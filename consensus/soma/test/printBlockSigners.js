const rlp = require('rlp')
const net = require('net');
const utils = require('ethereumjs-utils')
const Web3 = require('web3');

const secp256k1 = require('secp256k1')
const Buffer = require('safe-buffer').Buffer

const ecrecover = (msgHash, v, r, s, chainId) => {
  const signature = Buffer.concat([r,s], 64)
  const senderPubKey = secp256k1.recover(msgHash, signature, v)
  return secp256k1.publicKeyConvert(senderPubKey, false).slice(1)
}

/*
const HTTP_URI='http://localhost'
const HTTP_PORT=8501

const WS_URI='ws://localhost'
const WS_PORT=8545

const web3 = new Web3(WS_URI + ':' + WS_PORT)
*/

const sigHash = (header) => {
	const encodedHeader = rlp.encode([
		header.parentHash,
		header.sha3Uncles,
		header.miner,
		header.stateRoot,
		header.transactionsRoot,
		header.receiptsRoot,
		header.logsBloom,
		web3.utils.toBN(header.difficulty),
		web3.utils.toBN(header.number),
		header.gasLimit,
		header.gasUsed,
		web3.utils.toBN(header.timestamp),
		header.extraData,//[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.mixHash,
		header.nonce,
	])
	return web3.utils.sha3(encodedHeader)
}

const splitSignature = signatureStr => {
  const r = Buffer.from(signatureStr.slice(2).slice(0,32*2).padStart(32*2,0),'hex')
  const s = Buffer.from(signatureStr.slice(2).slice(32*2,(32*2)*2).padStart(32*2,0),'hex')
  const v = Number('0x' + signatureStr.slice(2).slice((32*2)*2))

  console.log(`SplitSignature():\n\tr: ${r.toString('hex')}`, `\n\ts: ${s.toString('hex')}`,`\n\tv: ${v}`)

  return { r, s, v }
}


const printHeader =  async header => {
  const extraVanity = 32
  const extraSeal = 65 // r(32 bytes) + s(32 bytes) + v(1 byte)
  const number = header.number
  const hash = header.hash

  const signature = '0x' + header.extraData.slice(-(extraSeal*2))
  const extraDataUnsigned = header.extraData.slice(0,header.extraData.length-(extraSeal*2))//.padEnd(header.extraData.length,0)

  // console.log('extraData:\t\t',header.extraData)
  // console.log('unsignedExtraData:\t',extraDataUnsigned)

  const blockHeaderNoSignature = Object.assign({},header, extraDataUnsigned)
  const blockHashNoSignature = sigHash(blockHeaderNoSignature)

  console.log(`HEADER ${number}`)

  const { r, s, v } = splitSignature(signature)
  const unsignedBlockBuffer = Buffer.from(blockHashNoSignature.slice(2),'hex')

  //const pubKey = ecrecover(Buffer.from(blockHashNoSignature.slice(2), 'hex'),v,r,s)
  //console.log('pubKey: ', pubKey.toString('hex'))

  //const signerAddress = await web3.eth.accounts.recover(blockHashNoSignature, signature, true)
  const sig = Buffer.concat([r,s,Buffer.from([v+27])])
  const signerAddress1 = await web3.eth.personal.ecRecover(unsignedBlockBuffer.toString(), '0x'+sig.toString('hex'))
  const pub = utils.ecrecover(unsignedBlockBuffer,v+27,r,s)
  const signerAddress = '0x' + utils.pubToAddress(pub).toString('hex')

  console.log(`\tSigner: ${signerAddress}`)
  console.log(`\tSigner: ${signerAddress1}`)
}

// Using the IPC provider in node.js
const web3 = new Web3('/home/doart3/soma_instance/data_node_0/geth.ipc', net);

const subscription = web3.eth.subscribe('newBlockHeaders')
  .on("data", printHeader)
  .on("error", console.error);

// unsubscribes the subscription
/*
subscription.unsubscribe(function(error, success){
    if (success) {
        console.log('Successfully unsubscribed!');
    }
});
*/
