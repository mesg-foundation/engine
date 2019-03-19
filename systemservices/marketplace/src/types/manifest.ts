export interface Manifest {
  version: string
  service: {
    hash: string
    hashVersion: string
    deployment: {
      type: string
      source: string
    }
    definition: {
      // basically mesg.yaml
      // name
      // description
      // sid
      // events
      // tasks
      // configuration
      // dependencies
    }
    readme: string
  }
}
