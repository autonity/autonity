#!/usr/bin/python3
try:
    from Crypto.Hash import keccak
    sha3_256 = lambda x: keccak.new(digest_bits=256, data=x).digest()
except:
    import sha3 as _sha3
    sha3_256 = lambda x: _sha3.sha3_256(x).digest()

import argparse
import subprocess
import re
import os
import json
import rlp
from eth_utils.hexadecimal import decode_hex

from time import sleep

from web3.auto import w3
from web3 import Web3
from solcx import compile_files
import solcx
from typing import List


def generate_genesis(addresses: List[str], enodes: List[str]):
    ##########################################################################################
    #   The following parameters should not be modified unless you know what you're doing.   #
    genesis = {
        "config": {
            "homesteadBlock": 0,
            "eip150Block": 0,
            "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "eip155Block": 0,
            "eip158Block": 0,
            "byzantinumBlock": 0,
            "constantinopleBlock": 0,
            "petersburgBlock": 0,
            "policy": 0,
            "block-period": 1,
            "tendermint": {
                "policy": 0,
                "block-period": 1,
            },
            "autonityContract": {
                "deployer": "",
                "bytecode": "",  # use default bin and abi from client side.
                "abi": "",
                "minGasPrice": 0,
                "users": [],
            }
        },
        "nonce": "0x0",
        "timestamp": "0x0",
        "gasLimit": "0xffffffff",
        "difficulty": "0x1",
        "coinbase": "0x0000000000000000000000000000000000000000",
        "number": "0x0",
        "gasUsed": "0x0",
        "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "mixHash": "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
    }
    #                                                                                        #
    ##########################################################################################

    # Default balance
    starting_balance = "0x000000000000000000100000000000000000000000000000000000000000000"
    genesis["alloc"] = {}
    genesis["validators"] = []
    # let operator and deployer to be the same account.
    genesis["config"]["autonityContract"]["deployer"] = f"0x{addresses[0]}"
    genesis["config"]["autonityContract"]["operator"] = f"0x{addresses[0]}"

    for i in range(0, len(addresses)):
        user = {
            "enode": enodes[i],
            "address": f"0x{addresses[i]}",
            "type": "validator",
            "stake": 2 if i == 0 else 1,
        }
        genesis["alloc"][f"0x{addresses[i]}"] = {"balance": starting_balance}
        genesis["config"]["autonityContract"]["users"].append(user)

    genesis["config"]["chainId"] = 1
    try:
        with open("./network-data/genesis.json", 'w') as out:
            out.write(json.dumps(genesis, indent=4) + '\n')
    except IOError as e:
        print("cannot create genesis.json ", e)
        raise e


def execute(cmd):
    try:
        print("[CMD] {}".format(cmd))
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, encoding="utf-8", shell=True)
        return process.communicate(timeout=5)
    except Exception as e:
        print("Error: ", e)
        raise e


def create_dir(dir_path):
    if execute("mkdir -p {}".format(dir_path)) is None:
        raise IOError("cannot create dir")


def remove_dir(dir_path):
    if execute("rm -rf {}".format(dir_path)) is None:
        raise IOError("cannot remove dir")


def create_network_dir():
    try:
        remove_dir("./network-data")
        create_dir("./network-data")
        for i in range(0, node_count):
            create_dir("./network-data/node{}".format(i))
    except IOError as e:
        print("cannot create dir ", e)
        raise e


def generate_new_accounts():
    try:
        addresses = []
        for node_id in range(node_count):
            execute("echo 123 > ./network-data/node{}/pass.txt".format(node_id))
            output = execute(
                '{} --datadir "./network-data/node{}/data" --password "./network-data/node{}/pass.txt" account new'
                    .format(autonity_path, node_id, node_id)
            )
            print(output)
            m = re.findall(r'0x(.{40})', output[0], re.MULTILINE)
            if len(m) == 0:
                print("Aborting - account creation failed")
            addresses.append(m[0])

        with open("./network-data/addresses.json", 'w') as out:
            out.write(json.dumps(addresses, indent=4) + '\n')
        return addresses
    except IOError as e:
        print('cannot generate accounts ', e)
        raise e


