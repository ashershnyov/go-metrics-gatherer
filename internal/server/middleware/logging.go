package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseData struct {
	status int
	size   int
	resp   []byte
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	r.responseData.resp = b
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Logging is a middleware to log uri, HTTP method, handling duration, reponse status and response size.
func Logging(logger *zap.SugaredLogger, h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		lrw := &loggingResponseWriter{
			w,
			&responseData{},
		}
		h.ServeHTTP(lrw, r)

		dur := time.Since(start)

		logger.Infoln(
			"uri", uri,
			"method", method,
			"duration", dur,
			"status", lrw.responseData.status,
			"response_size", lrw.responseData.size,
			"response_body", string(lrw.responseData.resp),
		)
	})
}
