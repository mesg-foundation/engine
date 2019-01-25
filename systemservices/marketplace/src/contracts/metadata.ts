import { Metadata } from "../types/service";
import Request from "request-promise-native"

const getMetadata = async (url: string): Promise<Metadata> => {
  return await Request.get(url, { json: true }) as Metadata
}

export { getMetadata }
