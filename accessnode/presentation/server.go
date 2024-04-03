package presentation

import (
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/presentation/index"
	"github.com/my-cloud/ruthenium/accessnode/presentation/transaction"
	"github.com/my-cloud/ruthenium/accessnode/presentation/transaction/info"
	"github.com/my-cloud/ruthenium/accessnode/presentation/transaction/output/progress"
	"github.com/my-cloud/ruthenium/accessnode/presentation/transactions"
	"github.com/my-cloud/ruthenium/accessnode/presentation/wallet/address"
	"github.com/my-cloud/ruthenium/accessnode/presentation/wallet/amount"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/clock"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/config"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
)

type Server struct {
	port string
}

func NewServer(port string, host *network.Neighbor, settings *config.Settings, templatePath string, watch *clock.Watch, logger *console.Logger) *Server {
	http.Handle("/", index.NewHandler(templatePath, logger))
	http.Handle("/transaction", transaction.NewHandler(host, logger))
	http.Handle("/transactions", transactions.NewHandler(host, logger))
	http.Handle("/transaction/info", info.NewHandler(host, settings, watch, logger))
	http.Handle("/transaction/output/progress", progress.NewHandler(host, settings, watch, logger))
	http.Handle("/wallet/address", address.NewHandler(logger))
	http.Handle("/wallet/amount", amount.NewHandler(host, settings, watch, logger))
	return &Server{port}
}

func (server *Server) Serve() error {
	return http.ListenAndServe("0.0.0.0:"+server.port, nil)
}
