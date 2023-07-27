package metrics

import (
	_ "embed"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	state *State
)

//go:embed state.yaml
var stateDefault string

type StateMetric struct {
	Key         string  `yaml:"key"`
	Description string  `yaml:"description"`
	Value       float64 `yaml:"value"`
}

type State struct {
	lock    *sync.Mutex
	Metrics []StateMetric `yaml:"metrics"`
}

func (s *State) Append(k, d string, v float64) {
	if s.lock == nil {
		s.lock = &sync.Mutex{}
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Metrics = append(s.Metrics, StateMetric{
		Key:         k,
		Description: d,
		Value:       v,
	})
}

func (s *State) SetValue(k string, v float64) bool {
	if s.lock == nil {
		s.lock = &sync.Mutex{}
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	for i, mtr := range s.Metrics {
		if mtr.Key == k {
			s.Metrics[i].Value = v
			return true
		}
	}
	return false
}

func loadState(file string) error {
	if !fileExists(file) {
		err := os.WriteFile(file, []byte(stateDefault), 0644)
		if err != nil {
			return fmt.Errorf("state file does not exist and could not be created")
		}
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	c := &State{}
	yaml.Unmarshal(data, c)
	state = c
	return nil
}

func saveState(file string) error {
	data, err := yaml.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}
