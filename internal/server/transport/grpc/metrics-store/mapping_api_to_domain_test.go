package metricsstore

import (
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_mapAPIToDomainMetricValue(t *testing.T) {
	testCases := []struct {
		desc           string
		in             *desc.Metric_Value
		outMetricValue models.MetricValue
		outMetricType  models.MetricType
		outError       error
	}{
		{
			desc:           "Positive case: counter",
			in:             &desc.Metric_Value{MetricValue: &desc.Metric_Value_Counter{Counter: &desc.Counter{Value: 111}}},
			outMetricValue: models.NewInt64MetricValue(111),
			outMetricType:  models.MetricTypeCounter,
			outError:       nil,
		},
		{
			desc:           "Positive case: gauge",
			in:             &desc.Metric_Value{MetricValue: &desc.Metric_Value_Gauge{Gauge: &desc.Gauge{Value: 222.5}}},
			outMetricValue: models.NewFloat64MetricValue(222.5),
			outMetricType:  models.MetricTypeGauge,
			outError:       nil,
		},
		{
			desc:     "Negative case: nil metric value",
			outError: fmt.Errorf("unknown metric type: %v", nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			outValue, outType, outErr := mapAPIToDomainMetricValue(tc.in)

			if tc.outError != nil {
				require.Error(t, outErr)
				require.EqualError(t, outErr, tc.outError.Error())
				return
			}

			require.Equal(t, tc.outMetricValue, outValue)
			require.Equal(t, tc.outMetricType, outType)
		})
	}
}

func Test_mapAPIToDomainMetricKey(t *testing.T) {
	testCases := []struct {
		desc         string
		inName       string
		inType       desc.MetricType
		outMetricKey models.MetricKey
		outError     error
	}{
		{
			desc:         "Positive case: counter",
			inName:       "test",
			inType:       desc.MetricType_METRIC_TYPE_COUNTER,
			outMetricKey: models.BuildMetricKey("test", models.MetricTypeCounter),
			outError:     nil,
		},
		{
			desc:         "Positive case: gauge",
			inName:       "test",
			inType:       desc.MetricType_METRIC_TYPE_GAUGE,
			outMetricKey: models.BuildMetricKey("test", models.MetricTypeGauge),
			outError:     nil,
		},
		{
			desc:     "Negative case: unspecified metric type",
			inType:   desc.MetricType_METRIC_TYPE_UNSPECIFIED,
			outError: fmt.Errorf("unknown metric type: %v", desc.MetricType_METRIC_TYPE_UNSPECIFIED),
		},
		{
			desc:     "Negative case: unknown metric type",
			inType:   desc.MetricType(666),
			outError: fmt.Errorf("unknown metric type: %v", 666),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			outKey, outErr := mapAPIToDomainMetricKey(tc.inName, tc.inType)

			if tc.outError != nil {
				require.Error(t, outErr)
				require.ErrorContains(t, outErr, tc.outError.Error())
				return
			}

			require.Equal(t, tc.outMetricKey, outKey)
		})
	}
}
