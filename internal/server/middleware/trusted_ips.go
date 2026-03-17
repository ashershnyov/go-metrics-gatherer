package middleware

import (
	"net"
	"net/http"
)

// CheckIP checks whether the client's ip is within the trusted subnet. If not, will return 403.
func CheckIP(subnet *net.IPNet) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if subnet == nil {
				next.ServeHTTP(w, r)
				return
			}

			clientIP := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(clientIP)
			if ip == nil {
				http.Error(w, "error parsing client's IP", http.StatusBadRequest)
				return
			}

			if !subnet.Contains(ip) {
				http.Error(w, "ip not trusted", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
