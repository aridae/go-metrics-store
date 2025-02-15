package metricpgrepo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *repo) GetByKey(ctx context.Context, key models.MetricKey) (*models.Metric, error) {
	queryable := r.txGetter.DefaultTrOrDB(ctx, r.db)

	qb := baseSelectQuery.Where(squirrel.Eq{keyColumn: key.String()})

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := queryable.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var dtos []metricDTO
	err = pgxscan.ScanAll(&dtos, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row into metricDTO: %w", err)
	}

	if len(dtos) == 0 {
		return nil, nil
	}

	d := dtos[0]
	metric, err := parseDTO(d)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metricDTO into metric business model: %w", err)
	}

	return &metric, nil
}

func (r *repo) GetAll(ctx context.Context) ([]models.Metric, error) {
	queryable := r.txGetter.DefaultTrOrDB(ctx, r.db)

	qb := baseSelectQuery.OrderBy(keyColumn)

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := queryable.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var dtos []metricDTO
	err = pgxscan.ScanAll(&dtos, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into dtos: %w", err)
	}

	if len(dtos) == 0 {
		return nil, nil
	}

	metrics, err := parseDTOs(dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dtos into metric business models: %w", err)
	}

	return metrics, nil
}
