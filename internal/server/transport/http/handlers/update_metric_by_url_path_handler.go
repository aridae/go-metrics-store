package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (rt *Router) updateMetricByURLPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	if paramsCount := strings.Split(strings.Trim(r.URL.Path, "/"), "/"); len(paramsCount) != 4 {
		http.Error(w, "Unknown URL path.", http.StatusNotFound)
		return
	}

	ctx := r.Context()

	metricTypeFromURL := chi.URLParam(r, urlParamMetricType)
	metricNameFromURL := chi.URLParam(r, urlParamMetricName)
	metricValueFromURL := chi.URLParam(r, urlParamMetricValue)

	metricFactory, err := resolveMetricFactoryForMetricType(metricTypeFromURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricValue, err := metricFactory.ParseMetricValue(metricValueFromURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricToRegister := metricFactory.CreateMetricUpsert(metricNameFromURL, metricValue)

	_, err = rt.useCasesController.UpsertMetric(ctx, metricToRegister)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
