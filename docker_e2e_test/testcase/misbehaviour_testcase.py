import log
from typing import List
from client.client import Client

class MisbehaviourTestCase:
    def __init__(self, test_case_conf, clients: List[Client]):
        self.test_case_conf = test_case_conf
        self.logger = log.get_logger()
        self.clients = {}
        for client in clients:
            self.clients[client.index] = client

    def run(self):
        # stop current running network
        for index, client in self.clients.items():
            if client.stop_client() is not True:
                return False

        # clean remote chain data
        for index, client in self.clients.items():
            if client.clean_chain_data() is not True:
                return False

        # re-generate local systemd file with fault simulator flag
        for index, client in self.clients.items():
            client.generate_system_service_file_with_fault_simulator(flag=flag, rule_id=id)

        # tar package
        for index, client in self.clients.items():
            client.generate_package()

        # deliver package
        for index, client in self.clients.items():
            if client.deliver_package() is not True:
                return False

        # run network
        for index, client in self.clients.items():
            if client.start_client() is not True:
                return False

        # check if on-chain proof is presented with timeout
        pass
