package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (g gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

// Gzip decompresses request and compresses response if necessary.
func Gzip(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.Reader

		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
		}

		body, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
