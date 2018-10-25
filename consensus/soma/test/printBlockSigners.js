const rlp = require('rlp')
const net = require('net');
const Web3 = require('web3');

const secp256k1 = require('secp256k1')
const Buffer = require('safe-buffer').Buffer

const CONTRACT_ADDRESS = '0xa4b794510e992936ae8c43bc3511a84ecc2e5d06'

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
		header.extraData,
		header.mixHash,
		header.nonce,
	])
	return web3.utils.sha3(encodedHeader)
}

const splitSignature = signatureStr => {
  const r = Buffer.from(signatureStr.slice(2).slice(0,32*2).padStart(32*2,0),'hex')
  const s = Buffer.from(signatureStr.slice(2).slice(32*2,(32*2)*2).padStart(32*2,0),'hex')
  const v = Number('0x' + signatureStr.slice(2).slice((32*2)*2))

  //console.log(`SplitSignature():\n\tr: ${r.toString('hex')}`, `\n\ts: ${s.toString('hex')}`,`\n\tv: ${v}`)

  return { r, s, v }
}

const getBlockSignerAddress = async header  => {
  const extraVanity = 32
  const extraSeal = 65 // r(32 bytes) + s(32 bytes) + v(1 byte)

  const signature = '0x' + header.extraData.slice(-(extraSeal*2))
  const extraDataUnsigned = header.extraData.slice(0,header.extraData.length-(extraSeal*2))//.padEnd(header.extraData.length,0)

  const blockHeaderNoSignature = Object.assign({},header, {extraData: extraDataUnsigned})
  const blockHashNoSignature = sigHash(blockHeaderNoSignature)

  const unsignedBlockBuffer = Buffer.from(blockHashNoSignature.slice(2),'hex')

  const signerAddress = await web3.eth.accounts.recover(blockHashNoSignature, signature, true)
  return signerAddress
}

const printHeader =  async header => {
  const signerAddress = await getBlockSignerAddress(header)

  const accounts = await web3.eth.getAccounts()
  const signerIndex = accounts.findIndex(acc => acc === signerAddress)


  console.log(`Block number: ${header.number.toString().padStart(5,0)} \
Signer(accounts[${signerIndex.toString().padStart(2,0)}]): ${signerAddress}`)
}

const getValidators = async (contractAddr) => {
  const validatorsFunc = web3.utils.sha3('validators(uint256)').slice(0,(4*2)+2)
  const validators = []
  for(let i =0; true; i++) {
    const argument = i.toString().padStart(32*2,0)
    const callRet = await web3.eth.call({ to: contractAddr, data: validatorsFunc+argument })
    if(callRet.length === 2) break
    else validators.push('0x'+callRet.slice(-20*2))
  }
  return validators
}

const listValidators = async (contractAddr) => {
  const validators = await getValidators(contractAddr)
  validators.forEach((acc,idx) => console.log(`(${idx.toString().padStart(2,0)}) ${acc}`))
}

const listAccounts = async () => {
  const accounts = await web3.eth.getAccounts()
  accounts.forEach((acc,idx) => console.log(`(${idx.toString().padStart(2,0)}) ${acc}`))
}

const castVote = async (contractAddr, from, candidate) => {
  const castVoteFunc = web3.utils.sha3('CastVote(address)').slice(0,(4*2)+2)
  const argument = candidate.slice(2).padStart(32*2,0)
  const data = castVoteFunc+argument
  const tx = { from, to: contractAddr, data, gas: "0xffffffff" }
  await web3.eth.personal.unlockAccount(from,'xxx',600)
  return web3.eth.sendTransaction(tx)
}

const removeValidator = async (contractAddr) => {
  const validators = await getValidators(contractAddr)
  if(validators.length <= 1) {
    console.log('1 or less validators in the validator list!')
    return
  }
  const leavingValidator = validators.slice(-1)[0]
  const votingValidators = validators.slice(0,-1)
  console.log(`Voting to remove validator: ${leavingValidator}\n\tVoters: ${votingValidators.join('\n\t\t')}`)
  const vontingPromi = votingValidators.map(v => castVote(contractAddr,v,leavingValidator))
  const tx = await Promise.all(vontingPromi)
  console.log(`Finished voting to remove validator: ${leavingValidator}\n${JSON.stringify(tx,' ',' ')}`)
}

const addValidator = async (contractAddr) => {
  const accounts = await web3.eth.getAccounts()
  const validators = await getValidators(contractAddr)
  if(validators.length >= accounts.length) {
    console.log('The validator list is longer or the same as the unlocked accounts in the node!')
    return
  }
  const candidates = accounts.map(acc => acc.toLowerCase()).filter(acc => !validators.includes(acc))
  const candidateValidator = candidates.slice(0,1)[0]
  const votingValidators = validators
  console.log(`Voting to add validator: ${candidateValidator}\n\tVoters: ${votingValidators.join('\n\t\t')}`)
  const vontingPromi = votingValidators.map(v => castVote(contractAddr,v,candidateValidator))
  const tx = await Promise.all(vontingPromi)
  console.log(`Finished voting to add validator: ${candidateValidator}\n${JSON.stringify(tx,' ',' ')}`)
}

const onKeyPress = async key => {
  switch (key.toString()) {
    case '+':
      addValidator(CONTRACT_ADDRESS)
      break
    case '-':
      removeValidator(CONTRACT_ADDRESS)
      break
    case 'v':
      console.log('\n================== Validators ==================')
      await listValidators(CONTRACT_ADDRESS)
      console.log('================================================\n')
      break
    case 'l':
      console.log('\n=================== Accounts ===================')
      await listAccounts()
      console.log('================================================\n')
      break
    case 'q':
      process.exit()
      break
    case 'h':
      console.log('\n================== Help ==================\n'
        + '(+): cast a vote to add validator (the first on accounts that is not a validator)\n'
        + '(-): cast a vote to remove validator (the last on the list of validators)\n'
        + '(v): print lis of validators\n'
        + '(l): print lis of accounts\n'
        + '(q): quit\n'
        + '===========================================\n'
      )
      break
    default:
      console.log(`UNKNOWN KEY! (${key})`)
      break
  }
}
onKeyPress('h')

process.stdin.setRawMode(true)
process.stdin.on('data', onKeyPress)

// Using the IPC provider in node.js
const web3 = new Web3('/home/doart3/soma_instance/data_node_0/geth.ipc', net);
// subscribe
web3.eth.subscribe('newBlockHeaders').on("data", printHeader).on("error", console.error)
