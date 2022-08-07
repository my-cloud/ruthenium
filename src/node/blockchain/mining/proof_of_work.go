package mining

import (
	"fmt"
	"strings"
)

const MiningDifficulty = 3

type ProofOfWork struct {
	hash [32]byte
}

func NewProofOfWork(hash [32]byte) *ProofOfWork {
	return &ProofOfWork{hash}
}

func (pow *ProofOfWork) IsInValid() bool {
	hashStr := fmt.Sprintf("%x", pow.hash)
	zeros := strings.Repeat("0", MiningDifficulty)
	return hashStr[:MiningDifficulty] != zeros
}
