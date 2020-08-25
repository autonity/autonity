## Autonity Hello World

Note all commands provided here assume the current directory is the helloworld
directory.

### tl;dr;

Run `docker-compose up -d` and off you go!

### What is this?

A simple script to start an Autonity network with IBFT.

### How do I run it?

You will need Docker and Docker-Compose. The versions we used while in development were:

```bash
$ docker -v
Docker version 18.09.2-ce, build 62479626f2

$ docker-compose -v
docker-compose version 1.23.1, build b02f1306
```

To deploy the network just run:

```bash
$ docker-compose up -d
```

### What should it look like when it is running?

When you first start the `docker-compose up -d` it should print out this information:

```bash
$ docker-compose up -d
WARNING: The Docker Engine you're using is running in swarm mode.

Compose does not use swarm mode to deploy services to multiple nodes in a swarm. All containers will be scheduled on the current node.

To deploy your application across the swarm, use `docker stack deploy`.

Creating network "helloworld_chainnet" with driver "bridge"
Creating autonity-node-2 ... done
Creating autonity-node-1 ... done
Creating autonity-node-5 ... done
Creating autonity-node-3 ... done
Creating autonity-node-4 ... done
Creating nodes-connector ... done
```

When the nodes have all been deployed and connected to each other, the `nodes-connector` should have exited. You can check this by doing the `ps` command:

```bash
$ docker-compose ps
     Name               Command         State                                                           Ports
----------------------------------------------------------------------------------------------------------------------------------------------------------------------
autonity-node-1   ./autonity-start.sh   Up      0.0.0.0:30313->30303/tcp, 0.0.0.0:30313->30303/udp, 0.0.0.0:8541->8545/tcp, 8546/tcp, 8547/tcp, 0.0.0.0:8641->8645/tcp
autonity-node-2   ./autonity-start.sh   Up      0.0.0.0:30323->30303/tcp, 0.0.0.0:30323->30303/udp, 0.0.0.0:8542->8545/tcp, 8546/tcp, 8547/tcp, 0.0.0.0:8642->8645/tcp
autonity-node-3   ./autonity-start.sh   Up      0.0.0.0:30333->30303/tcp, 0.0.0.0:30333->30303/udp, 0.0.0.0:8543->8545/tcp, 8546/tcp, 8547/tcp, 0.0.0.0:8643->8645/tcp
autonity-node-4   ./autonity-start.sh   Up      0.0.0.0:30343->30303/tcp, 0.0.0.0:30343->30303/udp, 0.0.0.0:8544->8545/tcp, 8546/tcp, 8547/tcp, 0.0.0.0:8644->8645/tcp
autonity-node-5   ./autonity-start.sh   Up      30303/tcp, 30303/udp, 8545/tcp, 8546/tcp, 8547/tcp
```

### How can I use the nodes?

You can connect to the nodes, through the autonity console all the RPC ports
are have been mapped to the host.

Here is an example of attaching a console to `autonity-node-1` via docker, note
that 172.25.0.11 is the ip address assigned to the container by docker and is
defined in `docker-compose.yml`

```bash
$
docker run --network helloworld_chainnet -ti --rm autonity attach http://172.25.0.11:8545
Welcome to the Autonity JavaScript console!

instance: Autonity/v0.5.0-b4d1f51f-20200812/linux-amd64/go1.14.7
coinbase: 0x850c1eb8d190e05845ad7f84ac95a318c8aab07f
at block: 414 (Tue Aug 25 2020 12:55:11 GMT+0000 (UTC))
 datadir: /autonity-data
 modules: admin:1.0 debug:1.0 eth:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 tendermint:1.0 txpool:1.0 web3:1.0

>
```

Here is an example of attaching a console to `autonity-node-1` from the host
machine, note that the host port for a container that corresponds to 8545 can
be found from the output of `docker-compose ps` in this case the port is 8541.

