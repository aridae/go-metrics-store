package handlers

import (
	"encoding/json"
	"net/http"

	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
)

func (rt *Router) updateMetricsBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var transportMetrics []httpmodels.Metric
	err := json.NewDecoder(r.Body).Decode(&transportMetrics)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusBadRequest)
		return
	}

	updatedMetrics := make([]httpmodels.Metric, 0, len(transportMetrics))
	for _, transportMetric := range transportMetrics {
		metricFactory, err := resolveMetricFactoryForMetricType(transportMetric.MType)
		if err != nil {
			mustWriteJSONError(w, err, http.StatusBadRequest)
			return
		}

		metric, err := buildMetricDomainModel(transportMetric, metricFactory)
		if err != nil {
			mustWriteJSONError(w, err, http.StatusBadRequest)
			return
		}

		metricUpsertStrategy := metricFactory.ProvideUpsertStrategy()

		newMetricState, err := rt.useCasesController.UpsertScalarMetric(ctx, metric, metricUpsertStrategy)
		if err != nil {
			mustWriteJSONError(w, err, http.StatusInternalServerError)
			return
		}

		transportMetricAfterUpsert, err := buildMetricTransportModel(newMetricState)
		if err != nil {
			mustWriteJSONError(w, err, http.StatusInternalServerError)
			return
		}

		updatedMetrics = append(updatedMetrics, transportMetricAfterUpsert)
	}

	err = json.NewEncoder(w).Encode(updatedMetrics)
	if err != nil {
		mustWriteJSONError(w, err, http.StatusInternalServerError)
		return
	}
}