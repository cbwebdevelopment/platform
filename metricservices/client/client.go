package client

import (
	"github.com/ant0ine/go-json-rest/rest"

	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/userservices/client"
)

type Context interface {
	Logger() log.Logger
	Request() *rest.Request
	AuthenticationDetails() client.AuthenticationDetails
}

type Client interface {
	RecordMetric(context Context, name string, data ...map[string]string) error
}
