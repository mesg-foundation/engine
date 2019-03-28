export interface Manifest {
  version: '1'
  service: {
    hash: string
    hashVersion: '1'
    deployment: {
      type: string
      source: 'ipfs'|'http'|'https'
    }
    definition: any,
    readme?: string
  }
}
