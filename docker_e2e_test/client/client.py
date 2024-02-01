import os
import re
import copy
import log
import utility
from web3.auto import w3
from fabric import Connection
from eth_rpc_client import Client as RpcClient
from invoke import Responder


AUTONITY_PATH = "/home/{}/network-data/autonity"
GENESIS_PATH = "/home/{}/network-data/genesis.json"
CHAIN_DATA_DIR = "/home/{}/network-data/{}/data/"
BOOT_KEY_FILE = "/home/{}/network-data/{}/boot.key"
KEY_PASSPHRASE_FILE = "/home/{}/network-data/{}/pass.txt"
PACKAGE_NAME = "./network-data/{}.tgz"
REMOTE_NAME = "/home/{}/{}.tgz"
SYSTEM_SERVICE_DIR = "/etc/systemd/system/"
DEPLOYMENT_DIR = '/home/{}/network-data'
SYSTEMD_START_CLIENT = 'sudo systemctl start autonity.service'
SYSTEMD_STOP_CLIENT = 'sudo systemctl stop autonity.service'

# use ip tables module of linux kernel which is common for all linux distributions to control peer connection.
CONNECT_PEER = "sudo iptables -j DROP -D INPUT -s {}"
IS_PEER_ALREADY_DISCONNECTED = "sudo iptables -j DROP -C INPUT -s {}"
DISCONNECT_PEER = "sudo iptables -j DROP -A INPUT -s {}"

# traffic control of network, please refer to https://wiki.linuxfoundation.org/networking/netem#delay_distribution
# use traffic control module of linux kernel which is common for all linux distributions to control out-coming delay.
SSH_DELAY_TX_COMMAND = "sudo tc qdisc add dev {} root netem delay {}ms loss {}% duplicate {}% reorder {}% corrupt {}%"
SSH_UN_DELAY_TX_COMMAND = "sudo tc qdisc del dev {} root netem"

# due to network module design of linux kernel, controlling in-coming delay need a stream redirection.
# to create virtual ethernet interface.
SSH_CREATE_IFB_MODULE = "sudo modprobe ifb"
SSH_UP_IFB_DEVICE_INTERFACE = "sudo ip link set dev ifb0 up"
# SSH_DOWN_IFB_DEVICE_INTERFACE = "sudo ip link set dev ifb0 down"

# attach and detach extra in-coming qdisc stream to public ethernet interface to simulate in-coming delay.
# qdisc need to be recycled after the test case done.
SSH_ADD_INCOMING_QUEUE_4_PUB_INTERFACE = "sudo tc qdisc add dev {} ingress"
SSH_DELETE_INCOMING_QUEUE_4_PUB_INTERFACE = "sudo tc qdisc del dev {} ingress"

# redirect the public ethernet interface in-coming stream into virtual ethernet interface.
SSH_REDIRECT_STREAM_TO_IFB = \
    "sudo tc filter add dev {} parent ffff: protocol ip u32 match u32 0 0 flowid 1:1 action mirred egress redirect dev ifb0"

# set delay meta data for in-coming data streams.
SSH_DELAY_RX_COMMAND = "sudo tc qdisc add dev ifb0 root netem delay {}ms loss {}% duplicate {}% reorder {}% corrupt {}%"
SSH_UN_DELAY_RX_COMMAND = "sudo tc qdisc del dev ifb0 root netem"

DEFAULT_DELAY = 1  # 1ms for default delay.

DEFAULT_PACKAGE_LOSS_RATE = 0.1  # 0.1%
DEFAULT_PACKAGE_DUPLICATE_RATE = 0.1  # 0.1%
DEFAULT_PACKAGE_REORDER_RATE = 0.1  # 0.1%
DEFAULT_PACKAGE_CORRUPT_RATE = 0.1  # 0.1%


