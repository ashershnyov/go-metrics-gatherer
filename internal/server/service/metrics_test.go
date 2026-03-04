package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/model"
	"github.com/ashershnyov/go-metrics-gatherer/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetric(t *testing.T) {
	t.Parallel()

	type testCase struct {
		metric      model.InternalMetric
		storageMock func(*gomock.Controller) *mocks.MockMetricStorage
		name        string
		wantError   bool
	}

	cases := []testCase{
		{
			name:   "Update gauge",
			metric: model.InternalMetric{Name: "TestGauge", Value: 12.5, Delta: 0, Type: model.Gauge},
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().UpdateGauge(gomock.Any(), "TestGauge", 12.5).Return(nil)
				return mock
			},
			wantError: false,
		},
		{
			name:   "Update counter",
			metric: model.InternalMetric{Name: "TestCounter", Value: 0, Delta: 12, Type: model.Counter},
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().UpdateCounter(gomock.Any(), "TestCounter", int64(12)).Return(nil)
				return mock
			},
			wantError: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := tt.storageMock(ctrl)
			service := &Service{
				storage: mockStorage,
			}

			ctx := context.Background()

			err := service.UpdateMetric(ctx, tt.metric)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetMetric(t *testing.T) {
	t.Parallel()

	type testCase[T float64 | int64] struct {
		wantMetricValue T
		storageMock     func(*gomock.Controller) *mocks.MockMetricStorage
		name            string
		metricName      string
		wantError       bool
	}

	casesGauges := []testCase[float64]{
		{
			name:            "Get gauge OK",
			metricName:      "TestGauge",
			wantMetricValue: 12.5,
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().GetGauge(gomock.Any(), "TestGauge").Return(12.5, nil)
				return mock
			},
			wantError: false,
		},
		{
			name:            "Get gauge Error",
			metricName:      "TestGaugeErr",
			wantMetricValue: 12.5,
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().GetGauge(gomock.Any(), "TestGaugeErr").Return(0.0, errors.New("not found"))
				return mock
			},
			wantError: true,
		},
	}

	for _, tt := range casesGauges {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := tt.storageMock(ctrl)
			service := &Service{
				storage: mockStorage,
			}

			ctx := context.Background()

			val, err := service.GetMetric(ctx, tt.metricName, "gauge")
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantMetricValue, val.Value)
			}
		})
	}

	casesCounters := []testCase[int64]{
		{
			name:            "Get counter OK",
			metricName:      "TestCounter",
			wantMetricValue: 12,
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().GetCounter(gomock.Any(), "TestCounter").Return(int64(12), nil)
				return mock
			},
			wantError: false,
		},
		{
			name:            "Get counter Error",
			metricName:      "TestCounterError",
			wantMetricValue: 12,
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().GetCounter(gomock.Any(), "TestCounterError").Return(int64(0), errors.New("not found"))
				return mock
			},
			wantError: true,
		},
	}

	for _, tt := range casesCounters {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := tt.storageMock(ctrl)
			service := &Service{
				storage: mockStorage,
			}

			ctx := context.Background()

			val, err := service.GetMetric(ctx, tt.metricName, "counter")
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantMetricValue, val.Delta)
			}
		})
	}
}

func TestUpdateMultipleMetrics(t *testing.T) {
	t.Parallel()

	type testCase struct {
		storageMock func(*gomock.Controller) *mocks.MockMetricStorage
		name        string
		metrics     []model.InternalMetric
		wantError   bool
	}

	cases := []testCase{
		{
			name: "Get gauge OK",
			metrics: []model.InternalMetric{
				model.InternalMetric{Name: "TestGauge", Value: 12.5, Delta: 0, Type: model.Gauge},
				model.InternalMetric{Name: "TestCounter", Value: 0, Delta: 12, Type: model.Gauge},
			},
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().UpdateMultipleMetrics(gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			wantError: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := tt.storageMock(ctrl)
			service := &Service{
				storage: mockStorage,
			}

			ctx := context.Background()

			err := service.UpdateMultipleMetrics(ctx, tt.metrics)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestListMetrics(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		storageMock func(*gomock.Controller) *mocks.MockMetricStorage
		wantMetrics []model.InternalMetric
		wantError   bool
	}

	cases := []testCase{
		{
			name: "Get List",
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().GetCounters(gomock.Any()).Return(map[string]int64{"TestCounter": 12}, nil)
				mock.EXPECT().GetGauges(gomock.Any()).Return(map[string]float64{"TestGauge": 12.5}, nil)
				return mock
			},
			wantMetrics: []model.InternalMetric{
				model.InternalMetric{Name: "TestGauge", Value: 12.5, Delta: 0, Type: model.Gauge},
				model.InternalMetric{Name: "TestCounter", Value: 0, Delta: 12, Type: model.Counter},
			},
			wantError: false,
		},
		{
			name: "Get List Counters Error",
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().GetCounters(gomock.Any()).Return(map[string]int64{}, errors.New("went wrong"))
				mock.EXPECT().GetGauges(gomock.Any()).Return(map[string]float64{"TestGauge": 12.5}, nil)
				return mock
			},
			wantMetrics: []model.InternalMetric{
				model.InternalMetric{Name: "TestGauge", Value: 12.5, Delta: 0, Type: model.Gauge},
				model.InternalMetric{Name: "TestCounter", Value: 0, Delta: 12, Type: model.Counter},
			},
			wantError: true,
		},
		{
			name: "Get List Gauges Error",
			storageMock: func(ctrl *gomock.Controller) *mocks.MockMetricStorage {
				mock := mocks.NewMockMetricStorage(ctrl)
				mock.EXPECT().GetGauges(gomock.Any()).Return(map[string]float64{}, errors.New("went wrong"))
				return mock
			},
			wantMetrics: []model.InternalMetric{
				model.InternalMetric{Name: "TestGauge", Value: 12.5, Delta: 0, Type: model.Gauge},
				model.InternalMetric{Name: "TestCounter", Value: 0, Delta: 12, Type: model.Counter},
			},
			wantError: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := tt.storageMock(ctrl)
			service := &Service{
				storage: mockStorage,
			}

			ctx := context.Background()

			metrics, err := service.ListMetrics(ctx)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, metrics, tt.wantMetrics)
			}
		})
	}
}
