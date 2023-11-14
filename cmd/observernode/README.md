# Observer Node
The user interface (UI) server lets to have a graphical user interface to easily communicate with a Ruthenium [host node](../validatornode/README.md).
Any other implementation of this UI server can communicate with a node using its [API](../validatornode/README.md#api).
In this repository, the UI is described in a simple `index.html`. Any other implementation of this UI can communicate with the UI server using its [API](#api). 

## Prerequisites
A Ruthenium node must be running.

## Launch
At root level (ruthenium folder), run the ui using the command `go run src/ui/main.go` with the add of some [program arguments](#program-arguments). For example:
```
go run src/ui/main.go -host-ip=0.0.0.0
```

## Program arguments:
```
-port: The TCP port number for the UI server (default: "8080")
-host-ip: The node host IP or DNS address (default: "127.0.0.1")
-host-port: The TCP port number of the host node (accepted values: "10600" for mainnet, "10601" to "10699" for testnet, default: "10600")
-templates-path: The UI templates path (default: "templates")
-log-level: The log level (accepted values: "debug", "info", "warn", "error", "fatal", default: "info")
```

Using a web browser, go to `http://localhost:8080` (If needed, replace `localhost` by the UI server IP address and `8080` by the TCP port number for the UI server)

## API
Base URL: `<UI server IP>:<UI server port>` (example: `localhost:8080`)

### Transactions pool
<details>
<summary><b>Add transaction</b></summary>

![POST](https://img.shields.io/badge/POST-seagreen?style=flat-square)
![/transaction](https://img.shields.io/badge//transaction-dimgray?style=flat-square)

*Description:* Add a transaction to the transactions pool.
* **parameters:** *none*
* **request body:** [TransactionRequest](#transactionrequest)
* **responses:**

  |Code|Description|
    |---|---|
  |201|Transaction added|
  |400|Bad request, if any request argument is invalid|
  |500|Internal server error, if an unexpected condition occurred|
</details>
<details>
<summary><b>Get transaction info</b></summary>

![GET](https://img.shields.io/badge/GET-steelblue?style=flat-square)
![/transaction/info](https://img.shields.io/badge//transaction/info-dimgray?style=flat-square)

*Description:* Get the transaction data needed for a transaction request.
* **parameters:**

  |Name|Description|Example|
      |---|---|---|
  |`address`|42 characters hexadecimal sender wallet address|`0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a`|
  |`value`|64 bits floating-point number value of the transaction|`0`|
* **request body:** *none*
* **responses:**

  |Code|Description|
      |---|---|
  |200|[TransactionInfo](#transactioninfo)|
  |400|Bad request, if any request argument is invalid|
  |405|Method not allowed, if the value exceeds the wallet amount for the given address|
  |500|Internal server error, if an unexpected condition occurred|
</details>
<details>
<summary><b>Get transactions</b></summary>

![GET](https://img.shields.io/badge/GET-steelblue?style=flat-square)
![/transactions](https://img.shields.io/badge//transactions-dimgray?style=flat-square)

*Description:* Get all the transactions of the current transactions pool.
* **parameters:** *none*
* **request body:** *none*
* **responses:**

  |Code|Description|
    |---|---|
  |200|Array of [transactions](#transaction)|
  |500|Internal server error, if an unexpected condition occurred|
</details>

### Wallet
<details>
<summary><b>Get wallet address</b></summary>

![GET](https://img.shields.io/badge/GET-steelblue?style=flat-square)
![/wallet/address](https://img.shields.io/badge//wallet/address-dimgray?style=flat-square)

*Description:* Get the wallet address depending on the given public key.
* **parameters:** *none*

  |Name|Description|Example|
    |---|---|---|
  |`publicKey`|132 characters hexadecimal public key|`0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782`|
* **request body:** *none*
* **responses:**

  |Code|Description|
    |---|---|
  |200|42 characters hexadecimal wallet address|
  |500|Internal server error, if an unexpected condition occurred|
</details>
<details>
<summary><b>Get wallet amount</b></summary>

![GET](https://img.shields.io/badge/GET-steelblue?style=flat-square)
![/wallet/amount](https://img.shields.io/badge//wallet/amount-dimgray?style=flat-square)

*Description:* Get the amount for the given wallet address.
* **parameters:**

  |Name|Description|Example|
    |---|---|---|
  |`address`|42 characters hexadecimal wallet address|`0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a`|
* **request body:** *none*
* **responses:**

  |Code|Description|
    |---|---|
  |200|64 bits floating-point number amount|
  |400|Bad request, if any request argument is invalid|
  |500|Internal server error, if an unexpected condition occurred|
</details>

---

### Schemas

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

#### InputInfo
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
}
```
</td>
<td>

```

The output index
The ID of the transaction holding the output

```
</td>
<td>

```
{
  "output_index": 0
  "transaction_id": "8ae72a72c0c99dc9d41c2b7d8ea67b5a2de25ff4463b1a53816ba179947ce77d"
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

#### TransactionInfo
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
  "inputs":    []InputInfo
  "rest":      uint64
  "timestamp": int64
}
```
</td>
<td>

```

The remaining amount to be used as a value for the output with the sender address
The utxos to be used as inputs of the transaction

```
</td>
<td>

```
{
  "inputs": []
  "rest": 0
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
