# Ethereum Wallet

Manage Ethereum accounts and sign transactions.

# Contents

- [Installation](#Installation)
- [Definitions](#Definitions)
  - [Tasks](#Tasks)
    - [Delete an account](#delete-an-account)
    - [Export an account](#export-an-account)
    - [Import an account](#import-an-account)
    - [List accounts](#list-accounts)
    - [Create a new account](#create-a-new-account)
    - [Sign transaction](#sign-transaction)
- [Test](#Test)

# Installation

## MESG Core

This service requires [MESG Core](https://github.com/mesg-foundation/core) to be installed first.

You can install MESG Core by running the following command or [follow the installation guide](https://docs.mesg.com/guide/start-here/installation.html).

```bash
bash <(curl -fsSL https://mesg.com/install)
```

## Service

Download the source code of this service, and then in the service's folder, run the following command:
```bash
mesg-core service deploy
```

# Definitions


# Tasks

## Delete an account

Task key: `delete`

Delete an account from the wallet. Need the address and its associated passphrase.

### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use to unlock the account. |

### Outputs

#### Error

Output key: `error`

Output when an error occurs.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Message** | `message` | `String` | The error message. |

#### Success

Output key: `success`

Output when the task executes successfully.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |


## Export an account

Task key: `export`

Export an existing account in order to backup it and import it in an other wallet. Respect the Web3 Secret Storage specification. See https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition for more information.

### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use to unlock the account. |

### Outputs

#### Error

Output key: `error`

Output when an error occurs.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Message** | `message` | `String` | The error message. |

#### Success

Output key: `success`

Output when the task executes successfully.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Account** | `account` | `Object` | The account in JSON format. |


## Import an account

Task key: `import`

Import an account. The account have to respect the Web3 Secret Storage specification. See https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition for more information.

### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Account** | `account` | `Object` | The account in JSON format. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use to unlock the account. |

### Outputs

#### Error

Output key: `error`

Output when an error occurs.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Message** | `message` | `String` | The error message. |

#### Success

Output key: `success`

Output when the task executes successfully.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |


## List accounts

Task key: `list`

Return the addresses of existing account.


### Outputs

#### Error

Output key: `error`

Output when an error occurs.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Message** | `message` | `String` | The error message. |

#### Success

Output key: `success`

Output when the task executes successfully.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Addresses** | `addresses` | `String` | List of addresses. |


## Create a new account

Task key: `create`

Create a new account with a passphrase. Make sure to backup the passphrase.

### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Passphrase** | `passphrase` | `String` | Passphrase to use to encrypt the account. |

### Outputs

#### Error

Output key: `error`

Output when an error occurs.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Message** | `message` | `String` | The error message. |

#### Success

Output key: `success`

Output when the task executes successfully.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |


## Sign transaction

Task key: `sign`

Sign a transaction with the specified account.

### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Address** | `address` | `String` | The public address of the account. |
| **Chain ID** | `chainID` | `Number` | The ID of the chain the transaction is for. See https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md#list-of-chain-ids. |
| **Passphrase** | `passphrase` | `String` | Passphrase to use to unlock the account. |
| **Transaction** | `transaction` | `Object` | The transaction to sign. |

### Outputs

#### Error

Output key: `error`

Output when an error occurs.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Message** | `message` | `String` | The error message. |

#### Success

Output key: `success`

Output when the task executes successfully.

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
| **Signed transaction** | `signedTransaction` | `Object` | The data of the signed transaction. |



# Test

A folder `test-data` contains test payloads to easily test the service. Adapt their content accordingly.

## create

```
mesg-core service execute ethwallet --task create --json ./test-data/create.json
```

## Lost

```
mesg-core service execute ethwallet --task list --json ./empty.json
```

## Delete

```
mesg-core service execute ethwallet --task delete --json ./test-data/delete.json
```

## Export

```
mesg-core service execute ethwallet --task export --json ./test-data/export.json
```

## Import

```
mesg-core service execute ethwallet --task import --json ./test-data/import.json
```

## Sign

```
mesg-core service execute ethwallet --task sign --json ./test-data/sign.json
```
