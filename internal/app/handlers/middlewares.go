package handlers

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type userCtxKey string

const (
	userCtx = userCtxKey("userId")
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipRead(r *http.Request) error {
	for _, array := range r.Header.Values("Content-Encoding") {
		for _, value := range strings.Split(array, ", ") {
			if strings.Contains(value, "gzip") {
				cr, err := gzip.NewReader(r.Body)
				if err != nil {
					return err
				}
				// меняем тело запроса на новое
				r.Body = cr
				defer cr.Close()
				// избегаем повторения операции
				return nil
			}
		}
	}
	return nil
}

func needGzipWriter(r *http.Request) bool {
	for _, array := range r.Header.Values("Accept-Encoding") {
		for _, value := range strings.Split(array, ", ") {
			if strings.Contains(value, "gzip") {
				return true
			}
		}
	}
	return false
}

func gzipReadWriterHandler(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// переопределяем writer, если клиента поддерживает gzip-сжатие
			wr := w
			if needGzipWriter(r) {
				gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
				if err != nil {
					logger.Error("ошибка при иницилизации gzip логгера", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
				}
				defer gz.Close()
				w.Header().Set("Content-Encoding", "gzip")
				wr = gzipWriter{ResponseWriter: w, Writer: gz}
			}
			// меняем тело запроса, если оно сжато
			if err := gzipRead(r); err != nil {
				logger.Error("ошибка при сжатии ответа", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(wr, r)
		})
	}
}

func userIdentity(services Service, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("authToken")
			if err != nil {
				val, err := services.GenerateToken(r.Context())
				if err != nil {
					logger.Error("ошибка при генерации токена", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				cookie = &http.Cookie{Name: "authToken", Value: val}
				http.SetCookie(w, cookie)
			}
			userID, err := services.ParseToken(cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctxuser := context.WithValue(r.Context(), userCtx, userID)
			next.ServeHTTP(w, r.WithContext(ctxuser))
		})
	}
}
