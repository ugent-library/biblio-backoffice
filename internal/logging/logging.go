package logging

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"
)

type MiddleWareLoggingHandler struct {
	logger  *zap.SugaredLogger
	handler http.Handler
}

func (h MiddleWareLoggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.MultipartForm != nil {
		req.MultipartForm.RemoveAll()
	}

	m := httpsnoop.CaptureMetrics(h.handler, w, req)

	h.logger.Infof("%s %s (code=%d dt=%s written=%d)",
		req.Method,
		req.URL,
		m.Code,
		m.Duration,
		m.Written,
	)
}

func HTTPHandler(logger *zap.SugaredLogger, h http.Handler) http.Handler {
	return MiddleWareLoggingHandler{logger, h}
}
