export interface Manifest {
  version: '1'
  service: {
    hash: string
    hashVersion: '1'
    deployment: {
      type: 'ipfs'|'http'|'https'
      source: string
    }
    definition: any,
    readme?: string
  }
}
