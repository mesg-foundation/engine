import { EventEmitter } from "events"
import Web3 from "web3"
import { NewBlockEventEmitterInterface } from "./interface"

export = async (web3: Web3, blockConfirmations: number, initialBlockNumber: number | null, pollingTime: number): Promise<NewBlockEventEmitterInterface> => {
  const newBlock: NewBlockEventEmitterInterface = new EventEmitter()
  let previousBlockNumber = initialBlockNumber
  const pollingBlockNumber = async () => {
    try {
      const latestBlockNumber = await web3.eth.getBlockNumber()
      const lastBlockNumber = latestBlockNumber - blockConfirmations
      if (previousBlockNumber === null) {
        previousBlockNumber = lastBlockNumber
      }
      for (let blockNumber = previousBlockNumber + 1; blockNumber <= lastBlockNumber; blockNumber++) {
        newBlock.emit('newBlock', blockNumber)
        previousBlockNumber = blockNumber
      }
    } catch (error) {
      console.error("catch polling", error)
    }
    return setTimeout(pollingBlockNumber, pollingTime)
  }
  pollingBlockNumber()
  return newBlock
}
