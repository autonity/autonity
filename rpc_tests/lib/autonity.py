import atexit
import json
import logging
import os
import pprint
import shutil
import signal
import socket
import tempfile
from collections import namedtuple
from os import path
from subprocess import STDOUT, Popen

from . import utils

ACCOUNT = "0x7d65dda8e1e3997d09542c62cc72ceadd186dc28"
OTHER_ACCOUNT = "0x1fea36661676734c9bd155311fe9449dc98de990"

_THIS_DIR = path.dirname(__file__)
BUILD_DIR = f"{_THIS_DIR}/../../build/bin"
CONFIG_DIR = f"{_THIS_DIR}/../config"
BASE_NET_PORT = 30303
BASE_RPC_PORT = 8545

Node = namedtuple("Node", ["process", "datadir", "logfile", "rpc_port"])

logger = logging.getLogger(__name__)


def launch(worker_no: int) -> Node:
    """Sets up and launces an Autonity node and returns a `Node` instance.

    `worker_no` is the number of the Pytest executor.
    """
    datadir = tempfile.mkdtemp(prefix="autonity.")
    genesis_file = f"{datadir}/genesis.json"
    logfile = tempfile.NamedTemporaryFile(mode="w", encoding="utf-8")

    # For parallel test execution
    net_port = BASE_NET_PORT + worker_no
    rpc_port = BASE_RPC_PORT + worker_no
    _create_genesis(genesis_file, net_port)

    cmd = [
        f"{BUILD_DIR}/autonity",
        "--datadir", datadir,
        "--genesis", genesis_file,
        "--keystore", f"{CONFIG_DIR}/keystore",
        "--nodekey", f"{CONFIG_DIR}/node.key",
        "--password", f"{CONFIG_DIR}/password.txt",
        "--syncmode", "full",
        "--port", str(net_port),
        "--http",
        "--http.port", str(rpc_port),
        "--http.addr", "127.0.0.1",
        "--http.corsdomain", "*",
        "--networkid", "1756",
        "--allow-insecure-unlock",
        "--unlock", ACCOUNT,
        "--debug",
        "--verbosity", "4",
        "--mine",
        "--miner.threads", "1",
        "--miner.etherbase", ACCOUNT,
    ]
    logger.debug(pprint.pformat(cmd))
    process = Popen(cmd, start_new_session=True, stdout=logfile, stderr=STDOUT)

    node = Node(process, datadir, logfile, rpc_port)
    # Pytest doesn't run the teardown fixture if the test run is interrupted
    atexit.register(terminate, node)

    _wait_for_port_to_open(net_port)
    _wait_for_port_to_open(rpc_port)
    _wait_for_mining_to_start(node)
    return node


def terminate(node: Node):
    """Kills an Autonity node and cleans up the workspace."""
    atexit.unregister(terminate)
    # Kill the entire process group to make sure all children are terminated
    os.killpg(node.process.pid, signal.SIGKILL)
    _wait_for_process_to_terminate(node.process)
    with open(node.logfile.name, encoding="utf-8") as f:
        print(f.readline())
    node.logfile.close()
    shutil.rmtree(node.datadir, ignore_errors=True)


def _create_genesis(genesis_file: str, net_port: int):
    genesis_template = f"{CONFIG_DIR}/genesis.template.json"
    with open(genesis_template) as f:
        genesis = json.load(f)
    user = genesis["config"]["autonityContract"]["users"][0]
    user["enode"] = user["enode"].replace("{{port}}", str(net_port))
    with open(genesis_file, "w") as f:
        json.dump(genesis, f)


def _wait_for_port_to_open(port: int, timeout: float = 5):
    if not utils.repeat_until(lambda: _is_port_used(port), timeout=timeout):
        raise TimeoutError(
            f"Autonity did not start listening on port {port} "
            f"within {timeout} seconds")


def _wait_for_mining_to_start(node: Node, timeout: float = 10):
    with open(node.logfile.name, encoding="utf-8") as f:
        if not utils.repeat_until(
                lambda: "Commit new mining work" in f.readline(),
                timeout=timeout,
                interval=0):
            raise TimeoutError(
                f"Autonity did not start mining within {timeout} seconds")


def _wait_for_process_to_terminate(process: Popen, timeout: float = 5):
    if not utils.repeat_until(
            lambda: process.poll() is not None, timeout=timeout):
        raise TimeoutError(
            f"Autonity did not terminate within {timeout} seconds")


def _is_port_used(port: int) -> bool:
    with socket.socket() as sock:
        return sock.connect_ex(("localhost", port)) != socket.errno.ECONNREFUSED
