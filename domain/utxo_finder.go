package domain

type UtxoFinder func(input InputInfoProvider) (UtxoInfoProvider, error)