def generate_enodes():
    try:
        enodes = []
        pubkeys = []
        for node_id in range(0, node_count):
            keystores_dir = "./network-data/node{}/data/keystore".format(node_id)
            keystore_file_path = keystores_dir + "/" + os.listdir(keystores_dir)[0]
            with open(keystore_file_path) as keyfile:
                encrypted_key = keyfile.read()
                account_private_key = w3.eth.account.decrypt(encrypted_key, "123").hex()[2:]
            with open("./network-data/node{}/boot.key".format(node_id), "w") as bootkey:
                bootkey.write(account_private_key)

            pub_key = \
            execute("{} -writeaddress -nodekey ./network-data/node{}/boot.key".format(bootnode_path, node_id))[
                0].rstrip()
            pubkeys.append(pub_key)
            port = 5000 + node_id
            enodes.append("enode://{}@127.0.0.1:{}".format(pub_key, port))
        return enodes
    except IOError as e:
        print("cannot generate enodes. ", e)
        raise e


def init_chains():
    try:
        for node_id in range(0, node_count):
            execute("""{} --datadir "./network-data/node{}/data/" init "./network-data/genesis.json" """
                    .format(autonity_path, node_id))
    except Exception as e:
        print("cannot init chain ", e)
        raise e


def reinit_chains():
    try:
        for node_id in range(0, node_count):
            remove_dir("./network-data/node{}/data/autonity".format(node_id))
            execute("""{} --datadir "./network-data/node{}/data/" init "./network-data/genesis.json" """
                    .format(autonity_path, node_id))
    except IOError as e:
        print("cannot reinit chain ", e)
        raise e


def tmux_start_clients(addresses, dont_start_id=None):
    try:
        execute("tmux new -s autonity -d")
        for node_id in range(0, node_count):
            if dont_start_id is not None and node_id == dont_start_id:
                continue
            execute("tmux new-window -t autonity:{} -n {}".format(node_id + 1, node_id))
            execute("""tmux send-keys -t autonity:{} "{}""".format(node_id + 1, autonity_path) +
                    """ --datadir ./network-data/node{}/data/""".format(node_id) +
                    """ --nodekey ./network-data/node{}/boot.key --syncmode 'full'""".format(node_id) +
                    """ --port {} --rpcport {} --rpc --rpcaddr '0.0.0.0' --ws --wsport {}""".format(5000 + node_id,
                                                                                                    6000 + node_id,
                                                                                                    7000 + node_id) +
                    """ --rpccorsdomain '*' --rpcapi 'personal,db,eth,net,web3,txpool,miner,tendermint,autonity,clique'"""
                    """ --networkid 1991  --gasprice '0' --allow-insecure-unlock --graphql --graphql.port {}""".format(
                        8000 + node_id) +
                    """ --unlock 0x{}""".format(addresses[node_id]) +
                    """ --password ./network-data/node{}/pass.txt --debug --mine --minerthreads '1'""".format(node_id) +
                    """ --etherbase 0x{} --verbosity 4" C-m """.format(addresses[node_id]))

            execute("""tmux split-window -h -t autonity:{}""".format(node_id + 1))
            execute("""tmux send-keys -t autonity:{} "sleep 10s" C-m""".format(node_id + 1))
            execute("""tmux send-keys -t autonity:{} "{} attach ipc:./network-data/node{}/data/autonity.ipc" C-m"""
                    .format(node_id + 1, autonity_path, node_id))
    except Exception as e:
        print("tmux cannot start clients ", e)
        raise e
    except IOError as e:
        print("cannot start client ", e)
        raise e


# install solc compiler.
def prepare_dependency():
    installed_solc = solcx.get_installed_solc_versions()
    if 'v0.5.11' in installed_solc:
        solcx.set_solc_version('v0.5.11')
    else:
        solcx.install_solc('v0.5.11')
        solcx.set_solc_version('v0.5.11')


