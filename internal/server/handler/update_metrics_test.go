package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
	"github.com/ashershnyov/go-metrics-gatherer/mocks"
	"github.com/golang/mock/gomock"
)

func TestUpdateMultipleJSONHandler(t *testing.T) {
	type testCase struct {
		auditMock func(*gomock.Controller) *mocks.MockauditLogger
		method    string
		body      string
		wantCode  int
	}

	cases := []testCase{
		{
			method: "POST",
			body:   `[{"ID":"testCounter","Type":"counter","Delta":12},{"ID":"testCounter2","Type":"counter","Delta":13}]`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			wantCode: http.StatusOK,
		},
		{
			method: "POST",
			body:   `[{"ID":"testCounter","Type":"counter"}]`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "POST",
			body:   `[{"ID":"testCounter","Type":"counter","Delta":12},{"ID":"testCounter2","Type":"counter","Delta":13}]`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("went wrong"))
				return mock
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			method: "POST",
			body:   `[{"ID":"testCounter","Type":"counter","Delta":12'},{"ID":"testCounter2","Type":"counter","Delta":13}]`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "POST",
			body:   `[{"ID":"testCounter","Type":"counter","Delta":12},{"ID":"testCounter2","Type":"gauge"}]`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "POST",
			body:   `[{"ID":"testCounter","Type":"counter"},{"ID":"testCounter2","Type":"counter","Delta":13}]`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "PUT",
			body:   `[{"ID":"testCounter","Type":"counter","Delta":12},{"ID":"testCounter2","Type":"counter","Delta":13}]`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest(tt.method, "/updates", bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()

			storage := storage.NewInMemory()
			service := service.NewService(storage)
			auditMock := tt.auditMock(ctrl)
			h := NewMetricsHandler(service, nil, auditMock)

			h.UpdateMultipleJSON().ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestUpdateMetricJSONHandler(t *testing.T) {
	type testCase struct {
		auditMock func(*gomock.Controller) *mocks.MockauditLogger
		method    string
		body      string
		wantCode  int
	}

	cases := []testCase{
		{
			method: "POST",
			body:   `{"ID":"testCounter","Type":"counter","Delta":12}`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			wantCode: http.StatusOK,
		},
		{
			method: "POST",
			body:   `{"ID":"testCounter","Type":"gauge","Value":12.5}`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			wantCode: http.StatusOK,
		},
		{
			method: "POST",
			body:   `{"ID":"testCounter","Type":"gauge","Value":12.5'}`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "POST",
			body:   `{"ID":"testCounter","Type":"gauge"}`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "PUT",
			body:   `{"ID":"testCounter","Type":"counter","Delta":12}`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			method: "POST",
			body:   `{"ID":"testCounter","Type":"counter","Delta":12}`,
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("went wrong"))
				return mock
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest(tt.method, "/update", bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()

			storage := storage.NewInMemory()
			service := service.NewService(storage)
			auditMock := tt.auditMock(ctrl)
			h := NewMetricsHandler(service, nil, auditMock)

			h.UpdateMetricJSON().ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantCode)
			}
		})
	}

}

func TestUpdateMetricHandler(t *testing.T) {
	type testCase struct {
		auditMock func(*gomock.Controller) *mocks.MockauditLogger
		method    string
		typ       string
		name      string
		val       string
		wantCode  int
	}

	cases := []testCase{
		{
			method: "POST",
			typ:    "gauge",
			name:   "Test1",
			val:    "23.5",
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			wantCode: http.StatusOK,
		},
		{
			method: "POST",
			typ:    "counter",
			name:   "Test2",
			val:    "25",
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			wantCode: http.StatusOK,
		},
		{
			method: "HEAD",
			typ:    "counter",
			name:   "Test3",
			val:    "25",
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			method: "POST",
			typ:    "counter",
			name:   "Test4",
			val:    "25.5",
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "POST",
			typ:    "counter",
			name:   "Test5",
			val:    "abc",
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "POST",
			typ:    "gauge",
			name:   "Test6",
			val:    "abc",
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "POST",
			typ:    "summary",
			name:   "Test7",
			val:    "abc",
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mock
			},
			wantCode: http.StatusBadRequest,
		},
		{
			method: "POST",
			typ:    "counter",
			name:   "Test2",
			val:    "25",
			auditMock: func(c *gomock.Controller) *mocks.MockauditLogger {
				mock := mocks.NewMockauditLogger(c)
				mock.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some error"))
				return mock
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest(tt.method, "/update/"+tt.typ+"/"+tt.name+"/"+tt.val, nil)
			w := httptest.NewRecorder()

			storage := storage.NewInMemory()
			service := service.NewService(storage)
			auditMock := tt.auditMock(ctrl)
			h := NewMetricsHandler(service, nil, auditMock)

			h.UpdateMetric().ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantCode {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantCode)
			}
		})
	}
}
