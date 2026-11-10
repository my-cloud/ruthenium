package history

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"net/http"
	"strconv"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type BlocksController struct {
	sender   application.Sender
	settings application.ProtocolSettingsProvider
	watch    application.TimeProvider
	logger   log.Logger
}

func NewBlocksController(sender application.Sender, settings application.ProtocolSettingsProvider, watch application.TimeProvider, logger log.Logger) *BlocksController {
	return &BlocksController{sender, settings, watch, logger}
}

func (controller *BlocksController) GetBlocks(writer http.ResponseWriter, req *http.Request) {
	response := io.NewResponse(writer, controller.logger)
	requestHeight := req.URL.Query().Get("height")
	startingBlockHeight, err := strconv.ParseUint(requestHeight, 10, 64)
	if err != nil {
		errorMessage := "failed to parse height value"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	blocks, err := controller.sender.GetBlocks(startingBlockHeight)
	if err != nil {
		errorMessage := "failed to get blocks"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	response.Write(http.StatusOK, string(blocks[:]))
}

func (controller *BlocksController) GetCurrentBlockHeight(writer http.ResponseWriter, _ *http.Request) {
	response := io.NewResponse(writer, controller.logger)
	firstBlockTimestamp, err := controller.sender.GetFirstBlockTimestamp()
	if err != nil {
		errorMessage := "failed to get first block timestamp"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusInternalServerError, errorMessage)
		return
	}
	var currentBlockHeight int64
	now := controller.watch.Now().UnixNano()
	if firstBlockTimestamp != 0 && now > firstBlockTimestamp {
		currentBlockHeight = (now - firstBlockTimestamp) / controller.settings.ValidationTimestamp()
	}
	response.WriteJson(http.StatusOK, currentBlockHeight)
}