def deploy_clients():
    global node_count
    global autonity_path
    global bootnode_path
    global addresses
    print("----------------------------------------------------")
    print("Autonity Local Network Deployment Utility")
    print("All rights reserved - Clearmatics Technologies Ltd.")
    print("----------------------------------------------------")

    try:
        parser = argparse.ArgumentParser()
        parser.add_argument("autonity", help='Autonity Binary Path', type=str, default='../build/bin/autonity')
        parser.add_argument("-n", help='Number of nodes', type=int, default=4)
        parser.add_argument("-r", help='Restart All', action="store_true")
        parser.add_argument("-o", help='Restart All except', type=int)
        parser.add_argument("-i", help='Reinit chains', action="store_true")
        args = parser.parse_args()

        node_count = args.n
        autonity_path = args.autonity

        bootnode_path = autonity_path.split("/")
        bootnode_path[len(bootnode_path) - 1] = "bootnode"
        bootnode_path = "/".join(bootnode_path)

        if args.o is not None:
            print("===== RESTART =====")
            execute("tmux kill-session -t autonity")
            print("===== REINIT CHAINS=====")
            reinit_chains()
            with open('./network-data/addresses.json') as f:
                addresses = json.load(f)
                print("===== STARTING CLIENTS =====")
                tmux_start_clients(addresses, args.o)
            return
        if args.r:
            print("===== KILL OLD SESSION=====")
            execute("tmux kill-session -t autonity")
            with open('./network-data/addresses.json') as f:
                addresses = json.load(f)
                print("===== STARTING CLIENTS =====")
                tmux_start_clients(addresses)
            return

        if args.i:
            print("===== KILL OLD SESSION=====")
            execute("tmux kill-session -t autonity")
            print("===== REINIT CHAINS=====")
            reinit_chains()
            return

        print("===== SETUP INITIALIZATION =====")
        create_network_dir()
        print("===== ACCOUNTS CREATION =====")
        addresses = generate_new_accounts()
        print(addresses)
        print("===== ENODES GENERATION =====")
        enodes = generate_enodes()
        print(enodes)
        print("===== GENESIS GENERATION =====")
        generate_genesis(addresses, enodes)
        print("===== CHAIN INITIALIZATION =====")
        init_chains()
        print("===== SETUP FINISHED =====")

        print("===== STARTING CLIENTS =====")
        tmux_start_clients(addresses)
    except Exception as e:
        print("cannot deploy the network, ", e)
        raise e


# connect to first node as default endpoint.
def get_http_end_point():
    return "http://0.0.0.0:6000"


# get autonity contract bin, abi, and a new version with v0.0.1.
def compile_contract():
    try:
        contract = compile_files(["../contracts/autonity/contract/contracts/Autonity.sol"])
        bin = contract["../contracts/autonity/contract/contracts/Autonity.sol:Autonity"]["bin"]
        abi = contract["../contracts/autonity/contract/contracts/Autonity.sol:Autonity"]["abi"]
        # to do get version from release note from Devop team.
        version = "v0.0.1"
        return bin, abi, version
    except Exception as e:
        print("cannot compile contract ", e)
        raise e


def to_string(value):
    if isinstance(value, bytes):
        return value
    if isinstance(value, str):
        return bytes(value, 'utf-8')
    if isinstance(value, int):
        return bytes(str(value), 'utf-8')


def sha3(seed):
    return sha3_256(to_string(seed))


def normalize_address(x, allow_blank=False):
    if allow_blank and x == '':
        return ''
    if len(x) in (42, 50) and x[:2] == '0x':
        x = x[2:]
    if len(x) in (40, 48):
        x = decode_hex(x)
    if len(x) == 24:
        assert len(x) == 24 and sha3(x[:20])[:4] == x[-4:]
        x = x[:20]
    if len(x) != 20:
        raise Exception("Invalid address format: %r" % x)
    return x


# get autonity contract address with deployer account and nonce 0.
# This is same as autonity set the address of autonity contract.
def get_autonity_contract_address():
    sender = get_system_deployer_account()
    nonce = 0
    return sha3(rlp.encode([normalize_address(sender), nonce]))[12:]


def get_system_operator_account():
    try:
        with open('./network-data/addresses.json') as address_file:
            data = json.load(address_file)
            return f"0x{data[0]}"
    except IOError as e:
        print("cannot read file ", e)
        raise e


def get_validators_from_config():
    try:
        with open('./network-data/addresses.json') as address_file:
            data = json.load(address_file)
            res = []
            for item in data:
                res.append("0x{}".format(item))
            return res
    except IOError as e:
        print("cannot read file ", e)
        raise e


