import Web3 from "web3"
import Service from "mesg-js/lib/service/service"

export = async (mesg: Service, web3: Web3, blockNumber: number) => {
  const block = await web3.eth.getBlock(blockNumber, true)
  block.transactions.forEach(async (transaction) => {
    try {
      const receipt = await web3.eth.getTransactionReceipt(transaction.hash)
      await mesg.emitEvent('transaction', {
        transactionHash: transaction.hash,
        transactionIndex: transaction.transactionIndex,
        blockHash: transaction.blockHash,
        blockNumber: transaction.blockNumber,
        from: transaction.from,
        to: transaction.to,
        status: receipt.status,
        value: transaction.value,
        gasPrice: transaction.gasPrice,
        gas: transaction.gas,
        gasUsed: receipt.gasUsed,
        input: transaction.input
      })
      if (receipt.logs !== undefined) {
        for (let i = 0; i < receipt.logs.length; i++) {
          const log = receipt.logs[i]
          await mesg.emitEvent('log', {
            address: log.address,
            data: log.data,
            topics: log.topics,
            logIndex: log.logIndex,
            transactionHash: log.transactionHash,
            transactionIndex: log.transactionIndex,
            blockHash: log.blockHash,
            blockNumber: log.blockNumber,
          })
        }
      }
    }
    catch (error) {
      console.error('catch transactions', error)
    }
  })
}