# Host node
In this repository, the host node is an implementation following the Ruthenium protocol. Any other implementation can contribute to run the network if it exposes the same [API](#api) and follows the protocol described in the Ruthenium [whitepaper](https://github.com/my-cloud/ruthenium/wiki/Whitepaper).

## Prerequisites
* A firewall port must be open. The port number will be the value of the `port` [program argument](#program-arguments).
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
-private-key: The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)
-infura-key: The infura key (required to check the proof of humanity)
-ip: The node IP or DNS address (detected if not provided)
-port: The TCP port number of the host node (accepted values: "10600" for mainnet, "10601" to "10699" for testnet, default: "10600")
-settings-path: The settings file path (default: "config/settings.json")
-seeds-path: The seeds file path (default: "config/seeds.json")
-log-level: The log level (accepted values: "debug", "info", "warn", "error", "fatal", default: "info")
```
  
## API
`base-url`: `<node IP>:<node port>` (example: seed-styx.ruthenium.my-cloud.me:10600)

Each request value or response value shall be marshaled to bytes or un-marshaled from bytes. All fields are required.

### Blockchain
<details>
<summary><b>Get blocks</b></summary>

*Route*: `base-url/blocks`

*Description*: Get the blocks starting from the given height (returned blocks array size is limited).
  * **request value:** 64 bits unsigned integer block height
  * **response value:** Array of [block responses](#block-response)
</details>
<details>
<summary><b>Get first block timestamp</b></summary>

*Route*: `base-url/first-block-timestamp`

*Description*: Get the first block timestamp.
  * **request value:** *none*
  * **response value:** 64 bits integer timestamp in nanoseconds
</details>

### Network
<details>
<summary><b>Share targets</b></summary>

*Route*: `base-url/targets`

*Description:* Share known validator node targets.
* **request value:** Array of string targets (IP and port, *e.g.* 0.0.0.0:0000)
* **response value:** *none*
</details>

### Transactions pool
<details>
<summary><b>Add transaction</b></summary>

*Route*: `base-url/transaction`

*Description:* Add a transaction to the transactions pool.
* **request value:** [Transaction request](#transaction-request)
* **response value:** *none*
</details>
<details>
<summary><b>Get transactions</b></summary>

*Route*: `base-url/transactions`

*Description:* Get all the transactions of the current transactions pool.
* **request value:** `GET TRANSACTIONS`
* **response value:** Array of [transaction responses](#transaction-response)
</details>

### Wallet
<details>
<summary><b>Get UTXOs</b></summary>

*Route*: `base-url/utxos`

*Description:* Get all the UTXOs for the given wallet address.
* **request value:** string wallet address
* **response value:** Array of [UTXO response](#utxo-response)
</details>

---

### Schemas

#### Block response
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
  Timestamp                  int64
  PreviousHash               [32]byte
  Transactions               []TransactionResponse
  AddedRegisteredAddresses   []string
  RemovedRegisteredAddresses []string
}
```
</td>
<td>

```
The data structure for block response
  The block timestamp
  The hash of the previous block in the chain
  The block transactions
  The added addresses registered in the PoH registry compared to the previous block
  The removed addresses registered in the PoH registry compared to the previous block

```
</td>
<td>

```
{
  "Timestamp":                  1667768884780639700
  "PreviousHash":               [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32]
  "Transactions":               []
  "AddedRegisteredAddresses":   [ 0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a ]
  "RemovedRegisteredAddresses": [ 0xb1477DcBBea001a339a92b031d14a011e36D008F ]
}
```
</td>
</tr>
</table>

#### Input
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
Input {
  OutputIndex   uint16
  TransactionId string
  PublicKey     string
  Signature     string
}
```
</td>
<td>

```
The input data structure
  The output index
  The ID of the transaction holding the output
  The output recipient public key
  The output signature

```
</td>
<td>

```
{
  "OutputIndex":   0
  "TransactionId": 8ae72a72c0c99dc9d41c2b7d8ea67b5a2de25ff4463b1a53816ba179947ce77d
  "PublicKey":     0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782
  "Signature":     4f3b24cbb4d2c13aaf60518fce70409fd29e1668db1c2109c0eac58427c203df59788bade6d5f3eb9df161b4ed3de451bac64f4c54e74578d69caf8cd401a38f
}
```
</td>
</tr>
</table>

#### Output
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
Output {
  Address   string
  HasReward bool
  HasIncome bool
  Value     uint64
}
```
</td>
<td>

```
The output data structure
  The address of this output recipient
  Whether this output contains a reward
  Whether this output should be used for income calculation
  The value at the transaction timestamp

```
</td>
<td>

```
{
  "Address":   0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a
  "HasReward": false
  "HasIncome": true
  "Value":     0
}
```
</td>
</tr>
</table>

#### Transaction request
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
  Inputs                       []InputRequest
  Outputs                      []OutputRequest
  Timestamp                    int64
  TransactionBroadcasterTarget string
}
```
</td>
<td>

```
The transaction data structure for request
  The inputs
  The outputs
  The timestamp
  The transaction broadcaster target

```
</td>
<td>

```
{
  "Inputs": []
  "Outputs": []
  "Timestamp":       1667768884780639700
  "TransactionBroadcasterTarget": 0.0.0.0:0000
}
```
</td>
</tr>
</table>

#### Transaction response
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
  Id        string
  Inputs    []*InputResponse
  Outputs   []*OutputResponse
  Timestamp int64
}
```
</td>
<td>

```
The transaction data structure for response
  The ID
  The inputs
  The outputs
  The timestamp

```
</td>
<td>

```
{
  "Id":            30148389df42b7cd0cb0d3ce951133da3f36ff4e1581d108da1ee05bacad64b7
  "Inputs": []
  "Outputs": []
  "Timestamp":     1667768884780639700
}
```
</td>
</tr>
</table>

#### UTXO response
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
type UtxoResponse struct {
  Address       string
  BlockHeight   int
  HasReward     bool
  HasIncome     bool
  OutputIndex   uint16
  TransactionId string
  Value         uint64
}
```
</td>
<td>

```
The data structure for UTXO response
  The address of the output recipient
  The output transaction block height
  Whether the output contains a reward
  Whether the output should be used for income calculation
  The output index
  The ID of the transaction holding the output
  The value at the transaction timestamp

```
</td>
<td>

```
{
  "Address":       0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a
  "BlockHeight":   0
  "HasReward":     false
  "HasIncome":     true
  "OutputIndex":   0
  "TransactionId": 8ae72a72c0c99dc9d41c2b7d8ea67b5a2de25ff4463b1a53816ba179947ce77d
  "Value":         0
}
```
</td>
</tr>
</table>

[1]: https://go.dev/blog/gob "Gobs official documentation"
