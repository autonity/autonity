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
        # clean remote chain data
        # re-generate local systemd file with fault simulator flag
        # tar package
        # deliver package
        # run network
        # check if on-chain proof is presented with timeout
        pass
