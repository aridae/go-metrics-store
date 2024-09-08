package usecases

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/models"
)

func (c *Controller) UpsertScalarMetric(ctx context.Context, updater models.ScalarMetricUpdater) error {
	switch updater.Type {
	case models.ScalarMetricTypeCounter:
		return c.counterUseCasesHandler.Upsert(ctx, updater)
	case models.ScalarMetricTypeGauge:
		return c.gaugeUseCasesHandler.Upsert(ctx, updater)
	default:
		return fmt.Errorf("unknown metric type: %v", updater.Type)
	}
}
