package metricsstore

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"testing"
)

func Test_GetMetric_Counter(t *testing.T) {
	tk := setUpTestKit(t)
	ctx := context.Background()

	domainMetric := models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: "counter-321",
			Val:   models.NewInt64MetricValue(321),
			Mtype: models.MetricTypeCounter,
		},
	}

	apiMetric := &desc.Metric{
		Name: "counter-321",
		Value: &desc.Metric_Value{MetricValue: &desc.Metric_Value_Counter{
			Counter: &desc.Counter{Value: 321},
		}},
	}

	expectedResp := &desc.GetMetricResponse{Metric: apiMetric}

	tk.useCasesController.EXPECT().GetMetricByKey(ctx, domainMetric.GetKey()).Return(&domainMetric, nil)

	resp, err := tk.service.GetMetric(ctx, &desc.GetMetricRequest{
		Name: "counter-321",
		Type: desc.MetricType_METRIC_TYPE_COUNTER,
	})
	require.NoError(t, err)

	if !proto.Equal(expectedResp, resp) {
		require.FailNow(t, fmt.Sprintf("got response %+v, expected response %+v", resp, expectedResp))
	}
}

func Test_GetMetric_Counter_NotFound(t *testing.T) {
	tk := setUpTestKit(t)
	ctx := context.Background()

	domainMetricKey := models.BuildMetricKey("counter-321", models.MetricTypeCounter)

	tk.useCasesController.EXPECT().GetMetricByKey(ctx, domainMetricKey).Return(nil, nil)

	_, err := tk.service.GetMetric(ctx, &desc.GetMetricRequest{
		Name: "counter-321",
		Type: desc.MetricType_METRIC_TYPE_COUNTER,
	})
	require.Error(t, err)

	st, ok := status.FromError(err)
	if !ok {
		require.FailNow(t, fmt.Sprintf("got error %+v not convertable to grpc status code", err))
	}

	require.Equal(t, codes.NotFound, st.Code())
}

func Test_GetMetric_Gauge(t *testing.T) {
	tk := setUpTestKit(t)
	ctx := context.Background()

	domainMetric := models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: "gauge-555.5",
			Val:   models.NewFloat64MetricValue(555.5),
			Mtype: models.MetricTypeGauge,
		},
	}

	apiMetric := &desc.Metric{
		Name: "gauge-555.5",
		Value: &desc.Metric_Value{MetricValue: &desc.Metric_Value_Gauge{
			Gauge: &desc.Gauge{Value: 555.5},
		}},
	}

	expectedResp := &desc.GetMetricResponse{Metric: apiMetric}

	tk.useCasesController.EXPECT().GetMetricByKey(ctx, domainMetric.GetKey()).Return(&domainMetric, nil)

	resp, err := tk.service.GetMetric(ctx, &desc.GetMetricRequest{
		Name: "gauge-555.5",
		Type: desc.MetricType_METRIC_TYPE_GAUGE,
	})
	require.NoError(t, err)

	if !proto.Equal(expectedResp, resp) {
		require.FailNow(t, fmt.Sprintf("got response %+v, expected response %+v", resp, expectedResp))
	}
}

func Test_GetMetric_Gauge_NotFound(t *testing.T) {
	tk := setUpTestKit(t)
	ctx := context.Background()

	domainMetricKey := models.BuildMetricKey("gauge-555.5", models.MetricTypeGauge)

	tk.useCasesController.EXPECT().GetMetricByKey(ctx, domainMetricKey).Return(nil, nil)

	_, err := tk.service.GetMetric(ctx, &desc.GetMetricRequest{
		Name: "gauge-555.5",
		Type: desc.MetricType_METRIC_TYPE_GAUGE,
	})
	require.Error(t, err)

	st, ok := status.FromError(err)
	if !ok {
		require.FailNow(t, fmt.Sprintf("got error %+v not convertable to grpc status code", err))
	}

	require.Equal(t, codes.NotFound, st.Code())
}
