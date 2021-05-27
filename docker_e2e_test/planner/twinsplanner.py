import json
from planner.networkplanner import NetworkPlanner
from client.client import Client
from client.twins_client import TwinsClient

"""
    Twins network planer, it takes autonity binary path, and the list of node index which is set as twins for each other.
    The steps of twins network planer is almost the same as the normal network planner. It get the IPs from the 
    configuration, setup the rpc clients for each nodes, create genesis block data for the chain, and upload chain settings
    for each nodes.
"""


def take_first(elem):
    return elem[0]


class TwinsNetworkPlanner(NetworkPlanner):
    def __init__(self, autonity_path=None, twins=None):
        super().__init__(autonity_path=autonity_path)
        # twins is a list of list, for example: [[0, 1], [2, 3], [4, 5]], we 3 pair of twins in the network.
        self.twins = []
        for _, l in enumerate(twins):
            l.sort()
            self.twins.append(l)

        self.twins.sort(key=take_first)

    def refer_to(self, index):
        for _, twins in enumerate(self.twins):
            if index == twins[0]:
                return None
            if index == twins[1]:
                return twins[0]

    def prepare_client_instances(self):
        self.logger.info("Prepare client instance in Twins network planner")
        for index, ip in enumerate(self.validator_ip_list):
            ref = self.refer_to(index)
            if ref is None:
                client = Client(host=ip, role="validator", autonity_path=self.autonity_path, bootnode_path=self.bootnode_path, index=index, hostname="client{}".format(index))
                self.clients.append(client)
            else:
                twins_client = TwinsClient(host=ip, role="validator", autonity_path=self.autonity_path, bootnode_path=self.bootnode_path, index=index, hostname="client{}".format(index), twins_client=self.clients[ref])
                self.clients.append(twins_client)

    def generate_genesis(self):
        self.logger.info("===== GENESIS GENERATION =====")
        #   The following parameters should not be modified unless you know what you're doing.   #
        genesis = {
            "config": {
                "homesteadBlock": 0,
                "eip150Block": 0,
                "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "eip155Block": 0,
                "eip158Block": 0,
                "byzantiumBlock": 0,
                "istanbulBlock": 0,
                "constantinopleBlock": 0,
                "petersburgBlock": 0,
                "tendermint": {
                    "policy": 1,
                    "block-period": 1,
                },
                "autonityContract": {
                    "bytecode": "",
                    "abi": "",
                    "minGasPrice": 5000,
                    "users": [],
                }
            },
            "nonce": "0x0",
            "timestamp": "0x0",
            "gasLimit": "10000000000",
            "difficulty": "0x1",
            "coinbase": "0x0000000000000000000000000000000000000000",
            "number": "0x0",
            "gasUsed": "0x0",
            "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "mixHash": "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
        }

        # Default balance
        starting_balance = "0x000000000000000000100000000000000000000000000000000000000000000"
        genesis["alloc"] = {}
        genesis["validators"] = []
        genesis["config"]["autonityContract"]["operator"] = "0x{}".format(self.clients[0].coin_base)
        genesis["config"]["autonityContract"]["deployer"] = "0x{}".format(self.clients[0].coin_base)
        genesis["config"]["chainId"] = 1  # can be assigned freely

        for index, client in enumerate(self.clients):
            ref = self.refer_to(index)
            if ref is None:
                coinbase = "0x{}".format(client.coin_base)
                user = {
                    "enode": client.e_node,
                    "address": coinbase,
                    "type": client.role,
                    "stake": 2 if client.role == "validator" else 1,
                }
                genesis["alloc"][coinbase] = {"balance": starting_balance}
                genesis["config"]["autonityContract"]["users"].append(user)

        with open("./network-data/genesis.json", 'w') as out:
            out.write(json.dumps(genesis, indent=4) + '\n')

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

    def set_ring_topology(self):
        self.clients[0].dis_connect_peer(self.clients[2].host, self.clients[2].p2p_port)
        self.clients[0].dis_connect_peer(self.clients[3].host, self.clients[3].p2p_port)
        self.clients[0].dis_connect_peer(self.clients[4].host, self.clients[4].p2p_port)

        self.clients[1].dis_connect_peer(self.clients[3].host, self.clients[3].p2p_port)
        self.clients[1].dis_connect_peer(self.clients[4].host, self.clients[4].p2p_port)
        self.clients[1].dis_connect_peer(self.clients[5].host, self.clients[5].p2p_port)

        self.clients[2].dis_connect_peer(self.clients[4].host, self.clients[4].p2p_port)
        self.clients[2].dis_connect_peer(self.clients[5].host, self.clients[5].p2p_port)

        self.clients[3].dis_connect_peer(self.clients[5].host, self.clients[5].p2p_port)
