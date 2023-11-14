package validatornode

type UtxoFinder func(input InputInfo) (Utxo, error)
