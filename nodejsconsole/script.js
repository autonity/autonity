let Web3 = require('web3');

// check that both url and port are specified
if (process.argv[2].split(":").length != 3){
  console.error("The supplied url is malformed. It should be in the format <protocol>://<url>:<port>")
  process.exit(1)
}

let url = process.argv[2]
const provider = new Web3.providers.WebsocketProvider(url, {
  reconnect: {
    auto: true,
    delay: 1000, // ms
    onTimeout: false,
    // maxAttempts:
  },
  timeout: 5000, // ms
  clientConfig: {
    maxReceivedFrameSize: 10000000000,
    maxReceivedMessageSize: 10000000000,
    keepalive: true,
    keepaliveInterval: 1000, // ms
    dropConnectionOnKeepaliveTimeout: true,
    keepaliveGracePeriod: 4000, // ms
  }
});
let web3 = new Web3(provider);
const liquidABI = [{"inputs":[{"internalType":"address","name":"_validator","type":"address"},{"internalType":"address payable","name":"_treasury","type":"address"},{"internalType":"uint256","name":"_commissionRate","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[],"name":"COMMISSION_RATE_PRECISION","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"FEE_FACTOR_UNIT_RECIP","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_owner","type":"address"},{"internalType":"address","name":"_spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_spender","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_delegator","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_account","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"burn","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"claimRewards","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_account","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"redistribute","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"_success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_sender","type":"address"},{"internalType":"address","name":"_recipient","type":"address"},{"internalType":"uint256","name":"_amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"_success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_account","type":"address"}],"name":"unclaimedRewards","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}];


const chequebook = function () {
  web3.extend({
    property: 'chequebook',
    methods: [
      new web3.extend.Method({
        name: 'deposit',
        call: 'chequebook_deposit',
        params: 1,
        inputFormatter: [null]
      }),
      new web3.extend.Method({
        name: 'balance',
        call: 'chequebook_balance',
        outputFormatter: web3.extend.utils.toDecimal
      }),
      new web3.extend.Method({
        name: 'cash',
        call: 'chequebook_cash',
        params: 1,
        inputFormatter: [null]
      }),
      new web3.extend.Method({
        name: 'issue',
        call: 'chequebook_issue',
        params: 2,
        inputFormatter: [null, null]
      }),
    ]
  })
};

const ethash = function () {
  web3.extend({
    property: 'ethash',
    methods: [
      new web3.extend.Method({
        name: 'getWork',
        call: 'ethash_getWork',
        params: 0
      }),
      new web3.extend.Method({
        name: 'getHashrate',
        call: 'ethash_getHashrate',
        params: 0
      }),
      new web3.extend.Method({
        name: 'submitWork',
        call: 'ethash_submitWork',
        params: 3,
      }),
      new web3.extend.Method({
        name: 'submitHashRate',
        call: 'ethash_submitHashRate',
        params: 2,
      }),
    ]
  })
};

const admin = function () {
  web3.extend({
    property: 'admin',
    methods: [
      new web3.extend.Method({
        name: 'addPeer',
        call: 'admin_addPeer',
        params: 1
      }),
      new web3.extend.Method({
        name: 'removePeer',
        call: 'admin_removePeer',
        params: 1
      }),
      new web3.extend.Method({
        name: 'addTrustedPeer',
        call: 'admin_addTrustedPeer',
        params: 1
      }),
      new web3.extend.Method({
        name: 'removeTrustedPeer',
        call: 'admin_removeTrustedPeer',
        params: 1
      }),
      new web3.extend.Method({
        name: 'exportChain',
        call: 'admin_exportChain',
        params: 3,
        inputFormatter: [null, null, null]
      }),
      new web3.extend.Method({
        name: 'importChain',
        call: 'admin_importChain',
        params: 1
      }),
      new web3.extend.Method({
        name: 'sleepBlocks',
        call: 'admin_sleepBlocks',
        params: 2
      }),
      new web3.extend.Method({
        name: 'startRPC',
        call: 'admin_startRPC',
        params: 4,
        inputFormatter: [null, null, null, null]
      }),
      new web3.extend.Method({
        name: 'stopRPC',
        call: 'admin_stopRPC'
      }),
      new web3.extend.Method({
        name: 'startWS',
        call: 'admin_startWS',
        params: 4,
        inputFormatter: [null, null, null, null]
      }),
      new web3.extend.Method({
        name: 'stopWS',
        call: 'admin_stopWS'
      }),
      new web3.extend.Method({
        name: 'nodeInfo',
        call: 'admin_nodeInfo'
      }),
      new web3.extend.Method({
        name: 'peers',
        call: 'admin_peers'
      }),
      new web3.extend.Method({
        name: 'datadir',
        call: 'admin_datadir'
      }),
    ]
  })
};

const debug = function () {
  web3.extend({
    property: 'debug',
    methods: [
      new web3.extend.Method({
        name: 'accountRange',
        call: 'debug_accountRange',
        params: 2
      }),
      new web3.extend.Method({
        name: 'printBlock',
        call: 'debug_printBlock',
        params: 1
      }),
      new web3.extend.Method({
        name: 'getBlockRlp',
        call: 'debug_getBlockRlp',
        params: 1
      }),
      new web3.extend.Method({
        name: 'setHead',
        call: 'debug_setHead',
        params: 1
      }),
      new web3.extend.Method({
        name: 'seedHash',
        call: 'debug_seedHash',
        params: 1
      }),
      new web3.extend.Method({
        name: 'dumpBlock',
        call: 'debug_dumpBlock',
        params: 1
      }),
      new web3.extend.Method({
        name: 'chaindbProperty',
        call: 'debug_chaindbProperty',
        params: 1,
        outputFormatter: console.log
      }),
      new web3.extend.Method({
        name: 'chaindbCompact',
        call: 'debug_chaindbCompact',
      }),
      new web3.extend.Method({
        name: 'verbosity',
        call: 'debug_verbosity',
        params: 1
      }),
      new web3.extend.Method({
        name: 'vmodule',
        call: 'debug_vmodule',
        params: 1
      }),
      new web3.extend.Method({
        name: 'backtraceAt',
        call: 'debug_backtraceAt',
        params: 1,
      }),
      new web3.extend.Method({
        name: 'stacks',
        call: 'debug_stacks',
        params: 0,
        outputFormatter: console.log
      }),
      new web3.extend.Method({
        name: 'freeOSMemory',
        call: 'debug_freeOSMemory',
        params: 0,
      }),
      new web3.extend.Method({
        name: 'setGCPercent',
        call: 'debug_setGCPercent',
        params: 1,
      }),
      new web3.extend.Method({
        name: 'memStats',
        call: 'debug_memStats',
        params: 0,
      }),
      new web3.extend.Method({
        name: 'gcStats',
        call: 'debug_gcStats',
        params: 0,
      }),
      new web3.extend.Method({
        name: 'cpuProfile',
        call: 'debug_cpuProfile',
        params: 2
      }),
      new web3.extend.Method({
        name: 'startCPUProfile',
        call: 'debug_startCPUProfile',
        params: 1
      }),
      new web3.extend.Method({
        name: 'stopCPUProfile',
        call: 'debug_stopCPUProfile',
        params: 0
      }),
      new web3.extend.Method({
        name: 'goTrace',
        call: 'debug_goTrace',
        params: 2
      }),
      new web3.extend.Method({
        name: 'startGoTrace',
        call: 'debug_startGoTrace',
        params: 1
      }),
      new web3.extend.Method({
        name: 'stopGoTrace',
        call: 'debug_stopGoTrace',
        params: 0
      }),
      new web3.extend.Method({
        name: 'blockProfile',
        call: 'debug_blockProfile',
        params: 2
      }),
      new web3.extend.Method({
        name: 'setBlockProfileRate',
        call: 'debug_setBlockProfileRate',
        params: 1
      }),
      new web3.extend.Method({
        name: 'writeBlockProfile',
        call: 'debug_writeBlockProfile',
        params: 1
      }),
      new web3.extend.Method({
        name: 'mutexProfile',
        call: 'debug_mutexProfile',
        params: 2
      }),
      new web3.extend.Method({
        name: 'setMutexProfileFraction',
        call: 'debug_setMutexProfileFraction',
        params: 1
      }),
      new web3.extend.Method({
        name: 'writeMutexProfile',
        call: 'debug_writeMutexProfile',
        params: 1
      }),
      new web3.extend.Method({
        name: 'writeMemProfile',
        call: 'debug_writeMemProfile',
        params: 1
      }),
      new web3.extend.Method({
        name: 'traceBlock',
        call: 'debug_traceBlock',
        params: 2,
        inputFormatter: [null, null]
      }),
      new web3.extend.Method({
        name: 'traceBlockFromFile',
        call: 'debug_traceBlockFromFile',
        params: 2,
        inputFormatter: [null, null]
      }),
      new web3.extend.Method({
        name: 'traceBadBlock',
        call: 'debug_traceBadBlock',
        params: 1,
        inputFormatter: [null]
      }),
      new web3.extend.Method({
        name: 'standardTraceBadBlockToFile',
        call: 'debug_standardTraceBadBlockToFile',
        params: 2,
        inputFormatter: [null, null]
      }),
      new web3.extend.Method({
        name: 'standardTraceBlockToFile',
        call: 'debug_standardTraceBlockToFile',
        params: 2,
        inputFormatter: [null, null]
      }),
      new web3.extend.Method({
        name: 'traceBlockByNumber',
        call: 'debug_traceBlockByNumber',
        params: 2,
        inputFormatter: [null, null]
      }),
      new web3.extend.Method({
        name: 'traceBlockByHash',
        call: 'debug_traceBlockByHash',
        params: 2,
        inputFormatter: [null, null]
      }),
      new web3.extend.Method({
        name: 'traceTransaction',
        call: 'debug_traceTransaction',
        params: 2,
        inputFormatter: [null, null]
      }),
      new web3.extend.Method({
        name: 'preimage',
        call: 'debug_preimage',
        params: 1,
        inputFormatter: [null]
      }),
      new web3.extend.Method({
        name: 'getBadBlocks',
        call: 'debug_getBadBlocks',
        params: 0,
      }),
      new web3.extend.Method({
        name: 'storageRangeAt',
        call: 'debug_storageRangeAt',
        params: 5,
      }),
      new web3.extend.Method({
        name: 'getModifiedAccountsByNumber',
        call: 'debug_getModifiedAccountsByNumber',
        params: 2,
        inputFormatter: [null, null],
      }),
      new web3.extend.Method({
        name: 'getModifiedAccountsByHash',
        call: 'debug_getModifiedAccountsByHash',
        params: 2,
        inputFormatter: [null, null],
      }),
      new web3.extend.Method({
        name: 'freezeClient',
        call: 'debug_freezeClient',
        params: 1,
      }),
    ],
  })
};

const eth = function () {
  web3.extend({
    property: 'eth',
    methods: [
      new web3.extend.Method({
        name: 'chainId',
        call: 'eth_chainId',
        params: 0
      }),
      new web3.extend.Method({
        name: 'sign',
        call: 'eth_sign',
        params: 2,
        inputFormatter: [web3.extend.formatters.inputAddressFormatter, null]
      }),
      new web3.extend.Method({
        name: 'resend',
        call: 'eth_resend',
        params: 3,
        inputFormatter: [web3.extend.formatters.inputTransactionFormatter, web3.extend.utils.fromDecimal, web3.extend.utils.fromDecimal]
      }),
      new web3.extend.Method({
        name: 'signTransaction',
        call: 'eth_signTransaction',
        params: 1,
        inputFormatter: [web3.extend.formatters.inputTransactionFormatter]
      }),
      new web3.extend.Method({
        name: 'submitTransaction',
        call: 'eth_submitTransaction',
        params: 1,
        inputFormatter: [web3.extend.formatters.inputTransactionFormatter]
      }),
      new web3.extend.Method({
        name: 'fillTransaction',
        call: 'eth_fillTransaction',
        params: 1,
        inputFormatter: [web3.extend.formatters.inputTransactionFormatter]
      }),
      new web3.extend.Method({
        name: 'getHeaderByNumber',
        call: 'eth_getHeaderByNumber',
        params: 1
      }),
      new web3.extend.Method({
        name: 'getHeaderByHash',
        call: 'eth_getHeaderByHash',
        params: 1
      }),
      new web3.extend.Method({
        name: 'getBlockByNumber',
        call: 'eth_getBlockByNumber',
        params: 2
      }),
      new web3.extend.Method({
        name: 'getBlockByHash',
        call: 'eth_getBlockByHash',
        params: 2
      }),
      new web3.extend.Method({
        name: 'getRawTransaction',
        call: 'eth_getRawTransactionByHash',
        params: 1
      }),
      new web3.extend.Method({
        name: 'getRawTransactionFromBlock',
        call: function (args) {
          return (web3.extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? 'eth_getRawTransactionByBlockHashAndIndex' : 'eth_getRawTransactionByBlockNumberAndIndex';
        },
        params: 2,
        inputFormatter: [web3.extend.formatters.inputBlockNumberFormatter, web3.extend.utils.toHex]
      }),
      new web3.extend.Method({
        name: 'getProof',
        call: 'eth_getProof',
        params: 3,
        inputFormatter: [web3.extend.formatters.inputAddressFormatter, null, web3.extend.formatters.inputBlockNumberFormatter]
      }),
      new web3.extend.Method({
        name: 'pendingTransactions',
        call: 'eth_pendingTransactions',
        outputFormatter: function (txs) {
          var formatted = [];
          for (var i = 0; i < txs.length; i++) {
            formatted.push(web3.extend.formatters.outputTransactionFormatter(txs[i]));
            formatted[i].blockHash = null;
          }
          return formatted;
        }
      }),
    ]
  })
};

const miner = function () {
  web3.extend({
    property: 'miner',
    methods: [
      new web3.extend.Method({
        name: 'start',
        call: 'miner_start',
        params: 1,
        inputFormatter: [null]
      }),
      new web3.extend.Method({
        name: 'stop',
        call: 'miner_stop'
      }),
      new web3.extend.Method({
        name: 'setEtherbase',
        call: 'miner_setEtherbase',
        params: 1,
        inputFormatter: [web3.extend.formatters.inputAddressFormatter]
      }),
      new web3.extend.Method({
        name: 'setExtra',
        call: 'miner_setExtra',
        params: 1
      }),
      new web3.extend.Method({
        name: 'setGasPrice',
        call: 'miner_setGasPrice',
        params: 1,
        inputFormatter: [web3.extend.utils.fromDecimal]
      }),
      new web3.extend.Method({
        name: 'setRecommitInterval',
        call: 'miner_setRecommitInterval',
        params: 1,
      }),
      new web3.extend.Method({
        name: 'getHashrate',
        call: 'miner_getHashrate'
      }),
    ],
  })
};

const net = function () {
  web3.extend({
    property: 'net',
    methods: [
      new web3.extend.Method({
        name: 'version',
        call: 'net_version'
      }),
    ]
  })
};

const personal = function () {
  web3.extend({
    property: 'personal',
    methods: [
      new web3.extend.Method({
        name: 'importRawKey',
        call: 'personal_importRawKey',
        params: 2
      }),
      new web3.extend.Method({
        name: 'sign',
        call: 'personal_sign',
        params: 3,
        inputFormatter: [null, web3.extend.formatters.inputAddressFormatter, null]
      }),
      new web3.extend.Method({
        name: 'ecRecover',
        call: 'personal_ecRecover',
        params: 2
      }),
      new web3.extend.Method({
        name: 'openWallet',
        call: 'personal_openWallet',
        params: 2
      }),
      new web3.extend.Method({
        name: 'deriveAccount',
        call: 'personal_deriveAccount',
        params: 3
      }),
      new web3.extend.Method({
        name: 'signTransaction',
        call: 'personal_signTransaction',
        params: 2,
        inputFormatter: [web3.extend.formatters.inputTransactionFormatter, null]
      }),
      new web3.extend.Method({
        name: 'unpair',
        call: 'personal_unpair',
        params: 2
      }),
      new web3.extend.Method({
        name: 'initializeWallet',
        call: 'personal_initializeWallet',
        params: 1
      }),
      new web3.extend.Method({
        name: 'listWallets',
        call: 'personal_listWallets'
      }),
      new web3.extend.Method({
        name: 'unlockAccount',
        call: 'personal_unlockAccount',
        params: 3,
        inputFormatter: [web3.extend.formatters.inputAddressFormatter, null, null]
      }),
    ]
  })
};

const rpc = function () {
  web3.extend({
    property: 'rpc',
    methods: [
      new web3.extend.Method({
        name: 'modules',
        call: 'rpc_modules',
      }),
      new web3.extend.Method({
        name: 'modules',
        call: 'rpc_modules'
      }),
    ]
  })
};

const shh = function () {
  web3.extend({
    property: 'shh',
    methods: [
      new web3.extend.Method({
        name: 'version',
        call: 'shh_version',
        outputFormatter: web3.extend.utils.toDecimal
      }),
      new web3.extend.Method({
        name: 'info',
        call: 'shh_info'
      }),
    ]
  })
};

const swarmfs = function () {
  web3.extend({
    property: 'swarmfs',
    methods:
      [
        new web3.extend.Method({
          name: 'mount',
          call: 'swarmfs_mount',
          params: 2
        }),
        new web3.extend.Method({
          name: 'unmount',
          call: 'swarmfs_unmount',
          params: 1
        }),
        new web3.extend.Method({
          name: 'listmounts',
          call: 'swarmfs_listmounts',
          params: 0
        }),
      ]
  })
};

const txpool = function () {
  web3.extend({
    property: 'txpool',
    methods: [
      new web3.extend.Method({
        name: 'content',
        call: 'txpool_content'
      }),
      new web3.extend.Method({
        name: 'inspect',
        call: 'txpool_inspect'
      }),
      new web3.extend.Method({
        name: 'status',
        call: 'txpool_status',
        outputFormatter: function (status) {
          status.pending = web3.extend.utils.toDecimal(status.pending);
          status.queued = web3.extend.utils.toDecimal(status.queued);
          return status;
        }
      }),
    ]
  })
};

const accounting = function () {
  web3.extend({
    property: 'accounting',
    methods: [
      new web3.extend.Method({
        name: 'balance',
        call: 'account_balance'
      }),
      new web3.extend.Method({
        name: 'balanceCredit',
        call: 'account_balanceCredit'
      }),
      new web3.extend.Method({
        name: 'balanceDebit',
        call: 'account_balanceDebit'
      }),
      new web3.extend.Method({
        name: 'bytesCredit',
        call: 'account_bytesCredit'
      }),
      new web3.extend.Method({
        name: 'bytesDebit',
        call: 'account_bytesDebit'
      }),
      new web3.extend.Method({
        name: 'msgCredit',
        call: 'account_msgCredit'
      }),
      new web3.extend.Method({
        name: 'msgDebit',
        call: 'account_msgDebit'
      }),
      new web3.extend.Method({
        name: 'peerDrops',
        call: 'account_peerDrops'
      }),
      new web3.extend.Method({
        name: 'selfDrops',
        call: 'account_selfDrops'
      }),
    ]
  })
};

const les = function () {
  web3.extend({
    property: 'les',
    methods:
      [
        new web3.extend.Method({
          name: 'getCheckpoint',
          call: 'les_getCheckpoint',
          params: 1
        }),
        new web3.extend.Method({
          name: 'clientInfo',
          call: 'les_clientInfo',
          params: 1
        }),
        new web3.extend.Method({
          name: 'priorityClientInfo',
          call: 'les_priorityClientInfo',
          params: 3
        }),
        new web3.extend.Method({
          name: 'setClientParams',
          call: 'les_setClientParams',
          params: 2
        }),
        new web3.extend.Method({
          name: 'setDefaultParams',
          call: 'les_setDefaultParams',
          params: 1
        }),
        new web3.extend.Method({
          name: 'addBalance',
          call: 'les_addBalance',
          params: 3
        }),
        new web3.extend.Method({
          name: 'latestCheckpoint',
          call: 'les_latestCheckpoint'
        }),
        new web3.extend.Method({
          name: 'checkpointContractAddress',
          call: 'les_getCheckpointContractAddress'
        }),
        new web3.extend.Method({
          name: 'serverInfo',
          call: 'les_serverInfo'
        }),
      ]
  })
};

const tendermint = function () {
  web3.extend({
    property: 'tendermint',
    methods:
      [
        new web3.extend.Method({
          name: 'getCommittee',
          call: 'tendermint_getCommittee',
          params: 1
        }),
        new web3.extend.Method({
          name: 'getCommitteeAtHash',
          call: 'tendermint_getCommitteeAtHash',
          params: 1
        }),
        new web3.extend.Method({
          name: 'getContractAddress',
          call: 'tendermint_getContractAddress',
          params: 0
        }),
        new web3.extend.Method({
          name: 'getContractABI',
          call: 'tendermint_getContractABI',
          params: 0
        }),
        new web3.extend.Method({
          name: 'getWhitelist',
          call: 'tendermint_getWhitelist',
          params: 0
        }),
        new web3.extend.Method({
          name: 'getCoreState',
          call: 'tendermint_getCoreState',
          params: 0
        })
      ]
  })
};

const lesPay = function () {
  web3.extend({
    property: 'lespay',
    methods:
      [
        new web3.extend.Method({
          name: 'distribution',
          call: 'lespay_distribution',
          params: 2
        }),
        new web3.extend.Method({
          name: 'timeout',
          call: 'lespay_timeout',
          params: 2
        }),
        new web3.extend.Method({
          name: 'value',
          call: 'lespay_value',
          params: 2
        }),
        new web3.extend.Method({
          name: 'requestStats',
          call: 'lespay_requestStats'
        }),
      ]
  })
};


let moduleInit = new Map();
moduleInit.set('accounting', accounting);
moduleInit.set('admin', admin);
moduleInit.set('chequebook', chequebook);
moduleInit.set('ethash', ethash);
moduleInit.set('debug', debug);
moduleInit.set('eth', eth);
moduleInit.set('miner', miner);
moduleInit.set('net', net);
moduleInit.set('personal', personal);
moduleInit.set('rpc', rpc);
moduleInit.set('shh', shh);
moduleInit.set('swarmfs', swarmfs);
moduleInit.set('txpool', txpool);
moduleInit.set('les', les);
moduleInit.set('tendermint', tendermint);
moduleInit.set('lespay', lesPay);

let WebSocket = require('ws');
const repl = require("repl");
let ws = new WebSocket(url);

let data = JSON.stringify({
  jsonrpc: '2.0',
  method: 'rpc_modules',
  id: 1
});

ws.addEventListener('open', () => {
  ws.send(data);
});

ws.addEventListener('message', setup);

async function setup(event) {
  let initMessage = 'modules: ';
  let parsed = JSON.parse(event.data);
  for (const [key, value] of Object.entries(parsed.result)) {
    if (moduleInit.has(key)) {
      moduleInit.get(key)()
      initMessage += key + ' ';
    }
  }
  ws.close(1000); // Successful close signal

  let hasAutonity = false;
  let contract;
  // If the tendermint module is loaded then load the contract bindings
  if (parsed.result.hasOwnProperty('tendermint')) {
    contract = await Promise.all([web3.tendermint.getContractABI(), web3.eth.getGasPrice(), web3.eth.getCoinbase()]).then((results) => {
      return new web3.eth.Contract(JSON.parse(results[0]), '0xbd770416a3345f91e4b34576cb804a576fa48eb1');
    });
    initMessage += 'autonity ';
    hasAutonity = true
  }

  console.log('Welcome to the Autonity node console');
  console.log(initMessage);
  console.log('Type "module".<Tab> to get started');

  // Start the repl
  const repl = require("repl");
  const server = repl.start();

  const defaultEval = server.eval;
  // await on returned promises by hooking default repl
  // https://github.com/nodejs/node/issues/13209#issuecomment-619526317
  server.eval = (cmd, context, filename, callback) => {
    defaultEval(cmd, context, filename, async (err, result) => {
      if (err) {
        // propagate errors from the eval
        callback(err);
      } else {
        // awaits the promise and either return result or error
        try {
          callback(null, await Promise.resolve(result));
        } catch (err) {
          callback(err);
        }
      }
    });
  };

  // Setting context values makes them available in the default scope
  // of the repl.
  for (let [key, val] of moduleInit) {
    server.context[key] = {}
    for (prop in web3[key]) {
      // This can be brittle..
      if (web3[key].hasOwnProperty(prop) && typeof web3[key][prop] == 'function') {
        server.context[key][prop] = web3[key][prop]
      }
    }
  }

  server.context.web3 = web3;
  server.on('exit', () => {
    process.exit();
  });

  if (hasAutonity){
    // Setting context values makes them available in the default scope
    // of the repl.
    server.context.contract = contract;
    server.context.autonity = contract.methods;

    // EXPERIMENTAL EXTENSIONS //

    async function getValidatorsData() {
      const vals = await contract.methods.getValidators().call();
      return Promise.all(vals.map(async val => {
        return contract.methods.getValidator(val).call();
      }));
    }
    /* Retrieve validators data */
    server.context.vals =  async () => {
      const data = await getValidatorsData();
      for(let i = 0; i< data.length; i++){
        console.log("\n\x1b[40m\x1b[37m__________________________Validator %d__________________________\x1b[0m", i)
        for (const key in data[i]) {
          if (!isNaN(key)) {
            continue
          }
          const separator = [...Array((19 - key.length)).keys()].reduce((a,b)=> a + " ","")
          console.log(`${key}:${separator} \x1b[1m${data[i][key]}\x1b[0m`);
        }
      }
    }

    /* Retrieve liquid staking wallet data */
    server.context.wal =  async (account) => {
      const vals = await getValidatorsData();
      const w = []
      for(let i = 0; i< vals.length; i++){
        const liquid = new web3.eth.Contract(liquidABI, vals[i].liquidContract);
        w.push({
          validator:     vals[i].addr,
          lnew:          await liquid.methods.balanceOf(account).call(),
          claimableRewards: await liquid.methods.unclaimedRewards(account).call()
        })
      }
      console.log("\n\x1b[40m\x1b[37m_____________________________________Staking Wallet_____________________________________\x1b[0m\n");
      console.log("\taccount: \t" + account + "\n");
      console.log("\tbalance: \t" + await web3.eth.getBalance(account) + " attoton\n");
      console.table(w);
    }

    /* Claim Rewards */
    server.context.rclm =  async (account, validator) => {
      const val = await contract.methods.getValidator(validator).call();
      console.log(`validator:      ${val.addr}`);
      console.log(`staker:         ${account}`);

      const liquid = new web3.eth.Contract(liquidABI, val.liquidContract);
      const claimable = await liquid.methods.unclaimedRewards(account).call()
      console.log(`claimable rewards: ${claimable}`);
      return await liquid.methods.claimRewards().send({from:account, gas:100000});
    }

    /* Claim rewards from all staked validators */
    server.context.rclm_a =  async (account) => {
      const vals = await getValidatorsData();
      let total = 0n;
      for(let i = 0; i< vals.length; i++){
        const liquid = new web3.eth.Contract(liquidABI, vals[i].liquidContract);
        const claimable = await liquid.methods.unclaimedRewards(account).call()
        if (BigInt(claimable) == 0n) {
          continue
        }
        console.log(`validator:          ${vals[i].addr}`);
        console.log(`claimable reward:   ${claimable}\n`);
        total += BigInt(claimable);
        await liquid.methods.claimRewards().send({from:account, gas: 100000});
      }
      console.log(`\n\x1b[1mtotal claimed:        ${total}\x1b[0m`);
    }

    /* Transfer Liquid Newton to recipient */
    server.context.lsend =  async ({val, from, to, value}) => {
      console.log(`Sending ${value} LNEW - Validator ${val}`);
      const validator = await contract.methods.getValidator(val).call();
      const liquid = new web3.eth.Contract(liquidABI, validator.liquidContract);
      return liquid.methods.transfer(to,value).send({from, gas: 300000})
    }

    /* liveness monitor */
    // this is the wrong approach - we should do better
    let mon_lastBlock
    server.context.livemon = async() => {
      web3.eth.subscribe("newBlockHeaders", (err,block) => {
        // skip first returned block
        if(mon_lastBlock == 0){
          mon_lastBlock = block.timestamp;
          return
        }
        const tdiff = block.timestamp - mon_lastBlock
        if(tdiff > 2){
          console.log("!!ALERT LIVENESS!!")
        }
        console.log(`#${block.number}\t ${tdiff}`);
        mon_lastBlock =  block.timestamp
      })
    }

  }
}
