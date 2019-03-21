# Marketplace service

If you get error like `error while downloading manifest ipfs :: QmUoE4fsthoirJRWtqosNMQF79DYBWbgtGw37Sgg2eMdre RequestError: Error: ESOCKETTIMEDOUT` it's because ipfs gateway remove the file from its cache. We need to host an IPFS node to guarantee availability.

### Start service

```
./dev-cli service dev ./systemservices/marketplace
```

### List service

```
./dev-cli service execute marketplace --task listServices --json ./systemservices/marketplace/test-data/empty.json
```

### Get service

```
./dev-cli service execute marketplace --task getService --json ./systemservices/marketplace/test-data/getService.json
```

### Publish service version

```
./dev-cli service execute marketplace --task publishServiceVersion --json ./systemservices/marketplace/test-data/publishServiceVersion.json
```

### Send sign transaction

```
./dev-cli service execute marketplace --task sendSignedTransaction --json ./systemservices/marketplace/test-data/sendSignedTransaction.json
```

### Create service offer

```
./dev-cli service execute marketplace --task createServiceOffer --json ./systemservices/marketplace/test-data/createServiceOffer.json
```

### Purchase

```
./dev-cli service execute marketplace --task purchase --json ./systemservices/marketplace/test-data/purchase.json
```

### Is Authorized

#### with hash
```
./dev-cli service execute marketplace --task isAuthorized --json ./systemservices/marketplace/test-data/isAuthorized.json
```

#### with sid
```
./dev-cli service execute marketplace --task isAuthorized --json ./systemservices/marketplace/test-data/isAuthorized-sid.json
```


### Note:

To send the transaction returned by some commands, execute:

```
./dev-cli service execute ethwallet --task sign --json sign.json
```

With `sign.json` file like (you need to replace the data):
```json
{
  "address": "0xcAB79fA69c68CB4C65fa5C6E05BC4dBa5FB57D11",
  "passphrase": "1",
  "transaction": {
    "chainID": 3,
    "data": "0x095ea7b300000000000000000000000094f4cb92fe9f547574aec617b1594b13abd47ad300000000000000000000000000000000000000000000003635c9adc5dea00000",
    "gas": 1000000,
    "gasPrice": "1000000000",
    "nonce": 21,
    "to": "0x5861B3DC52339d4f976B7fa5d80dB6cd6f477F1B",
    "value": "0"
  }
}
```

Then, publish transaction, copy past after `signedTransaction=` the tx outputted:

```
./dev-cli service execute marketplace --task sendSignedTransaction --data signedTransaction=
```

### Decode transaction receipt

#### decodeLogCreateServiceOffer

```
./dev-cli service execute marketplace --task decodeLogCreateServiceOffer --json ./systemservices/marketplace/test-data/decodeLogCreateServiceOffer.json
```

#### decodeLogPublishServiceVersion

```
./dev-cli service execute marketplace --task decodeLogPublishServiceVersion --json ./systemservices/marketplace/test-data/decodeLogPublishServiceVersion.json 
```

#### decodeLogPurchase

```
./dev-cli service execute marketplace --task decodeLogPurchase --json ./systemservices/marketplace/test-data/decodeLogPurchase.json
```
