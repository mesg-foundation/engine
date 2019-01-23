import Web3 from "web3"
import Service from "mesg-js/lib/service/service"

export = async (mesg: Service, web3: Web3, blockNumber: number) => {
  const block = await web3.eth.getBlock(blockNumber, false)
  await mesg.emitEvent('block', {
    number: block.number,
    hash: block.hash,
    parentHash: block.parentHash,
    sha3Uncles: block.sha3Uncles,
    logsBloom: block.logsBloom,
    stateRoot: block.stateRoot,
    miner: block.miner,
    extraData: block.extraData,
    gasLimit: block.gasLimit,
    gasUsed: block.gasUsed,
    timestamp: block.timestamp,
    size: block.size,
    difficulty: block.difficulty
  })
}