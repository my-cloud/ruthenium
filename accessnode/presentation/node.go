package presentation

import (
	"github.com/gin-gonic/gin"
	"github.com/my-cloud/ruthenium/accessnode/presentation/api"
	"github.com/my-cloud/ruthenium/accessnode/presentation/api/payment"
	"github.com/my-cloud/ruthenium/accessnode/presentation/api/wallet"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"github.com/my-cloud/ruthenium/validatornode/domain/clock"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
)

type Node struct {
	port   string
	rooter *gin.Engine
}

func NewNode(port string, sender application.Sender, settings application.ProtocolSettingsProvider, templatePath string, watch *clock.Watch, logger *console.Logger) *Node {
	rooter := gin.Default()
	indexController := api.NewIndexController(templatePath, logger)
	transactionController := payment.NewTransactionController(sender, logger)
	transactionsController := payment.NewTransactionsController(sender, logger)
	infoController := payment.NewInfoController(sender, settings, watch, logger)
	progressController := payment.NewProgressController(sender, settings, watch, logger)
	addressController := wallet.NewAddressController(logger)
	amountController := wallet.NewAmountController(sender, settings, watch, logger)
	rooter.GET("/", func(c *gin.Context) { indexController.GetIndex(c.Writer, c.Request) })
	rooter.POST("/transaction", func(c *gin.Context) { transactionController.PostTransaction(c.Writer, c.Request) })
	rooter.GET("/transactions", func(c *gin.Context) { transactionsController.GetTransactions(c.Writer, c.Request) })
	rooter.GET("/transaction/info", func(c *gin.Context) { infoController.GetTransactionInfo(c.Writer, c.Request) })
	rooter.PUT("/transaction/output/progress", func(c *gin.Context) { progressController.GetTransactionProgress(c.Writer, c.Request) })
	rooter.GET("/wallet/address", func(c *gin.Context) { addressController.GetWalletAddress(c.Writer, c.Request) })
	rooter.GET("/wallet/amount", func(c *gin.Context) { amountController.GetWalletAmount(c.Writer, c.Request) })
	return &Node{port, rooter}
}

func (node *Node) Run() error {
	return node.rooter.Run(":" + node.port)
}
