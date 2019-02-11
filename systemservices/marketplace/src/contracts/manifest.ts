import { Manifest } from "../types/service";
import Request from "request-promise-native"

const getManifest = async (protocol: string, source: string): Promise<Manifest|undefined> => {
  switch (protocol) {
    case 'https':
    case 'http':
      // return await Request.get(source, { json: true }) as Manifest
  }
  console.warn('protocol ' + protocol + ' is not compatible with this service')
  return undefined
}

export { getManifest }
