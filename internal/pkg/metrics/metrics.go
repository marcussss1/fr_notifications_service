package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	CountActiveGoroutines  *prometheus.GaugeVec
	CountMessages          prometheus.Counter
	CountFailedMessages    prometheus.Counter
	CountSucceededMessages prometheus.Counter
}

func NewMetricsService() (*Metrics, error) {
	metrics := &Metrics{
		CountActiveGoroutines: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "count_active_goroutines",
			Help: "счетчик активных горутин",
		}, []string{}),
		CountMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "count_messages",
			Help: "счетчик общего количество сообщений",
		}),
		CountFailedMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "count_failed_messages",
			Help: "счетчик не доставленных сообщений",
		}),
		CountSucceededMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "count_succeeded_messages",
			Help: "счетчик доставленных сообщений",
		}),
	}

	if err := prometheus.Register(metrics.CountMessages); err != nil {
		return nil, err
	}
	if err := prometheus.Register(metrics.CountFailedMessages); err != nil {
		return nil, err
	}
	if err := prometheus.Register(metrics.CountSucceededMessages); err != nil {
		return nil, err
	}
	if err := prometheus.Register(metrics.CountActiveGoroutines); err != nil {
		return nil, err
	}
	go func() {
		for {
			numGoroutines := runtime.NumGoroutine()
			metrics.CountActiveGoroutines.WithLabelValues().Set(float64(numGoroutines))
			time.Sleep(time.Second)
		}
	}()

	return metrics, nil
}
