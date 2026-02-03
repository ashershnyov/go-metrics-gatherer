package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
)

func TestUpdateMetricHandler(t *testing.T) {
	cases := []struct {
		method   string
		typ      string
		name     string
		val      string
		wantCode int
	}{
		{
			method:   "POST",
			typ:      "gauge",
			name:     "Test1",
			val:      "23.5",
			wantCode: http.StatusOK,
		},
		{
			method:   "POST",
			typ:      "counter",
			name:     "Test2",
			val:      "25",
			wantCode: http.StatusOK,
		},
		{
			method:   "HEAD",
			typ:      "counter",
			name:     "Test3",
			val:      "25",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			method:   "POST",
			typ:      "counter",
			name:     "Test4",
			val:      "25.5",
			wantCode: http.StatusBadRequest,
		},
		{
			method:   "POST",
			typ:      "counter",
			name:     "Test5",
			val:      "abc",
			wantCode: http.StatusBadRequest,
		},
		{
			method:   "POST",
			typ:      "gauge",
			name:     "Test6",
			val:      "abc",
			wantCode: http.StatusBadRequest,
		},
		{
			method:   "POST",
			typ:      "summary",
			name:     "Test7",
			val:      "abc",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/update/"+tt.typ+"/"+tt.name+"/"+tt.val, nil)
			w := httptest.NewRecorder()
			storage := storage.NewInMemory()
			service := service.NewService(storage)
			h := NewMetricsHandler(service, nil, nil)
			h.UpdateMetric().ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantCode)
			}
		})
	}
}