```bash
$
../build/bin/autonity attach http://0.0.0.0:8541
Welcome to the Autonity JavaScript console!

instance: Autonity/v0.5.0-b4d1f51f-20200812/linux-amd64/go1.14.7
coinbase: 0x850c1eb8d190e05845ad7f84ac95a318c8aab07f
at block: 1181 (Tue Aug 25 2020 14:07:58 GMT+0100 (BST))
 datadir: /autonity-data
 modules: admin:1.0 debug:1.0 eth:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 tendermint:1.0 txpool:1.0 web3:1.0

>
```

You can also run a simple Javascript command without having an interactive console:

```bash
$ docker run --network helloworld_chainnet -ti --rm autonity attach http://172.25.0.11:8545 --exec '[eth.coinbase, eth.getBlock("latest").number, eth.getBlock("latest").hash, eth.mining]'
["0x850c1eb8d190e05845ad7f84ac95a318c8aab07f", 493, "0xcffb1c661b4bd87430079656fa8b233fb0a0585250282f46506b7d44151560f0", true]
```

### What are all these files in the `helloword` directory?

The files in the `helloworld` directory are used to deploy and run the network, you can alter them and reploy to see how the changes affected the network. Here is the file list:

```bash
$ ls -lh
total 60K
-rwxr-xr-x 1 clearmatics clearmatics 2.2K Feb 13 15:12 autonity-connect.sh
-rwxr-xr-x 1 clearmatics clearmatics  577 Feb 13 15:13 autonity-start.sh
-rw-r--r-- 1 clearmatics clearmatics 3.1K Feb 13 00:15 docker-compose.yml
-rw-r--r-- 1 clearmatics clearmatics  410 Feb 12 16:35 Dockerfile
-rw-r--r-- 1 clearmatics clearmatics 1.4K Feb 13 13:41 genesis-clique.json
-rw-r--r-- 1 clearmatics clearmatics 2.6K Feb 12 12:15 genesis-ibft.json
drwx------ 2 clearmatics clearmatics 4.0K Feb 12 14:19 keystore
-rw-r--r-- 1 clearmatics clearmatics   65 Feb 12 23:40 nodekey1
-rw-r--r-- 1 clearmatics clearmatics   65 Feb 12 23:40 nodekey2
-rw-r--r-- 1 clearmatics clearmatics   65 Feb 12 23:40 nodekey3
-rw-r--r-- 1 clearmatics clearmatics   65 Feb 12 23:41 nodekey4
-rw-r--r-- 1 clearmatics clearmatics   65 Feb 12 23:41 nodekey5
-rw-r--r-- 1 clearmatics clearmatics 3.3K Feb 13 15:35 README.md
```

* `Dockerfile` is used by Docker to build the image, that will be reused everytime you deploy a container
* `docker-compose.yml` is used by Docker-Compose and it describes how the nodes should be deployed (what are the cointaner names, what images should be used, what is the order of deployment)
* `autonity-start.sh` script to start an autonity node, used evertime a container is deployed
* `autonity-connect.sh` script run everytime the `autonity-connector` container is started (it connects 5 nodes to ech other, sets the coinbase value, and starts the miner)
* `keystore` directory with all the keystores (keystores are used to keep the private keys of the accounts, our keystores all use the password `test`)
* `nodekey1` file containing Node Key used to generate ENode (this way the enodes never change, although it is not relevant for the Clique Hello World, it will be used in the future for the IBFT Hello World)

### How can the validator set be changed?

There are two ways to update the validator set:

1. Update the Soma and Glienicke smart contracts
2. Update the `nodekey` files
3. Change the `genesis-ibft.json`

#### Update Glienicke and Soma contract

The _Glienick_ contract is responsible for making sure that only nodes in its list are able to connect to the Autonity client.

