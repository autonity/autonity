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

    args = parser.parse_args()
    autonity_path = args.autonity
    start_id = args.id

    # Generate 3 consecutive test IDs
    batch_size = 3
    test_ids = range(start_id, start_id + batch_size)
    LG.debug(f"Running tests: {list(test_ids)}")

    conf.load_project_conf()
    network_planner = NetworkPlanner(autonity_path)
    network_planner.plan()
    network_planner.deploy_all_nodes()
    network_planner.start_all_nodes()
    clients = network_planner.get_clients()

    all_passed = True
    test_set = conf.get_test_case_conf()

    for test_id in test_ids:
        try:
            if test_id >= len(test_set["playbook"]["testcases"]):
                break

            LG.debug("\n" + "="*50)
            LG.debug(f"Starting Test ID: {test_id}")

            test_case = test_set["playbook"]["testcases"][test_id]
            test = TestCase(test_case, clients)

            LG.info("Starting test case: %s", test_case)
            result = test.start_test()

            if result:
                LG.info(f"[TEST {test_id} PASSED]")
            else:
                LG.error(f"[TEST {test_id} FAILED]")
                all_passed = False

            network_planner.re_genesis_network()

        except (KeyError, TypeError, IndexError) as e:
            LG.error(f"Invalid test ID {test_id}: {e}")
            all_passed = False
        except Exception as e:
            LG.error(f"Error in test ID {test_id}: {e}")
            all_passed = False

    if all_passed:
        LG.info("[TEST PASSED]")
        exit(0)
    else:
        LG.info("[TEST FAILED]")
        LG.info("Log collecting for failed tests...")
        time.sleep(180)
        exit(1)
