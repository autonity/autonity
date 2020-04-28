from client.client import Client
from typing import List
import log
import sched
import threading
import time
import copy


class Scheduler(object):
    def __init__(self, test_case_conf, clients: List[Client]):
        self.test_case_conf = test_case_conf
        self.clients = {}
        for client in clients:
            self.clients[client.index] = client
        self.logger = log.get_logger()
        self.scheduler = sched.scheduler(time.time, time.sleep)
        self.thread = threading.Thread(target=self.scheduler.run)

    def schedule(self):
        try:
            case_name = self.test_case_conf["name"]
            crash_nodes = self.test_case_conf["condition"]["crashNodes"]
            scenario = self.test_case_conf["condition"]["scenario"]
        except KeyError as e:
            self.logger.error("parameter missing in testcaseconf.yml. %s", e)
            return None
        except TypeError as e:
            self.logger.error("parameter with wrong type in testcaseconf.yml. %s", e)
            return None

        # checking the parameters.
        if self.is_valid_parameters(crash_nodes, scenario) is False:
            return None

        # schedule the disaster actions.
        if self.schedule_actions(case_name, scenario) is not True:
            return None

        # start scheduling thread.
        self.start_scheduling_thread()
        return True

    def is_scheduling_events(self):
        if self.thread.is_alive() or self.scheduler.empty() is False:
            return True
        return False

    def start_scheduling_thread(self):
        self.thread.start()

    def stop_scheduling_events(self):
        event_list = self.scheduler.queue
        for event in event_list:
            try:
                self.scheduler.cancel(event)
            except ValueError as e:
                self.logger.debug('event not in queue, %s', e)
                continue

    def try_join(self):
        if self.thread.is_alive():
            self.thread.join()

    def is_valid_parameters(self, crash_nodes, scenario):
        # to do check nodes in scenario happens in crash_nodes list.
        try:
            for step in scenario:
                # skip none disaster steps.
                if step["action"] != "stop":
                    continue

                for target in step["target"]:
                    if target not in crash_nodes:
                        return False
        except (TypeError, KeyError) as e:
            self.logger.error("Wrong configuration. %s", e)
            return False
        self.logger.debug("crash_nodes: %s", crash_nodes)
        # to do check crash nodes presents in the test bed.
        try:
            for crash_node in crash_nodes:
                if crash_node not in self.clients:
                    return False
        except (TypeError, KeyError) as e:
            self.logger.error("Wrong configuration for crash nodes. %s", e)
            return False

        return True

    def schedule_actions(self, case_name, scenario):
        self.logger.debug('Start disaster for test case: %s, disaster impacts: %s.', case_name, scenario)
        try:
            if scenario is []:
                self.logger.debug("There is no disaster defined for case, skip scheduling.")
                return True

            for step in scenario:
                if self.schedule_action(step) is not True:
                    return None
        except (TypeError, KeyError) as e:
            self.logger.error("Wrong configuration. %s", e)
            return None
        return True

    def cold_deploy_network(self):
        failed = False
        for index, client in self.clients.items():
            if client.stop_client() is not True:
                failed = True
            if client.clean_chain_data() is not True:
                failed = True
            if client.deliver_package() is not True:
                failed = True
        return True if not failed else False

    def connect_peers(self, test, peers):
        for peer in peers:
            if len(peer) is not 2:
                self.logger.warning('Wrong peer configuration, skip the connection control %s', peer)
                continue
            if peer[0] not in self.clients and peer[1] not in self.clients:
                self.logger.warning("wrong peer configuration. skip the connection control %s.", peer)
                continue
            is_disconnected = self.clients[peer[0]].is_peer_disconnected(self.clients[peer[1]].host,
                                                                         self.clients[peer[1]].p2p_port)
            if is_disconnected is False:
                continue
            self.clients[peer[0]].connect_peer(self.clients[peer[1]].host, self.clients[peer[1]].p2p_port)

    def dis_connect_peers(self, test, peers):
        for peer in peers:
            if len(peer) is not 2:
                self.logger.warning('Wrong peer configuration, skip the connection control %s', peer)
                continue
            if peer[0] not in self.clients and peer[1] not in self.clients:
                self.logger.warning("wrong peer configuration. skip the connection control %s.", peer)
                continue
            if self.clients[peer[0]].is_peer_disconnected(self.clients[peer[1]].host,
                                                                         self.clients[peer[1]].p2p_port):
                continue
            self.clients[peer[0]].dis_connect_peer(self.clients[peer[1]].host, self.clients[peer[1]].p2p_port)

    def stop_clients(self, test, nodes):
        for index in nodes:
            if index not in self.clients:
                self.logger.warning("wrong node index in test case. skip the crash action. %s", index)
                continue
            self.clients[index].stop_client()

    def start_clients(self, test, nodes):
        for index in nodes:
            if index not in self.clients:
                self.logger.warning("wrong node index in test case. skip the crash action. %s", index)
                continue
            self.clients[index].start_client()

    def set_delays(self, test, step):
        try:
            for delay_meta in step["latency"]:
                if delay_meta["host"] not in self.clients:
                    self.logger.warning("wrong node index in delay meta data. skip the delay simulation. %s",
                                        delay_meta["host"])
                    continue
                if delay_meta["uplink"] is not None:
                    if self.clients[delay_meta["host"]].set_up_link_delay(delay_meta["uplink"]) is not True:
                        self.logger.warning("Set uplink delay failed: %s, %s.", self.clients[delay_meta["host"], delay_meta])
                if delay_meta["downlink"] is not None:
                    if self.clients[delay_meta["host"]].set_down_link_delay(delay_meta["downlink"]) is not True:
                        self.logger.warning("Set downlink delay failed: %s, %s.", self.clients[delay_meta["host"], delay_meta])

        except Exception as e:
            self.logger.warning("cannot simulate delays. %s", e)
            return False
        return True

    def clear_delays(self, test, step):
        try:
            for index in step["target"]:
                if index not in self.clients:
                    self.logger.warning("wrong node index in delay meta data. skip the delay clearance. %s", index)
                    continue
                self.clients[index].cancel_up_link_delay()
                self.clients[index].cancel_down_link_delay()
        except Exception as e:
            self.logger.warning("cannot clear delays. %s", e)
            return False
        return True

    def schedule_action(self, step):
        try:
            self.logger.debug("Schedule disaster action within %ds for step: %s.", step["delay"], step)
            if step['action'] == 're-deploy' and 're-deploy' in self.test_case_conf \
                    and self.test_case_conf['re-deploy'] is True:
                self.scheduler.enter(step['delay'], 1, self.cold_deploy_network)
                return True

            if step["action"] == "connect":
                peers = copy.deepcopy(step["peers"])
                self.scheduler.enter(step["delay"], 1, self.connect_peers, argument=(1, peers))
                return True

            if step["action"] == "disconnect":
                peers = copy.deepcopy(step["peers"])
                self.scheduler.enter(step["delay"], 1, self.dis_connect_peers, argument=(1, peers))
                return True

            if step["action"] == "stop":
                target = copy.deepcopy(step["target"])
                self.scheduler.enter(step["delay"], 1, self.stop_clients, argument=(1, target))
                return True

            if step["action"] == "start":
                target = copy.deepcopy(step["target"])
                self.scheduler.enter(step["delay"], 1, self.start_clients, argument=(1, target))
                return True

            if step["action"] == "delay":
                step_copy = copy.deepcopy(step)
                self.scheduler.enter(step["delay"], 1, self.set_delays, argument=(1, step_copy))
                return True

            if step["action"] == "un-delay":
                step_copy = copy.deepcopy(step)
                self.scheduler.enter(step["delay"], 1, self.clear_delays, argument=(1, step_copy))
                return True

        except (KeyError, TypeError) as e:
            self.logger.error("Wrong configuration. %s", e)
            return None
        except Exception as e:
            self.logger.error("Exception happens when schedule actions. %s", e)
            return None
        return True
