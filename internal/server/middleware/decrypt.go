package middleware

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"
)

// Decrypt decrypts the request body using the provided private key.
func Decrypt(key *rsa.PrivateKey) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encrypted := r.Header.Get("Encryption")
			if key == nil || encrypted != "rsa" {
				next.ServeHTTP(w, r)
				return
			}

			data, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "error reading request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			decryptedMessage, err := rsa.DecryptPKCS1v15(nil, key, data)
			if err != nil {
				http.Error(w, "error decrypting request body", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(decryptedMessage))

			next.ServeHTTP(w, r)
		})
	}
}
