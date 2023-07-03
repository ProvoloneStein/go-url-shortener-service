package server

import (
	"fmt"
	"net/http"
)

const (
	serverReadTimeout  = 10
	serverWriteTimeout = 10
	maxHeaderBytes     = 1 << 20 // 1 MB
)

func Run(addr string, handler http.Handler) error {
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    serverReadTimeout,
		WriteTimeout:   serverWriteTimeout,
	}
	if err := httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("server run:  %w", err)
	}
	return nil
}
