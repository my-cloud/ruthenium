# Validator Node
In this repository, the validator node is one implementation of the Ruthenium protocol. Any other implementation can contribute to run the network if it exposes the same [API](#api) and follows the protocol described in the Ruthenium [whitepaper](https://github.com/my-cloud/ruthenium/wiki/Whitepaper).

## Prerequisites
* A firewall port must be open. The port number will be the value of the `port` [program argument](#program-arguments).
* In order to [validate](https://github.com/my-cloud/ruthenium/wiki/Whitepaper#validation) [blocks](https://github.com/my-cloud/ruthenium/wiki/Whitepaper#block) or get an [income](https://github.com/my-cloud/ruthenium/wiki/Whitepaper#income), the node wallet address must be registered in the [Proof of Humanity](https://github.com/my-cloud/ruthenium/wiki/Whitepaper#proof-of-humanity) registry.

## Launch
At root level (ruthenium folder), run the validator node using the command `go run validatornode/main.go` with the add of some [program argument](#program-arguments). For example:
```
go run validatornode/main.go -private-key=0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd
```

## Program arguments:
```
-mnemonic: The mnemonic (required if the private key is not provided)
-derivation-path: The derivation path (unused if the mnemonic is omitted, default: "m/44'/60'/0'/0/0")
-password: The mnemonic password (unused if the mnemonic is omitted)
-private-key: The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)
-infura-key: The infura key (required to check the proof of humanity)
-ip: The validator node IP or DNS address (detected if not provided)
-port: The TCP port number of the validator node (accepted values: "10600" for mainnet, "10601" to "10699" for testnet, default: "10600")
-settings-path: The settings file path (default: "config/settings.json")
-seeds-path: The seeds file path (default: "config/seeds.json")
-log-level: The log level (accepted values: "debug", "info", "warn", "error", "fatal", default: "info")
```
  
## API
Base URL: `<validator node IP>:<validator node port>` (example: seed-styx.ruthenium.my-cloud.me:10600)

Each request value or response value shall be marshaled to bytes or un-marshaled from bytes. All fields are required.

### History
<details>
<summary><b>Get blocks</b></summary>

![/blocks](https://img.shields.io/badge//blocks-dimgray?style=flat-square)

*Description*: Get the blocks starting from the given height (returned blocks array size is limited).
  * **request value:** 64 bits unsigned integer block height
  * **response value:** Array of [blocks](#block)
</details>
<details>
<summary><b>Get first block timestamp</b></summary>

![/first-block-timestamp](https://img.shields.io/badge//first--block--timestamp-dimgray?style=flat-square)

*Description*: Get the first block timestamp.
  * **request value:** *none*
  * **response value:** 64 bits integer timestamp in nanoseconds
</details>

### Network
<details>
<summary><b>Share targets</b></summary>

![/targets](https://img.shields.io/badge//targets-dimgray?style=flat-square)

*Description:* Share known validator node targets.
* **request value:** Array of target strings (IP and port, *e.g.* ["0.0.0.0:0000", "1.1.1.1:1111"])
* **response value:** *none*
</details>

### Payment
<details>
<summary><b>Add transaction</b></summary>

![/transaction](https://img.shields.io/badge//transaction-dimgray?style=flat-square)

*Description:* Add a transaction to the transactions pool.
* **request value:** [TransactionRequest](#transactionrequest)
* **response value:** *none*
</details>
<details>
<summary><b>Get transactions</b></summary>

![/transactions](https://img.shields.io/badge//transactions-dimgray?style=flat-square)

*Description:* Get all the transactions of the current transactions pool.
* **request value:** *none*
* **response value:** Array of [transactions](#transaction)
</details>

### Protocol
<details>
<summary><b>Get settings</b></summary>

![/settings](https://img.shields.io/badge//settings-dimgray?style=flat-square)

*Description:* Get protocol settings.
* **request value:** [settings](#settings)
* **response value:** *none*
</details>

### Wallet
<details>
<summary><b>Get UTXOs</b></summary>

![/utxos](https://img.shields.io/badge//utxos-dimgray?style=flat-square)

*Description:* Get all the UTXOs for the given wallet address.
* **request value:** wallet address string
* **response value:** Array of [UTXOs](#utxo)
</details>

---

### Schemas

#### Block
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
{
  "previous_hash":                [32]byte
  "added_registered_addresses":   []string
  "removed_registered_addresses": []string
  "timestamp":                    int64
  "transactions":                 []Transaction
}
```
</td>
<td>

```

The hash of the previous block in the chain
The added addresses registered in the PoH registry compared to the previous block
The removed addresses registered in the PoH registry compared to the previous block
The block timestamp
The block transactions

```
</td>
<td>

```
{
  "previous_hash": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32]
  "added_registered_addresses": ["0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a"]
  "removed_registered_addresses": ["0xb1477DcBBea001a339a92b031d14a011e36D008F"]
  "timestamp": 1667768884780639700
  "transactions": []
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
{
  "output_index":   uint16
  "transaction_id": string
  "public_key":     string
  "signature":      string
}
```
</td>
<td>

```

The output index
The ID of the transaction holding the output
The output recipient public key
The output signature

```
</td>
<td>

```
{
  "output_index": 0
  "transaction_id": "8ae72a72c0c99dc9d41c2b7d8ea67b5a2de25ff4463b1a53816ba179947ce77d"
  "public_key": "0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782"
  "signature": "4f3b24cbb4d2c13aaf60518fce70409fd29e1668db1c2109c0eac58427c203df59788bade6d5f3eb9df161b4ed3de451bac64f4c54e74578d69caf8cd401a38f"
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
{
  "address":     string
  "is_yielding": bool
  "value":       uint64
}
```
</td>
<td>

```

The address of this output recipient
Whether this output should be used for income calculation
The value at the transaction timestamp

```
</td>
<td>

```
{
  "address": "0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a"
  "is_yielding": true
  "value": 0
}
```
</td>
</tr>
</table>

#### Settings
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
{
	"blocksCountLimit":                 uint64
	"genesisAmount":                    uint64
	"halfLifeInDays":                   float64
	"incomeBase":                       uint64
	"incomeLimit":                      uint64
	"maxOutboundsCount":                int
	"minimalTransactionFee":            uint64
	"smallestUnitsPerCoin":             uint64
	"synchronizationIntervalInSeconds": int
	"validationIntervalInSeconds":      int64
	"validationTimeoutInSeconds":       int64
	"verificationsCountPerValidation":  int64
}
```
</td>
<td>

```

	The maximum blocks count returned by a blocks request
	The genesis amount in smallest unit
	The half-life in days
	The income base in smallest unit
	The income limit in smallest unit
	The maximum node outbounds count
	The minimal transaction fee in smallest unit
	The number of smallest uints per coin
	The synchronization interval in seconds
	The validation interval in seconds
	The validation timeout in seconds
	The verifications count per validation

```
</td>
<td>

```
{
  "blocksCountLimit": 1440,
  "genesisAmount": 5000000000000,
  "halfLifeInDays": 373.59,
  "incomeBase": 100000000000,
  "incomeLimit": 5000000000000,
  "maxOutboundsCount": 8,
  "minimalTransactionFee": 1000,
  "smallestUnitsPerCoin": 100000000,
  "synchronizationIntervalInSeconds": 10,
  "validationIntervalInSeconds": 60,
  "validationTimeoutInSeconds": 5,
  "verificationsCountPerValidation": 6
}
```
</td>
</tr>
</table>

#### Transaction
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
{
  "id":        string
  "inputs":    []Input
  "outputs":   []Output
  "timestamp": int64
}
```
</td>
<td>

```

The ID
The inputs
The outputs
The timestamp

```
</td>
<td>

```
{
  "id": "30148389df42b7cd0cb0d3ce951133da3f36ff4e1581d108da1ee05bacad64b7"
  "inputs": []
  "outputs": []
  "timestamp": 1667768884780639700
}
```
</td>
</tr>
</table>

#### TransactionRequest
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
{
  "transaction":                    Transaction
  "transaction_broadcaster_target": string
}
```
</td>
<td>

```

The transaction
The transaction broadcaster target

```
</td>
<td>

```
{
  "transaction": {}
  "transaction_broadcaster_target": "0.0.0.0:0000"
}
```
</td>
</tr>
</table>

#### UTXO
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
{
  "address":        string
  "block_height":   int
  "is_yielding":     bool
  "output_index":   uint16
  "transaction_id": string
  "value":          uint64
}
```
</td>
<td>

```

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
  "address": "0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a"
  "block_height": 0
  "is_yielding": true
  "output_index": 0
  "transaction_id": "8ae72a72c0c99dc9d41c2b7d8ea67b5a2de25ff4463b1a53816ba179947ce77d"
  "value": 0
}
```
</td>
</tr>
</table>

[1]: https://go.dev/blog/gob "Gobs official documentation"
