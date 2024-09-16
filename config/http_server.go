package config

import (
	"context"
	"net/http"
	"time"
)

type HTTPServer struct {
	address string
	handler http.Handler
	option  *httpServerOptionValue
}

type httpServerOptionValue struct {
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

func NewHTTPServer(address string, handler http.Handler, opts ...HttpServerOptionFn) (*HTTPServer, error) {
	var option httpServerOptionValue

	for _, opFn := range opts {
		if err := opFn(&option); err != nil {
			return nil, err
		}
	}

	return &HTTPServer{address, handler, &option}, nil
}

func (h *HTTPServer) Create(_ context.Context) (*http.Server, error) {
	return &http.Server{
		Addr:              h.address,
		Handler:           h.handler,
		IdleTimeout:       h.option.IdleTimeout,
		ReadHeaderTimeout: h.option.ReadHeaderTimeout,
		ReadTimeout:       h.option.ReadTimeout,
		WriteTimeout:      h.option.WriteTimeout,
	}, nil
}

// # HTTP Server OPTIONS

type HttpServerOptionFn func(in *httpServerOptionValue) error

type HttpServerOption struct{}

func NewHttpServerOption() HttpServerOption {
	return HttpServerOption{}
}

func (HttpServerOption) WithIdleTimeout(value time.Duration) HttpServerOptionFn {
	return func(in *httpServerOptionValue) error {
		in.IdleTimeout = value
		return nil
	}
}

func (HttpServerOption) WithReadHeaderTimeout(value time.Duration) HttpServerOptionFn {
	return func(in *httpServerOptionValue) error {
		in.ReadHeaderTimeout = value
		return nil
	}
}

func (HttpServerOption) WithReadTimeout(value time.Duration) HttpServerOptionFn {
	return func(in *httpServerOptionValue) error {
		in.ReadTimeout = value
		return nil
	}
}

func (HttpServerOption) WithWriteTimeout(value time.Duration) HttpServerOptionFn {
	return func(in *httpServerOptionValue) error {
		in.WriteTimeout = value
		return nil
	}
}
