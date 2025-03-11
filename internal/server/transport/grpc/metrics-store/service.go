package metricsstore

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	desc "github.com/aridae/go-metrics-store/pkg/pb/metrics-store"
)

type useCasesController interface {
	UpsertMetricsBatch(context.Context, []models.MetricUpsert) ([]models.Metric, error)
	UpsertMetric(context.Context, models.MetricUpsert) (models.Metric, error)
	GetMetricByKey(context.Context, models.MetricKey) (*models.Metric, error)
	GetAllMetrics(context.Context) ([]models.Metric, error)
}

type Implementation struct {
	desc.UnimplementedMetricsStoreAPIServer

	useCasesController useCasesController
}

func NewAPI(useCasesController useCasesController) *Implementation {
	return &Implementation{useCasesController: useCasesController}
}
