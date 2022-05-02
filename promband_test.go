package promband_test

import (
	"context"
	promband "github.com/collierconsulting/traefik-plugin-prometheus-bandwidth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPromBand(t *testing.T) {
	cfg := promband.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := promband.New(ctx, next, cfg, "prometheus-bandwidth-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	httpClient := http.Client{
		Timeout: time.Second * 1,
	}
	metricsURI := "http://localhost:9666/metrics"
	metricsReq, err := http.NewRequest("GET", metricsURI, nil)
	metricsResp, err := httpClient.Do(metricsReq)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, metricsResp.StatusCode)
}
