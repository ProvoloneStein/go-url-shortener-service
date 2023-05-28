package handlers

import (
	"compress/gzip"
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

func gzipWriterHandler(next http.Handler) http.Handler {
	// сжимаем gzip
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO если успею, то надо добавить провеку размера ответа (пока не понял как))
		// проверяем, что входящие данные поддерживаемого формата
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "text/plain") && !strings.HasPrefix(ct, "application/json") {
			next.ServeHTTP(w, r)
			return
		}
		// проверяем, что клиент поддерживает gzip-сжатие
		for _, array := range r.Header.Values("Accept-Encoding") {
			for _, value := range strings.Split(array, ", ") {
				if strings.Contains(value, "gzip") {
					// TODO подумать как использовать gzip.Reset
					gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
					if err != nil {
						io.WriteString(w, err.Error())
						return
					}
					defer gz.Close()
					w.Header().Set("Content-Encoding", "gzip")
					// передаём обработчику страницы переменную типа gzipWriter для вывода данных
					next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func gzipReaderHandler(next http.Handler) http.Handler {
	// распаковываем gzip тело
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	head:
		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		for _, array := range r.Header.Values("Content-Encoding") {
			for _, value := range strings.Split(array, ", ") {
				if strings.Contains(value, "gzip") {
					cr, err := gzip.NewReader(r.Body)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					// меняем тело запроса на новое
					r.Body = cr
					defer cr.Close()
					//избегаем повторения операции
					break head
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
