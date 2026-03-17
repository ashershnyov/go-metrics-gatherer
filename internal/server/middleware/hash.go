package middleware

import (
	"bytes"
	"io"
	"net/http"

	hg "github.com/ashershnyov/go-metrics-gatherer/pkg/hasher"
)

type hasher interface {
	Hash([]byte) string
}

const hashHeaderKey = "HashSHA256"

// Hashing computes sha256 hash of the incoming data and compares to the one in the request headers.
func Hashing(key string) middleware {
	hasher := hg.NewHasher(key)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedHash := r.Header.Get("HashSHA256")
			if expectedHash == "" || key == "" {
				next.ServeHTTP(w, r)
				return
			}

			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			actualHash := hasher.Hash(buf.Bytes())
			r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
			if actualHash != expectedHash {
				http.Error(w, "request data is corrupt as hashes don't match"+"actual:"+actualHash+"excpected:"+expectedHash, http.StatusBadRequest)
				return
			}

			w.Header().Set(hashHeaderKey, actualHash)
			next.ServeHTTP(w, r)
		})
	}
}
