package index

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"html/template"
	"net/http"
	"path"
)

type Handler struct {
	templatesPath string
	logger        log.Logger
}

func NewHandler(templatesPath string, logger log.Logger) *Handler {
	return &Handler{templatesPath, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(handler.templatesPath, "index.html"))
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
