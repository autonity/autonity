import docker
import utility
import ipaddress
import argparse
import threading
import signal
import os

TEST_ENGINE_IMAGE_NAME = "enginehost{}/ubuntu"
TEST_ENGINE_DOCKER_FILE = "./Dockerfile"
CLIENT_IMAGE_NAME = "clienthost/ubuntu:2404"
CLIENT_DOCKER_FILE = "./clientDockerFile"
BUILDER_IMAGE_NAME = "go-builder/ubuntu"
NUM_OF_CLIENT = 6
NODE_NAME = "Node{}_{}"
ENGINE_NAME = "Engine{}"
VALIDATOR_IP_LIST_FILE = "./etc/validator.ip"
COMMAND_START_TEST = "python3 e2etestengine.py ./test_bin/autonity -id {}"
FAILED_TEST_LOGS = "./JOB_{}.tar"
SYSTEM_LOG_PATH = "/system_log"
JOB_ID = ""
COMMIT_HASH = ""


def check_environment():
    try:
        client = docker.from_env()
    except Exception as e:
        print("dockerd is not find: ", e)
        raise Exception("please check if docker is installed and dockerd is started")
    print("docker version: ", client.version())


def init():
    try:
        check_environment()
    except Exception as e:
        print("check environment: check to install docker.", e)
        if is_docker_installed() is False:
            install_docker()


def is_docker_installed():
    try:
        result = utility.execute("docker --version")
        if result[1] == "":
            return True
    except Exception as e:
        print("checking docker failed: ", e)
        return False
    print("docker is not installed yet, going to install docker.")
    return False


def install_docker():
    try:
        utility.execute("sudo apt-get update")
        utility.execute("sudo apt-get install --yes docker.io")
        print("docker is been installed.")
    except Exception as e:
        print("cannot install docker: ", e)


def check_docker_daemon():
    result = ("", "")
    try:
        result = utility.execute("pidof dockerd")
    except Exception as e:
        print("unknown state of dockerd. ", e)
        utility.execute("sudo service docker start")
    if result[0] == "" and result[1] == "":
        utility.execute("sudo service docker start")
        print("docker daemon is started by this script.")


def check_to_build_client_images():
    client_image_found = False
    try:
        client = docker.from_env()
        image_list = client.images.list()
        for image in image_list:
            if CLIENT_IMAGE_NAME in image.tags:
                print("image name ", image.tags)
                client_image_found = True
    except Exception as e:
        print("check", e)
    if not client_image_found:
        print("client image is not founded, going to build it.")
        create_image(CLIENT_IMAGE_NAME, CLIENT_DOCKER_FILE)


def create_image(tag, docker_file):
    path = "."
    network_mode = "bridge"
    try:
        client = docker.from_env()
        client.images.build(path=path, network_mode=network_mode, tag=tag, dockerfile=docker_file)
    except Exception as e:
        print("cannot build image: ", e)
        exit(1)
    print("image is been built: ", tag)


def start_client_containers(job_id):
    ip_set = set()
    try:
        client = docker.from_env()
        for i in range(0, NUM_OF_CLIENT):
            node_name = NODE_NAME.format(job_id, i)
            container = client.containers.run(CLIENT_IMAGE_NAME, name=node_name,
                                              detach=True, cap_add=["SYS_ADMIN", "NET_ADMIN"],
                                              volumes={"/sys/fs/cgroup": {"bind": "/sys/fs/cgroup", "mode": "ro"}})
            print("create new container: ", container.id)
            container.logs()
            result = utility.execute("sudo docker inspect -f \'{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}\' " + node_name)
            if result[1] != "":
                print("cannot get container ip: ", result[1])
                continue

            print("get container ip: ", result[0])
            for part in result[0].split():
                try:
                    a = ipaddress.ip_network(part)
                except ValueError:
                    pass
                else:
                    if a.is_private:
                        ip_set.add(str(a.network_address))
    except Exception as e:
        print("create container failed: ", e)
    finally:
        print("test bed was created: ", ip_set)
        return sorted(ip_set)


def prune_unused_images():
    try:
        client = docker.from_env()
        result = client.images.prune(filters={"dangling": True})
        print("docker prune unused, un-tagged images: ", result)
    except Exception as e:
        print("prune unused images: ", e)


def prune_unused_volumes():
    try:
        client = docker.from_env()
        result = client.volumes.prune()
        print("docker prune unused volumes: ", result)
    except Exception as e:
        print("prune unused volumes: ", e)


def prune_unused_network():
    try:
        client = docker.from_env()
        result = client.networks.prune()
        print("docker prune unused network: ", result)
    except Exception as e:
        print("prune unused network: ", e)


def remove_test_engine_image(commit_hash):
    try:
        client = docker.from_env()
        result = client.api.remove_image(TEST_ENGINE_IMAGE_NAME.format(commit_hash), force=True, noprune=True)
        print("docker remove image: ", result)
    except Exception as e:
        print("docker remove image: ", e)


# to stop and remove client and test engine containers.
def clean_test_bed_containers(job_id):
    print("start clean up current test context: test bed and test engine.")
    client = docker.from_env()
    for i in range(0, NUM_OF_CLIENT):
        try:
            node_name = NODE_NAME.format(job_id, i)
            container = client.containers.get(node_name)
            container.stop()
            container.remove()
            print("remove container: ", node_name)
        except Exception as e:
            print("remove container: ", e)
            continue
    try:
        # stop and remove test engine container.
        container = client.containers.get(ENGINE_NAME.format(job_id))
        container.stop()
        container.remove()
        print("remove test engine: ", ENGINE_NAME.format(job_id))
    except Exception as err:
        print("stop and remove test engine container: ", err)


