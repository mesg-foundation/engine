export interface Manifest {
  version: number
  service: {
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
