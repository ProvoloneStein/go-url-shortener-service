package server

import (
	"net/http"
	"time"
)

func Run(addr string, handler http.Handler) error {
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return httpServer.ListenAndServe()
}
