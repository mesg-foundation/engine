import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"
import { Marketplace } from "../contracts/Marketplace"
import BigNumber from "bignumber.js"


const hexToAscii = (x: any) => Web3.utils.hexToAscii(x).replace(/\u0000/g, '')

export default (contract: Marketplace) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const servicesCount = new BigNumber(await contract.methods.getServicesCount().call())
    const services = []
    for (let i = new BigNumber(0); servicesCount.isGreaterThan(i); i = i.plus(1)) {
      const service = await contract.methods.services(i.toString()).call()
      const versionsCount = new BigNumber(await contract.methods.getServiceVersionsCount(i.toString()).call())
      const versions = []
      for (let j = new BigNumber(0); versionsCount.isGreaterThan(j); j = j.plus(1)) {
        const version = await contract.methods.getServiceVersion(i.toString(), j.toString()).call()
        versions.push({
          hash: version.hash,
          url: hexToAscii(version.url),
        })
      }
      services.push({
        sid: service.sid,
        owner: service.owner,
        price: service.price,
        versions: versions
      })
    }
    console.log('services', services)
    outputs.success({ services })
  }
  catch (error) {
    console.error('error im listServices', error)
    outputs.error({ message: error.toString() })
  }
}
