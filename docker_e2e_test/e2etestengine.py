#!/usr/bin/env/ python3
import argparse

import log
from conf import conf
from testcase.testcase import TestCase
from planner.networkplanner import NetworkPlanner
import time

LG = log.get_logger()


if __name__ == '__main__':
    LG.debug("##########################################")
    LG.debug("\n\nTest Engine start.")

    parser = argparse.ArgumentParser()
    parser.add_argument("autonity", help='Autonity Binary Path')
    parser.add_argument("-id", help='Starting testcase index', type=int, required=True, default=0)
    # IPs in certain format: 172.17.0.2,172.17.0.3,172.17.0.4,172.17.0.5,172.17.0.6,172.17.0.7
    parser.add_argument("-ips", help='Validator IPs', type=str, required=True, default="")

    args = parser.parse_args()
    autonity_path = args.autonity
    test_id = args.id
    validator_ips = args.ips.split(",")

    LG.info(f"Running test {test_id} with validator nodes {args.ips}")

    conf.load_project_conf()
    network_planner = NetworkPlanner(autonity_path, validator_ips, test_id+1)
    network_planner.plan()
    network_planner.deploy_all_nodes()
    network_planner.start_all_nodes()
    clients = network_planner.get_clients()

    passed = True
    test_set = conf.get_test_case_conf()

    try:
        if test_id >= len(test_set["playbook"]["testcases"]):
            exit(0)

        LG.debug("\n" + "="*50)
        LG.debug(f"Starting Test ID: {test_id}")

        test_case = test_set["playbook"]["testcases"][test_id]
        test = TestCase(test_case, clients)

        LG.info("Starting test case: %s", test_case)
        result = test.start_test()

        if result:
            LG.info(f"[TEST {test_id+1} PASSED]")
        else:
            LG.error(f"[TEST {test_id+1} FAILED]")
            passed = False

    except (KeyError, TypeError, IndexError) as e:
        LG.error(f"Invalid test ID {test_id}: {e}")
        passed = False
    except Exception as e:
        LG.error(f"Error in test ID {test_id}: {e}")
        passed = False

    if passed:
        LG.info("[TEST PASSED]")
        exit(0)
    else:
        LG.info("[TEST FAILED]")
        exit(1)