def check_to_build_engine_image(commit_hash):
    print("try to build test engine image")
    image_found = False
    try:
        client = docker.from_env()
        image_list = client.images.list()
        for image in image_list:
            if TEST_ENGINE_IMAGE_NAME.format(commit_hash) in image.tags:
                print("image name ", image.tags)
                image_found = True
    except Exception as e:
        print("check", e)
    if not image_found:
        print("engine image is not founded, going to build it.")
        create_image(TEST_ENGINE_IMAGE_NAME.format(commit_hash), TEST_ENGINE_DOCKER_FILE)


def dump_ips_to_engine_conf(ips):
    try:
        with open(VALIDATOR_IP_LIST_FILE, 'w+') as f:
            for ip in ips:
                f.write("{}\n".format(ip))
        print("test bed IP list was configured into test engine.")
    except Exception as e:
        print("failed to dump ip into test engine validator.ip file. ", e)


def start_test_engine_container(commit_hash, job_id, test_id):
    print("start test engine container, the testcase will be run in it.")
    try:
        print("test engine is going to start:")
        client = docker.from_env()
        cmd = COMMAND_START_TEST.format(test_id)
        container = client.containers.run(TEST_ENGINE_IMAGE_NAME.format(commit_hash), command=cmd,
                                          name=ENGINE_NAME.format(job_id), detach=True, privileged=True)
        print("test engine is started.")
        return container
    except Exception as e:
        print("create test engine container failed: ", e)


def clean_up(job_id, commit_hash):
    clean_test_bed_containers(job_id)
    remove_test_engine_image(commit_hash)
    prune_unused_images()
    prune_unused_volumes()
    prune_unused_network()


def thread_func_copy_system_logs(job_id, path):
    try:
        print("***: start collecting logs from test engine container.")
        with open(FAILED_TEST_LOGS.format(job_id), 'wb') as f:
            bits, stat = container.get_archive(path)
            print(stat)
            for chunk in bits:
                f.write(chunk)
    except Exception as e:
        print("***: collecting system logs failed. ", e)
    finally:
        print("***: log was collected at: ", FAILED_TEST_LOGS.format(job_id))


def receive_signal(signal_number, frame):
    print('Signal Received: ', signal_number)
    clean_up(JOB_ID, COMMIT_HASH)
    raise SystemExit('Exiting')
    return


if __name__ == "__main__":
    exit_code = 1
    parser = argparse.ArgumentParser()
    parser.add_argument("autonity", help="Autonity WorkDir Path")
    # Adding test case ID for the target test case to be run.
    parser.add_argument("-id", help='Test Case ID', type=int, required=True, default=0)
    parser.add_argument("-hash", help='Commit hash', type=str, required=True, default="0xaabcdef")

    args = parser.parse_args()
    autonity_path = os.path.abspath(args.autonity)
    bootnode_bin= os.path.join(autonity_path,"build/bin/bootnode")
    autonity_bin= os.path.join(autonity_path,"build/bin/autonity")
    key_inspector_bin= os.path.join(autonity_path,"build/bin/ethkey")

    test_id = args.id
    COMMIT_HASH = args.hash
    print("start docker test for commit hash: %s", COMMIT_HASH)

    JOB_ID = "hash_{}_case_{}".format(COMMIT_HASH, test_id)

    # cleanup in case of test is killed by ci.
    signal.signal(signal.SIGTERM, receive_signal)

    try:
        # check to install docker, start docker daemon.
        init()
        check_docker_daemon()

        # release unused docker resources if there exists.
        prune_unused_images()
        prune_unused_volumes()
        prune_unused_network()

        # copy binary to binary dir for image building.
        utility.execute("cp {} ./test_bin/".format(autonity_bin))
        utility.execute("cp {} ./test_bin/".format(bootnode_bin))
        utility.execute("cp {} ./test_bin/".format(key_inspector_bin))

        # tyr to build autonity client image if the image hasn't being built.
        check_to_build_client_images()

        # create the network infra for the test.
        ips = start_client_containers(JOB_ID)

        # collect infra information and try to build test engine image.
        dump_ips_to_engine_conf(ips)
        check_to_build_engine_image(COMMIT_HASH)

        # start the test engine with specific test case id.
        container = start_test_engine_container(COMMIT_HASH, JOB_ID, test_id)

        if container is not None:
            thd = None
            for line in container.logs(stdout=True, stderr=True, stream=True):
                print(line.decode())
                if line == b"INFO - [TEST PASSED]\n":
                    exit_code = 0
                if line == b"INFO - [TEST FAILED]\n":
                    exit_code = 1
                    thd = threading.Thread(target=thread_func_copy_system_logs, args=(JOB_ID, SYSTEM_LOG_PATH))
                    thd.start()
            # wait if log collection is not finished.
            if thd is not None:
                thd.join(timeout=300)

    except Exception as e:
        print("e2e testing failed: ", e)
    finally:
        clean_up(JOB_ID, COMMIT_HASH)
        exit(exit_code)
