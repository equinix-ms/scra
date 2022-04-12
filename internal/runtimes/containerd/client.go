package containerd

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

type Client struct {
	client *containerd.Client
}

func NewClient(address string) (*Client, error) {
	client, err := containerd.New(address)
	if err != nil {
		return nil, fmt.Errorf("error connecting to containerd: %v", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) ListNamespaces() ([]string, error) {
	ctx := context.Background()
	namespaces := c.client.NamespaceService()
	nss, err := namespaces.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing namespaces: %v", err)
	}

	r := make([]string, len(nss))
	for i, ns := range nss {
		r[i] = ns
	}

	return r, nil
}

func (c *Client) ListContainers(namespace string) ([]containerd.Container, error) {
	ctx := namespaces.WithNamespace(context.Background(), namespace)
	containers, err := c.client.Containers(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %v", err)
	}

	return containers, nil
}
