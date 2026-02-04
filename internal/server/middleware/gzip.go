package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"
)

var writerPool = sync.Pool{
	New: func() any {
		w, err := gzip.NewWriterLevel(nil, gzip.BestSpeed)
		if err != nil {
			panic(err)
		}
		return w
	},
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

// Gzip decompresses request and compresses response if necessary.
func Gzip() middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Encoding") == "gzip" {
				gzReader, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer gzReader.Close()
				r.Body = gzReader
			}

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				h.ServeHTTP(w, r)
				return
			}

			gzWriter := writerPool.Get().(*gzip.Writer)
			gzWriter.Reset(w)
			defer func() {
				gzWriter.Close()
				writerPool.Put(gzWriter)
			}()

			w.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gzWriter}, r)
		})
	}
}
