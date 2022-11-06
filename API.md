<details open>
<summary>v1.0</summary>

# Node host
This version supports [Gobs][1] only.

### End points
<details>
<summary>Get blocks</summary>

request value: `GET BLOCKS REQUEST`  
response value: [`[]BlockResponse`](#block-response)
</details>
<details>
<summary>Get transactions</summary>

request value: `GET TRANSACTIONS REQUEST`  
response value: [`[]TransactionResponse`](#transaction-response)
</details>
<details>
<summary>Mine</summary>

request value: `MINE REQUEST`  
</details>
<details>
<summary>Start mining</summary>

request value: `START MINING REQUEST`  
</details>
<details>
<summary>Stop mining</summary>

request value: `STOP MINING REQUEST`  
</details>

### Models
<details>
<summary>Block</summary>

#### Block response
```
type BlockResponse struct {
    Timestamp           int64
    PreviousHash        [32]byte
    Transactions        []*TransactionResponse
    RegisteredAddresses []string
}
```
</details>
<details>
<summary>Transaction</summary>

#### Transaction response
```
type TransactionResponse struct {
    RecipientAddress string
    SenderAddress    string
    SenderPublicKey  string
    Signature        string
    Timestamp        int64
    Value            uint64
    Fee              uint64
}
```
</details>

# UI server
<details>
<summary>Add transaction</summary>

![Wiki](https://img.shields.io/badge/POST-brightgreen) `localhost:8080/transaction`  
![Wiki](https://img.shields.io/badge/POST-brightgreen?style=for-the-badge) `localhost:8080/transaction`  
![Wiki](https://img.shields.io/badge/POST-brightgreen?style=plastic) `localhost:8080/transaction`  
parameters:   
request body:   
responses:   
</details>
</details>

[1]: https://go.dev/blog/gob "Gobs official documentation"
