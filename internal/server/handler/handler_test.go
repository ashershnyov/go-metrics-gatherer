package handler

import (
	"io"
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
			storage := storage.NewMetricStorage()
			service := service.NewService(storage)
			h := NewMetricsHandler(service, nil)
			h.UpdateMetric().ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestGetMetricHandler(t *testing.T) {
	cases := []struct {
		method   string
		typ      string
		name     string
		wantVal  string
		wantCode int
	}{
		{
			method:   "GET",
			typ:      "gauge",
			name:     "Test1",
			wantVal:  "23.5",
			wantCode: http.StatusOK,
		},
		{
			method:   "GET",
			typ:      "counter",
			name:     "Test2",
			wantVal:  "123",
			wantCode: http.StatusOK,
		},
		{
			method:   "POST",
			typ:      "counter",
			name:     "Test2",
			wantVal:  "",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			method:   "GET",
			typ:      "counter",
			name:     "Test3",
			wantVal:  "",
			wantCode: http.StatusNotFound,
		},
		{
			method:   "GET",
			typ:      "Sumamry",
			name:     "Test3",
			wantVal:  "",
			wantCode: http.StatusBadRequest,
		},
	}

	storage := storage.NewMetricStorage()
	storage.UpdateGauge("Test1", 23.5)
	storage.UpdateCounter("Test2", 123)
	service := service.NewService(storage)
	h := NewMetricsHandler(service, nil)

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/value/"+tt.typ+"/"+tt.name, nil)
			w := httptest.NewRecorder()
			h.GetMetric().ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantCode)
			}
			if res.StatusCode == http.StatusOK {
				b, _ := io.ReadAll(res.Body)
				if string(b) != tt.wantVal {
					t.Errorf("got value %v, want %v", string(b), tt.wantVal)
				}
			}
		})
	}
}
