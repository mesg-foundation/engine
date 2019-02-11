import BigNumber from "bignumber.js";

interface Service {
  owner: string;
  sid: string;
  versions: Version[];
  offers: Offer[];
  purchases: Purchase[];
}

interface Version {
  hash: string;
  manifestSource: string;
  manifestProtocol: string;
  manifest: Manifest|undefined;
}

interface Offer {
  index: BigNumber;
  price: BigNumber;
  duration: BigNumber;
  active: boolean;
}

interface Purchase {
  purchaser: string;
  expire: BigNumber;
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

export { Service, Version, Offer, Purchase, Manifest }