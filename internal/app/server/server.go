package server

import (
	"net/http"
	"time"
)

const (
	serverReadTimeout  = 15 * time.Second
	serverWriteTimeout = 15 * time.Second
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
	//if err := httpServer.ListenAndServe(); err != nil {
	//	return fmt.Errorf("server run:  %w", err)
	//}
	return httpServer.ListenAndServe()
}
