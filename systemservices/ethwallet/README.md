# Ethereum Wallet (ID: ethwallet)

Manage Ethereum accounts and sign transactions.

## Contents

- [Installation](#Installation)
  - [MESG SDK](#MESG-SDK)
  - [Deploy the Service](#Service)
- [Definitions](#Definitions)
  - [Tasks](#Tasks)
    - [List accounts](#list)
    - [Create a new account](#create)
    - [Delete an account](#delete)
    - [Export an account](#export)
    - [Import an account](#import)
    - [Sign transaction](#sign)
    - [Import an account from a private key](#importFromPrivateKey)

## Installation

### MESG SDK

This service requires [MESG SDK](https://github.com/mesg-foundation/engine) to be installed first.

You can install MESG SDK by running the following command or [follow the installation guide](https://docs.mesg.com/guide/start-here/installation.html).

```bash
npm install -g mesg-cli
```

### Deploy the Service

To deploy this service, go to [this service page](https://marketplace.mesg.com/services/ethwallet) on the [MESG Marketplace](https://marketplace.mesg.com) and click the button "get/buy this service".

## Definitions


### Tasks

<h4 id="list">List accounts</h4>

Task key: `list`

Return the addresses of existing account.

  
##### Outputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Addresses** | `addresses` | `String` | List of addresses. |
<h4 id="create">Create a new account</h4>

Task key: `create`

Create a new account with a passphrase. Make sure to backup the passphrase.

##### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Passphrase** | `passphrase` | `String` | Passphrase to use with the account. |
  
##### Outputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
<h4 id="delete">Delete an account</h4>

Task key: `delete`

Delete an account from the wallet. Need the address and its associated passphrase.

##### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use with the account. |
  
##### Outputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
<h4 id="export">Export an account</h4>

Task key: `export`

Export an existing account in order to backup it and import it in an other wallet. Respect the Web3 Secret Storage specification. See https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition for more information.

##### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use with the account. |
  
##### Outputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
| **ID** | `id` | `String` | The id of the account. |
| **Version** | `version` | `Number` | The version used to export the account. |
| **Crypto** | `crypto` | `Object` | The encrypted account. |
<h4 id="import">Import an account</h4>

Task key: `import`

Import an account. The account have to respect the Web3 Secret Storage specification. See https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition for more information.

##### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Account** | `account` | `Object` | The JSON encoded account. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use with the account. |
  
##### Outputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
<h4 id="sign">Sign transaction</h4>

Task key: `sign`

Sign a transaction with the specified account.

##### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use with the account. |
| **Transaction** | `transaction` | `Object` | The transaction to sign. |
  
##### Outputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Signed transaction** | `signedTransaction` | `String` | The signed transaction. |
<h4 id="importFromPrivateKey">Import an account from a private key</h4>

Task key: `importFromPrivateKey`

Import an account from a private key.

##### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Private key** | `privateKey` | `String` | The private key to import. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use with the account. |
  
##### Outputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
