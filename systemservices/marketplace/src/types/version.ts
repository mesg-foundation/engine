import { Manifest } from "./manifest";

export interface Version {
  hash: string;
  manifest: string;
  manifestProtocol: string;
  manifestData: Manifest|undefined;
  createTime: Date;
}
