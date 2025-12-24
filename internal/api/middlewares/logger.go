package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rww := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		slog.Info("request started",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("user_agent", r.UserAgent()),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(rww, r)

		slog.Info("request completed",
			slog.Int("status", rww.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("size", rww.size),
		)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(data []byte) (int, error) {
	rw.size += int64(len(data))
	return rw.ResponseWriter.Write(data)
}