In the default Docker Compose deployment the contract can be found at the `0x522B3294E6d06aA25Ad0f1B8891242E335D3B459` address. You can find the contract deployed in the Autonity code in the [`contracts`](https://github.com/clearmatics/autonity/tree/master/contracts/Glienicke) directory.

The _Soma_ contract allows anyone to vote on the IBFT set of validators.

In the default Docker Compose deployment the contract can be found at the `0xc3d854209eF19803954916F2fe4712448094363e` address. You can find the contract deployed in the Autonity code in the [`contracts`](https://github.com/clearmatics/autonity/tree/master/contracts/Soma) directory.

#### Change the `genesis-ibft.json` and update the `nodekey` files

_The Autonity Hello World limits the amount of validators to 4, but in a real world application you can have more validators_

It is possible update the set of validators by updating the genesis file and the nodekey files, the steps needed are:

1. Update the `nodekey1` file (or 2,3,4) with the private key of the validator
2. Update the `enodeWhitelist` property in the genesis file. Enode address can be a few formats:
* Ethereum enodeV4 `enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303`
* with domain instead of IP `enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@domain.com:30303`
* any of the above without port `enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@domain.com`
* by default, if it's not specified port `30303` will be used.
3. Update the `validators` property in the genesis file by with a proper node ID, eg:

```
"validators": [
    "0x850C1Eb8D190e05845ad7F84ac95a318C8AaB07f",
    "0x4AD219b58a5b46A1D9662BeAa6a70DB9F570deA5",
    "0x4B07239Bd581d21AEfcdEe0c6dB38070F9A5FD2D",
    "0xc443C6c6AE98F5110702921138D840e77dA67702",
    "0x09428e8674496e2d1e965402f33a9520c5fcbbe2"
]
```

The `validators` has higher priority compare to `extraData` and if both are specified, than `extraData` will be rewritten.

### What are the keystore passwords?

All the keystores use the same password: `test` (*please do not use in any production enviroment*)

## Tutorial

So you have nice and running 5 node cluster. Let's examine it and try _Soma_ and _Glienicke_ features.
First of all we need to connect to one of the nodes: `autonity attach http://0.0.0.0:8541`

In the Genesis file `genesis-ibft.json` we've defined the while list of nodes that have permission to connect to the network:
```bash
    [
        "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
        "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
        "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
        "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
        "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
    ]
```

We have 5 validators so we expect to see 4 connected peers of each node:

```bash
net.peerCount
    4
```

Let's check what are this peers:
```bash
admin.peers.map(function(peer) {return peer.enode;});
    [  
       "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:59360",
       "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
       "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303",
       "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303"
    ]
```

### Validators
We can get the list of actual validator for any block:
```bash
// for 10th block
istanbul.getValidators(web3.toHex(10))

// for current block
istanbul.getValidators(web3.toHex(eth.blockNumber))

    [
       "0x4ad219b58a5b46a1d9662beaa6a70db9f570dea5",
       "0x4b07239bd581d21aefcdee0c6db38070f9a5fd2d",
       "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
       "0xc443c6c6ae98f5110702921138d840e77da67702"
    ]
```

If you want you can get the same information by hash:
```bash
//for current block
istanbul.getValidatorsAtHash(eth.getBlock(eth.blockNumber).hash)
```

#### Adding and removing from validators
First of all we need to get _Soma_ contract ABI and address:

```bash
var somaContractAbi = eth.contract(JSON.parse(istanbul.getSomaContractABI()));
var somaContract = somaContractAbi.at(istanbul.getSomaContractAddress());

// check somaContract
somaContract

    ...
    AddValidator: function(),
    RemoveValidator: function(),
    allEvents: function(),
    getValidators: function(),
    validators: function()
    ...
    
```

Now it's possible to use `somaContract` object to call _Soma_.

##### Add a validator

```
// getValidators
somaContract.getValidators();

    [
       "0x4ad219b58a5b46a1d9662beaa6a70db9f570dea5",
       "0x4b07239bd581d21aefcdee0c6db38070f9a5fd2d",
       "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
       "0xc443c6c6ae98f5110702921138d840e77da67702"
    ]

// AddValidator
web3.personal.unlockAccount(eth.accounts[0], 'test');
somaContract.AddValidator("0x000000000000000000000000000000", {from: eth.accounts[0]});

somaContract.getValidators()
    [
       "0x4ad219b58a5b46a1d9662beaa6a70db9f570dea5",
       "0x4b07239bd581d21aefcdee0c6db38070f9a5fd2d",
       "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
       "0xc443c6c6ae98f5110702921138d840e77da67702",
       "0x0000000000000000000000000000000000000000"
    ]
```

A new validator `0x0000000000000000000000000000000000000000` has been added.

If you try to add an incorrect node ID, you get an error:
```bash
somaContract.AddValidator("incorrect_ID", {from: eth.accounts[0]});

    Error: new BigNumber() not a number: incorrect_ID
```

##### Remove a validator
```bash
somaContract.RemoveValidator("0x0000000000000000000000000000000000000000", {from: eth.accounts[0]});

somaContract.getValidators()

    [
       "0x4ad219b58a5b46a1d9662beaa6a70db9f570dea5",
       "0x4b07239bd581d21aefcdee0c6db38070f9a5fd2d",
       "0x850c1eb8d190e05845ad7f84ac95a318c8aab07f",
       "0xc443c6c6ae98f5110702921138d840e77da67702"
    ]
```

### Permissioned network
As it was done for _Soma_ we need to get _Glienicke_ contract:
``` 
var glienickeContractAbi = eth.contract(JSON.parse(istanbul.getGlienickeContractABI()));
var glienickeContract = glienickeContractAbi.at(istanbul.getGlienickeContractAddress());

glienickeContract;

    ...
    transactionHash: null,
    AddEnode: function(),
    RemoveEnode: function(),
    allEvents: function(),
    compareStringsbyBytes: function(),
    enodes: function(),
    getWhitelist: function()
    ...
```

#### Remove and add a user to the white list
The current white list can be gotten:
```
istanbul.getWhitelist();

    [
       "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
       "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
       "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
       "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
       "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
    ]
```

Lets remove one peer from white list and check that the node will be dropped.
At the moment we have 4 connections on each node:
```
// current network connections
net.peerCount;

    4

// list of connected peers
admin.peers.map(function(peer) {return peer.enode;});

    [
       "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:57262",
       "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
       "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:60654",
       "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303"
    ]

// current white list
istanbul.getWhitelist();

    [
       "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
       "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
       "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
       "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
       "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
    ]
    
```

To remove a peer from white list we should use `glienickeContract`:
```
glienickeContract.RemoveEnode("enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303", {from: eth.accounts[0], gas: 100000000});

net.peerCount;

    3

admin.peers.map(function(peer) {return peer.enode;});

    [
       "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
       "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:40640",
       "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:53914"
    ]

istanbul.getWhitelist();

    [
       "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
       "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
       "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
       "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
    ]

```

The connection with the removed peer won't be established until the peer will be added to white list again. If this happened, the connection will be established in few seconds:

```
glienickeContract.AddEnode("enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303", {from: eth.accounts[0], gas: 100000000});

net.peerCount;

    4

admin.peers.map(function(peer) {return peer.enode;});

    [
       "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
       "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:40640",
       "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:53914",
       "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:33168"
    ]


istanbul.getWhitelist();

    [
       "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
       "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
       "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
       "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303",
       "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303"
    ]

```

#### Error handling

If we try to add an incorrect enode:
```bash
glienickeContract.AddEnode("incorrect_Enode", {from: eth.accounts[0]});
```

The error should be logged in cluster. To get logs run the command `docker-compose logs | grep "ERROR"`:
```
ERROR[04-11|08:48:09.034] Invalid whitelisted enode                returned enode=incorrect_Enode error="invalid URL scheme, want \"enode\""
```
