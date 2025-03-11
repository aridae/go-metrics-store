package metricsstore

import (
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
)

func mapAPIToDomainMetricKey(name string, mtype desc.MetricType) (models.MetricKey, error) {
	domainMetricType, err := mapAPIToDomainMetricType(mtype)
	if err != nil {
		return "", fmt.Errorf("mapAPIToDomainMetricType: %w", err)
	}

	domainMetricKey := models.BuildMetricKey(name, domainMetricType)

	return domainMetricKey, nil
}

func mapAPIToDomainMetricType(mtype desc.MetricType) (models.MetricType, error) {
	apiToDomainMetricType := map[desc.MetricType]models.MetricType{
		desc.MetricType_METRIC_TYPE_COUNTER: models.MetricTypeCounter,
		desc.MetricType_METRIC_TYPE_GAUGE:   models.MetricTypeGauge,
	}

	domainMetricType, ok := apiToDomainMetricType[mtype]
	if !ok {
		return "", fmt.Errorf("unknown metric type: %v", mtype)
	}

	return domainMetricType, nil
}

func mapAPIToDomainMetricUpsert(upsert *desc.Metric) (models.MetricUpsert, error) {
	mvalue, mtype, err := mapAPIToDomainMetricValue(upsert.GetValue())
	if err != nil {
		return models.MetricUpsert{}, fmt.Errorf("mapAPIToDomainMetricValue: %w", err)
	}

	return models.MetricUpsert{
		Val:   mvalue,
		Mtype: mtype,
		MName: upsert.GetName(),
	}, nil
}

func mapAPIToDomainMetricValue(value *desc.Metric_Value) (models.MetricValue, models.MetricType, error) {
	switch mvalue := value.GetMetricValue().(type) {
	case *desc.Metric_Value_Counter:
		return models.NewInt64MetricValue(mvalue.Counter.GetValue()), models.MetricTypeCounter, nil
	case *desc.Metric_Value_Gauge:
		return models.NewFloat64MetricValue(float64(mvalue.Gauge.GetValue())), models.MetricTypeGauge, nil
	default:
		return nil, "", fmt.Errorf("unknown metric type: %v", mvalue)
	}
}
