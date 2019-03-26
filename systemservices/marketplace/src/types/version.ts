import { Manifest } from "./manifest";

export interface Version {
  versionHash: string;
  manifest: Manifest;
  createTime: Date;
}
