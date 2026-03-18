package http

import "net/http"

// PingDB checks connection to the DB.
func (h *MetricsHandler) PingDB() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			if h.db == nil {
				http.Error(w, "Could not connect to the database", http.StatusInternalServerError)
				return
			}

			err := h.db.PingContext(r.Context())
			if err != nil {
				http.Error(w, "Could not connect to the database", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		},
	)
}
