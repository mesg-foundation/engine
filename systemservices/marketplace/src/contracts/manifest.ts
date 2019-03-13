import Request from "request-promise-native"
import { Manifest } from "../types/manifest";

const getIpfs = async (source: string): Promise<any> => {
  return await Request.get('https://gateway.ipfs.io/ipfs/' + source, { json: true, timeout: 10000 })
}

const getHttp = async (source: string): Promise<any> => {
  return await Request.get(source, { json: true, timeout: 10000 })
}

const get: {[key: string]: (source: string) => Promise<any>} = {
  'ipfs': getIpfs,
  'https': getHttp,
  'http': getHttp,
}

const getManifest = async (protocol: string, source: string): Promise<Manifest|undefined> => {
  try {
    if (!get[protocol]) {
      console.warn('protocol', protocol, 'is not compatible with this service')
      return
    }
    const manifest = await get[protocol](source)
    if (typeof manifest === 'object') {
      return manifest as Manifest
    }
    console.warn('manifest ', protocol, '::', source, 'is not a valid manifest')
    return
  }
  catch (error) {
    console.warn('error while downloading manifest', protocol, '::', source, error.toString())
    return
  }
}

export { getManifest }
