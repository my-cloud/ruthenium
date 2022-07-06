package chain

import (
	"fmt"
	"math/big"
)

type Signature struct {
	// Public key x coordinate
	R *big.Int

	// Can be computed by referring to information
	// like the transactions hash and the temporary public key
	// for generating signature
	S *big.Int
}

func (signature *Signature) String() string {
	return fmt.Sprintf("%x%x", signature.R, signature.S)
}
