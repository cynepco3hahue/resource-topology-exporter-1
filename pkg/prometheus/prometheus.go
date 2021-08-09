package prometheus

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	utilwait "k8s.io/apimachinery/pkg/util/wait"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const prometheusDefaultPort = "2112"

var (
	ApiCallsFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "topology_updater_api_call_failures_total",
		Help: "The total number of api calls that failed by the updater",
	})

	OperationDelay = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "topology_updater_hw_topology_update_operation_measurement",
			Help: "Represent the latency of the update operation",
		})
)

func InitPrometheus() error {
	var port = prometheusDefaultPort
	if envValue, ok := os.LookupEnv("METRICS_PORT"); ok {
		if _, err := strconv.Atoi(envValue); err != nil {
			return fmt.Errorf("the env variable PROMETHEUS_PORT has inccorrect value %q; err %v", envValue, err)
		}
		port = envValue
	}

	http.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	go utilwait.Until(func() {
		if err :=  http.ListenAndServe(addr, nil); err != nil {
			utilruntime.HandleError(fmt.Errorf("failed to run prometheus server; %v", err))
		}
	}, 5*time.Second, utilwait.NeverStop)

	return nil
}
