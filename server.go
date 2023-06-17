package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	errMissing = &fiber.Error{
		Code:    403000,
		Message: "Missing API key",
	}
	errInvalid = &fiber.Error{
		Code:    403001,
		Message: "Invalid API key",
	}
)

type Server struct {
	addr      string
	stateFile string
	api       *fiber.App
	apiKeys   []string
	data      map[string]*metric
	lock      *sync.Mutex
}

func (srv *Server) Create(key, description string, value interface{}) bool {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if _, ok := srv.data[key]; !ok {
		srv.data[key] = newMetric(key, description)
		srv.data[key].set(value)
		return true
	}
	return false
}

func (srv *Server) CreateUpdate(key, description string, value interface{}) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if _, ok := srv.data[key]; !ok {
		srv.data[key] = newMetric(key, description)
	}
	srv.data[key].set(value)
}

func (srv *Server) Read(key string) (float64, bool) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if mtr, ok := srv.data[key]; ok {
		return mtr.get(), true
	}
	return 0, false
}

func (srv *Server) Update(key string, value interface{}) bool {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if _, ok := srv.data[key]; !ok {
		return false
	}
	srv.data[key].set(value)
	return true
}

func (srv *Server) Delete(key string) bool {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if m, ok := srv.data[key]; ok {
		prometheus.DefaultRegisterer.Unregister(m.gauge)
		delete(srv.data, key)
		return true
	}
	return false
}

func (srv *Server) Increment(key string) bool {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if m, ok := srv.data[key]; ok {
		m.inc()
		return true
	}
	return false
}

func (srv *Server) Decrement(key string) bool {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if m, ok := srv.data[key]; ok {
		m.dec()
		return true
	}
	return false
}

func (srv *Server) Add(key string, v interface{}) bool {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if m, ok := srv.data[key]; ok {
		if f, ok := interfaceToFloat64(v); ok {
			m.add(f)
			return true
		}
	}
	return false
}

func (srv *Server) Sub(key string, v interface{}) bool {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	key = sanitizeKey(key)
	if m, ok := srv.data[key]; ok {
		if f, ok := interfaceToFloat64(v); ok {
			m.sub(f)
			return true
		}
	}
	return false
}

func (srv *Server) initMiddlewares() {
	srv.api.Use(idempotency.New())
	srv.api.Use(
		keyauth.New(keyauth.Config{
			KeyLookup: "header:x-api-key",
			Validator: func(ctx *fiber.Ctx, s string) (bool, error) {
				if s == "" {
					return false, errMissing
				}
				for _, k := range srv.apiKeys {
					if s == k {
						return true, nil
					}
				}

				return false, errInvalid
			},
		}),
	)
}

func (srv *Server) initAPI() {
	// PROMETHEUS handler
	srv.api.Get("/__metrics", func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(c.Context())
		return nil
	})

	// CREATE handler
	srv.api.Post("/:metric", func(c *fiber.Ctx) error {
		if srv.Create(c.Params("metric"), string(c.Body()), 0.0) {
			return c.SendStatus(fiber.StatusCreated)
		}
		// if we get here the metric already exists
		return c.SendStatus(fiber.StatusOK)
	})

	// READ handler
	srv.api.Get("/:metric", func(c *fiber.Ctx) error {
		if v, ok := srv.Read(c.Params("metric")); ok {
			return c.SendString(fmt.Sprint(v))
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	// UPDATE handler
	srv.api.Put("/:metric", func(c *fiber.Ctx) error {
		if srv.Update(c.Params("metric"), string(c.Body())) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	// INCREMENT handler
	srv.api.Put("/:metric/inc", func(c *fiber.Ctx) error {
		if srv.Increment(c.Params("metric")) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		// if we get here the metric already exists
		return c.SendStatus(fiber.StatusOK)
	})

	// DECREMENT handler
	srv.api.Put("/:metric/dec", func(c *fiber.Ctx) error {
		if srv.Decrement(c.Params("metric")) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		// if we get here the metric already exists
		return c.SendStatus(fiber.StatusOK)
	})

	// ADD handler
	srv.api.Put("/:metric/add", func(c *fiber.Ctx) error {
		if srv.Add(c.Params("metric"), string(c.Body())) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		// if we get here the metric already exists
		return c.SendStatus(fiber.StatusOK)
	})

	// SUB handler
	srv.api.Put("/:metric/sub", func(c *fiber.Ctx) error {
		if srv.Sub(c.Params("metric"), string(c.Body())) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		// if we get here the metric already exists
		return c.SendStatus(fiber.StatusOK)
	})

	// DELETE handler
	srv.api.Delete("/:metric", func(c *fiber.Ctx) error {
		if srv.Delete(c.Params("metric")) {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

}

func (srv *Server) AddAPIKey(key string) {
	for _, k := range srv.apiKeys {
		if k == key {
			return
		}
	}
	srv.apiKeys = append(srv.apiKeys, key)
}

func (srv *Server) Start(keyFile, certFile string) error {
	err := loadState(srv.stateFile)
	if err != nil {
		return err
	}
	for _, mtr := range state.Metrics {
		srv.Create(mtr.Key, mtr.Description, mtr.Value)
	}
	go func() {
		for {
			time.Sleep(time.Minute)
			_ = saveState(srv.stateFile)
		}
	}()

	srv.api = fiber.New()
	srv.initMiddlewares()
	srv.initAPI()

	keyFile, certFile, err = generateSelfSignedCertificate("local.nexus", "metric-nexus", keyFile, certFile)
	if err != nil && err.Error() != "files already exist" {
		return err
	}

	return srv.api.ListenTLS(srv.addr, certFile, keyFile)
}

func NewServer(host string, port int, stateFile string) *Server {
	srv := &Server{
		addr:      fmt.Sprintf("%s:%d", host, port),
		stateFile: stateFile,
		data:      map[string]*metric{},
		lock:      &sync.Mutex{},
	}
	return srv
}
