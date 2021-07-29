#!/usr/bin/env/ python3
import argparse

import log
from conf import conf
from testcase.testcase import TestCase
from testcase.misbehaviour_testcase import MisbehaviourTestCase
from testcase.misbehaviour_testcase import TwinsTestCase
from planner.networkplanner import NetworkPlanner
from planner.twinsplanner import TwinsNetworkPlanner
import time

LG = log.get_logger()


def network_disaster_tests(binary_path, conf, num_of_cases, passed_testcases, failed_testcases):
    # return num_of_cases, passed_testcases, failed_testcases
    # Deploy will create brand new configurations then bootstrap entire network from genesis block.
    network_planner = NetworkPlanner(autonity_path=binary_path)
    network_planner.plan()
    network_planner.deploy()
    network_planner.start_all_nodes()
    clients = network_planner.get_clients()
    try:
        # load test case view, and start testing one by one.
        test_set = conf.get_test_case_conf()
        num_of_cases += len(test_set["playbook"]["testcases"])
        for test_case in test_set["playbook"]["testcases"]:
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
                LG.info('TEST CASE PASSED: %s', test_case)
                passed_testcases.append(test_case)
            if result is False:
                LG.error('TEST CASE FAILED: %s', test_case)
                failed_testcases.append(test_case)

    except (KeyError, TypeError) as e:
        LG.error("Wrong configuration. %s", e)
    except Exception as e:
        LG.error("Get error: %s", e)
    return num_of_cases, passed_testcases, failed_testcases


def malicious_behaviour_tests(binary_path, conf, num_of_cases, passed_testcases, failed_testcases):
    # return num_of_cases, passed_testcases, failed_testcases
    network_planner = NetworkPlanner(autonity_path=binary_path)
    network_planner.plan()
    network_planner.deploy()
    network_planner.start_all_nodes()
    clients = network_planner.get_clients()
    # run misbehaviour test cases
    try:
        # load test case view, and start testing one by one.
        test_set = conf.get_misbehaviour_test_case_conf()
        num_of_cases += len(test_set["playbook"]["malicious_tests"])
        for test_case in test_set["playbook"]["malicious_tests"]:
            test = MisbehaviourTestCase(test_case, clients)
            LG.debug("")
            LG.debug("")
            LG.info("start test case: %s", test_case)
            LG.debug("")
            LG.debug("")
            result = test.run()
            if result is True:
                LG.info('TEST CASE PASSED: %s', test_case)
                passed_testcases.append(test_case)
            if result is False:
                LG.error('TEST CASE FAILED: %s', test_case)
                failed_testcases.append(test_case)
    except (KeyError, TypeError) as e:
        LG.error("Wrong configuration. %s", e)
    except Exception as e:
        LG.error("Get error: %s", e)
    return num_of_cases, passed_testcases, failed_testcases


def twins_test(binary_path, conf, num_of_cases, passed_testcases, failed_testcases):
    # twins in autonity is not allowed since the user management in Autonity contract does not support duplicated user
    # with different enode endpoint, by the removal of the duplicated user checker, then the twins can run.
    # set node_0 and node_3 as twins for the time being, todo: get it from config.
    twins = [[0, 3]]
    twins_planner = TwinsNetworkPlanner(autonity_path=binary_path, twins=twins)
    twins_planner.plan()
    twins_planner.deploy()
    twins_planner.set_ring_topology()
    twins_planner.start_all_nodes()
    clients = twins_planner.get_clients()
    # run twins test cases
    try:
        # load test case view, and start testing one by one.
        test_set = conf.get_misbehaviour_test_case_conf()
        num_of_cases += len(test_set["playbook"]["twins_tests"])
        for test_case in test_set["playbook"]["twins_tests"]:
            test = TwinsTestCase(test_case, clients, twins)
            LG.debug("")
            LG.debug("")
            LG.info("start test case: %s", test_case)
            LG.debug("")
            LG.debug("")
            result = test.run()
            if result is True:
                LG.info('TEST CASE PASSED: %s', test_case)
                passed_testcases.append(test_case)
            if result is False:
                LG.error('TEST CASE FAILED: %s', test_case)
                failed_testcases.append(test_case)
    except (KeyError, TypeError) as e:
        LG.error("Wrong configuration. %s", e)
    except Exception as e:
        LG.error("Get error: %s", e)
    return num_of_cases, passed_testcases, failed_testcases


if __name__ == '__main__':
    LG.debug("##########################################")
    LG.debug("")
    LG.debug("")
    LG.debug("Test Engine start.")

    parser = argparse.ArgumentParser()
    parser.add_argument("autonity", help='Autonity Binary Path')
    parser.add_argument("-d", help='Start deploy remote network with brand new configurations.', type=bool, default=True)
    parser.add_argument("-t", help='Start test remote network.', type=bool, default=True)
    args = parser.parse_args()

    autonity_path = args.autonity

    conf.load_project_conf()
    passed_testcases = []
    failed_testcases = []

    num_of_cases = 0

    # run disaster tests
    num_of_cases, passed_testcases, failed_testcases = \
        network_disaster_tests(autonity_path, conf=conf, num_of_cases=num_of_cases, passed_testcases=passed_testcases, failed_testcases=failed_testcases)

    # run malicious behaviour tests
    num_of_cases, passed_testcases, failed_testcases = \
        malicious_behaviour_tests(autonity_path, conf=conf, num_of_cases=num_of_cases, passed_testcases=passed_testcases, failed_testcases=failed_testcases)

    # run twins tests
    num_of_cases, passed_testcases, failed_testcases = \
        twins_test(autonity_path, conf=conf, num_of_cases=num_of_cases, passed_testcases=passed_testcases, failed_testcases=failed_testcases)

    # generate an overview of the test report.
    if len(passed_testcases) == num_of_cases:
        LG.info("[TEST PASSED]")

    LG.info("[PASS] %d/%d cases were passed.", len(passed_testcases), num_of_cases)
    for case in passed_testcases:
        LG.info("[PASS] %s", case["name"])

    if len(failed_testcases) > 0:
        exit_code = 1
        LG.info("[TEST FAILED]")
        LG.info("[FAILED] %d/%d cases were failed.", len(failed_testcases), num_of_cases)
        for case in failed_testcases:
            LG.info("[ERROR] %s", case["name"])

        for i in range(0, len(failed_testcases)):
            LG.info("Log collecting...")
            time.sleep(180)
        exit(1)

    exit(0)
