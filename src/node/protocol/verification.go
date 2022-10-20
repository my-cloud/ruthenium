package protocol

import (
	"time"
)

const verificationsCountPerValidation = 6

type Verification struct {
	blockchain *Blockchain
	pool       *Pool
	network    *Network
}

func NewVerification(blockchain *Blockchain, pool *Pool, network *Network) *Verification {
	return &Verification{blockchain, pool, network}
}

func (verification *Verification) Start() {
	timer := validationIntervalInSeconds * time.Second / verificationsCountPerValidation
	ticker := time.NewTicker(timer)
	go func() {
		for {
			for i := 0; i < verificationsCountPerValidation; i++ {
				if i > 0 {
					go verification.verifyBlockchain()
				}
				<-ticker.C
			}
		}
	}()
}

func (verification *Verification) verifyBlockchain() {
	neighbors := verification.network.Neighbors()
	// FIXME lock the transactions pool
	verification.blockchain.Verify(neighbors)
	if verification.blockchain.IsReplaced() {
		verification.pool.Clear()
	}
}
