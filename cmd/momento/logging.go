package main

import (
	"log"
	"net/http"
	"time"
)

func WithLogging(h http.Handler) http.Handler {
	logger := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%s %s %s panic: %s", r.RemoteAddr,
				r.RequestURI, r.Method, err)
			}
		}()

		start := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("%s %s %s %s", r.RemoteAddr, r.RequestURI,
		r.Method, duration.String())
	}

	return http.HandlerFunc(logger)
}
