package agent

import (
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/context"
)

type Agent struct {
	context.Context

	config *configuration.Config
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Infoln("agent running")
	defer context.GetLogger(agent).Infoln("agent shutting down")
	return nil
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	return &Agent{
		Context: ctx,
		config:  config,
	}, nil
}