class Client(object):
    def __init__(self, host=None, p2p_port=None, acn_port=None, rpc_port=None, ws_port=None, net_interface=None,
                 coin_base=None, ssh_user=None, ssh_pass=None, ssh_key=None, sudo_pass=None, autonity_path=None,
                 bootnode_path=None, role=None, index=None, e_node=None):
        self.autonity_path = autonity_path
        self.bootnode_path = bootnode_path
        self.host = host
        self.p2p_port = p2p_port
        self.acn_port = acn_port
        self.rpc_port = rpc_port
        self.ws_port = ws_port
        self.net_interface = net_interface
        self.ssh_user = ssh_user
        self.ssh_pass = ssh_pass
        self.ssh_key = ssh_key
        self.sudo_pass = sudo_pass
        self.coin_base = coin_base
        self.role = role
        self.index = index
        self.e_node = e_node
        self.rpc_client = None
        self.logger = log.get_logger()
        self.disconnected_peers = []
        self.client_stopped = False
        self.up_link_delayed = False
        self.down_link_delayed = False
        self.is_local_address = False

    def create_work_dir(self, data_dir):
        work_dir = "{}/{}".format(data_dir, self.host)
        utility.create_dir(work_dir)

    def generate_new_account(self):
        folder = self.host
        utility.execute("echo 123 > ./network-data/{}/pass.txt".format(folder))
        output = utility.execute(
            '{} --datadir "./network-data/{}/data" --password "./network-data/{}/pass.txt" account new'
            .format(self.autonity_path, folder, folder)
        )
        self.logger.debug(output)
        m = re.findall(r'0x(.{40})', output[0], re.MULTILINE)
        if len(m) == 0:
            self.logger.error("Aborting - account creation failed")
            return None
        else:
            self.coin_base = m[0]
            return self.coin_base

    def generate_enode(self):
        folder = self.host

        keystores_dir = "./network-data/{}/data/keystore".format(folder)
        keystore_file_path = keystores_dir + "/" + os.listdir(keystores_dir)[0]
        with open(keystore_file_path) as keyfile:
            encrypted_key = keyfile.read()
            account_private_key = w3.eth.account.decrypt(encrypted_key, "123").hex()[2:]
        with open("./network-data/{}/boot.key".format(folder), "w") as bootkey:
            ## todo, generate autonity keys from cli, now it used a fix key for testing.
            autonity_keys = account_private_key + "3ab975b09167b550d25f8f0b31f4e3ebaf5ea73b3cc0eb1ca2c0957a5331de2d"
            bootkey.write(autonity_keys)

        pub_key = \
            utility.execute("{} -writeaddress -nodekey ./network-data/{}/boot.key".
                            format(self.bootnode_path, folder))[0].rstrip()
        #self.e_node = "enode://{}@{}:{}".format(pub_key, self.host, self.p2p_port)
        # new patern: "enode://pubKey:host:port?discPort=30303&acnep=host:port"
        self.e_node = "enode://{}@{}:{}?discPort={}&acnep={}:{}".format(pub_key, self.host, self.p2p_port, self.p2p_port, self.host, self.acn_port)
        return self.e_node

    def generate_system_service_file(self):
        template_remote = "[Unit]\n" \
                   "Description=Clearmatics Autonity Client server\n" \
                   "After=syslog.target network.target\n" \
                   "[Service]\n" \
                   "Type=simple\n" \
                   "ExecStart={} --genesis {} --datadir {} --autonitykeys {} --syncmode 'full' --port {} --consensus.port {} " \
                   "--http.port {} --http --http.addr '0.0.0.0' --ws --ws.port {} --http.corsdomain '*' "\
                   "--http.api 'personal,debug,db,eth,net,web3,txpool,miner,tendermint,clique' --networkid 1991  " \
                   "--allow-insecure-unlock --graphql " \
                   "--unlock 0x{} --password {} " \
                   "--mine --miner.threads '1' --verbosity 4 --miner.gaslimit 10000000000 \n" \
                   "KillMode=process\n" \
                   "KillSignal=SIGINT\n" \
                   "TimeoutStopSec=1\n" \
                   "Restart=on-failure\n" \
                   "RestartSec=1s\n" \
                   "[Install]\n" \
                   "Alias=autonity.service\n"\
                   "WantedBy=multi-user.target"

        folder = self.host

        print("prepare autonity systemd service file for node: %s", self.host)
        bin_path = AUTONITY_PATH.format(self.ssh_user)
        genesis_path = GENESIS_PATH.format(self.ssh_user)
        data_dir = CHAIN_DATA_DIR.format(self.ssh_user, folder)
        boot_key_file = BOOT_KEY_FILE.format(self.ssh_user, folder)
        p2p_port = self.p2p_port
        acn_port = self.acn_port
        rpc_port = self.rpc_port
        ws_port = self.ws_port
        coin_base = self.coin_base
        password_file = KEY_PASSPHRASE_FILE.format(self.ssh_user, folder)

        content = template_remote.format(bin_path, genesis_path, data_dir, boot_key_file, p2p_port, acn_port, rpc_port, ws_port,
                                         coin_base, password_file)
        with open("./network-data/{}/autonity.service".format(folder), 'w') as out:
            out.write(content)

    def generate_package(self):
        folder = self.host
        utility.execute('cp {} ./network-data/'.format(self.autonity_path))
        utility.execute('tar -zcvf ./network-data/{}.tgz ./network-data/{}/ ./network-data/genesis.json ./network-data/autonity'.format(folder, folder))

    def deliver_package(self):
        try:
            with Connection(self.host, user=self.ssh_user, connect_kwargs={
                #"key_filename": self.ssh_key,
                "password": self.ssh_pass
            }) as c:
                sudopass = Responder(
                    pattern=r'\[sudo\] password for ' + self.ssh_user + ':',
                    response=self.sudo_pass + '\n'
                )
                c.put(PACKAGE_NAME.format(self.host), REMOTE_NAME.format(self.ssh_user, self.host))
                self.logger.info('Chain package was uploaded to %s.', self.host)
                result = c.run('sudo tar -C /home/{} -zxvf {}'.format(self.ssh_user, REMOTE_NAME.format(self.ssh_user, self.host)), pty=True,
                               watchers=[sudopass], warn=True, hide=True)
                if result and result.exited == 0 and result.ok:
                    self.logger.info('Chain package was unpacked to %s.', self.host)
                    return True
                else:
                    self.logger.error('chain package fail to unpacked to %s.', self.host)
        except Exception as e:
            self.logger.error("cannot deliver package to host. %s, %s", self.host, e)
        return False

    def load_systemd_file(self):
        try:
            with Connection(self.host, user=self.ssh_user, connect_kwargs={
                #"key_filename": self.ssh_key,
                "password": self.ssh_pass
            }) as c:
                sudopass = Responder(
                    pattern=r'\[sudo\] password for ' + self.ssh_user + ':',
                    response=self.sudo_pass + '\n'
                )
                src = '/home/{}/network-data/{}/autonity.service'.format(self.ssh_user, self.host)
                result = c.run('sudo cp {} {}'.format(src, SYSTEM_SERVICE_DIR), pty=True, watchers=[sudopass],
                               warn=True, hide=True)
                if result and result.exited == 0 and result.ok:
                    self.logger.info('system service loaded. %s', self.host)
                else:
                    self.logger.error('systemd file loading failed. %s', self.host)
        except Exception as e:
            self.logger.error("cannot load systemd file. %s, %s", self.host, e)

    def start_client(self):
        try:
            with Connection(self.host, user=self.ssh_user, connect_kwargs={
                #"key_filename": self.ssh_key,
                "password": self.ssh_pass
            }) as c:
                sudopass = Responder(
                    pattern=r'\[sudo\] password for ' + self.ssh_user + ':',
                    response=self.sudo_pass + '\n'
                )
                #cmd = self.generate_start_cmd()
                cmd = SYSTEMD_START_CLIENT
                result = c.run(cmd, pty=True, watchers=[sudopass],
                               warn=True, hide=True)
                if result and result.exited == 0 and result.ok:
                    self.logger.info('system service started. %s', self.host)
                    self.client_stopped = False
                    return True
                else:
                    self.logger.error('systemd service starting failed. %s', self.host)
        except Exception as e:
            self.logger.error("cannot start service. %s, %s.", self.host, e)
        return False

    def deploy_client(self):
        self.deliver_package()
        self.load_systemd_file()

    def stop_client(self):
        try:
            with Connection(self.host, user=self.ssh_user, connect_kwargs={
                #"key_filename": self.ssh_key,
                "password": self.ssh_pass
            }) as c:
                sudopass = Responder(
                    pattern=r'\[sudo\] password for ' + self.ssh_user + ':',
                    response=self.sudo_pass + '\n'
                )
                #cmd = self.generate_stop_cmd()
                cmd = SYSTEMD_STOP_CLIENT
                result = c.run(cmd, pty=True, watchers=[sudopass],
                               warn=True, hide=True)
                if result and result.exited == 0 and result.ok:
                    self.logger.info('system service stopped. %s', self.host)
                    self.client_stopped = True
                    return True
                else:
                    self.logger.error('system service stopping failed. %s', self.host)
        except Exception as e:
            self.logger.error("cannot stop client, %s, %s", self.host, e)
        return False

    def clean_chain_data(self):
        try:
            with Connection(self.host, user=self.ssh_user, connect_kwargs={
                #"key_filename": self.ssh_key,
                "password": self.ssh_pass
            }) as c:
                sudopass = Responder(
                    pattern=r'\[sudo\] password for ' + self.ssh_user + ':',
                    response=self.sudo_pass + '\n'
                )
                result = c.run('sudo rm -rf {}'.format(DEPLOYMENT_DIR.format(self.ssh_user)), pty=True, watchers=[sudopass],
                               warn=True, hide=True)
                if result and result.exited == 0 and result.ok:
                    self.logger.info('chain data cleaned. %s', self.host)
                else:
                    self.logger.error('chain data cleaning failed. %s', self.host)
        except Exception as e:
            self.logger.error("cannot clean chain data. %s, %s.", self.host, e)
            return False
        return True

    def redirect_system_log(self, log_folder):
        try:
            zip_file = "{}/{}.tgz".format(log_folder, self.host)
            log_file = "{}/{}.log".format(log_folder, self.host)
            # untar file,
            utility.execute("tar -zxvf {} --directory {}".format(zip_file, log_folder))
            # read file and print into log file.
            self.logger.info("\t\t\t **** node_%s logs started from here. **** \n\n\n", self.host)
            with open(log_file, "r", encoding="utf-8") as fp:
                for _, line in enumerate(fp):
                    self.logger.info("NODE_%s_%s: %s", self.index, self.host, line.encode("utf-8"))
            # remove file.
            utility.execute("rm -f {}".format(log_file))
        except Exception as e:
            self.logger.error('Exception happens. %s', e)

    def download_system_log(self, log_folder):
        try:
            with Connection(self.host, user=self.ssh_user, connect_kwargs={
                #"key_filename": self.ssh_key,
                "password": self.ssh_pass,
            }) as c:
                sudopass = Responder(
                    pattern=r'\[sudo\] password for ' + self.ssh_user + ':',
                    response=self.sudo_pass + '\n'
                )
                # dump log from previous boot at remote node.
                file_name = "./{}.log".format(self.host)
                cmd = 'sudo journalctl -u autonity.service -b > {}'.format(file_name)
                result = c.run(cmd, pty=True, watchers=[sudopass], warn=True, hide=True)
                if result and result.exited == 0 and result.ok:
                    self.logger.info('log was dump on host: %s', self.host)
                    # tar logs for remote node.
                    tar_file = "./{}.log.tgz".format(self.host)
                    cmd = "sudo tar -zcvf {} {}".format(tar_file, file_name)
                    result = c.run(cmd, pty=True, watchers=[sudopass], warn=True, hide=True)
                    if result and result.exited == 0 and result.ok:
                        self.logger.info('log was zip on host: %s', self.host)
                        # download logs.
                        local_dir = log_folder
                        local_file = "{}/{}.tgz".format(local_dir, self.host)
                        c.get(tar_file, local=local_file)
                        self.logger.info('log files was saved to %s.', local_dir)
                    else:
                        self.logger.error('cannot zip log file at host: %s', self.host)
                else:
                    self.logger.error('Cannot dump logs from autonity.service. %s', self.host)

        except (KeyError, TypeError) as e:
            self.logger.error('wrong configuration file. %s', e)
        except Exception as e:
            self.logger.error('Exception happens. %s', e)

    def send_transaction(self, to=None, gas=None, gas_price=None, value=0, data=None):
        try:
            if self.rpc_client is None:
                self.rpc_client = RpcClient(host=self.host, port=self.rpc_port)
                self.rpc_client.session.headers.update({"Content-type": "application/json"})
            # send transaction
            tx_hash = self.rpc_client.send_transaction(_from="0x{}".format(self.coin_base), to=to,
                                                       gas=gas, value=value, data=data)
            return tx_hash
        except Exception as e:
            self.logger.warn("send tx failed due to exception: %s", e)
            return None

    def get_balance(self):
        try:
            if self.rpc_client is None:
                self.rpc_client = RpcClient(host=self.host, port=self.rpc_port)
                self.rpc_client.session.headers.update({"Content-type": "application/json"})
            balance = self.rpc_client.get_balance("0x{}".format(self.coin_base))
            return balance
        except Exception as e:
            self.logger.error("Cannot get balance due to exception. %s", e)
            return None

    def get_block_hash_by_height(self, height):
        try:
            if self.rpc_client is None:
                self.rpc_client = RpcClient(host=self.host, port=self.rpc_port)
                self.rpc_client.session.headers.update({"Content-type": "application/json"})
            block = self.rpc_client.get_block_by_number(height)
            if block is None:
                self.logger.error("Cannot find block with height: %d at host: %s", height, self.host)
                return None
            return block["hash"]
        except IOError as e:
            self.logger.error("Cannot access RPC API from remote. %s", e)
            return None
        except Exception as e:
            self.logger.error("Exception happens: %s", e)
            return None

    def get_transaction_by_hash(self, hash):
        try:
            if self.rpc_client is None:
                self.rpc_client = RpcClient(host=self.host, port=self.rpc_port)
                self.rpc_client.session.headers.update({"Content-type": "application/json"})
            result = self.rpc_client.get_transaction_by_hash(hash)
            return result
        except Exception as e:
            self.logger.error("Cannot get balance due to exception. %s", e)
            return None

    def get_chain_height(self):
        try:
            if self.rpc_client is None:
                self.rpc_client = RpcClient(host=self.host, port=self.rpc_port)
                self.rpc_client.session.headers.update({"Content-type": "application/json"})
            height = self.rpc_client.get_block_number()
            self.logger.debug("get height: %d, %s.", height, self.host)
            return height
        except Exception as e:
            self.logger.debug("cannot get chain height yet, please retry. %s", e)
            return None

    def execute_ssh_cmd(self, cmd):
        self.logger.debug("ssh cmd: %s ", cmd)
        try:
            with Connection(self.host, user=self.ssh_user, connect_kwargs={
                #"key_filename": self.ssh_key,
                "password": self.ssh_pass,
            }) as c:
                sudopass = Responder(
                    pattern=r'\[sudo\] password for ' + self.ssh_user + ':',
                    response=self.sudo_pass + '\n'
                )
                result = c.run(cmd, pty=True, watchers=[sudopass], warn=True, hide=True)
                if result and result.exited == 0 and result.ok:
                    self.logger.debug('SSH executed fine. %s for node: %s', cmd, self.host)
                    return True
                else:
                    self.logger.debug('%s', result)
                    return result.stdout.strip()
        except IOError as e:
            self.logger.error("Cannot connect to node: %s via network. %s", self.host, e)
            return None
        except ValueError as e:
            self.logger.error("Wrong parameter for fabric lib. %s", e)
            return None
        except Exception as e:
            self.logger.error('Command get exception: %s', e)
            return None

    def is_peer_disconnected(self, host, port):
        cmd = IS_PEER_ALREADY_DISCONNECTED.format(host)
        is_disconnected = self.execute_ssh_cmd(cmd)
        if is_disconnected is True:
            return True
        return False

    def connect_peer(self, host, port):
        if self.is_peer_disconnected(host, port) is False:
            return True
        # connect_peer
        cmd = CONNECT_PEER.format(host)
        if self.execute_ssh_cmd(cmd) is not True:
            return False

        peer = "{}:{}".format(host, port)
        if peer in self.disconnected_peers:
            self.logger.info("peer connected: %s to %s", self.host, host)
            self.disconnected_peers.remove(peer)
        return True

    def dis_connect_peer(self, host, port):
        if self.is_peer_disconnected(host, port):
            return True
        cmd = DISCONNECT_PEER.format(host)
        if self.execute_ssh_cmd(cmd) is not True:
            return False
        peer = "{}:{}".format(host, port)
        if peer not in self.disconnected_peers:
            self.logger.info("peer disconnected: %s to %s", self.host, host)
            self.disconnected_peers.append(peer)
        return True

    def set_up_link_delay(self, up_link_delay_meta):
        try:
            delay = \
                DEFAULT_DELAY if up_link_delay_meta['delay'] is None else up_link_delay_meta['delay']
            loss_rate = \
                DEFAULT_PACKAGE_LOSS_RATE if up_link_delay_meta['lossRate'] is None else up_link_delay_meta['lossRate']
            duplicate_rate = \
                DEFAULT_PACKAGE_DUPLICATE_RATE \
                if up_link_delay_meta['duplicateRate'] is None else up_link_delay_meta['duplicateRate']
            reorder_rate = \
                DEFAULT_PACKAGE_REORDER_RATE \
                if up_link_delay_meta['reorderRate'] is None else up_link_delay_meta['reorderRate']
            corrupt_rate = \
                DEFAULT_PACKAGE_CORRUPT_RATE \
                if up_link_delay_meta['corruptRate'] is None else up_link_delay_meta['corruptRate']

            # to do checking parameters before formatting command.
            delay = delay if isinstance(delay, (int, float)) and not isinstance(delay, bool) else DEFAULT_DELAY
            loss_rate = loss_rate \
                if isinstance(loss_rate, (int, float)) and not isinstance(loss_rate, bool) \
                else DEFAULT_PACKAGE_LOSS_RATE
            duplicate_rate = duplicate_rate if isinstance(duplicate_rate, (int, float)) \
                and not isinstance(duplicate_rate, bool) else DEFAULT_PACKAGE_DUPLICATE_RATE
            reorder_rate = reorder_rate \
                if isinstance(reorder_rate, (int, float)) and not isinstance(reorder_rate, bool)\
                else DEFAULT_PACKAGE_REORDER_RATE
            corrupt_rate = corrupt_rate \
                if isinstance(corrupt_rate, (int, float)) and not isinstance(corrupt_rate, bool)\
                else DEFAULT_PACKAGE_CORRUPT_RATE

            ether_id = self.net_interface
            if ether_id is None:
                self.logger.error('Cannot find host ethernet interface id.')
                return None
            # to do formatting shell command.
            command = SSH_DELAY_TX_COMMAND.format(ether_id, delay, loss_rate, duplicate_rate, reorder_rate, corrupt_rate)
            result = self.execute_ssh_cmd(command)
            if result is True:
                self.up_link_delayed = True
                self.logger.info("network traffic control rule applied to host: %s", self.host)
            else:
                self.logger.error('SSH ERR: %s', result)
                return None
        except (KeyError, TypeError) as e:
            self.logger.error('Wrong configuration. %s', e)
            return None
        except Exception as e:
            self.logger.error('Exception happens: %s', e)
            return None
        return True

    def cancel_up_link_delay(self):
        try:
            ether_id = self.net_interface
            if ether_id is None:
                self.logger.error('Cannot find host ethernet interface id.')
                return None
            command = SSH_UN_DELAY_TX_COMMAND.format(ether_id)
            result = self.execute_ssh_cmd(command)
            if result is True:
                self.logger.info("network traffic control rule canceled to host: %s", self.host)
                self.up_link_delayed = False
            else:
                self.logger.error('undelay up-link failed for host: %s, error: %s', self.host, result)
                return None

        except (KeyError, TypeError) as e:
            self.logger.error('wrong configuration. %s', e)
            return None
        except Exception as e:
            self.logger.error('Exception happens: %s', e)
            return None
        self.logger.info('SSH command execute fine. %s', command)
        return True

    def set_down_link_delay(self, down_link_delay_meta):
        # since docker is a LXC base container solution, it does not involved the kernel module
        # which our TC (traffic control) depends on to simulate the traffic delays, so
        # we remove this from docker based testbed, while in VM based testing, TC is
        # still valid.
        return True
        try:
            delay = \
                DEFAULT_DELAY if down_link_delay_meta['delay'] is None else down_link_delay_meta['delay']
            loss_rate = \
                DEFAULT_PACKAGE_LOSS_RATE if down_link_delay_meta['lossRate'] \
                is None else down_link_delay_meta['lossRate']
            duplicate_rate = \
                DEFAULT_PACKAGE_DUPLICATE_RATE \
                if down_link_delay_meta['duplicateRate'] is None else down_link_delay_meta['duplicateRate']
            reorder_rate = \
                DEFAULT_PACKAGE_REORDER_RATE \
                if down_link_delay_meta['reorderRate'] is None else down_link_delay_meta['reorderRate']
            corrupt_rate = \
                DEFAULT_PACKAGE_CORRUPT_RATE \
                if down_link_delay_meta['corruptRate'] is None else down_link_delay_meta['corruptRate']

            # to do checking parameters before formatting command.
            delay = delay if isinstance(delay, (int, float)) and not isinstance(delay, bool) else DEFAULT_DELAY
            loss_rate = loss_rate \
                if isinstance(loss_rate, (int, float)) and not isinstance(loss_rate, bool) \
                else DEFAULT_PACKAGE_LOSS_RATE
            duplicate_rate = duplicate_rate if isinstance(duplicate_rate, (int, float)) \
                and not isinstance(duplicate_rate, bool) else DEFAULT_PACKAGE_DUPLICATE_RATE
            reorder_rate = reorder_rate \
                if isinstance(reorder_rate, (int, float)) and not isinstance(reorder_rate, bool)\
                else DEFAULT_PACKAGE_REORDER_RATE
            corrupt_rate = corrupt_rate \
                if isinstance(corrupt_rate, (int, float)) and not isinstance(corrupt_rate, bool)\
                else DEFAULT_PACKAGE_CORRUPT_RATE

            # get ip from host name.
            ether_id = self.net_interface
            if ether_id is None:
                self.logger.error('cannot get ethernet interface id')
                return None
            # step 1, create virtual network interface.
            result = \
                self.execute_ssh_cmd(SSH_CREATE_IFB_MODULE)
            if result is not True:
                self.logger.error('create virtual network interface failed. %s', result)
                return None
            # step 2, start up virtual network interface.
            result = self.execute_ssh_cmd(SSH_UP_IFB_DEVICE_INTERFACE)
            if result is not True:
                self.logger.error('start up virtual network interface failed. %s', result)
                return None
            # step 3, add in-coming queue to network interface.
            command = SSH_ADD_INCOMING_QUEUE_4_PUB_INTERFACE.format(ether_id)
            result = self.execute_ssh_cmd(command)
            if result is not True:
                self.logger.error('create redirect stream failed. %s', result)
                return None

            # step 4, redirect in-coming stream into virtual network interface.
            command = SSH_REDIRECT_STREAM_TO_IFB.format(ether_id)
            result = self.execute_ssh_cmd(command)
            if result is not True:
                self.logger.error('redirect in-coming stream into virtual network interface failed. %s', result)
                return None

            # step 5, apply delay in in-coming stream.
            command = SSH_DELAY_RX_COMMAND.format(delay, loss_rate, duplicate_rate, reorder_rate, corrupt_rate)
            result = self.execute_ssh_cmd(command)
            if result is True:
                self.down_link_delayed = True
            else:
                self.logger.error('SSH ERR: %s', result)
                return None
        except (KeyError, TypeError) as e:
            self.logger.error('Wrong configuration. %s', e)
            return None
        except Exception as e:
            self.logger.error('Exception happens: %s', e)
            return None
        return True

    def cancel_down_link_delay(self):
        # since docker is a LXC base container solution, it does not involved the kernel module
        # which our TC (traffic control) depends on to simulate the traffic delays, so
        # we remove this from docker based testbed, while in VM based testing, TC is
        # still valid.
        return True
        try:
            ether_id = self.net_interface
            if ether_id is None:
                self.logger.error('Cannot find host ethernet interface id.')
                return None
            # step 1, un-apply delay on in-coming data stream.
            result = self.execute_ssh_cmd(SSH_UN_DELAY_RX_COMMAND)
            if result is not True:
                self.logger.error('un-apply delay on in-coming data stream failed. %s', result)

            # step 2, recycle the redirection stream.
            command = SSH_DELETE_INCOMING_QUEUE_4_PUB_INTERFACE.format(ether_id)
            result = self.execute_ssh_cmd(command)
            if result is True:
                self.down_link_delayed = False
            else:
                self.logger.error('undelay up-link failed for host: %s, error: %s', self.host, result)
                return None

        except (KeyError, TypeError) as e:
            self.logger.error('wrong configuration. %s', e)
            return None
        except Exception as e:
            self.logger.error('Exception happens: %s', e)
            return None
        self.logger.info('SSH command executed fine. %s', command)
        return True

    def heal_from_disaster(self):
        failed = False
        if self.client_stopped:
            if self.start_client() is not True:
                failed = True

        disconnected_peers = copy.deepcopy(self.disconnected_peers)
        if len(disconnected_peers) > 0:
            for peer in disconnected_peers:
                endpoint = peer.split(":")
                if self.connect_peer(endpoint[0], endpoint[1]) is not True:
                    failed = True

        if self.up_link_delayed:
            if self.cancel_up_link_delay() is not True:
                failed = True
        if self.down_link_delayed:
            if self.cancel_down_link_delay() is not True:
                failed = True
        return True if not failed else False
