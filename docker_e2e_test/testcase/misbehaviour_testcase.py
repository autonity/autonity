import log
import os
import time
from timeit import default_timer as timer

TIME_OUT = 60 * 30  # 30 minutes
SYSTEM_LOG_DIR = './system_log/'
TEST_CASE_SYSTEM_LOG_DIR = './system_log/{}'


class MisbehaviourTestCase:
    def __init__(self, test_case_conf, clients):
        self.start_time = time.time()
        self.test_case_conf = test_case_conf
        self.logger = log.get_logger()
        self.clients = {}
        for client in clients:
            self.clients[client.index] = client

    def clean_up_test_bed(self):
        # stop current running network
        for index, client in self.clients.items():
            if client.stop_client() is not True:
                return False

        # clean remote chain data
        for index, client in self.clients.items():
            if client.clean_chain_data() is not True:
                return False
        return True

    def generate_system_service_file(self):
        try:
            flag = self.test_case_conf["flag"]
            rule_id = self.test_case_conf["value"]
        except Exception as e:
            self.logger.error("re-generate local systemd file with fault simulator flag failed: %s", e)
            return False
        # re-generate local systemd file with fault simulator flag
        for index, client in self.clients.items():
            client.generate_system_service_file_with_fault_simulator(flag=flag, rule_id=rule_id)
        return True

    def setup_test_bed(self):
        # tar package
        for index, client in self.clients.items():
            client.generate_package()

        # deliver package
        for index, client in self.clients.items():
            if client.deliver_package() is not True:
                return False

        # load autonity systemd service file
        for index, client in self.clients.items():
            client.load_systemd_file()

        # run network
        for index, client in self.clients.items():
            if client.start_client() is not True:
                return False
        return True

    def verify_test(self):
        try:
            name = self.test_case_conf["name"]
            flag = self.test_case_conf["flag"]
            rule_id = self.test_case_conf["value"]
        except Exception as e:
            self.logger.error("re-generate local systemd file with fault simulator flag failed: %s", e)
            return False
        # check if on-chain proof is presented with timeout
        start = timer()
        while (timer() - start) < TIME_OUT:
            for index, client in self.clients.items():
                height = client.get_chain_height()
                if height is not None:
                    self.logger.debug("chain height: %d, node: %s", height, client.host)
                if client.is_proof_presented(flag, rule_id) is True:
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

    def run(self):
        if self.clean_up_test_bed() is not True:
            return False

        if self.generate_system_service_file() is not True:
            return False

        if self.setup_test_bed() is not True:
            return False

        return self.verify_test()


class TwinsTestCase(MisbehaviourTestCase):
    def __init__(self, test_case_conf, clients, twins):
        super().__init__(test_case_conf, clients)
        self.logger.info("twins test case is created: %s", twins)
        self.twins = []
        for _, l in enumerate(twins):
            l.sort()
            self.twins.append(l)

    def refer_to(self, index):
        for _, twins in enumerate(self.twins):
            if index == twins[0]:
                return None
            if index == twins[1]:
                return twins[0]

    def generate_system_service_file(self):
        # set fault simulator config for twins node only.
        try:
            flag = self.test_case_conf["flag"]
            rule_id = self.test_case_conf["value"]
        except Exception as e:
            self.logger.error("re-generate local systemd file with fault simulator flag failed: %s", e)
            return False
        # re-generate local systemd file with fault simulator flag
        for index, client in self.clients.items():
            ref = self.refer_to(index)
            if ref is None:
                client.generate_system_service_file()
            else:
                client.generate_system_service_file_with_fault_simulator(flag=flag, rule_id=rule_id)
                self.logger.info("system service file generated for twins client: node_%s", client.host)
        return True

    def run(self):
        if self.clean_up_test_bed() is not True:
            return False

        if self.generate_system_service_file() is not True:
            return False

        if self.setup_test_bed() is not True:
            return False

        return self.verify_test()
