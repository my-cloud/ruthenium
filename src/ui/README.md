# UI server
The user interface (UI) server lets to have a graphical user interface to easily communicate with a Ruthenium [host node](../node/README.md).
Any other implementation of this UI server can communicate with a node using its [API](../node/README.md#api).
In this repository, the UI is described in a simple `index.html`. Any other implementation of this UI can communicate with the UI server using its [API](#api). 

## Prerequisites
A Ruthenium node must be running.

## Launch
At root level (ruthenium folder), run the ui using the command `go run src/ui/main.go` with the add of some [program arguments](#program-arguments). For example:
```
go run src/ui/main.go -host-ip=0.0.0.0 -private-key=0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd
```

## Program arguments:
```
-mnemonic: The mnemonic (required if the private key is not provided)
-derivation-path: The derivation path (unused if the mnemonic is omitted, default: "m/44'/60'/0'/0/0")
-password: The mnemonic password (unused if the mnemonic is omitted)
-private-key: The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)
-port: The TCP port number for the UI server (default: 8080)
-host-ip: The node host IP address
-host-port: The TCP port number of the host node (default: 8106)
-templates-path: The UI templates path (default: "templates")
-log-level: The log level (default: "info")
-delay-startup: The delay prior to start up the server. Valid unit times: ms, s, m (default: 0s)
```

Using a web browser, go to `http://localhost:8080` (If needed, replace `localhost` by the UI server IP address and `8080` by the TCP port number for the UI server)

## API
Base url: `<UI server IP>:<UI server port>` (example: `localhost:8080`)

### Transactions pool
<details>
<summary><b>Add transaction</b></summary>

![POST](https://img.shields.io/badge/POST-seagreen?style=flat-square)
![Transaction](https://img.shields.io/badge//transaction-dimgray?style=flat-square)

*Description:* Add a transaction to the transactions pool.
* **parameters:** *none*
* **request body:** [Transaction request](#transaction-request)
* **responses:**

  |Code|Description|
    |---|---|
  |201|Transaction added|
  |400|Bad request|
  |500|Internal server error|
</details>
<details>
<summary><b>Get transactions</b></summary>

![GET](https://img.shields.io/badge/GET-steelblue?style=flat-square)
![Transactions](https://img.shields.io/badge//transactions-dimgray?style=flat-square)

*Description:* Get all the transactions of the current transactions pool.
* **parameters:** *none*
* **request body:** *none*
* **responses:**

  |Code|Description|
    |---|---|
  |200|Array of [transaction responses](#transaction-response)|
  |500|Internal server error|
</details>

### Validation
<details>
<summary><b>Start validation</b></summary>

![POST](https://img.shields.io/badge/POST-seagreen?style=flat-square)
![Validation start](https://img.shields.io/badge//validation/start-dimgray?style=flat-square)

*Description:* Start validating one block per minute.
* **parameters:** *none*
* **request body:** *none*
* **responses:**

  |Code|Description|
    |---|---|
  |200|Validation started|
  |500|Internal server error|
</details>
<details>
<summary><b>Stop validation</b></summary>

![POST](https://img.shields.io/badge/POST-seagreen?style=flat-square)
![Validation stop](https://img.shields.io/badge//validation/stop-dimgray?style=flat-square)

*Description:* Stop validating one block per minute.
* **parameters:** *none*
* **request body:** *none*
* **responses:**

  |Code|Description|
    |---|---|
  |200|Validation stopped|
  |500|Internal server error|
</details>

### Wallet
<details>
<summary><b>Create wallet</b></summary>

![POST](https://img.shields.io/badge/POST-seagreen?style=flat-square)
![Wallet](https://img.shields.io/badge//wallet-dimgray?style=flat-square)

*Description:* Create a new wallet instance with the arguments provided at UI server program launch.
* **parameters:** *none*
* **request body:** *none*
* **responses:**

  |Code|Description|
    |---|---|
  |201|[Wallet response](#wallet-response)|
  |500|Internal server error|
</details>
<details>
<summary><b>Get wallet amount</b></summary>

![GET](https://img.shields.io/badge/GET-steelblue?style=flat-square)
![Wallet amount](https://img.shields.io/badge//wallet/amount-dimgray?style=flat-square)

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
  |400|Bad request|
  |500|Internal server error|
</details>

### Schemas

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
  SenderPrivateKey string
  SenderAddress    string
  RecipientAddress string
  Value            string
}
```
</td>
<td>

```
The transaction data structure for request
The sender wallet private key
The sender wallet address
The recipient wallet address
The value

```
</td>
<td>

```
{
  "SenderPrivateKey": 0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd
  "SenderAddress":    0x9C69443c3Ec0D660e257934ffc1754EB9aD039CB
  "RecipientAddress": 0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a
  "Value":            100000000
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

#### Wallet response
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
Wallet {
  PrivateKey string
  PublicKey  string
  Address    string
}
```
</td>
<td>

```
The wallet data structure
The wallet private key
The wallet public key
The wallet address

```
</td>
<td>

```
{
  "PrivateKey": 0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd
  "PublicKey":  0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782
  "Address":    0x9C69443c3Ec0D660e257934ffc1754EB9aD039CB
}
```
</td>
</tr>
</table>
