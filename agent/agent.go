package agent

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/MustWin/cmeter/collector"
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/containers"
	containersFactory "github.com/MustWin/cmeter/containers/factory"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/reporting"
	reportingFactory "github.com/MustWin/cmeter/reporting/factory"
	"github.com/MustWin/cmeter/shared/disposer"
)

type Agent struct {
	context.Context

	config *configuration.Config

	collector *collector.Collector

	machineCollector *collector.MachineCollector

	containers containers.Driver

	registry *containers.Registry

	reporting reporting.Driver

	dispose *disposer.Disposer
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Info("starting agent")
	defer context.GetLogger(agent).Info("shutting down agent")
	err := agent.InitializeContainers()
	if err != nil {
		return fmt.Errorf("error initializing container states: %v", err)
	}

	agent.bootstrapSignalHandler()

	go agent.ProcessSamples(agent.dispose.Quitter())
	go agent.ProcessHostSamples(agent.dispose.Quitter())
	go agent.ProcessEvents(agent.dispose.Quitter())

	agent.dispose.Wait()
	agent.Shutdown()
	return nil
}

func (agent *Agent) Shutdown() error {
	for _, c := range agent.registry.List() {
		agent.ProcessStateChange(&containers.StateChange{
			Container: c,
			Source: &containers.Event{
				Container: c,
				Timestamp: time.Now().Unix(),
				Type:      containers.EventMeterShutdown,
			},
			State: containers.StateStopped,
		}, true)
	}

	agent.dispose.QuitAll()
	return nil
}

func (agent *Agent) bootstrapSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		sig := <-c
		context.GetLogger(agent).Infof("detected signal %v: shutting down", sig)
		agent.dispose.Dispose()
	}()
}

func (agent *Agent) InitializeContainers() error {
	feed, err := agent.containers.GetMachineUsage(agent)
	if err != nil {
		return err
	}

	agent.machineCollector = collector.NewMachine(agent, feed, agent.config.Collector)
	active, err := agent.containers.GetContainers(agent)
	if err != nil {
		return err
	}

	context.GetLogger(agent).Infof("found %d active containers", len(active))
	for _, containerInfo := range active {
		c := &containers.StateChange{
			State: containers.StateRunning,
			Source: &containers.Event{
				Type:      containers.EventContainerExisted,
				Timestamp: time.Now().Unix(),
			},
			Container: containerInfo,
		}

		go agent.ProcessStateChange(c, false)
	}

	return nil
}

// TODO: break-out
func (agent *Agent) ProcessHostSamples(quitCh <-chan struct{}) {
	context.GetLogger(agent).Info("machine usage processor started")
	defer context.GetLogger(agent).Info("machine usage processor stopped")

	if err := agent.machineCollector.Start(); err != nil {
		context.GetLogger(agent).Errorf("error starting machine usage collection: %v", err)
		return
	}

	for {
		select {
		case <-quitCh:
			return
		case sample := <-agent.machineCollector.GetChannel():
			e := reporting.Generate(agent, reporting.EventMachineSample, sample)
			go func() {
				_, err := agent.reporting.Report(agent, e)
				if err != nil {
					context.GetLogger(agent).Errorf("error reporting machine usage: %v", err)
				} else {
					context.GetLogger(agent).Debug("machine usage reported")
				}
			}()

		}
	}
}

// TODO: break-out
func (agent *Agent) ProcessStateChange(c *containers.StateChange, registered bool) {
	if !registered {
		if err := agent.registry.Register(agent, c.Container); err != nil {
			if err != containers.ErrNotTrackable {
				context.GetLogger(agent).Errorf("error registering container: %v", err)
			}

			return
		}
	}

	if c.State == containers.StateStopped {
		if err := agent.registry.Drop(agent, c.Container.Name); err != nil {
			context.GetLogger(agent).Errorf("error dropping container: %v", err)
			return
		}

		if _, err := agent.collector.Stop(agent, c.Container); err != nil {
			context.GetLogger(agent).Errorf("error stopping container usage collection: %v", err)
			return
		}
	} else {
		ch, err := agent.containers.GetContainerUsage(agent, c.Container.Name)
		if err != nil {
			context.GetLogger(agent).Errorf("error opening usage channel: %v", err)
			return
		} else if err = agent.collector.Collect(agent, ch); err != nil {
			context.GetLogger(agent).Errorf("error starting container usage collection: %v", err)
			return
		}
	}

	e := reporting.Generate(agent, reporting.EventStateChange, c)
	go func() {
		_, err := agent.reporting.Report(agent, e)
		if err != nil {
			context.GetLogger(agent).Errorf("error reporting state change: %v", err)
		} else {
			context.GetLogger(agent).Debug("state change reported")
		}
	}()
}

