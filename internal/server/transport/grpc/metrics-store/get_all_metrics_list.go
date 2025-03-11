package metricsstore

import (
	"context"
	"fmt"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
	"github.com/aridae/go-metrics-store/pkg/slice"
)

func (i *Implementation) GetAllMetricsList(ctx context.Context, _ *desc.GetAllMetricsListRequest) (*desc.GetAllMetricsListResponse, error) {
	domainMetrics, err := i.useCasesController.GetAllMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("useCasesController.GetAllMetrics: %w", err)
	}

	apiMetrics, err := slice.MapBatch(domainMetrics, mapDomainToAPIMetric)
	if err != nil {
		return nil, fmt.Errorf("batch mapDomainToAPIMetric: %w", err)
	}

	return &desc.GetAllMetricsListResponse{Metrics: apiMetrics}, nil
}
