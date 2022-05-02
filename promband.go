package promband

import (
	"context"
	"github.com/emicklei/go-restful"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

// Config the plugin configuration
type Config struct {
}

// CreateConfig creates the default plugin configuration
func CreateConfig() *Config {
	return &Config{}
}

// Metrics prometheus metrics
type Metrics struct {
	RequestTotal prometheus.Counter
}

// PromBand prometheus bandwidth plugin
type PromBand struct {
	metrics           *Metrics
	metricsHTTPServer *http.Server
	next              http.Handler
	name              string
	logger            log.Logger
}

// New create a new PromBand plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	logger := *log.New(os.Stdout, "plugin:prometheusBandwidth ", log.Ldate|log.Ltime)
	// Create web service to handle prometheus endpoint
	wsContainer := restful.NewContainer()
	wsContainer.DoNotRecover(false)

	// Create metrics for prometheus
	metrics := Metrics{
		RequestTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "promband",
			Name:      "request_total",
			Help:      "The total number of requests through promband",
		}),
	}

	wsContainer.Handle("/metrics", promhttp.Handler())
	metricsHTTPServer := &http.Server{Addr: ":9666", Handler: wsContainer}

	go func() {
		err := metricsHTTPServer.ListenAndServe()
		if err != nil {
			logger.Fatalf("promband http server failed, %v", err)
		}
	}()

	// Todo check config values here
	return &PromBand{
		metrics:           &metrics,
		metricsHTTPServer: metricsHTTPServer,
		next:              next,
		name:              name,
		logger:            logger,
	}, nil
}

func (a *PromBand) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.metrics.RequestTotal.Inc()
	a.next.ServeHTTP(rw, req)
}
