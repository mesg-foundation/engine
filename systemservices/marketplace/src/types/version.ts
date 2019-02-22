import { Manifest } from "./manifest";

export interface Version {
  hash: string;
  manifestSource: string;
  manifestProtocol: string;
  manifest: Manifest|undefined;
  createTime: Date;
}
