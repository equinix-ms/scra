package containerd

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

type Client struct {
	client  *containerd.Client
	context context.Context
}

func NewClient(address string, ctx context.Context) (*Client, error) {
	client, err := containerd.New(address)
	if err != nil {
		return nil, fmt.Errorf("error connecting to containerd: %v", err)
	}

	return &Client{client: client, context: ctx}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) ListNamespaces() ([]string, error) {
	namespaces := c.client.NamespaceService()
	nss, err := namespaces.List(c.context)
	if err != nil {
		return nil, fmt.Errorf("error listing namespaces: %v", err)
	}

	return nss, nil
}

func (c *Client) ListContainers(namespace string) ([]containerd.Container, error) {
	ctx := namespaces.WithNamespace(c.context, namespace)
	containers, err := c.client.Containers(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %v", err)
	}

	return containers, nil
}
