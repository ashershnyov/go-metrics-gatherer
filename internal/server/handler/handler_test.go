package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricUpdateHandler(t *testing.T) {
	cases := []struct {
		method     string
		metricType string
		metricName string
		metricVal  string
		wantCode   int
	}{
		{
			method:     "POST",
			metricType: "gauge",
			metricName: "Test1",
			metricVal:  "23.5",
			wantCode:   http.StatusOK,
		},
		{
			method:     "POST",
			metricType: "counter",
			metricName: "Test2",
			metricVal:  "25",
			wantCode:   http.StatusOK,
		},
		{
			method:     "HEAD",
			metricType: "counter",
			metricName: "Test3",
			metricVal:  "25",
			wantCode:   http.StatusMethodNotAllowed,
		},
		{
			method:     "POST",
			metricType: "counter",
			metricName: "Test4",
			metricVal:  "25.5",
			wantCode:   http.StatusBadRequest,
		},
		{
			method:     "POST",
			metricType: "counter",
			metricName: "Test5",
			metricVal:  "abc",
			wantCode:   http.StatusBadRequest,
		},
		{
			method:     "POST",
			metricType: "gauge",
			metricName: "Test6",
			metricVal:  "abc",
			wantCode:   http.StatusBadRequest,
		},
		{
			method:     "POST",
			metricType: "summary",
			metricName: "Test7",
			metricVal:  "abc",
			wantCode:   http.StatusBadRequest,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/update/"+tt.metricType+"/"+tt.metricName+"/"+tt.metricVal, nil)
			w := httptest.NewRecorder()
			h := NewMetricUpdateHandler()
			h.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantCode)
			}
		})
	}
}
