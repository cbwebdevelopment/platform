package client

import (
	"github.com/tidepool-org/platform/auth"
	"github.com/tidepool-org/platform/client"
	"github.com/tidepool-org/platform/task"
)

type Client struct {
	client *client.Client
}

func New(cfg *client.Config) (*Client, error) {
	clnt, err := client.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: clnt,
	}, nil
}

func (c *Client) GetStatus(ctx auth.Context) (*task.Status, error) {
	sts := &task.Status{}
	if err := c.client.SendRequestWithServerToken(ctx, "GET", c.client.BuildURL("status"), nil, sts); err != nil {
		return nil, err
	}

	return sts, nil
}