/**
 * Use this file to configure your truffle project. It's seeded with some
 * common settings for different networks and features like migrations,
 * compilation and testing. Uncomment the ones you need or modify
 * them to suit your project as necessary.
 *
 * More information about configuration can be found at:
 *
 * truffleframework.com/docs/advanced/configuration
 *
 * To deploy via Infura you'll need a wallet provider (like truffle-hdwallet-provider)
 * to sign your transactions before they're sent to a remote public node. Infura accounts
 * are available for free at: infura.io/register.
 *
 * You'll also need a mnemonic - the twelve word phrase the wallet uses to generate
 * public/private key pairs. If you're publishing your code to GitHub make sure you load this
 * phrase from a file you've .gitignored so it doesn't accidentally become public.
 *
 */

// module.exports = {
//  /**
//   * Networks define how you connect to your ethereum client and let you set the
//   * defaults web3 uses to send transactions. If you don't specify one truffle
//   * will spin up a development blockchain for you on port 9545 when you
//   * run `develop` or `test`. You can ask a truffle command to use a specific
//   * network from the command line, e.g
//   *
//   * $ truffle test --network <network-name>
//   */
//
//  networks: {
//    // Useful for testing. The `development` name is special - truffle uses it by default
//    // if it's defined here and no other network is specified at the command line.
//    // You should run a client (like ganache-cli, geth or parity) in a separate terminal
//    // tab if you use this network and you must also set the `host`, `port` and `network_id`
//    // options below to some value.
//    autonity: {
//      host: "127.0.0.1",     // Localhost (default: none)
//      port: 8545,            // Standard Ethereum port (default: none)
//      network_id: "*",       // Any network (default: none)
//      gas: 46123880
//    },
//
//    // Another network with more advanced options...
//    // advanced: {
//    // port: 8777,             // Custom port
//    // network_id: 1342,       // Custom network
//    // gas: 8500000,           // Gas sent with each transaction (default: ~6700000)
//    // gasPrice: 20000000000,  // 20 gwei (in wei) (default: 100 gwei)
//    // from: <address>,        // Account to send txs from (default: accounts[0])
//    // websockets: true        // Enable EventEmitter interface for web3 (default: false)
//    // },
//
//    // Useful for deploying to a public network.
//    // NB: It's important to wrap the provider as a function.
//    // ropsten: {
//    // provider: () => new HDWalletProvider(mnemonic, `https://ropsten.infura.io/v3/YOUR-PROJECT-ID`),
//    // network_id: 3,       // Ropsten's id
//    // gas: 5500000,        // Ropsten has a lower block limit than mainnet
//    // confirmations: 2,    // # of confs to wait between deployments. (default: 0)
//    // timeoutBlocks: 200,  // # of blocks before a deployment times out  (minimum/default: 50)
//    // skipDryRun: true     // Skip dry run before migrations? (default: false for public nets )
//    // },
//
//    // Useful for private networks
//    // private: {
//    // provider: () => new HDWalletProvider(mnemonic, `https://network.io`),
//    // network_id: 2111,   // This network is yours, in the cloud.
//    // production: true    // Treats this network as if it was a public net. (default: false)
//    // }
//  },
//
//  // Set default mocha options here, use special reporters etc.
//  mocha: {
//    // timeout: 100000
//  },
//
//  // Configure your compilers
//  compilers: {
//    solc: {
//      version: "0.8.19",    // Fetch exact version from solc-bin (default: truffle's version)
//      // docker: true,        // Use "0.5.1" you've installed locally with docker (default: false)
//      optimizer: {
//        enabled: true,
//        runs: 1
//      },
//    }
//  }
//}

const HDWalletProvider = require("@truffle/hdwallet-provider");

// Private keys for the 10 test accounts, ordered to match the original
// --unlock order so that accounts[N] indices are preserved.
const TEST_PRIVATE_KEYS = [
  "a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91", // 0x850c1eb8d190e05845ad7f84ac95a318c8aab07f
  "aa4b77b1305f8f265e81599587c623d8950624f3e1bd9c121ef2461a7a1e7527", // 0x4ad219b58a5b46a1d9662beaa6a70db9f570dea5
  "9e19e8f26a59eb5465cf39dde39161eafc84ba934fb641bb0078ffcf8393e133", // 0x4b07239bd581d21aefcdee0c6db38070f9a5fd2d
  "4ec99383dc50aa3f3117fcbfba7b69188ba60d3418185fb353c9a69d066e55d9", // 0xc443c6c6ae98f5110702921138d840e77da67702
  "0c8698f456533170fe07c6dcb753d47bef8bedd46443efa57a859c989887b56b", // 0x09428e8674496e2d1e965402f33a9520c5fcbbe2
  "82b34a8613c090fd991f7c6e2122f28752fd7066a26e5105576f6ef78d49b7d5", // 0x64852003fc0b84d6c49c5cb3dfcd17922affddc1
  "3073d06125723c03a4b2317b73280878e2d5e7916731e1dc27f7650e13b03d5c", // 0x4839950a5f07d6d6cd82f933d1de8574c48d6e74
  "fa59cb05e77d39425f16d17277bafdfb742e1057177476173d9c79576c6bf779", // 0x160bc705bf2e5871557722c9332cfa185c02b765
  "58951d75562e20501fdcbc8fa6d36b6a10e87aa429ea0e0d302cc0718973f9f2", // 0xe12b43b69e57ed6acdd8721eb092bf7c8d41df41
  "e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b", // 0xde03b7806f885ae79d2aa56568b77cadb0de073e
];

module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",
      port: 8545,
      network_id: "*",
      gas: 56123880
    },
    autonity: {
      provider: () => new HDWalletProvider({
        privateKeys: TEST_PRIVATE_KEYS,
        providerOrUrl: "http://127.0.0.1:8545",
      }),
      network_id: "*",
      gas: 56123880,
    },
  },

  mocha: {
     timeout: 900000,
  },

  compilers: {
    solc: {
      version: "0.8.30",
      evmVersion: "london",
      // docker: true,
      optimizer: {
        enabled: true,
        runs: 1
      },
    }
  }
}
