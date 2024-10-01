package handlers

import (
	"encoding/json"
	httpmodels "github.com/aridae/go-metrics-store/internal/server/transport/http/models"
	"net/http"
)

func (rt *Router) updateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	transportMetric := httpmodels.Metric{}
	err := json.NewDecoder(r.Body).Decode(&transportMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricFactory, err := resolveMetricFactoryForMetricType(transportMetric.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := buildMetricDomainModel(transportMetric, metricFactory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricUpsertStrategy := metricFactory.ProvideUpsertStrategy()

	newMetricState, err := rt.useCasesController.UpsertScalarMetric(ctx, metric, metricUpsertStrategy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transportMetricAfterUpsert, err := buildMetricTransportModel(newMetricState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(transportMetricAfterUpsert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}