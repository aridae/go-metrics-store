package metricsstore

import (
	"context"
	"fmt"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
)

func (i *Implementation) UpsertMetric(ctx context.Context, req *desc.UpsertMetricRequest) (*desc.UpsertMetricResponse, error) {
	domainUpsert, err := mapAPIToDomainMetricUpsert(req.GetMetric())
	if err != nil {
		return nil, fmt.Errorf("mapAPIToDomainMetricUpsert: %w", err)
	}

	upsertedMetric, err := i.useCasesController.UpsertMetric(ctx, domainUpsert)
	if err != nil {
		return nil, fmt.Errorf("useCasesController.UpsertMetric: %w", err)
	}

	apiMetric, err := mapDomainToAPIMetric(upsertedMetric)
	if err != nil {
		return nil, fmt.Errorf("mapDomainToAPIMetric: %w", err)
	}

	return &desc.UpsertMetricResponse{UpsertedMetric: apiMetric}, nil
}
