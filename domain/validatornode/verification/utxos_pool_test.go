package verification

//func Test_UtxosByAddress_UnknownAddress_ReturnsEmptyArray(t *testing.T) {
//	// Arrange
//	logger := log.NewLoggerMock()
//	genesisValidatorAddress := ""
//	var genesisAmount uint64 = 0
//	settings := new(validatornode.SettingsProviderMock)
//	settings.GenesisAmountFunc = func() uint64 { return genesisAmount }
//	blockchain := NewBlockchain(nil, settings, nil, logger)
//
//	// Act
//	utxosBytes := blockchain.Utxos(genesisValidatorAddress)
//
//	// Assert
//	var utxos []*ledger.Utxo
//	_ = json.Unmarshal(utxosBytes, &utxos)
//	test.Assert(t, len(utxos) == 0, "utxos should be empty")
//}
//
//func Test_Utxos_UtxoExists_ReturnsUtxo(t *testing.T) {
//	// Arrange
//	registry := new(validatornode.RegistrationsManagerMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := log.NewLoggerMock()
//	neighborsManagerMock := new(network.NeighborsManagerMock)
//	var validationInterval int64 = 1
//	settings := new(validatornode.SettingsProviderMock)
//	settings.GenesisAmountFunc = func() uint64 { return 1 }
//	blockchain := NewBlockchain(registry, settings, neighborsManagerMock, logger)
//	registeredAddress := ""
//	var expectedValue uint64 = 1
//	var genesisTimestamp int64 = 0
//	transaction, _ := ledger.NewRewardTransaction(registeredAddress, true, genesisTimestamp+validationInterval, expectedValue)
//	transactions := []*ledger.Transaction{transaction}
//	transactionsBytes, _ := json.Marshal(transactions)
//	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
//	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, transactionsBytes, []string{registeredAddress})
//	_ = blockchain.AddBlock(genesisTimestamp+2*validationInterval, nil, nil)
//
//	// Act
//	utxosBytes := blockchain.Utxos(registeredAddress)
//
//	// Assert
//	var utxos []*ledger.Utxo
//	_ = json.Unmarshal(utxosBytes, &utxos)
//	actualValue := utxos[0].Value(genesisTimestamp+2*validationInterval, 1, 1, 1)
//	test.Assert(t, actualValue == expectedValue, fmt.Sprintf("utxo amount is %d whereas it should be %d", actualValue, expectedValue))
//}
