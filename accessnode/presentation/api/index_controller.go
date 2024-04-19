package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type IndexController struct {
	templatePath string
	logger       log.Logger
}

func NewIndexController(templatePath string, logger log.Logger) *IndexController {
	return &IndexController{templatePath, logger}
}

func (controller *IndexController) GetIndex(writer http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles(controller.templatePath)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to parse the template: %w", err).Error())
		return
	}
	if err = t.Execute(writer, ""); err != nil {
		controller.logger.Error(fmt.Errorf("failed to execute the template: %w", err).Error())
	}
}
