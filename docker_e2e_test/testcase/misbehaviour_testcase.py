import log
import os
from typing import List
from client.client import Client
import time
from timeit import default_timer as timer

TIME_OUT = 60 * 5 # 5 minutes
SYSTEM_LOG_DIR = './system_log/'
TEST_CASE_SYSTEM_LOG_DIR = './system_log/{}'

class MisbehaviourTestCase:
    def __init__(self, test_case_conf, clients: List[Client]):
        self.start_time = time.time()
        self.test_case_conf = test_case_conf
        self.logger = log.get_logger()
        self.clients = {}
        for client in clients:
            self.clients[client.index] = client

    def run(self):
        try:
            name = self.test_case_conf["name"]
            flag = self.test_case_conf["flag"]
            id = self.test_case_conf["value"]
        except Exception as e:
            self.logger.error("re-generate local systemd file with fault simulator flag failed: %s", e)
            return False

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
        start = timer()
        while (timer() - start) < TIME_OUT:
            for index, client in self.clients.items():
                if client.is_proof_presented(flag, id) is True:
                    self.logger.info("misbehaviour test case passed: %s", name)
                    return True
            time.sleep(1)

        self.logger.error("misbehaviour test case timeout: %s", name)
        # start collect system logs
        try:
            # try to create dirs
            os.makedirs(SYSTEM_LOG_DIR, exist_ok=True)  # It never fail even if the dir is existed.
            os.makedirs(TEST_CASE_SYSTEM_LOG_DIR.format(self.start_time), exist_ok=True)
            for index, client in self.clients.items():
                client.collect_system_log(TEST_CASE_SYSTEM_LOG_DIR.format(self.start_time))
        except Exception as e:
            self.logger.error('Cannot fetch logs from node. %s.', e)
            return False
        try:
            # redirect client logs into test engine's logger.
            for index, client in self.clients.items():
                client.redirect_system_log(TEST_CASE_SYSTEM_LOG_DIR.format(self.start_time))
        except Exception as e:
            self.logger.error('Cannot redirect system logs from client into test engine log file %s.', e)
            return False
        return False
