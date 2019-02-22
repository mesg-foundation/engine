import BigNumber from "bignumber.js";

export interface Offer {
  index: BigNumber;
  price: BigNumber;
  duration: BigNumber;
  active: boolean;
  createTime: Date
}
