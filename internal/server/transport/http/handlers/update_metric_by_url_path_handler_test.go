package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers/_stub"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func Test_getUpdateMetricByURLPathHandler_TableTest(t *testing.T) {
	t.Parallel()

	type prereq struct {
		httpMethod  string
		urlEndpoint string
		chiParams   map[string]string
	}

	type want struct {
		httpCode int
	}

	testCases := []struct {
		desc   string
		prereq prereq
		want   want
	}{
		{
			desc: "negative: invalid method: get",
			prereq: prereq{
				httpMethod:  http.MethodGet,
				urlEndpoint: "/update/counter/name/123",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: put",
			prereq: prereq{
				httpMethod:  http.MethodPut,
				urlEndpoint: "/update/counter/name/123",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: patch",
			prereq: prereq{
				httpMethod:  http.MethodPatch,
				urlEndpoint: "/update/counter/name/123",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid url: value absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update/counter/testName",
			},
			want: want{
				httpCode: http.StatusNotFound,
			},
		},
		{
			desc: "negative: invalid url: name absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update/counter",
			},
			want: want{
				httpCode: http.StatusNotFound,
			},
		},
		{
			desc: "negative: invalid url: type absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
			},
			want: want{
				httpCode: http.StatusNotFound,
			},
		},
		{
			desc: "negative: invalid url: type unknown",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update/unknown/testName/123",
				chiParams: map[string]string{
					urlParamMetricType:  "unknown",
					urlParamMetricName:  "testName",
					urlParamMetricValue: "123",
				},
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid url: value not castable - float counter",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update/counter/testName/123.666",
				chiParams: map[string]string{
					urlParamMetricType:  "counter",
					urlParamMetricName:  "testName",
					urlParamMetricValue: "123.666",
				},
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid url: value not castable - string counter",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update/counter/testName/lalala",
				chiParams: map[string]string{
					urlParamMetricType:  "counter",
					urlParamMetricName:  "testName",
					urlParamMetricValue: "lalala",
				},
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid url: value not castable - string gauge",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update/gauge/testName/lalala",
				chiParams: map[string]string{
					urlParamMetricType:  "gauge",
					urlParamMetricName:  "testName",
					urlParamMetricValue: "lalala",
				},
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "positive: ok counter",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update/counter/testName/123",
				chiParams: map[string]string{
					urlParamMetricType:  "counter",
					urlParamMetricName:  "testName",
					urlParamMetricValue: "123",
				},
			},
			want: want{
				httpCode: http.StatusOK,
			},
		},
		{
			desc: "positive: ok gauge",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update/gauge/testName/123.5",
				chiParams: map[string]string{
					urlParamMetricType:  "gauge",
					urlParamMetricName:  "testName",
					urlParamMetricValue: "123.5",
				},
			},
			want: want{
				httpCode: http.StatusOK,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			req := httptest.NewRequest(test.prereq.httpMethod, test.prereq.urlEndpoint, nil)

			rctx := chi.NewRouteContext()
			for k, v := range test.prereq.chiParams {
				rctx.URLParams.Add(k, v)
			}

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			router := NewRouter(&_stub.ControllerNoErrStub{})
			router.updateMetricByURLPathHandler(w, req)

			resp := w.Result()
			_ = resp.Body.Close()

			assert.Equal(t, test.want.httpCode, resp.StatusCode)
		})
	}
}