def get_system_deployer_account():
    try:
        with open('./network-data/addresses.json') as address_file:
            data = json.load(address_file)
            return f"0x{data[0]}"
    except IOError as e:
        print("cannot read file ", e)
        raise e


# bootstrap timer after network setup and before test case starts.
def wait_for_network_bootstrap():
    i = 10
    while i:
        sleep(1)
        i -= 1
        print("waiting for {} seconds for network bootstrap.".format(i))


# get autontiy client and contract object
def get_autonity_and_contract():
    try:
        # connect to autonity.
        end_point = get_http_end_point()
        autonity = Web3(Web3.HTTPProvider(end_point, request_kwargs={'timeout': 60}))
    except Exception as e:
        print("cannot connect to endpoint. ", e)
        raise e
    try:
        # construct contract object.
        addr = get_autonity_contract_address()
        byte_code, abi, version = compile_contract()
        autonity_contract = autonity.eth.contract(address=addr, abi=abi)
    except Exception as e:
        print("cannot create contract object from client. ", e)
        raise e
    return autonity, autonity_contract


#####################################################################
#          Your new e2e test cases implements from here             #
def test_upgrade_autonity_contract():
    try:
        autonity, contract = get_autonity_and_contract()
    except Exception as e:
        print("cannot connect to client. ", e)
        raise Exception("test_upgrade_autonity_contract failed.")

    try:
        start_height = autonity.eth.getBlock("latest")['number']
        print("test started at height: ", start_height)
        # get contract version from chain before upgrade.
        previous_version = contract.functions.getVersion().call()
        print("get contract previous version: ", previous_version)
        byte_code, abi, version = compile_contract()
        abi = json.dumps(abi)
        gas_estimate = contract.functions.upgradeContract(byte_code, abi, version).estimateGas()
        print("gas estimate to transact with ugradeContract: {}".format(gas_estimate))
        deployer = get_system_deployer_account()
        contract.functions.upgradeContract(byte_code, abi, version).transact({'from': deployer, 'gas': 999999999})
        # wait TX to be mined.
        sleep(5)
        version_new = contract.functions.getVersion().call()
        print("get contrace new version: ", version_new)
        if version != version_new:
            raise Exception("test_upgrade_autonity_contract failed")
        end_height = autonity.eth.getBlock("latest")['number']
        print("test endedat at height: ", end_height)
        if end_height <= start_height:
            print("chain on-hold after contract upgrade")
            raise Exception("test_upgrade_autonity_contract failed")
        print("test_upgrade_autonity_contract passed.")

    except Exception as e:
        print("test case failed: ", e)
        raise Exception("test_upgrade_autonity_contract failed.")


def test_get_validators():
    try:
        autonity, contract = get_autonity_and_contract()
        start_height = autonity.eth.getBlock("latest")['number']
        print("test started at height: ", start_height)
        # get validators from contract
        validators = contract.functions.getValidators().call()
        print("validators: ", validators)

        # get validators from configuration
        expected_validators = get_validators_from_config()
        print("expected validators: ", expected_validators)

        # compare two list
        result = set(validators) & set(expected_validators)
        if len(result) != len(validators):
            print("wrong validators set expected.")
            raise Exception("test_get_validators failed.")

        end_height = autonity.eth.getBlock("latest")['number']
        print("test started at height: ", end_height)
        print("test_get_validators passed.")
    except Exception as e:
        print("cannot connect to client. ", e)
        raise Exception("test_get_validators failed.")


def run_tests():
    try:
        test_upgrade_autonity_contract()
        test_get_validators()
    except Exception as e:
        raise e


if __name__ == "__main__":

    # before test, trying to clean up test-bed context.
    execute("tmux kill-session -t autonity")

    try:
        deploy_clients()
        wait_for_network_bootstrap()
    except Exception as e:
        print("cannot deploy network", e)
        exit(1)

    try:
        prepare_dependency()
    except Exception as e:
        print("cannot install dependencies ", e)
        exit(1)

    try:
        run_tests()
    except Exception as e:
        print("test case failed: ", e)
        exit(1)
