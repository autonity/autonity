from client.client import Client, AUTONITY_PATH, GENESIS_PATH, CHAIN_DATA_DIR, BOOT_KEY_FILE, KEY_PASSPHRASE_FILE
import utility


class TwinsClient(Client):
    def __init__(self, host=None, p2p_port=None, rpc_port=None, ws_port=None, net_interface=None,
                 coin_base=None, ssh_user=None, ssh_pass=None, ssh_key=None, sudo_pass=None, autonity_path=None,
                 bootnode_path=None, role=None, index=None, e_node=None, twins_client=None, hostname=None):
        super().__init__(host=host, p2p_port=p2p_port, rpc_port=rpc_port, ws_port=ws_port,
                         net_interface=net_interface, coin_base=coin_base, ssh_user=ssh_user, ssh_pass=ssh_pass,
                         ssh_key=ssh_key, sudo_pass=sudo_pass, autonity_path=autonity_path, bootnode_path=bootnode_path,
                         role=role, index=index, e_node=e_node, hostname=hostname)
        self.refer_client = twins_client

    def generate_new_account(self):
        self.logger.info("generate new account in twins client.")
        folder = self.host
        utility.execute("echo 123 > ./network-data/{}/pass.txt".format(folder))

        utility.execute("mkdir ./network-data/{}/data".format(self.host))
        # copy twins' keystore.
        utility.execute("cp -rf ./network-data/{}/data/ ./network-data/{}/".format(self.refer_client.host, self.host))
        self.coin_base = self.refer_client.coin_base
        return self.coin_base
