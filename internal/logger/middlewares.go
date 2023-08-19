package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	// Берём структуру для хранения сведений об ответе.
	responseData struct {
		status int
		size   int
	}
	// Добавляем реализацию http.ResponseWriter.
	loggingResponseWriter struct {
		http.ResponseWriter // Встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// Записываем ответ, используя оригинальный http.ResponseWriter.
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return 0, fmt.Errorf("ResponseWriter: %w", err)
	}
	r.responseData.size += size // Захватываем размер
	return size, nil
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// Записываем код статуса, используя оригинальный http.ResponseWriter.
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // Захватываем код статуса
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: w, // Встраиваем оригинальный http.ResponseWriter.
				responseData:   respData,
			}
			start := time.Now()
			h.ServeHTTP(&lw, r)
			duration := time.Since(start)
			logger.Info("got incoming HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
				zap.Int("status", respData.status),
				zap.Int("size", respData.size),
			)
		})
	}
}
