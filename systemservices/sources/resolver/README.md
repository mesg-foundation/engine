# System Resolver Service

## Commands

### Add peers command

```
mesg-core service execute SERVICE_ID -t addPeers -d addresses="[\"PEER_IP:50052\"]"
```

### Resolve command

```
mesg-core service execute SERVICE_ID -t resolve -d serviceID=SERVICE_ID_TO_RESOLVE
```



# Tasks

## addPeers

Task key: `addPeers`



### Inputs

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **addresses** | `Object` |  |


### Outputs

##### error

Output key: `error`

Output when error

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **message** | `String` | The error&#39;s message |

##### success

Output key: `success`

Output when success

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **addresses** | `Object` |  |




## resolve

Task key: `resolve`



### Inputs

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **serviceID** | `String` |  |


### Outputs

##### error

Output key: `error`

Output when error

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **message** | `String` | The error&#39;s message |

##### found

Output key: `found`

A peer matches

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **address** | `String` |  |
| **serviceID** | `String` |  |

##### notFound

Output key: `notFound`

No peers have been found

| **Key** | **Type** | **Description** |
| --- | --- | --- |
| **serviceID** | `String` |  |




