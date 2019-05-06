export interface Manifest {
  version: '1'
  service: {
    deployment: {
      type: 'ipfs'|'http'|'https'
      source: string
    }
    definition: any, // TODO: add ts definition or json schema validation
    readme?: string
  }
}
