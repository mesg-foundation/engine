import { Manifest } from "../types/service";
import Request, { RequestPromise } from "request-promise-native"

const getManifest = async (protocol: string, source: string): Promise<Manifest|undefined> => {
  const data: any = {
    'ipfs': getIpfs,
    'https': getHttp,
    'http': getHttp,
  }
  if (!data[protocol]) {
    console.warn('protocol ' + protocol + ' is not compatible with this service')
    return
  }
  const manifest = await data[protocol](source)
  if (typeof manifest === 'object') {
    return manifest as Manifest
  }
  console.warn('source ' + source + ' is not an object')
}

const getIpfs = async (source: string): Promise<any> => {
  return await Request.get('https://ipfs.io/ipfs/' + source, { json: true })
}
const getHttp = async (source: string): Promise<any> => {
  return await Request.get(source, { json: true })
}

export { getManifest }
