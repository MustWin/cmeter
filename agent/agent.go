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
)

type Agent struct {
	context.Context

	config *configuration.Config

	pipeline pipeline.Pipeline

	containers containers.ContainersDriver
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
	containers, err := agent.monitor.GetContainers()
	if err != nil {
		return err
	}

	for _, containerId := range containers {
		m := registercontainer.NewMessage(containerId)
		if err := agent.pipeline.Send(agent, m); err != nil {
			return err
		}
	}

	return nil
}

func (agent *Agent) ProcessEvents() error {
	eventChan := agent.monitor.OpenEventChannel(monitor.EVENTS_CREATION, monitor.EVENTS_DELETION)
	for event := range <-eventChan {

	}
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	filters := []pipeline.Filter{
		logFilter.New(),
	}

	monitorParams := config.Monitor.Parameters()
	if monitorParams == nil {
		monitorParams = make(configuration.Parameters)
	}

	monitor, err := monitorFactory.Create(config.Monitor.Type(), monitorParams)
	if err != nil {
		return nil, err
	}

	context.GetLogger(ctx).Debugf("using %q monitor driver", config.Monitor.Type())
	return &Agent{
		Context:  ctx,
		config:   config,
		monitor:  monitor,
		pipeline: pipeline.New(filters...),
	}, nil
}
