package metricsstore

import (
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
)

func mapDomainToAPIMetric(domainMetric models.Metric) (*desc.Metric, error) {
	apiMetricValue, err := mapDomainToAPIMetricValue(domainMetric.GetType(), domainMetric.GetValue())
	if err != nil {
		return nil, fmt.Errorf("mapDomainToAPIMetricValue: %w", err)
	}

	return &desc.Metric{
		Name:  domainMetric.GetName(),
		Value: apiMetricValue,
	}, nil
}

func mapDomainToAPIMetricValue(domainMetricType models.MetricType, domainMetricValue models.MetricValue) (*desc.Metric_Value, error) {
	switch domainMetricType {
	case models.MetricTypeCounter:
		return &desc.Metric_Value{MetricValue: &desc.Metric_Value_Counter{Counter: &desc.Counter{Value: domainMetricValue.UnsafeCastInt()}}}, nil
	case models.MetricTypeGauge:
		return &desc.Metric_Value{MetricValue: &desc.Metric_Value_Gauge{Gauge: &desc.Gauge{Value: float32(domainMetricValue.UnsafeCastFloat())}}}, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %v", domainMetricType)
	}
}
