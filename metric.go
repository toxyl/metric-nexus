package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metric struct {
	gauge       prometheus.Gauge
	lock        *sync.Mutex
	key         string
	description string
	value       float64
}

func (m *metric) set(v interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if f, ok := interfaceToFloat64(v); ok {
		m.gauge.Set(f)
		m.value = f
		if !state.SetValue(m.key, f) {
			state.Append(m.key, m.description, m.value)
		}
	}
}

func (m *metric) add(v float64) float64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	v = m.value + v
	m.gauge.Set(v)
	m.value = v

	_ = state.SetValue(m.key, v)
	return v
}

func (m *metric) sub(v float64) float64 {
	return m.add(-v)
}

func (m *metric) inc() float64 {
	return m.add(1)
}

func (m *metric) dec() float64 {
	return m.sub(1)
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
		lock:        &sync.Mutex{},
		key:         key,
		description: description,
		value:       0.0,
	}
	return mc
}
