package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metric struct {
	gauge       prometheus.Gauge
	key         string
	description string
	value       float64
	lock        *sync.Mutex
}

func (m *metric) set(v interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if f, ok := interfaceToFloat64(v); ok {
		m.gauge.Set(f)
		m.value = f
		for i, mtr := range state.Metrics {
			if mtr.Key == m.key {
				state.Metrics[i].Value = f
				return
			}
		}
		state.Metrics = append(state.Metrics, StateMetric{
			Key:         m.key,
			Description: m.description,
			Value:       m.value,
		})
	}
}

func (m *metric) get() float64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.value
}

func newMetric(key, description string) *metric {
	mc := &metric{
		gauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: key,
			Help: description,
		}),
		key:         key,
		description: description,
		value:       0.0,
		lock:        &sync.Mutex{},
	}
	return mc
}
