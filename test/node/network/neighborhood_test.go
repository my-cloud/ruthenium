package network

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Synchronize_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	configurationPath := "../../../config"
	jsonFile, _ := os.Open(configurationPath + "/seeds.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	_ = jsonFile.Close()
	var seedsIps []string
	_ = json.Unmarshal(byteValue, &seedsIps)
	watch := node.NewWatchMock()
	logger := log.NewLogger(log.Fatal)
	neighborhood := network.NewNeighborhood("", 0, watch, configurationPath, logger)

	// Act
	neighborhood.Synchronize()

	// Assert
	neighborhood.Wait()
	neighbors := neighborhood.Neighbors()
	expectedNeighborsCount := len(seedsIps)
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
