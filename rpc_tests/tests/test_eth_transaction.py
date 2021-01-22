import eth_utils

from lib import autonity, utils, validators


def test__eth_signTransaction(web3):
    """https://eth.wiki/json-rpc/API#eth_signtransaction"""
    result = web3.request("eth_signTransaction", [{
        "from": autonity.ACCOUNT,
        "to": autonity.OTHER_ACCOUNT,
        "gas": eth_utils.to_hex(1000000),
        "gasPrice": eth_utils.to_hex(5000),
        "value": eth_utils.to_hex(1000000),
        "nonce": "0x0",
    }])
    validators.HexString().validate(result)


def test__eth_sendTransaction(web3):
    """https://eth.wiki/json-rpc/API#eth_sendtransaction"""
    result = web3.request("eth_sendTransaction", [{
        "from": autonity.ACCOUNT,
        "to": autonity.OTHER_ACCOUNT,
        "gas": eth_utils.to_hex(1000000),
        "gasPrice": eth_utils.to_hex(5000),
        "value": eth_utils.to_hex(1000000),
        "nonce": "0x0",
    }])
    validators.HexString().validate(result)


def test__eth_getTransactionReceipt(web3, tx_hash):
    """https://eth.wiki/json-rpc/API#eth_gettransactionreceipt"""
    result = utils.repeat_until(
        lambda: web3.request("eth_getTransactionReceipt", [tx_hash]), 3)
    validators.Object({
        "transactionHash": validators.HexString(32),
        "transactionIndex": validators.HexString(),
        "blockHash": validators.HexString(32),
        "blockNumber": validators.HexString(),
        "from": validators.HexString(20),
        "to": validators.HexString(20),
        "cumulativeGasUsed": validators.HexString(),
        "gasUsed": validators.HexString(),
        "contractAddress": validators.Null(),
        "logs": validators.Array(),
        "logsBloom": validators.HexString(256),
    }).validate(result)


def test__eth_getTransactionReceipt__no_receipt(web3):
    """https://eth.wiki/json-rpc/API#eth_gettransactionreceipt"""
    result = web3.request("eth_getTransactionReceipt", [
        "0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"
    ])
    validators.Null().validate(result)
