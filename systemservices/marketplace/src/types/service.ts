interface Service {
  owner: string;
  sid: string;
  versions: Version[];
}

interface Version {
  hash: string;
  manifestSource: string;
  manifestProtocol: string;
  manifest: Manifest;
}

interface Manifest {
  version: number
  service: {
    deployment: {
      source: string
      // env
      // etc...
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
    // author: string
    // logo
    // tags
    
  }
}

export { Service, Version, Manifest }