package metricsstore

import (
	"context"
	"fmt"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
	"github.com/aridae/go-metrics-store/pkg/slice"
)

func (i *Implementation) UpsertMetricsBatch(ctx context.Context, req *desc.UpsertMetricsBatchRequest) (*desc.UpsertMetricsBatchResponse, error) {
	domainUpserts, err := slice.MapBatch(req.GetMetrics(), mapAPIToDomainMetricUpsert)
	if err != nil {
		return nil, fmt.Errorf("batch mapAPIToDomainMetricUpsert: %w", err)
	}

	upsertedMetrics, err := i.useCasesController.UpsertMetricsBatch(ctx, domainUpserts)
	if err != nil {
		return nil, fmt.Errorf("useCasesController.UpsertMetricsBatch: %w", err)
	}

	apiMetrics, err := slice.MapBatch(upsertedMetrics, mapDomainToAPIMetric)
	if err != nil {
		return nil, fmt.Errorf("batch mapDomainToAPIMetric: %w", err)
	}

	return &desc.UpsertMetricsBatchResponse{UpsertedMetrics: apiMetrics}, nil
}
