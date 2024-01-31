import os
from typing import List
from scheduler.scheduler import Scheduler
from conf import conf
import threading
import log
import time
from client.client import Client
from timeit import default_timer as timer

HEAL_TIME_OUT = 60 * 5  # 5 minutes
TX_HISTORY_FILE = './TXs_Per_TC_{}'
BLOCK_CONSISTENT_CHECKING_DURATION = 2
ENGINE_STATE_CHECKING_DURATION = 60
TEST_CASE_CONTEXT_FILE_NAME = './system_log/failed_{}_context/test_case_context.log'
SYSTEM_LOG_DIR = './system_log/'
TEST_CASE_SYSTEM_LOG_DIR = './system_log/testcase_{}_{}'


class TestCase:
    """A TestCase define the meta data of test case, includes condition, input, output."""

    def __init__(self, test_case_conf, clients: List[Client]):
        self.passed = False
        self.test_case_conf = test_case_conf
        self.id = test_case_conf['name'].split(':')[0]
        self.logger = log.get_logger()
        self.clients = {}
        for client in clients:
            self.clients[client.index] = client
        self.scheduler = Scheduler(test_case_conf, clients)
        # TC statistics
        self.start_chain_height = 0
        self.end_chain_height_before_recover = 0
        self.end_chain_height_after_recover = 0
        self.start_time = time.time()
        self.start_recover_time = None
        self.end_recover_time = None

        self.tx_history_file = TX_HISTORY_FILE.format(self.start_time)
        # TX issuing statistics
        self.tx_start_chain_height = 0
        self.tx_end_chain_height = 0
        self.tx_start_time = None
        self.tx_end_time = None
        self.tx_sent = 0
        self.tx_mined = 0
        self.balance_mined_by_the_test = 0
        self.sender_before_balance = 0
        self.receiver_before_balance = 0

    def __del__(self):
        try:
            if os.path.exists(self.tx_history_file):
                os.remove(self.tx_history_file)
        except Exception as e:
            self.logger.warning('remove file failed %s', e)

    def get_chain_height(self):
        best_height = 0
        for index, client in self.clients.items():
            height = client.get_chain_height()
            if height is None:
                continue
            else:
                best_height = height if height > best_height else best_height
        return best_height

    def do_context_clean_up(self):
        # clean up scheduled events.
        self.scheduler.stop_scheduling_events()
        # recover disasters simulated in the test bed.
        self.start_recover_time = time.time()
        self.end_chain_height_before_recover = self.get_chain_height()
        self.recover()
        # checking if disaster is healed with block synced again with alive nodes within a specified duration.
        if self.is_healed() is True:
            self.end_chain_height_after_recover = self.get_chain_height()
            self.end_recover_time = time.time()
            #self.generate_report()

    def tx_send(self):
        try:
            if 'startAt' in self.test_case_conf['input']:
                start_point = self.test_case_conf['input']['startAt']
                start = timer()
                while True:
                    gap = start_point - (timer() - start)
                    if gap > 0:
                        time.sleep(1)
                        self.logger.debug('Waiting for %ds to start sending TX.', gap)
                    else:
                        break

        except (KeyError, TypeError) as e:
            self.logger.error("Wrong configuration file. %s", e)
            return None

        self.tx_start_time = time.time()
        self.tx_start_chain_height = self.get_chain_height()
        try:
            start = timer()
            duration = self.test_case_conf["input"]["duration"]
            sender_index = self.test_case_conf["input"]["senderNode"]
            receiver_index = self.test_case_conf["input"]["receiverNode"]
            amount_per_tx = self.test_case_conf["input"]["amountperTX"]

            if sender_index not in self.clients or receiver_index not in self.clients:
                return None
            self.sender_before_balance = self.clients[sender_index].get_balance()
            self.receiver_before_balance = self.clients[receiver_index].get_balance()
            while (timer() - start) < duration or self.scheduler.is_scheduling_events():
                time.sleep(1)
                try:
                    txn_hash = self.clients[sender_index].send_transaction(
                        to="0x{}".format(self.clients[receiver_index].coin_base), value=amount_per_tx, gas_price=5000)
                    if txn_hash is not None:
                        with open(self.tx_history_file, 'a+') as f:
                            f.write('{}\n'.format(txn_hash))
                        self.tx_sent += 1
                except Exception as e:
                    self.logger.error("Send TX failed due to exception: %s.", e)

            self.tx_end_time = time.time()
            self.tx_end_chain_height = self.get_chain_height()
        except Exception as e:
            self.logger.error("cannot access remote RPC endpoint: %s", e)
            return None
        return True

    def is_balance_okay(self):
        """Verify balance base on test_case_conf between sender and receiver."""
        self.logger.debug("Before test, sender have: %d tokens", self.sender_before_balance)
        self.logger.debug("Before test, receiver have: %d tokens", self.receiver_before_balance)
        amount_per_tx = self.test_case_conf["input"]["amountperTX"]
        sender_index = self.test_case_conf["input"]["senderNode"]
        receiver_index = self.test_case_conf["input"]["receiverNode"]
        try:
            with open(self.tx_history_file, 'r') as reader:
                for tx_hash in reader:
                    # check if TX is mined, then calculate balance between sender and receiver.
                    result = self.clients[sender_index].get_transaction_by_hash(tx_hash.strip('\n'))
                    # TX was mined, count the expected balance
                    if result["blockHash"] is not None:
                        self.tx_mined += 1
                        self.balance_mined_by_the_test += amount_per_tx
        except IOError as e:
            self.logger.error("Cannot get TX via RPC api: %s", e)
        except (KeyError, TypeError) as e:
            self.logger.error("Cannot find blockHash from result, something wrong from RPC service: %s", e)
        except Exception as e:
            self.logger.error("Something wrong happens at balance validation. %s", e)

        sender_after_balance = self.clients[sender_index].get_balance()
        receiver_after_balance = self.clients[receiver_index].get_balance()

        if sender_after_balance is None or receiver_after_balance is None:
            return False

        # checking balance if sending tokens to self.
        if sender_index == receiver_index:
            self.logger.debug("sender balance: %d, receiver balance: %d", sender_after_balance, receiver_after_balance)
            if sender_after_balance != receiver_after_balance:
                return False

        # checking sender's balance.
        if self.sender_before_balance - self.balance_mined_by_the_test == sender_after_balance is False:
            return False

        # checking receiver's balance.
        if self.receiver_before_balance + self.balance_mined_by_the_test == receiver_after_balance is False:
            return False
        return True

    def get_dead_validators(self):
        try:
            dead_nodes = self.test_case_conf["condition"]["crashNodes"]
            self.logger.debug("get dead node: %s", dead_nodes)
            return dead_nodes
        except (KeyError, TypeError) as e:
            self.logger.error("Wrong configuration for test case. %s", e)
            return None

    def get_alive_validators(self):
        try:
            alive_nodes = []
            dead_nodes = self.get_dead_validators()
            for index, client in self.clients.items():
                if client.index not in dead_nodes:
                    alive_nodes.append(client.index)
        except (KeyError, TypeError) as e:
            self.logger.error("Wrong configuration for test case. %s", e)
            return None
        self.logger.debug("Get alive nodes: %s", alive_nodes)
        return alive_nodes

    def is_block_in_consistent_state(self):
        # to do tracking the block hash via RPC through out the validators.
        try:
            map_height_hash = {}
            alive_nodes = self.get_alive_validators()
            if alive_nodes is None:
                return False
            for r in range(1, BLOCK_CONSISTENT_CHECKING_DURATION):
                reference_height = self.get_chain_height()
                for node in alive_nodes:
                    block_hash = self.clients[node].get_block_hash_by_height(reference_height)
                    if block_hash is None:
                        continue
                    if reference_height in map_height_hash:
                        if map_height_hash.get(reference_height) != block_hash and block_hash is not None:
                            self.logger.error("BLOCK INCONSISTENT in round: %d, at height: %d", r, reference_height)
                            return False
                        else:
                            self.logger.debug('checking consistence in height: %d, hash %s of node: %s port: %s.',
                                              reference_height, block_hash if block_hash is not None else 'NULL',
                                              self.clients[node].host, self.clients[node].p2p_port)
                    else:
                        map_height_hash[reference_height] = block_hash
                time.sleep(1)
            self.logger.debug("Verified block consistent with heights: %s", map_height_hash)
            return True
        except (KeyError, TypeError) as e:
            self.logger.error("Wrong configuration file. %s", e)
            return False

    def get_expected_engine_state(self):
        try:
            if self.test_case_conf["output"]["engineAlive"] is True:
                return True
            return False
        except (KeyError, TypeError) as e:
            self.logger.error("Wrong configuration for test case. %s", e)
            return False

    def is_engine_state_expected(self):
        # get expected state from test conf.
        should_engine_produce_block = self.get_expected_engine_state()

        try:
            on_hold_counter = 0
            height = self.get_chain_height()
            self.logger.debug("latest chain height %d", height)
            for r in range(1, ENGINE_STATE_CHECKING_DURATION):
                time.sleep(1)
                new_height = self.get_chain_height()
                if height < new_height:
                    # reset on-holding counter.
                    on_hold_counter = 0
                    self.logger.info("Consensus engine is keeping produce blocks: %d.", new_height)
                    height = new_height
                    if should_engine_produce_block is False:
                        return False
                    else:
                        break
                else:
                    on_hold_counter += 1
                    self.logger.info('Consensus engine is on-holding for %ds', on_hold_counter)
                    if on_hold_counter == (ENGINE_STATE_CHECKING_DURATION - 1):
                        if should_engine_produce_block is True:
                            return False

        except (KeyError, TypeError) as e:
            self.logger.error("Wrong configuration file. %s ", e)
            return False
        return True

    def start_test(self):
        self.passed = self.run()
        if self.passed is False:
            self.collect_test_case_context_log()
            # save system log for failed case, do not print them into engine's std output.
            self.save_system_log()
            #self.prints_system_log()
            return False
        else:
            # save autonity clients logs as well for passed test case.
            self.save_system_log()
            return True

    def run(self):
        """run the test case, and tear down the test case with network recovery."""
        self.logger.debug("before running test case, thread: %d.", threading.active_count())
        self.start_chain_height = self.get_chain_height()
        self.logger.debug("start schedule events...")
        if self.scheduler.schedule() is not True:
            return False
        if self.tx_send() is not True:
            self.do_context_clean_up()
            return False
        if self.is_balance_okay() is not True:
            self.do_context_clean_up()
            return False
        if self.is_engine_state_expected() is not True:
            self.do_context_clean_up()
            return False
        if self.is_block_in_consistent_state() is not True:
            self.do_context_clean_up()
            return False

        self.start_recover_time = time.time()
        self.end_chain_height_before_recover = self.get_chain_height()

        if self.recover() is not True:
            self.scheduler.stop_scheduling_events()
            return False
        self.logger.debug("After disaster recover, thread: %d.", threading.active_count())

        # checking if disaster is healed with block synced again with alive nodes within a specified duration.
        if self.is_healed() is True:
            self.end_chain_height_after_recover = self.get_chain_height()
            self.end_recover_time = time.time()
            self.logger.info("TESTCASE: %s is passed.", self.test_case_conf["name"])
            #self.generate_report()
            self.scheduler.try_join()
            return True

        self.scheduler.try_join()
        self.logger.info('Recovering timeout happens.')
        return False

    def generate_report(self):
        if self.tx_start_chain_height > self.tx_end_chain_height:
            self.logger.info('Blockchain was re-initialized from scratch.')
            self.tx_start_chain_height = 0

        if self.start_chain_height > self.end_chain_height_after_recover:
            self.logger.info('Blockchain was re-initialized from scratch.')
            self.start_chain_height = 0

        self.logger.info('statistics: $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$')
        self.logger.info('statistics: Name: %s', self.test_case_conf['name'])
        self.logger.info('statistics: TC: start time:            %s', time.asctime(time.localtime(self.start_time)))
        self.logger.info('statistics: TC: duration:              %ds', self.start_recover_time - self.start_time)
        self.logger.info('statistics: TC: start recover time:    %s', time.asctime(time.localtime(self.start_recover_time)))
        self.logger.info('statistics: TC: duration:              %ds', self.end_recover_time - self.start_recover_time)
        self.logger.info('statistics: TC: end recover time:      %s', time.asctime(time.localtime(self.end_recover_time)))
        self.logger.info('statistics: TC: start chain height:    %d', self.start_chain_height)
        self.logger.info('statistics: TC: before recover height: %d', self.end_chain_height_before_recover)
        self.logger.info('statistics: TC: after recover height:  %d', self.end_chain_height_after_recover)
        self.logger.info('statistics: TC: Block producing speed: %1.3f block/s',
                         (self.end_chain_height_after_recover - self.start_chain_height) /
                         (self.end_recover_time - self.start_time))
        self.logger.info('statistics: -------------------------------------------------')
        self.logger.info('statistics: TX: start:                 %s', time.asctime(time.localtime(self.tx_start_time)))
        self.logger.info('statistics: TX: end:                   %s', time.asctime(time.localtime(self.tx_end_time)))
        self.logger.info('statistics: TX: duration:              %ds', self.tx_end_time - self.tx_start_time)
        self.logger.info('statistics: TX: start height:          %d', self.tx_start_chain_height)
        self.logger.info('statistics: TX: end   height:          %d', self.tx_end_chain_height)
        self.logger.info('statistics: TX: Block producing speed: %1.3f block/s.',
                         (self.tx_end_chain_height - self.tx_start_chain_height) / (self.tx_end_time - self.tx_start_time))
        self.logger.info('statistics: TX: %d of %d was mined in ledger, TPL: %.2f%%. Token delivered: %d.',
                         self.tx_mined, self.tx_sent, (self.tx_mined/self.tx_sent)*100, self.balance_mined_by_the_test)
        self.logger.info('statistics: $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$')

    def collect_test_case_context_log(self):
        try:
            # try to create dirs
            os.makedirs(SYSTEM_LOG_DIR, exist_ok=True)  # It never fail even if the dir is existed.
            log_dir = TEST_CASE_SYSTEM_LOG_DIR.format(self.id, "passed" if self.passed else "failed")
            os.makedirs(log_dir, exist_ok=True)

            self.test_case_conf['testcase_start_time'] = time.ctime(self.start_time)
            self.test_case_conf['testcase_start_height'] = self.start_chain_height
            self.test_case_conf['testcase_end_height'] = self.end_chain_height_before_recover
            self.test_case_conf['testcase_end_time'] = time.ctime(time.time())
            self.test_case_conf['ip_mapping'] = []
            for index, client in self.clients.items():
                self.test_case_conf['ip_mapping'].append("{}:{}".format(index, client.host))

            ret = conf.write_yaml(TEST_CASE_CONTEXT_FILE_NAME.format(self.id), self.test_case_conf)
            if ret is not True:
                self.logger.warning('cannot save test case context.')

            # print test case context in log file.
            self.logger.info("\n\n\n")
            self.logger.info("The failed test case context is collected as below:")
            self.logger.info(self.test_case_conf)
            self.logger.info("\n\n\n")
        except Exception as e:
            self.logger.error("Cannot collect test case context logs. %s", e)

    def recover(self):
        failed = False
        for index, client in self.clients.items():
            if client.heal_from_disaster() is not True:
                failed = True
        return True if not failed else False

    def is_healed(self):
        # measure the best height.
        best_height = self.get_chain_height()
        healed_clients = {}
        start = timer()
        while (timer() - start) < HEAL_TIME_OUT:
            self.logger.debug("IsHeal, current thread count: %d", threading.active_count())
            if len(healed_clients) == len(self.clients):
                return True
            for index, client in self.clients.items():
                height = client.get_chain_height()
                if height is None:
                    continue
                if height >= best_height:
                    healed_clients[index] = client
            time.sleep(1)
        self.logger.warning('Disaster recovering timeout. 5 minutes!')
        return False

    def save_system_log(self):
        try:
            # try to create dirs
            os.makedirs(SYSTEM_LOG_DIR, exist_ok=True)  # It never fail even if the dir is existed.
            log_dir = TEST_CASE_SYSTEM_LOG_DIR.format(self.id, "passed" if self.passed else "failed")
            os.makedirs(log_dir, exist_ok=True)
            for index, client in self.clients.items():
                client.download_system_log(log_dir)
        except Exception as e:
            self.logger.error('Cannot fetch logs from node. %s.', e)
            return None
        return True

    def prints_system_log(self):
        if self.save_system_log() is None:
            return None

        # redirect autonity client logs into test engine's logger, this is used for github CI to dump logs in report.
        try:
            for index, client in self.clients.items():
                log_dir = TEST_CASE_SYSTEM_LOG_DIR.format(self.id, "passed" if self.passed else "failed")
                client.redirect_system_log(log_dir)
        except Exception as e:
            self.logger.error('Cannot redirect system logs from client into test engine log file %s.', e)
            return None
        return True
