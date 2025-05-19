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
    parser.add_argument("-id", help='testcase index', type=int, required=True, default=0)

    args = parser.parse_args()
    autonity_path = args.autonity
    test_id = args.id

    LG.debug(f"testcase index: {test_id}")

    conf.load_project_conf()

    # plan the autonity network configurations, and start all autonity clients.
    network_planner = NetworkPlanner(autonity_path)
    network_planner.plan()
    network_planner.deploy()
    network_planner.start_all_nodes()

    clients = network_planner.get_clients()
    try:
        # load test case view, and start the specific test
        test_set = conf.get_test_case_conf()
        test_case = test_set["playbook"]["testcases"][test_id]
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
    except Exception as e:
        LG.error("Get error: %s", e)

    exit(1)
