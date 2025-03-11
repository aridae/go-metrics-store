package metricsstore

import (
	"context"
	"fmt"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetMetric(ctx context.Context, req *desc.GetMetricRequest) (*desc.GetMetricResponse, error) {
	domainMetricKey, err := mapAPIToDomainMetricKey(req.GetName(), req.GetType())
	if err != nil {
		return nil, fmt.Errorf("mapAPIToDomainMetricKey: %w", err)
	}

	domainMetric, err := i.useCasesController.GetMetricByKey(ctx, domainMetricKey)
	if err != nil {
		return nil, fmt.Errorf("useCasesController.GetAllMetrics: %w", err)
	}
	if domainMetric == nil {
		return &desc.GetMetricResponse{}, status.Errorf(codes.NotFound, fmt.Sprintf("Metric with name %s, type %s not found", req.GetName(), req.GetType().String()))
	}

	apiMetric, err := mapDomainToAPIMetric(*domainMetric)
	if err != nil {
		return nil, fmt.Errorf("mapDomainToAPIMetric: %w", err)
	}

	return &desc.GetMetricResponse{Metric: apiMetric}, nil
}