// TODO: break-out
func (agent *Agent) ProcessEvents(quitCh <-chan struct{}) {
	eventTypes := []containers.EventType{
		containers.EventContainerCreation,
		containers.EventContainerDeletion,
		containers.EventContainerOom,
		containers.EventContainerOomKill,
	}

	eventChan, err := agent.containers.WatchEvents(agent, eventTypes...)
	if err != nil {
		context.GetLogger(agent).Panicf("error opening event channel: %v", err)
		return
	}

	context.GetLogger(agent).Info("event monitor started")
	defer context.GetLogger(agent).Info("event monitor stopped")

	for {
		select {
		case <-quitCh:
			return
		case event := <-eventChan.GetChannel():
			var c *containers.ContainerInfo
			registered := false
			state := containers.StateFromEvent(event.Type)

			if cc, found := agent.registry.Get(event.Container.Name); found {
				c = cc
				registered = true
			} else if state != containers.StateStopped {
				c, err = agent.containers.GetContainer(agent, event.Container.Name)
				if err != nil {
					if err == containers.ErrContainerNotFound {
						context.GetLogger(agent).Warnf("info for container %q not available", event.Container.Name)
					} else {
						context.GetLogger(agent).Errorf("error getting event container info: %v", err)
					}

					continue
				}
			} else {
				continue
			}

			change := &containers.StateChange{
				Container: c,
				State:     state,
				Source:    event,
			}

			go agent.ProcessStateChange(change, registered)
		}
	}
}

// TODO: break-out
func (agent *Agent) ProcessSamples(quitCh <-chan struct{}) {
	context.GetLogger(agent).Info("sample collector started")
	defer context.GetLogger(agent).Info("sample collector stopped")
	for {
		select {
		case <-quitCh:
			return
		case sample := <-agent.collector.GetChannel():
			e := reporting.Generate(agent, reporting.EventSample, sample)
			go func() {
				_, err := agent.reporting.Report(agent, e)
				if err != nil {
					context.GetLogger(agent).Errorf("error reporting usage: %v", err)
				} else {
					context.GetLogger(agent).Debug("usage reported")
				}
			}()
		}
	}
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	ctx, err := configureLogging(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error configuring logging: %v", err)
	}

	log := context.GetLogger(ctx)
	log.Info("initializing agent")

	containersParams := config.Containers.Parameters()
	if containersParams == nil {
		containersParams = make(configuration.Parameters)
	}

	containersDriver, err := containersFactory.Create(config.Containers.Type(), containersParams)
	if err != nil {
		return nil, err
	}

	reportingParams := config.Reporting.Parameters()
	if reportingParams == nil {
		reportingParams = make(configuration.Parameters)
	}

	reportingDriver, err := reportingFactory.Create(config.Reporting.Type(), reportingParams)
	if err != nil {
		return nil, err
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	log.Infof("using %q containers driver", config.Containers.Type())
	log.Infof("using %q reporting driver", config.Reporting.Type())
	log.Infof("monitoring containers with a %q label", config.Tracking.Marker.Label)
	if config.Tracking.Marker.Env != "" {
		log.Infof("monitoring containers with a %q env variable", config.Tracking.Marker.Env)
	}

	return &Agent{
		Context:    ctx,
		config:     config,
		dispose:    disposer.New(),
		containers: containersDriver,
		collector:  collector.New(config.Collector),
		//machineCollector: collector.NewMachineCollector(config.Collector),
		reporting: reportingDriver,
		registry:  containers.NewRegistry(config.Tracking.Marker),
	}, nil
}

func configureLogging(ctx context.Context, config *configuration.Config) (context.Context, error) {
	log.SetLevel(logLevel(config.Log.Level))
	formatter := config.Log.Formatter
	if formatter == "" {
		formatter = "text"
	}

	switch formatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	case "text":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	default:
		if config.Log.Formatter != "" {
			return ctx, fmt.Errorf("unsupported log formatter: %q", config.Log.Formatter)
		}
	}

	if len(config.Log.Fields) > 0 {
		var fields []interface{}
		for k := range config.Log.Fields {
			fields = append(fields, k)
		}

		ctx = context.WithValues(ctx, config.Log.Fields)
		ctx = context.WithLogger(ctx, context.GetLogger(ctx, fields...))
	}

	ctx = context.WithLogger(ctx, context.GetLogger(ctx))
	return ctx, nil
}

func logLevel(level configuration.LogLevel) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q", level, err, l)
	}

	return l
}
