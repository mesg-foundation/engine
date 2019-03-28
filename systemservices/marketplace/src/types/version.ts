import { Manifest } from "./manifest";

export interface Version {
  versionHash: string;
  manifest: string;
  manifestProtocol: string;
  manifestData: Manifest|null;
  createTime: Date;
}
