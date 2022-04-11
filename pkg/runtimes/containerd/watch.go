package containerd

import (
	"context"

	apievents "github.com/containerd/containerd/api/events"
	"github.com/containerd/containerd/events"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/typeurl"

	"go.uber.org/zap"
)

func (a *Auditor) Watch() error {
	ctx := context.Background()
	eventStream, errC := a.containerdClient.client.EventService().Subscribe(ctx, `topic=="/containers/create"`)
	for {
		var (
			event *events.Envelope
			err   error
		)
		select {
		case err = <-errC:
			if err != nil {
				a.logger.Warn("received error", zap.Error(err))
				continue
			}
		case event = <-eventStream:
			if event.Event == nil {
				a.logger.Debug("invalid (nil) event")
				continue
			}

			e, err := typeurl.UnmarshalAny(event.Event)
			if err != nil {
				a.logger.Warn("failed to unmarshall event", zap.Error(err))
				continue
			}

			switch t := e.(type) {
			case *apievents.ContainerCreate:
				nsCtx := namespaces.WithNamespace(ctx, event.Namespace)

				container, err := a.containerdClient.client.LoadContainer(nsCtx, t.ID)
				if err != nil {
					a.logger.Warn("error getting container details", zap.Error(err))
					continue
				}

				err = a.auditContainer(event.Namespace, container)
				if err != nil {
					a.logger.Warn("could not audit container", zap.Error(err))
				}

			default:
				a.logger.Debug("received unknown event")
			}

		}
	}
}
