#!/usr/bin/env/ python3
import argparse

import log
from conf import conf
from testcase.testcase import TestCase
from planner.networkplanner import NetworkPlanner
from client.client import Client
import time

LG = log.get_logger()


if __name__ == '__main__':
    LG.debug("##########################################")
    LG.debug("")
    LG.debug("")
    LG.debug("Test Engine start.")

    parser = argparse.ArgumentParser()
    parser.add_argument("autonity", help='Autonity Binary Path')
    parser.add_argument("-d", help='Start deploy remote network with brand new configurations.', type=bool, default=True)
    parser.add_argument("-t", help='Start test remote network.', type=bool, default=True)
    parser.add_argument("-id", help='testcase index', type=int, required=True, default=0)

    args = parser.parse_args()

    is_deploy = args.d
    is_testing = args.t
    autonity_path = args.autonity
    id = args.id

    LG.debug(f"testcase index: {id}")

    conf.load_project_conf()
    network_planner = None

    exit_code = 0

    network_planner = NetworkPlanner(autonity_path)
    network_planner.plan()
    network_planner.deploy()
    network_planner.start_all_nodes()

    if is_testing:
        clients = None
        if network_planner:
            clients = network_planner.get_clients()
        else:
            # load network view from generated testbed.conf
            clients = []
            test_bed = conf.get_test_bed_conf()
            try:
                for node in test_bed["targetNetwork"]["nodes"]:
                    client = Client(host=node["name"], p2p_port=node["p2pPort"], rpc_port=node["rpcPort"], ws_port=node["wsPort"],
                                    net_interface=node["ethernetInterfaceID"], coin_base=node["coinBase"][2:],
                                    ssh_user=node["sshCredential"]["sshUser"], ssh_pass=node["sshCredential"]["sshPass"],
                                    ssh_key=node["sshCredential"]["sshKey"], sudo_pass=node["sshCredential"]["sudoPass"],
                                    role=node["role"], index=node["index"])
                    clients.append(client)
            except Exception as e:
                LG.error("Process exit with cannot conf from test bed conf.", e)
                exit_code = 1

        try:
            # load test case view, and start the specific test
            test_set = conf.get_test_case_conf()
            for index, test_case in enumerate(test_set["playbook"]["testcases"]):
                if index != id:
                    continue
                playbook = conf.get_test_case_conf()
                if playbook["playbook"]["stop"] is True:
                    LG.info("Playbook is stopped by user configuration: testcaseconf.yml/playbook/stop: true.")
                    break
                test = TestCase(test_case, clients)
                LG.debug("")
                LG.debug("")
                LG.info("start test case: %s", test_case)
                LG.debug("")
                LG.debug("")
                result = test.start_test()
                if result is True:
                    LG.info("[TEST PASSED]")
                    exit(0)
                if result is False:
                    LG.info("[TEST FAILED]")
                    LG.info("Log collecting for failed test ......")
                    time.sleep(180)
                    exit(1)
        except (KeyError, TypeError) as e:
            LG.error("Wrong configuration. %s", e)
            exit_code = 1
        except Exception as e:
            LG.error("Get error: %s", e)
            exit_code = 1

    exit(exit_code)
