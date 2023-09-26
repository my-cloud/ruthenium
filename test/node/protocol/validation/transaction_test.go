package validation

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"testing"
)

func Test(t *testing.T) {
	marshaledTransactionForId, _ := json.Marshal(struct {
		Inputs    []*network.InputResponse  `json:"inputs"`
		Outputs   []*network.OutputResponse `json:"outputs"`
		Timestamp int64                     `json:"timestamp"`
	}{
		Timestamp: 0,
	})
	transactionHash := sha256.Sum256(marshaledTransactionForId)
	transactionResponse := &network.TransactionResponse{
		Id:        fmt.Sprintf("%x", transactionHash),
		Timestamp: 0,
	}
	var transaction *validation.Transaction
	marshaledTransaction, _ := json.Marshal(transactionResponse)
	_ = json.Unmarshal(marshaledTransaction, &transaction)

	fmt.Println("tx1 id", transaction.Id())

	marshaledTransaction2ForId, _ := json.Marshal(struct {
		Inputs    []*network.InputResponse  `json:"inputs"`
		Outputs   []*network.OutputResponse `json:"outputs"`
		Timestamp int64                     `json:"timestamp"`
	}{
		Timestamp: 1,
	})
	transactionHash2 := sha256.Sum256(marshaledTransaction2ForId)
	transactionResponse2 := &network.TransactionResponse{
		Id:        fmt.Sprintf("%x", transactionHash2),
		Timestamp: 1,
	}
	var transaction2 *validation.Transaction
	marshaledTransaction2, _ := json.Marshal(transactionResponse2)
	_ = json.Unmarshal(marshaledTransaction2, &transaction2)

	fmt.Println("tx2 id", transaction2.Id())

	var transactionResponses []*network.TransactionResponse
	transactionResponses = append(transactionResponses, transactionResponse)
	transactionResponses = append(transactionResponses, transactionResponse2)

	var transactions []*validation.Transaction
	marshaledTransactions, _ := json.Marshal(transactionResponses)
	err := json.Unmarshal(marshaledTransactions, &transactions)

	fmt.Println(err)

	fmt.Println(transactions)
	//fmt.Println("tx1 id", transactions[0].Id())
	//fmt.Println("tx2 id", transactions[1].Id())

}

func Test2(t *testing.T) {
	var addresses map[string]bool
	addresses["hello"] = true
	delete(addresses, "jnl")
	addresses["hello"] = false
	addresses["new"] = false
	delete(addresses, "hello")
	fmt.Println(addresses)
}

func Test3(t *testing.T) {
	originalData := []byte{25, 65, 84, 56} // This can be any data you want to format and then parse.
	formattedString := fmt.Sprintf("%x", originalData)

	var parsedData []byte
	_, err := fmt.Sscanf(formattedString, "%x", &parsedData)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Original Data:", originalData)
	fmt.Println("Parsed Data:", parsedData)
}

func Test4(t *testing.T) {
	originalData := []byte{25, 65, 84, 56} // This can be any data you want to format and then parse.
	formattedString := fmt.Sprintf("%x", originalData)

	var parsedData []byte
	_, err := fmt.Sscanf(formattedString, "%x", &parsedData)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Original Data:", originalData)
	fmt.Println("Parsed Data:", parsedData)

	bytes, _ := json.Marshal(originalData)
	var unmarshalled []byte
	_ = json.Unmarshal(bytes, &unmarshalled)
	formattedString = fmt.Sprintf("%x", unmarshalled)

	fmt.Println("Marshalled bytes:", bytes)
}
