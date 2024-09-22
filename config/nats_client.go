package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/Mind2Screen-Dev-Team/go-skeleton/app/registry"
	"github.com/Mind2Screen-Dev-Team/go-skeleton/gen/appconfig"
	"github.com/Mind2Screen-Dev-Team/go-skeleton/pkg/lazy"
)

type natsClient struct{}

func NewNatsClient() *natsClient {
	return &natsClient{}
}

func (n *natsClient) Create(_ context.Context, cfg *appconfig.AppConfig) (*nats.Conn, error) {
	if !cfg.Nats.Enabled {
		return nil, errors.New("nats client message broker is disabled")
	}

	var options []nats.Option
	if cfg.Nats.Auth.Enabled {
		options = append(options, nats.UserInfo(
			cfg.Nats.Auth.Username,
			cfg.Nats.Auth.Password,
		))
	}

	return nats.Connect(
		fmt.Sprintf(
			"nats://%s:%d",
			cfg.Nats.Host,
			cfg.Nats.Port,
		),
		options...,
	)
}

func (n *natsClient) Loader(ctx context.Context, cfg *appconfig.AppConfig, app *registry.AppDependency) {
	app.NatsConn = lazy.New(func() (*nats.Conn, error) {
		return n.Create(ctx, cfg)
	})

	app.NatsJetStreamConn = lazy.New(func() (jetstream.JetStream, error) {
		return jetstream.New(app.NatsConn.Value())
	})
}
