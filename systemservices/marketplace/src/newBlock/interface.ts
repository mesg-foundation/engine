import { EventEmitter } from "events"

export interface NewBlockEventEmitterInterface extends EventEmitter {
  on(event: 'newBlock', listener: (blockNumber: number) => void): this
  on(event: 'error', listener: (error: Error) => void): this
  
  emit(event: 'newBlock', blockNumber: number): boolean
  emit(event: 'error', blockHeader: Error): boolean
}