package protocol

type UtxoFinder func(input InputInfo) (Utxo, error)
