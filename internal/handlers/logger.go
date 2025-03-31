package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jantytgat/go-kit/pkg/slogd"
)

func NewLogger(h http.Handler) *Logger {
	return &Logger{
		handler: h,
	}
}

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)

	var clientIp string
	if clientIp = r.Header.Get("X-Forwarded-For"); clientIp == "" {
		clientIp = r.RemoteAddr
	}
	slogd.Logger().LogAttrs(r.Context(), slogd.LevelInfo, "handler served", slog.String("method", r.Method), slog.String("hostname", r.Host), slog.String("uri", r.RequestURI), slog.String("client", clientIp), slog.Duration("duration", time.Since(start)))
}
