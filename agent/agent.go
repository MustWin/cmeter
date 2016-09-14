package agent

import (
	"fmt"

	"github.com/MustWin/cmeter/collector"
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/containers"
	containersFactory "github.com/MustWin/cmeter/containers/factory"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	collectionsEmitter "github.com/MustWin/cmeter/pipeline/emitters/collections"
	eventsEmitter "github.com/MustWin/cmeter/pipeline/emitters/events"
	logFilter "github.com/MustWin/cmeter/pipeline/filters/logger"
	registryFilter "github.com/MustWin/cmeter/pipeline/filters/registry"
	resolveContainerFilter "github.com/MustWin/cmeter/pipeline/filters/resolvecontainer"
	resolveServiceFilter "github.com/MustWin/cmeter/pipeline/filters/resolveservice"
	"github.com/MustWin/cmeter/pipeline/messages/containerdiscovery"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
)

type Agent struct {
	context.Context

	config *configuration.Config

	collector collector.Collector

	pipeline pipeline.Pipeline

	containers containers.Driver

	registry *containers.Registry

	emitters []pipeline.Emitter
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Info("starting agent")
	defer context.GetLogger(agent).Info("shutting down agent")
	err := agent.InitializeContainers()
	if err != nil {
		return fmt.Errorf("error initializing container states: %v", err)
	}

	return agent.ProcessEvents()
}

func (agent *Agent) InitializeContainers() error {
	containers, err := agent.containers.GetContainers(agent)
	if err != nil {
		return err
	}

	context.GetLogger(agent).Infof("found %d active containers", len(containers))
	for _, containerInfo := range containers {
		m := containerdiscovery.NewMessage(containerInfo)
		if err := agent.pipeline.Send(agent, m); err != nil {
			return err
		}
	}

	return nil
}

func (agent *Agent) ProcessEvents() error {
	eventChan, err := agent.containers.WatchEvents(agent, containers.EventContainerCreation, containers.EventContainerDeletion)
	if err != nil {
		return fmt.Errorf("error opening event channel: %v", err)
	}

	context.GetLogger(agent).Info("event monitor started")
	defer context.GetLogger(agent).Info("event monitor stopped")
	for event := range eventChan.GetChannel() {
		m := statechange.NewMessage(event)
		agent.pipeline.Send(agent, m)
	}

	return nil
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	context.GetLogger(ctx).Info("initializing agent")

	registry := containers.NewRegistry()
	collector := collector.New(agent)

	containersParams := config.Containers.Parameters()
	if containersParams == nil {
		containersParams = make(configuration.Parameters)
	}

	containersDriver, err := containersFactory.Create(config.Containers.Type(), containersParams)
	if err != nil {
		return nil, err
	}

	context.GetLogger(ctx).Infof("using %q containers driver", config.Containers.Type())

	filters := []pipeline.Filter{
		logFilter.New(),
		resolveContainerFilter.New(containersDriver),
		registryFilter.New(registry, config.Tracking.TrackingLabel),
		resolveServiceFilter.New(registry, config.Tracking.ServiceKeyLabel),
	}

	pipeline := pipeline.New(filters...)
	emitters := []pipeline.Emitter{
		collectionsEmitter.New(collector, pipeline),
		eventsEmitter.New(containers, pipeline),
	}

	return &Agent{
		Context:    ctx,
		config:     config,
		containers: containersDriver,
		collector:  collector,
		pipeline:   pipeline,
		registry:   registry,
		emitters:   emitters,
	}, nil
}
