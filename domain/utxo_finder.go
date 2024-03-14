package domain

type UtxoFinder func(input InputInfo) (Utxo, error)
