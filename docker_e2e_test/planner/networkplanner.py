import copy
import json
import utility
import log
from client.client import Client
from conf import conf


class NetworkPlanner(object):
    def __init__(self, autonity_path):
        self.logger = log.get_logger()
        self.autonity_path = autonity_path
        path_list = autonity_path.split("/")
        path_list[len(path_list) - 1] = "bootnode"
        self.bootnode_path = "/".join(path_list)
        path_list[len(path_list) - 1] = "ethkey"
        self.key_inspector_path = "/".join(path_list)
        self.validator_ip_list = []
        self.participant_ip_list = []
        self.clients = []

    def get_clients(self):
        return self.clients

    def prepare_network_ips(self):
        engine_conf = conf.get_engine_conf()
        if engine_conf is None:
            return
        if engine_conf["local_mode"]:
            for i in range(0, engine_conf["default_scalability"]):
                # use loop-back addresses for local different clients.
                self.validator_ip_list.append("127.0.0.{}".format(i+1))
            return
        if engine_conf["local_mode"] is False:
            self.validator_ip_list, self.participant_ip_list = conf.get_client_ips()

    def prepare_client_instances(self):
        for index, ip in enumerate(self.validator_ip_list):
            self.clients.append(Client(ip, role="validator", autonity_path=self.autonity_path,
                                       bootnode_path=self.bootnode_path, key_inspector_path=self.key_inspector_path, index=index))
        for index, ip in enumerate(self.participant_ip_list):
            self.clients.append(Client(ip, role="participant", autonity_path=self.autonity_path,
                                       bootnode_path=self.bootnode_path, key_inspector_path=self.key_inspector_path, index=index+len(self.validator_ip_list)))

    def create_work_dir(self):
        self.logger.info("===== SETUP INITIALIZATION =====")
        engine_conf = conf.get_engine_conf()
        if engine_conf is None:
            return
        data_dir = engine_conf["network_data_dir"]
        utility.remove_dir(data_dir)
        utility.create_dir(data_dir)

        for client in self.clients:
            client.create_work_dir(data_dir)

    def generate_accounts(self):
        self.logger.info("===== ACCOUNTS CREATION =====")
        accounts = []
        for client in self.clients:
            account = client.generate_new_account()
            if account:
                accounts.append(account)
        self.logger.info(accounts)

    def generate_testbed_conf(self):
        self.logger.info("===== TEST BED CONF GENERATION =====")
        template = conf.get_testbed_template()
        try:
            node_template = template["targetNetwork"]["nodes"].pop(0)
            node = copy.deepcopy(node_template)
            nodes_to_apply = []

            for index, client in enumerate(self.clients):
                # sync template data to client instance.
                client.p2p_port = node_template["p2pPort"]
                client.acn_port = node_template["acnPort"]
                client.rpc_port = node_template["rpcPort"]
                client.ws_port = node_template["wsPort"]
                client.net_interface = node_template["ethernetInterfaceID"]
                client.ssh_key = node_template["sshCredential"]["sshKey"]
                client.ssh_pass = node_template["sshCredential"]["sshPass"]
                client.ssh_user = node_template["sshCredential"]["sshUser"]
                client.sudo_pass = node_template["sshCredential"]["sudoPass"]
                client.index = index

                # sync template data to testbedconf.yml
                node["index"] = client.index
                node["role"] = client.role
                node["name"] = client.host
                node["coinBase"] = "0x{}".format(client.coin_base)
                node["p2pPort"] = client.p2p_port
                node["acnPort"] = client.acn_port
                node["rpcPort"] = client.rpc_port
                node["wsPort"] = client.ws_port
                node["ethernetInterfaceID"] = client.net_interface
                node["sshCredential"]["sshKey"] = client.ssh_key
                node["sshCredential"]["sshPass"] = client.ssh_pass
                node["sshCredential"]["sshUser"] = client.ssh_user
                node["sshCredential"]["sudoPass"] = client.sudo_pass

                nodes_to_apply.append(copy.deepcopy(node))

            template["targetNetwork"]["nodes"] = nodes_to_apply
            if conf.dump_test_bed_conf(template):
                self.logger.info("test bed conf file was generated.")
                return True
            else:
                self.logger.error("cannot generate test bed conf file.")
        except (IOError, OSError) as e:
            self.logger.error("cannot generate test bed conf %s", e)
        except (KeyError, ValueError) as e:
            self.logger.error("Wrong configuration from template file. %s", e)

    def generate_enodes(self):
        self.logger.info("===== ENODES GENERATION =====")
        enodes = []
        for client in self.clients:
            enode = client.generate_enode()
            if enode:
                enodes.append(enode)
        self.logger.info(enodes)

    def generate_genesis(self):
        self.logger.info("===== GENESIS GENERATION =====")
        #   The following parameters should not be modified unless you know what you're doing.   #
        genesis = {
            "config": {
                "chainId": 1,
                "autonity": {
                    "minBaseFee": 5000,
                    "delegationRate": 1000,
                    "blockPeriod": 1,
                    "maxCommitteeSize": 7,
                    "unbondingPeriod": 120,
                    "epochPeriod": 30,
                    "treasuryFee": 150000000,
                    "validators": [],
                },
                "oracle": {"votePeriod": 10},
                "accountability": {
                    "innocenceProofSubmissionWindow": 30,
                    "baseSlashingRateLow": 500,
                    "baseSlashingRateMid": 1000,
                    "collusionFactor": 550,
                    "historyFactor": 750,
                    "jailFactor": 60,
                    "slashingRatePrecision": 10000
                },
            },
            "nonce": "0x0",
            "timestamp": "0x0",
            "baseFee": "15000000000",
            "gasLimit": "10000000000",
            "difficulty": "0x0",
            "coinbase": "0x0000000000000000000000000000000000000000",
            "number": "0x0",
            "gasUsed": "0x0",
            "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "mixHash": "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
        }

        # Default balance
        starting_balance = "0x000000000000000000100000000000000000000000000000000000000000000"
        genesis["alloc"] = {}
        genesis["config"]["autonity"]["operator"] = "0x{}".format(self.clients[0].coin_base)
        genesis["config"]["autonity"]["treasury"] = "0x{}".format(self.clients[0].coin_base)

        for index, client in enumerate(self.clients):
            coinbase = "0x{}".format(client.coin_base)
            validator = {
                "treasury": coinbase,
                "enode": client.e_node,
                "bondedStake": 10000  if client.role  == "validator" else 5000,
                "oracleAddress": coinbase,
                "consensusKey": client.consensus_pub_key,
            }
            genesis["alloc"][coinbase] = {"balance": starting_balance}
            genesis["config"]["autonity"]["validators"].append(validator)

        with open("./network-data/genesis.json", 'w') as out:
            out.write(json.dumps(genesis, indent=4) + '\n')

    def generate_systemd_service_file(self):
        self.logger.info("===== SYSTEMD SERVICE FILE GENERATION =====")
        for client in self.clients:
            client.generate_system_service_file()

    def generate_package(self):
        self.logger.info("===== PACKAGE GENERATION =====")
        for client in self.clients:
            client.generate_package()

    def plan(self):
        self.prepare_network_ips()
        self.prepare_client_instances()
        self.create_work_dir()
        self.generate_accounts()
        self.generate_testbed_conf()
        self.generate_enodes()
        self.generate_genesis()
        self.generate_systemd_service_file()
        self.generate_package()
        self.logger.info("===== SETUP FINISHED =====")

    def deploy(self):
        for client in self.clients:
            client.deploy_client()

    def stop_all_nodes(self):
        for client in self.clients:
            client.stop_client()

    def start_all_nodes(self):
        for client in self.clients:
            if client.start_client() is not True:
                return False
        return True
