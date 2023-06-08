package handlers

import (
	"compress/gzip"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func newGzipWriter(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, error) {
	for _, array := range r.Header.Values("Accept-Encoding") {
		for _, value := range strings.Split(array, ", ") {
			if strings.Contains(value, "gzip") {
				// TODO подумать как использовать gzip.Reset
				gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
				if err != nil {
					return w, err
				}
				defer gz.Close()
				w.Header().Set("Content-Encoding", "gzip")
				return gzipWriter{ResponseWriter: w, Writer: gz}, nil
			}
		}
	}
	return w, nil
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
				//избегаем повторения операции
				return nil
			}
		}
	}
	return nil
}

func gzipReadWriterHandler(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO если успею, то надо добавить провеку размера ответа (пока не понял как))
			// переопределяем writer, если клиента поддерживает gzip-сжатие
			wr, err := newGzipWriter(w, r)
			if err != nil {
				logger.Error("ошибка при иницилизации gzip логгера", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
			// меняем тело запроса, если оно сжато
			gzipRead(r)
			next.ServeHTTP(wr, r)
		})
	}
}
