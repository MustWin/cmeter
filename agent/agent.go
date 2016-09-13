package agent

import (
	"fmt"

	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/containers"
	containersFactory "github.com/MustWin/cmeter/containers/factory"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	logFilter "github.com/MustWin/cmeter/pipeline/filters/logger"
	"github.com/MustWin/cmeter/pipeline/messages/registercontainer"
	"github.com/MustWin/cmeter/pipeline/messages/statechange"
)

type Agent struct {
	context.Context

	config *configuration.Config

	pipeline pipeline.Pipeline

	containers containers.Driver

	tracker *containers.Tracker
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Infoln("agent running")
	defer context.GetLogger(agent).Infoln("agent shutting down")

	err := agent.InitializeContainers()
	if err != nil {
		return fmt.Errorf("error initializing container state: %v", err)
	}

	return agent.ProcessEvents()
}

func (agent *Agent) InitializeContainers() error {
	containers, err := agent.containers.GetContainers()
	if err != nil {
		return err
	}

	for _, containerInfo := range containers {
		m := registercontainer.NewMessage(containerInfo)
		if err := agent.pipeline.Send(agent, m); err != nil {
			return err
		}
	}

	return nil
}

func (agent *Agent) ProcessEvents() error {
	eventChan, err := agent.containers.WatchEvents(containers.EventContainerCreation, containers.EventContainerDeletion)
	if err != nil {
		return fmt.Errorf("error processing events: %v", err)
	}

	for event := range eventChan.GetChannel() {
		m := statechange.NewMessage(event)
		agent.pipeline.Send(agent, m)
	}

	return nil
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	filters := []pipeline.Filter{
		logFilter.New(),
	}

	containersParams := config.Containers.Parameters()
	if containersParams == nil {
		containersParams = make(configuration.Parameters)
	}

	containersDriver, err := containersFactory.Create(config.Containers.Type(), containersParams)
	if err != nil {
		return nil, err
	}

	context.GetLogger(ctx).Debugf("using %q containers driver", config.Containers.Type())
	return &Agent{
		Context:    ctx,
		config:     config,
		containers: containersDriver,
		pipeline:   pipeline.New(filters...),
	}, nil
}
