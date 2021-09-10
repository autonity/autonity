import eth_utils
from pytest import fixture

from lib import autonity, json_rpc, utils


@fixture
def node(worker_id):
    if worker_id == "master":
        worker_no = 0
    else:
        # Workers are named as 'gw0', 'gw1', ...
        worker_no = int(worker_id.replace("gw", ""))

    autonity_node = autonity.launch(worker_no)
    yield autonity_node
    autonity.terminate(autonity_node)


@fixture
def web3(node):
    return json_rpc.Client(node.rpc_port)


@fixture
def tx_hash(web3):
    return web3.request("eth_sendTransaction", [{
        "from": autonity.ACCOUNT,
        "to": autonity.OTHER_ACCOUNT,
        "gas": eth_utils.to_hex(1000000),
        "gasPrice": eth_utils.to_hex(5000),
        "value": eth_utils.to_hex(1000000),
        "nonce": "0x0",
    }])


@fixture
def tx_receipt(web3, tx_hash):
    return utils.repeat_until(
        lambda: web3.request("eth_getTransactionReceipt", [tx_hash]), 3)
