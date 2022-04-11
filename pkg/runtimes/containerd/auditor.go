package containerd

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/equinix-ms/scra/pkg/audit"

	"go.uber.org/zap"
)

type Auditor struct {
	containerdClient *Client
	rootPrefix       string
	logger           *zap.Logger
}

func NewAuditor(address string, prefix string) (*Auditor, error) {
	client, err := NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("error setting up zap logging: %v", err)
	}

	a := &Auditor{
		containerdClient: client,
		rootPrefix:       prefix,
		logger:           logger,
	}

	return a, nil
}

func (a *Auditor) AuditOnce() error {
	nss, err := a.containerdClient.ListNamespaces()
	if err != nil {
		return fmt.Errorf("error listing namespaces: %v", err)
	}

	for _, ns := range nss {
		containers, err := a.containerdClient.ListContainers(ns)
		if err != nil {
			return fmt.Errorf("error listing containers: %v", err)
		}

		for _, container := range containers {
			err = a.auditContainer(ns, container)
			if err != nil {
				return fmt.Errorf("error auditing container %s: %v", container.ID(), err)
			}
		}
	}

	return nil
}

func (a *Auditor) auditContainer(namespace string, container containerd.Container) error {
	ctx := namespaces.WithNamespace(context.Background(), namespace)
	spec, err := container.Spec(ctx)
	if err != nil {
		return fmt.Errorf("error getting spec: %v", err)
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		task = nil // non running tasks error, apparantly. Flag with nil task
	}

	containerStatus := containerd.Unknown
	if task != nil {
		status, err := task.Status(ctx)
		if err != nil {
			return fmt.Errorf("error getting status: %v", err)
		}

		containerStatus = status.Status
	}

	image, err := container.Image(ctx)
	if err != nil {
		return fmt.Errorf("error getting image information: %v", err)
	}

	info, err := container.Info(ctx, containerd.WithoutRefreshedMetadata)
	if err != nil {
		return fmt.Errorf("error getting container info: %v", err)
	}

	mounts := make([]string, len(spec.Mounts))
	for i, mount := range spec.Mounts {
		mounts[i] = mount.Source
	}

	namespaces, err := getNamespaceInfo(spec, a.rootPrefix)
	if err != nil {
		return fmt.Errorf("error getting namespace info: %v", err)
	}

	networks, err := getNetworkInfo(spec)
	if err != nil {
		return fmt.Errorf("error getting network info: %v", err)
	}

	devices := make([]*audit.LinuxDevice, len(spec.Linux.Devices))
	for i, dev := range spec.Linux.Devices {
		devices[i] = &audit.LinuxDevice{&dev}
	}

	r := audit.Report{
		Runtime:      info.Runtime.Name,
		ID:           container.ID(),
		Image:        image.Name(),
		PID:          -1, // we will fill it in later, if we can find it
		Namespaces:   *namespaces,
		Networks:     networks,
		Created:      info.CreatedAt,
		Mounts:       mounts,
		CgroupsPath:  spec.Linux.CgroupsPath,
		Status:       string(containerStatus),
		Capabilities: &audit.LinuxCapabilities{spec.Process.Capabilities},
		Devices:      devices,
	}

	if task != nil {
		r.PID = int(task.Pid())
	}

	a.logger.Info("container audit report", zap.Object("report", &r))

	return nil
}
