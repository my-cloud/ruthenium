package api

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/application"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_Serve_NeighborConnectionError_LogsWithTraceLevel(t *testing.T) {
	// Arrange
	port := findFreePort(t)
	logger := log.NewLoggerMock()
	traced := make(chan string, 1)
	logger.TraceFunc = func(msg string) {
		select {
		case traced <- msg:
		default:
		}
	}
	host, err := NewHost(new(application.BlocksManagerMock), new(application.SendersManagerMock), new(application.TransactionsManagerMock), new(application.UtxosManagerMock), port, []byte{}, time.Second, logger)
	test.Assert(t, err == nil, "failed to create host: %v", err)
	go func() { _ = host.Serve() }()
	conn := dialUntilServing(t, port)

	// Act
	_ = conn.Close()

	// Assert
	select {
	case <-traced:
	case <-time.After(5 * time.Second):
		test.Assert(t, false, "neighbor connection error is not logged with trace level whereas it should be")
	}
	test.Assert(t, len(logger.ErrorCalls()) == 0, "neighbor connection error is logged with error level whereas it should not")
}

func findFreePort(t *testing.T) string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find a free port: %v", err)
	}
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	_ = listener.Close()
	return port
}

func dialUntilServing(t *testing.T, port string) net.Conn {
	// The gp2p server exposes no readiness signal, so the connection is retried within a bounded delay
	deadline := time.Now().Add(5 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, time.Second)
		if err == nil {
			return conn
		}
		test.Assert(t, time.Now().Before(deadline), "failed to connect to host: %v", err)
		time.Sleep(10 * time.Millisecond)
	}
}
