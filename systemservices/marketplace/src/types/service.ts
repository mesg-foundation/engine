interface Service {
  owner: string;
  sid: string;
  price: string;
  versions: Version[];
}

interface Version {
  hash: String;
  metadata: Metadata;
}

interface Metadata {
  version: Number
  service: {
    deployment: {
      source: String
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
    readme: String
    // author: String
    // logo
    // tags
    
  }
}

export { Service, Version, Metadata }