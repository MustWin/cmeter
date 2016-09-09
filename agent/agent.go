package agent

import (
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	logFilter "github.com/MustWin/cmeter/pipeline/filters/logger"
)

type Agent struct {
	context.Context

	config *configuration.Config

	pipeline pipeline.Pipeline
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Infoln("agent running")
	defer context.GetLogger(agent).Infoln("agent shutting down")
	return nil
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	filters := []pipeline.Filter{
		logFilter.New(),
	}

	return &Agent{
		Context:  ctx,
		config:   config,
		pipeline: pipeline.New(filters...),
	}, nil
}
