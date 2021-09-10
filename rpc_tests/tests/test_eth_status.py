import eth_utils

from lib import autonity, validators


def test__eth_protocolVersion(web3):
    """https://eth.wiki/json-rpc/API#eth_protocolversion"""
    result = web3.request("eth_protocolVersion", [])
    assert isinstance(result, str)
    assert result != ""


def test__eth_syncing(web3):
    """https://eth.wiki/json-rpc/API#eth_syncing"""
    result = web3.request("eth_syncing", [])
    validators.Boolean().validate(result)
    assert result is False


def test__eth_coinbase(web3):
    """https://eth.wiki/json-rpc/API#eth_coinbase"""
    result = web3.request("eth_coinbase", [])
    validators.HexString(20).validate(result)
    assert result == autonity.ACCOUNT


def test__eth_mining(web3):
    """https://eth.wiki/json-rpc/API#eth_mining"""
    result = web3.request("eth_mining", [])
    validators.Boolean().validate(result)
    assert result is True


def test__eth_hashrate(web3):
    """https://eth.wiki/json-rpc/API#eth_hashrate"""
    result = web3.request("eth_hashrate", [])
    validators.HexString().validate(result)
    assert result == eth_utils.to_hex(0)


def test__eth_gasPrice(web3):
    """https://eth.wiki/json-rpc/API#eth_gasprice"""
    result = web3.request("eth_gasPrice", [])
    validators.HexString().validate(result)
    assert result == eth_utils.to_hex(5000)


def test__eth_accounts(web3):
    """https://eth.wiki/json-rpc/API#eth_accounts"""
    result = web3.request("eth_accounts", [])
    validators.Array(validators.HexString(20), 1).validate(result)
    assert result[0] == autonity.ACCOUNT


def test__eth_blockNumber(web3):
    """https://eth.wiki/json-rpc/API#eth_blocknumber"""
    result = web3.request("eth_blockNumber", [])
    validators.HexString().validate(result)
    assert result == eth_utils.to_hex(1)


def test__eth_getBalance(web3):
    """https://eth.wiki/json-rpc/API#eth_getbalance"""
    result = web3.request("eth_getBalance", [autonity.ACCOUNT, "latest"])
    validators.HexString().validate(result)
    assert result == "0x200000000000000000000000000000000000000000000000000000000000000"


def test__eth_getStorageAt(web3):
    """https://eth.wiki/json-rpc/API#eth_getstorageat"""
    result = web3.request(
        "eth_getStorageAt", [autonity.ACCOUNT, "0x0", "latest"])
    validators.HexString().validate(result)
    assert result == "0x0000000000000000000000000000000000000000000000000000000000000000"


def test__eth_getTransactionCount(web3):
    """https://eth.wiki/json-rpc/API#eth_gettransactioncount"""
    result = web3.request(
        "eth_getTransactionCount", [autonity.ACCOUNT, "latest"])
    validators.HexString().validate(result)
    assert result == "0x0"


def test__eth_getCode(web3):
    """https://eth.wiki/json-rpc/API#eth_getcode"""
    result = web3.request("eth_getCode", [autonity.ACCOUNT, "latest"])
    validators.HexString().validate(result)
    assert result == "0x"


def _validate_block_object_fields(result):
    validators.Object({
        "number": validators.HexString(),
        "hash": validators.HexString(32),
        "parentHash": validators.HexString(32),
        "nonce": validators.HexString(8),
        "sha3Uncles": validators.HexString(32),
        "logsBloom": validators.HexString(256),
        "transactionsRoot": validators.HexString(32),
        "stateRoot": validators.HexString(32),
        "receiptsRoot": validators.HexString(32),
        "miner": validators.HexString(20),
        "difficulty": validators.HexString(),
        "totalDifficulty": validators.HexString(),
        "extraData": validators.HexString(),
        "size": validators.HexString(),
        "gasLimit": validators.HexString(),
        "gasUsed": validators.HexString(),
        "timestamp": validators.HexString(),
        "transactions": validators.Array(
            validators.HexString(32), 1),
        "uncles": validators.Array(),
    }).validate(result)


def test_eth_getBlockByNumber__transaction_hashes(web3, tx_receipt):
    """https://eth.wiki/json-rpc/API#eth_getblockbynumber"""
    result = web3.request(
        "eth_getBlockByNumber", [tx_receipt["blockNumber"], False])
    _validate_block_object_fields(result)


def test_eth_getBlockByNumber__no_such_block(web3):
    """https://eth.wiki/json-rpc/API#eth_getblockbynumber"""
    result = web3.request("eth_getBlockByNumber", ["0xF", True])
    validators.Null().validate(result)


def test_eth_getBlockByHash__transaction_hashes(web3, tx_receipt):
    result = web3.request(
        "eth_getBlockByHash", [tx_receipt["blockHash"], False])
    _validate_block_object_fields(result)


def test_eth_getBlockByHash__no_such_hash(web3):
    result = web3.request(
        "eth_getBlockByHash",
        ["0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae",
         True])
    validators.Null().validate(result)
