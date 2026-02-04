package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/audit"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
)

func ExampleMetricsHandler_UpdateMetricJSON() {
	jsonBody := `{"id": "TestCounter", "type": "counter", "delta": 12}`
	req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	storage := storage.NewInMemory()
	service := service.NewService(storage)
	auditLogger := audit.NewLogger(nil)

	handler := NewMetricsHandler(service, nil, auditLogger)

	handler.UpdateMetricJSON().ServeHTTP(rec, req)

	fmt.Println(rec.Code)

	// Output:
	// 200
}

func ExampleMetricsHandler_UpdateMultipleJSON() {
	jsonBody := `[{"id": "TestCounter1", "type": "counter", "delta": 12}, {"id": "TestCounter2", "type": "counter", "delta": 24}]`
	req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	storage := storage.NewInMemory()
	service := service.NewService(storage)
	auditLogger := audit.NewLogger(nil)

	handler := NewMetricsHandler(service, nil, auditLogger)

	handler.UpdateMultipleJSON().ServeHTTP(rec, req)

	fmt.Println(rec.Code)

	// Output:
	// 200
}

func ExampleMetricsHandler_UpdateMetric() {
	req := httptest.NewRequest(http.MethodPost, "/update/counter/Test/12", nil)

	rec := httptest.NewRecorder()

	storage := storage.NewInMemory()
	service := service.NewService(storage)
	auditLogger := audit.NewLogger(nil)

	handler := NewMetricsHandler(service, nil, auditLogger)

	handler.UpdateMetric().ServeHTTP(rec, req)

	fmt.Println(rec.Code)

	// Output:
	// 200
}
