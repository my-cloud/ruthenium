package presentation

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"net/http"
	"testing"
)

func Test_CalculateFee_UnknownTransactionId_ReturnsError(t *testing.T) {
	// Arrange
	// Act
	node := NewNode("", nil, nil, "", nil, nil)

	// Assert
	test.Assert(t, node != nil, "node is nil whereas it should not")
}

func Test_NewNode_RoutesAreRegistered(t *testing.T) {
	// Arrange
	expectedRoutes := []struct{ method, path string }{
		{http.MethodGet, "/"},
		{http.MethodPost, "/transactions"},
		{http.MethodGet, "/transactions"},
		{http.MethodGet, "/transactions/info"},
		{http.MethodGet, "/transactions/output/progress"},
		{http.MethodGet, "/wallet/address"},
		{http.MethodGet, "/wallet/amount"},
	}

	// Act
	node := NewNode("", nil, nil, "", nil, nil)

	// Assert
	actualRoutes := node.rooter.Routes()
	test.Assert(t, len(actualRoutes) == len(expectedRoutes), fmt.Sprintf("Wrong routes count. expected: %d actual: %d", len(expectedRoutes), len(actualRoutes)))
	for _, expectedRoute := range expectedRoutes {
		var isRouteRegistered bool
		for _, actualRoute := range actualRoutes {
			if actualRoute.Method == expectedRoute.method && actualRoute.Path == expectedRoute.path {
				isRouteRegistered = true
				break
			}
		}
		test.Assert(t, isRouteRegistered, fmt.Sprintf("Route %s %s is not registered whereas it should be.", expectedRoute.method, expectedRoute.path))
	}
}
