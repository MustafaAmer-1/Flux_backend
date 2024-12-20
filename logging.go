package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type responseCapture struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseCapture) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseCapture) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log Request
		start := time.Now()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
		} else {
			log.Printf("Request: %s %s, Body: %s", r.Method, r.URL.Path, string(body))
		}

		// Restore the body for further use
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Capture Response
		responseRecorder := &responseCapture{ResponseWriter: w, body: bytes.NewBuffer(nil)}
		next.ServeHTTP(responseRecorder, r)

		// Log Response
		log.Printf(
			"Response: %s %s, Status: %d, Body: %s, Time: %v",
			r.Method, r.URL.Path, responseRecorder.statusCode, responseRecorder.body.String(), time.Since(start),
		)
	})
}
