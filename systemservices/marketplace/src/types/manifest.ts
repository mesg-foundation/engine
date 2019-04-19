export interface Manifest {
  version: '1'
  service: {
    hashVersion: '1'
    deployment: {
      type: 'ipfs'|'http'|'https'
      source: string
    }
    definition: any, // TODO: add ts definition or json schema validation
    readme?: string
  }
}
