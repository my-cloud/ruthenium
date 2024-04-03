package index

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type Handler struct {
	templatePath string
	logger       log.Logger
}

func NewHandler(templatePath string, logger log.Logger) *Handler {
	return &Handler{templatePath, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(handler.templatePath)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to parse the template: %w", err).Error())
			return
		}
		if err = t.Execute(writer, ""); err != nil {
			handler.logger.Error(fmt.Errorf("failed to execute the template: %w", err).Error())
		}
	default:
		handler.logger.Error("invalid HTTP method")
	}
}
