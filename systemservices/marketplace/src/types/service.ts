import { Version } from "./version";
import { Offer } from "./offer";
import { Purchase } from "./purchase";

export interface Service {
  owner: string;
  sid: string;
  versions: Version[];
  offers: Offer[];
  purchases: Purchase[];
  createTime: Date;
}
