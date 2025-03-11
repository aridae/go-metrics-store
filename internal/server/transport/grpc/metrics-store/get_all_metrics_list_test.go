package metricsstore

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"testing"
)

func Test_GetAllMetricsList(t *testing.T) {
	tk := setUpTestKit(t)
	ctx := context.Background()

	domainMetrics := []models.Metric{
		{MetricUpsert: models.MetricUpsert{
			MName: "counter-321",
			Val:   models.NewInt64MetricValue(321),
			Mtype: models.MetricTypeCounter,
		}},
		{MetricUpsert: models.MetricUpsert{
			MName: "gauge-555.5",
			Val:   models.NewFloat64MetricValue(555.5),
			Mtype: models.MetricTypeGauge,
		}},
	}

	apiMetrics := []*desc.Metric{
		{
			Name: "counter-321",
			Value: &desc.Metric_Value{MetricValue: &desc.Metric_Value_Counter{
				Counter: &desc.Counter{Value: 321},
			}},
		},
		{
			Name: "gauge-555.5",
			Value: &desc.Metric_Value{MetricValue: &desc.Metric_Value_Gauge{
				Gauge: &desc.Gauge{Value: 555.5},
			}},
		},
	}

	expectedResp := &desc.GetAllMetricsListResponse{Metrics: apiMetrics}

	tk.useCasesController.EXPECT().GetAllMetrics(ctx).Return(domainMetrics, nil)

	resp, err := tk.service.GetAllMetricsList(ctx, &desc.GetAllMetricsListRequest{})
	require.NoError(t, err)

	if !proto.Equal(expectedResp, resp) {
		require.FailNow(t, fmt.Sprintf("got response %+v, expected response %+v", resp, expectedResp))
	}
}
