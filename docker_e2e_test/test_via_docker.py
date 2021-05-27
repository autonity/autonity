import docker
import utility
import ipaddress
import argparse
import threading
import signal
import time
import os

TEST_ENGINE_IMAGE_NAME = "enginehost{}/ubuntu"
TEST_ENGINE_DOCKER_FILE = "./Dockerfile"
CLIENT_IMAGE_LATEST = "clienthost/ubuntu:latest"
CLIENT_IMAGE_NAME = "clienthost/ubuntu"
CLIENT_DOCKER_FILE = "./clientDockerFile"
BUILDER_IMAGE_NAME = "go-builder/ubuntu"
BUILDER_DOCKER_FILE = "./builderDockerfile"
NUM_OF_CLIENT = 6
NODE_NAME = "Node{}_{}"
ENGINE_NAME = "Engine{}"
VALIDATOR_IP_LIST_FILE = "./etc/validator.ip"
COMMAND_START_TEST = "python3 e2etestengine.py ./bin/autonity"
FAILED_TEST_LOGS = "./JOB_{}.tar"
SYSTEM_LOG_PATH = "/system_log"
JOB_ID = ""


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
            if CLIENT_IMAGE_LATEST in image.tags:
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


def create_test_bed(job_id):
    ip_set = set()
    try:
        client = docker.from_env()
        for i in range(0, NUM_OF_CLIENT):
            node_name = NODE_NAME.format(job_id, i)
            # todo: remove the hard coding twins pair.
            hostname = "client{}".format(i)
            if i == 3:
                hostname = "client{}".format(0)

            container = client.containers.run(CLIENT_IMAGE_NAME, name=node_name, hostname=hostname,
                                              detach=True, privileged=True,
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


def remove_test_engine_image(job_id):
    try:
        client = docker.from_env()
        result = client.api.remove_image(TEST_ENGINE_IMAGE_NAME.format(job_id), force=True, noprune=True)
        print("docker remove image: ", result)
    except Exception as e:
        print("docker remove image: ", e)


# to stop and remove client and test engine containers.
def clean_test_bed_containers(job_id):
    print("start clean up current test context: test bed, test engine.")
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


def create_test_engine_image_per_run(job_id):
    print("start to build test engine image")
    create_image(TEST_ENGINE_IMAGE_NAME.format(job_id), TEST_ENGINE_DOCKER_FILE)


def dump_ips_to_engine_conf(ips):
    try:
        with open(VALIDATOR_IP_LIST_FILE, 'w+') as f:
            for ip in ips:
                f.write("{}\n".format(ip))
        print("test bed IP list was configured into test engine.")
    except Exception as e:
        print("failed to dump ip into test engine validator.ip file. ", e)


def start_test_engine_container(job_id):
    print("start test engine container, the testcase will be run in it.")
    try:
        print("test engine is going to start:")
        client = docker.from_env()
        container = client.containers.run(TEST_ENGINE_IMAGE_NAME.format(job_id), command=COMMAND_START_TEST,
                                          name=ENGINE_NAME.format(job_id), detach=True, privileged=True)
        print("test engine is started.")
        return container
    except Exception as e:
        print("create test engine container failed: ", e)


def clean_up(job_id):
    clean_test_bed_containers(job_id)
    prune_unused_images()
    prune_unused_volumes()
    prune_unused_network()
    remove_test_engine_image(job_id)


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
    clean_up(JOB_ID)
    raise SystemExit('Exiting')
    return


if __name__ == "__main__":
    exit_code = 1
    parser = argparse.ArgumentParser()
    parser.add_argument("autonity", help="Autonity WorkDir Path")
    args = parser.parse_args()
    job_id = str(time.time())
    JOB_ID = job_id
    autonity_path = os.path.abspath(args.autonity)
    bootnode_bin= os.path.join(autonity_path,"build/bin/bootnode")
    autonity_bin= os.path.join(autonity_path,"build/bin/autonity")

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


        # Build builder image to removed since docker in docker build fails in CI.
        #create_image(BUILDER_IMAGE_NAME, BUILDER_DOCKER_FILE)
        #container = docker.from_env().containers.run(BUILDER_IMAGE_NAME, name="go-builder",
        #    detach=False, remove=True, command='bash -c "cd autonity && make all"',
        #    volumes={autonity_path: {"bind": "/autonity", "mode": "rw"}})

        # copy binary to binary dir for image building.
        utility.execute("cp {} ./bin/".format(autonity_bin))
        utility.execute("cp {} ./bin/".format(bootnode_bin))

        # build autonity client image.
        check_to_build_client_images()

        # create test bed for the latter deployment.
        ips = create_test_bed(job_id)

        # prepare test engine image.
        dump_ips_to_engine_conf(ips)
        create_test_engine_image_per_run(job_id)

        # start the e2e testing.
        container = start_test_engine_container(job_id)

        if container is not None:
            thd = None
            for line in container.logs(stdout=True, stderr=True, stream=True):
                print(line.decode())
                if line == b"INFO - [TEST PASSED]\n":
                    exit_code = 0
                if line == b"INFO - [TEST FAILED]\n":
                    exit_code = 1
                    thd = threading.Thread(target=thread_func_copy_system_logs, args=(job_id, SYSTEM_LOG_PATH))
                    thd.start()
            # wait if log collecting is not finished.
            if thd is not None:
                thd.join(timeout=300)

    except Exception as e:
        print("e2e testing failed: ", e)
    finally:
        clean_up(job_id)
        exit(exit_code)
