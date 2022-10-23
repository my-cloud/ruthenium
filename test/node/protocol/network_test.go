package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/clock"
	"io/ioutil"
	"os"
	"testing"
)

func Test_SynchronizeNeighbors_OneNeighbor_NeighborAdded(t *testing.T) {
	// Arrange
	configurationPath := "../../../config"
	jsonFile, _ := os.Open(configurationPath + "/seeds.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	_ = jsonFile.Close()
	var seedsIps []string
	_ = json.Unmarshal(byteValue, &seedsIps)
	watch := clock.NewWatch()
	logger := log.NewLogger(log.Fatal)
	network := protocol.NewNetwork("", 0, watch, configurationPath, logger)

	// Act
	network.SynchronizeNeighbors()

	// Assert
	network.Wait()
	neighbors := network.Neighbors()
	expectedNeighborsCount := len(seedsIps)
	test.Assert(t, len(neighbors) == expectedNeighborsCount, fmt.Sprintf("Wrong neighbors count. Expected: %d - Actual: %d", expectedNeighborsCount, len(neighbors)))
}
