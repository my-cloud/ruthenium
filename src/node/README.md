# Host node
In this repository, the host node is an implementation following the Ruthenium protocol. Any other implementation can contribute to run the network if it exposes the same [API](#api) and follows the protocol described in the Ruthenium [whitepaper](https://github.com/my-cloud/ruthenium/wiki/Whitepaper).

## Prerequisites
* A DNS port must be open. The port number will be the value of the `port` [program argument](#program-arguments).
* If you want to [validate](https://github.com/my-cloud/ruthenium/wiki/Whitepaper#validation) [blocks](https://github.com/my-cloud/ruthenium/wiki/Whitepaper#block) or get an [income](https://github.com/my-cloud/ruthenium/wiki/Whitepaper#income), you must be registered in the [Proof of Humanity](https://github.com/my-cloud/ruthenium/wiki/Whitepaper#proof-of-humanity) registry with an Ethereum wallet address for which you are the owner of the private key.

## Launch
At root level (ruthenium folder), run the node using the command `go run src/node/main.go` with the add of some [program argument](#program-arguments). For example:
```
go run src/node/main.go -private-key=0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd
```

## Program arguments:
```
-mnemonic: The mnemonic (required if the private key is not provided)
-derivation-path: The derivation path (unused if the mnemonic is omitted, default: "m/44'/60'/0'/0/0")
-password: The mnemonic password (unused if the mnemonic is omitted)
-privateKey: The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)
-port: TCP port number of the host node (default: 8106)
-configuration-path: The configuration files path (default: "config")
-log-level: The log level (default: "info")
```
  
## API
This version supports [Gobs][1] only.

### Blockchain
<details>
<summary><b>Get blocks</b></summary>

*Description*: Get all the blocks of the blockchain for the current timestamp.
  * **request value:** `GET BLOCKS`  
  * **response value:** Array of [block responses](#block-response)
</details>

### Network
<details>
<summary><b>Share targets</b></summary>

*Description:* Share known validator node targets (IP and port).
* **request value:** Array of [target requests](#target-request)  
* **response value:** no response
</details>

### Transactions pool
<details>
<summary><b>Add transaction</b></summary>

*Description:* Add a transaction to the transactions pool.
* **request value:** [TransactionRequest](#transaction-request)  
* **response value:** *none*
</details>
<details>
<summary><b>Get transactions</b></summary>

*Description:* Get all the transactions of the current transactions pool.
* **request value:** `GET TRANSACTIONS`  
* **response value:** Array of [transaction responses](#transaction-response)
</details>

### Validation
<details>
<summary><b>Start validation</b></summary>

*Description:* Start validating one block per minute.
* **request value:** `START VALIDATION`  
* **response value:** *none*
</details>
<details>
<summary><b>Stop validation</b></summary>

*Description:* Stop validating one block per minute.
* **request value:** `STOP VALIDATION`  
* **response value:** *none*
</details>

### Wallet
<details>
<summary><b>Get amount</b></summary>

*Description:* Get the amount for the given wallet address.
* **request value:** [Amount request](#amount-request)  
* **response value:** [Amount response](#amount-response)
</details>

---
<details open>
<summary style="font-size:24px"><b>Schemas</b></summary>

### Amount request
<table>
<th>
Schema
</th>
<th>
Description
</th>
<th>
Example
</th>
<tr>
<td>

```
AmountRequest {
  Address string
}
```
</td>
<td>

```
The amount data structure for request
The wallet address for which to get the amount

```
</td>
<td>

```
{
  "Address": 0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a
}
```
</td>
</tr>
</table>

### Amount response
<table>
<th>
Schema
</th>
<th>
Description
</th>
<th>
Example
</th>
<tr>
<td>

```
AmountResponse {
  Amount uint64
}
```
</td>
<td>

```
The amount data structure for response
The amount

```
</td>
<td>

```
{
  "Amount": 100000000
}
```
</td>
</tr>
</table>

### Block response
<table>
<th>
Schema
</th>
<th>
Description
</th>
<th>
Example
</th>
<tr>
<td>

```
BlockResponse {
  Timestamp           int64
  PreviousHash        [32]byte
  Transactions        []TransactionResponse
  RegisteredAddresses []string
}
```
</td>
<td>

```
The block data structure for response
The block timestamp
The hash of the previous block in the chain
The block transactions
The addresses registered in the PoH registry

```
</td>
<td>

```
{
  "Timestamp": 1667768884780639700
  "PreviousHash": [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32]
  "Transactions":        []
  "RegisteredAddresses": [ 0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a ]
}
```
</td>
</tr>
</table>

### Target request
<table>
<th>
Schema
</th>
<th>
Description
</th>
<th>
Example
</th>
<tr>
<td>

```
TargetRequest {
    Ip   string
    Port uint16
}
```
</td>
<td>

```
The target data structure for request
The IP
The port

```
</td>
<td>

```
{
  "Ip":   0.0.0.0
  "Port": 8106
}
```
</td>
</tr>
</table>

### Transaction request

<table>
<th>
Schema
</th>
<th>
Description
</th>
<th>
Example
</th>
<tr>
<td>

```
TransactionRequest {
  RecipientAddress string
  SenderAddress    string
  SenderPublicKey  string
  Signature        string
  Timestamp        int64
  Value            uint64
  Fee              uint64
}
```
</td>
<td>

```
The transaction data structure for request
The recipient wallet address
The sender wallet address
The sender wallet public key
The signature generated by both the private and public keys
The timestamp
The value
The fee

```
</td>
<td>

```
{
  "RecipientAddress": 0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a
  "SenderAddress":    0x9C69443c3Ec0D660e257934ffc1754EB9aD039CB
  "SenderPublicKey":  0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782
  "Signature":        4f3b24cbb4d2c13aaf60518fce70409fd29e1668db1c2109c0eac58427c203df59788bade6d5f3eb9df161b4ed3de451bac64f4c54e74578d69caf8cd401a38f
  "Timestamp":        1667768884780639700
  "Value":            100000000
  "Fee":              1000
}
```
</td>
</tr>
</table>

### Transaction response

<table>
<th>
Schema
</th>
<th>
Description
</th>
<th>
Example
</th>
<tr>
<td>

```
TransactionResponse {
  RecipientAddress string
  SenderAddress    string
  SenderPublicKey  string
  Signature        string
  Timestamp        int64
  Value            uint64
  Fee              uint64
}
```
</td>
<td>

```
The transaction data structure for response
The recipient wallet address
The sender wallet address
The sender wallet public key
The signature generated by both the private and public keys
The timestamp
The value
The fee

```
</td>
<td>

```
{
  "RecipientAddress": 0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a
  "SenderAddress":    0x9C69443c3Ec0D660e257934ffc1754EB9aD039CB
  "SenderPublicKey":  0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782
  "Signature":        4f3b24cbb4d2c13aaf60518fce70409fd29e1668db1c2109c0eac58427c203df59788bade6d5f3eb9df161b4ed3de451bac64f4c54e74578d69caf8cd401a38f
  "Timestamp":        1667768884780639700
  "Value":            100000000
  "Fee":              1000
}
```
</td>
</tr>
</table>
</details>

[1]: https://go.dev/blog/gob "Gobs official documentation"
