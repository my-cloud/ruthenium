package p2p

import (
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type NeighborLogger struct {
	logger log.Logger
}

func NewNeighborLogger(logger log.Logger) *NeighborLogger {
	return &NeighborLogger{logger}
}

func (neighborLogger *NeighborLogger) Info(msg string) {
	neighborLogger.logger.Trace(msg)
}

func (neighborLogger *NeighborLogger) Warn(msg string) {
	neighborLogger.logger.Trace(msg)
}

func (neighborLogger *NeighborLogger) Error(msg string) {
	neighborLogger.logger.Trace(msg)
}
