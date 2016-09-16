package agent

import (
	"fmt"
	"sync"

	"github.com/MustWin/cmeter/api"
	"github.com/MustWin/cmeter/collector"
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/containers"
	containersFactory "github.com/MustWin/cmeter/containers/factory"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	sampleCollectionFilter "github.com/MustWin/cmeter/pipeline/filters/collector"
	logFilter "github.com/MustWin/cmeter/pipeline/filters/logger"
	registryFilter "github.com/MustWin/cmeter/pipeline/filters/registry"
	resolveContainerFilter "github.com/MustWin/cmeter/pipeline/filters/resolvecontainer"
	resolveServiceFilter "github.com/MustWin/cmeter/pipeline/filters/resolveservice"
	uplinkFilter "github.com/MustWin/cmeter/pipeline/filters/uplink"
	"github.com/MustWin/cmeter/pipeline/messages/containerdiscovery"
	"github.com/MustWin/cmeter/pipeline/messages/containersample"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
)

type Agent struct {
	context.Context

	config *configuration.Config

	collector *collector.Collector

	pipeline pipeline.Pipeline

	containers containers.Driver

	registry *containers.Registry

	api api.Client
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Info("starting agent")
	defer context.GetLogger(agent).Info("shutting down agent")
	err := agent.InitializeContainers()
	if err != nil {
		return fmt.Errorf("error initializing container states: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go agent.ProcessSamples(wg)
	go agent.ProcessEvents(wg)
	wg.Wait()
	return nil
}

func (agent *Agent) InitializeContainers() error {
	containers, err := agent.containers.GetContainers(agent)
	if err != nil {
		return err
	}

	context.GetLogger(agent).Infof("found %d active containers", len(containers))
	for _, containerInfo := range containers {
		m := containerdiscovery.NewMessage(containerInfo)
		go agent.pipeline.Send(agent, m)
	}

	return nil
}

func (agent *Agent) ProcessEvents(wg sync.WaitGroup) {
	defer wg.Done()
	eventChan, err := agent.containers.WatchEvents(agent, containers.EventContainerCreation, containers.EventContainerDeletion)
	if err != nil {
		context.GetLogger(agent).Panicf("error opening event channel: %v", err)
	}

	context.GetLogger(agent).Info("event monitor started")
	defer context.GetLogger(agent).Info("event monitor stopped")
	for event := range eventChan.GetChannel() {
		m := statechange.NewMessage(event)
		go agent.pipeline.Send(agent, m)
	}
}

func (agent *Agent) ProcessSamples(wg sync.WaitGroup) {
	defer wg.Done()

	for sample := range agent.collector.GetChannel() {
		m := containersample.NewMessage(sample)
		go agent.pipeline.Send(agent, m)
	}
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	context.GetLogger(ctx).Info("initializing agent")

	registry := containers.NewRegistry()
	collector := collector.New(config.Collector)
	apiClient := api.NewClient(config.Api)

	containersParams := config.Containers.Parameters()
	if containersParams == nil {
		containersParams = make(configuration.Parameters)
	}

	containersDriver, err := containersFactory.Create(config.Containers.Type(), containersParams)
	if err != nil {
		return nil, err
	}

	log := context.GetLogger(ctx)
	log.Infof("using %q containers driver", config.Containers.Type())
	log.Infof("tracking containers with %q label", config.Tracking.TrackingLabel)

	filters := []pipeline.Filter{
		logFilter.New(),
		resolveContainerFilter.New(containersDriver, registry),
		registryFilter.New(registry, config.Tracking.TrackingLabel),
		resolveServiceFilter.New(registry, config.Tracking.ServiceKeyLabel),
		sampleCollectionFilter.New(containersDriver, collector),
		uplinkFilter.New(apiClient),
	}

	pipeline := pipeline.New(filters...)

	return &Agent{
		Context:    ctx,
		config:     config,
		containers: containersDriver,
		collector:  collector,
		pipeline:   pipeline,
		registry:   registry,
		api:        apiClient,
	}, nil
}
