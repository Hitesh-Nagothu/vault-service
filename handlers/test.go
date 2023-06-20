package handlers

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Test struct {
	logger *zap.Logger
}

func NewTest(l *zap.Logger) *Test {
	return &Test{l}
}

func (t *Test) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.logger.Log(zapcore.InfoLevel, "hello from test world")
	fmt.Fprintln(w, "hello from test world")
	w.WriteHeader(http.StatusNotImplemented)
}
