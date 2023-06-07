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

func gzipReadWriterHandler(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wr := w
			// TODO если успею, то надо добавить провеку размера ответа (пока не понял как))
			// проверяем, что клиент поддерживает gzip-сжатие
		writer:
			for _, array := range r.Header.Values("Accept-Encoding") {
				for _, value := range strings.Split(array, ", ") {
					if strings.Contains(value, "gzip") {
						// TODO подумать как использовать gzip.Reset
						gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
						if err != nil {
							logger.Error("ошибка при иницилизации gzip логгера", zap.Error(err))
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						defer gz.Close()
						w.Header().Set("Content-Encoding", "gzip")
						// переопределяем writer
						wr = gzipWriter{ResponseWriter: w, Writer: gz}
						break writer
					}
				}
			}
		reader:
			// проверяем, что клиент отправил серверу сжатые данные в формате gzip
			for _, array := range r.Header.Values("Content-Encoding") {
				for _, value := range strings.Split(array, ", ") {
					if strings.Contains(value, "gzip") {
						cr, err := gzip.NewReader(r.Body)
						if err != nil {
							logger.Error("ошибка при сжатии ответа", zap.Error(err))
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						// меняем тело запроса на новое
						r.Body = cr
						defer cr.Close()
						//избегаем повторения операции
						break reader
					}
				}
			}
			next.ServeHTTP(wr, r)
		})
	}
}
