package metrics

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Client struct {
	addr            string
	apiKey          string
	allowSelfSigned bool
}

func (c *Client) Create(key, description string) error {
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.Set("x-api-key", c.apiKey)
	req.SetRequestURI(fmt.Sprintf("%s/%s", c.addr, key))
	req.SetBodyString(description)

	if err := a.Parse(); err != nil {
		return err
	}

	if c.allowSelfSigned {
		a = a.InsecureSkipVerify()
	}
	code, _, errs := a.Bytes()

	if code != fiber.StatusCreated && code != fiber.StatusOK {
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
		return errors.New("failed to create metric")
	}

	return nil
}

func (c *Client) Update(key string, value interface{}) error {
	v, ok := interfaceToFloat64(value)
	if !ok {
		return errors.New("could not parse value")
	}
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodPut)
	req.Header.Set("x-api-key", c.apiKey)
	req.SetRequestURI(fmt.Sprintf("%s/%s", c.addr, key))
	req.SetBodyString(fmt.Sprint(v))

	if err := a.Parse(); err != nil {
		return err
	}

	if c.allowSelfSigned {
		a = a.InsecureSkipVerify()
	}
	code, _, errs := a.Bytes()

	if code != fiber.StatusNoContent {
		if len(errs) > 0 {
			return errors.Join(errs...)
		}
		return errors.New("failed to update metric")
	}

	return nil
}

func (c *Client) CreateUpdate(key, description string, value interface{}) error {
	_ = c.Create(key, description)
	return c.Update(key, value)
}

func (c *Client) Read(key string) (float64, error) {
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Set("x-api-key", c.apiKey)
	req.SetRequestURI(fmt.Sprintf("%s/%s", c.addr, key))

	if err := a.Parse(); err != nil {
		return 0, err
	}

	if c.allowSelfSigned {
		a = a.InsecureSkipVerify()
	}
	code, body, errs := a.Bytes()

	if code != fiber.StatusOK {
		if len(errs) > 0 {
			return 0, errors.Join(errs...)
		}

		return 0, errors.New("failed to read metric")
	}

	if v, ok := interfaceToFloat64(body); ok {
		return v, nil
	}

	return 0, errors.New("failed to parse read metric")
}

func (c *Client) Delete(key string) error {
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodDelete)
	req.Header.Set("x-api-key", c.apiKey)
	req.SetRequestURI(fmt.Sprintf("%s/%s", c.addr, key))

	if err := a.Parse(); err != nil {
		return err
	}

	if c.allowSelfSigned {
		a = a.InsecureSkipVerify()
	}
	code, _, errs := a.Bytes()

	if code != fiber.StatusNoContent {
		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return errors.New("failed to delete metric")
	}

	return nil
}

func NewClient(host string, port int, apiKey string, allowSelfSigned bool) *Client {
	c := &Client{
		addr:            fmt.Sprintf("https://%s:%d", host, port),
		apiKey:          apiKey,
		allowSelfSigned: allowSelfSigned,
	}
	return c
}
