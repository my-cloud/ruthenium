package protocol

type UtxoFinder func(input Input) (Utxo, error)
