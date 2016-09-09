package configuration

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

func (version *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var versionString string
	err := unmarshal(&versionString)
	if err != nil {
		return err
	}

	newVersion := Version(versionString)
	if _, err := newVersion.major(); err != nil {
		return err
	}

	if _, err := newVersion.minor(); err != nil {
		return err
	}

	*version = newVersion
	return nil
}

type Parameters map[string]interface{}

type Monitor map[string]Parameters

func (monitor Monitor) Type() string {
	var monitorType []string

	for k := range monitor {
		monitorType = append(monitorType, k)
	}

	if len(monitorType) > 1 {
		panic("multiple monitor drivers specified in the configuration or environment: " + strings.Join(monitorType, ", "))
	}

	if len(monitorType) == 1 {
		return monitorType[0]
	}

	return ""
}

func (monitor Monitor) Parameters() Parameters {
	return monitor[monitor.Type()]
}

func (monitor Monitor) setParameter(key string, value interface{}) {
	monitor[monitor.Type()][key] = value
}

func (monitor *Monitor) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var monitorMap map[string]Parameters
	err := unmarshal(&monitorMap)
	if err == nil && len(monitorMap) > 1 {
		types := make([]string, 0, len(monitorMap))
		for k := range monitorMap {
			types = append(types, k)
		}

		if len(types) > 1 {
			return fmt.Errorf("Must provide exactly one monitor type. provided: %v", types)
		}

		*monitor = monitorMap
		return nil
	}

	var monitorType string
	if err = unmarshal(&monitorType); err != nil {
		return err
	}

	*monitor = Monitor{
		monitorType: Parameters{},
	}

	return nil
}

func (monitor Monitor) MarshalYAML() (interface{}, error) {
	if monitor.Parameters() == nil {
		return monitor.Type(), nil
	}

	return map[string]Parameters(monitor), nil
}

type MockApiConfig struct {
	Addr string
}

type LogLevel string

func (logLevel *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strLogLevel string
	err := unmarshal(&strLogLevel)
	if err != nil {
		return err
	}

	strLogLevel = strings.ToLower(strLogLevel)
	switch strLogLevel {
	case "error", "warn", "info", "debug":
	default:
		return fmt.Errorf("Invalid log level %s. Must be one of [error, warn, info, debug]", strLogLevel)
	}

	*logLevel = LogLevel(strLogLevel)
	return nil
}

type LogConfig struct {
	Level     LogLevel               `yaml:"level,omitempty"`
	Formatter string                 `yaml:"formatter,omitempty"`
	Fields    map[string]interface{} `yaml:"fields,omitempty"`
}

type CollectorConfig struct {
	Rate          int64  `yaml:"rate"`
	KeyLabel      string `yaml:"key_label,omitempty"`
	TrackingLabel string `yaml:"tracking_label,omitempty"`
}

type Config struct {
	Log       LogConfig       `yaml:"log"`
	Monitor   Monitor         `yaml:"monitor"`
	MockApi   MockApiConfig   `yaml:"mockapi"`
	Collector CollectorConfig `yaml:"collector"`
}

type ApiConfig struct {
	Addr string `yaml:"addr,omitempty"`
}

type v0_1Config Config

func newConfig() *Config {
	config := &Config{
		Log: LogConfig{
			Level:     "debug",
			Formatter: "text",
			Fields:    make(map[string]interface{}),
		},

		Monitor: make(Monitor),

		MockApi: MockApiConfig{
			Addr: ":9090",
		},

		Collector: CollectorConfig{
			Rate:          10000,
			KeyLabel:      "com.cmeter.service",
			TrackingLabel: "com.cmeter.track",
		},
	}

	return config
}

func Parse(rd io.Reader) (*Config, error) {
	in, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	p := NewParser("cmeter", []VersionedParseInfo{
		{
			Version: MajorMinorVersion(0, 1),
			ParseAs: reflect.TypeOf(v0_1Config{}),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				if v0_1, ok := c.(*v0_1Config); ok {
					if v0_1.Monitor.Type() == "" {
						return nil, fmt.Errorf("no monitor configuration provided")
					}

					return (*Config)(v0_1), nil
				}

				return nil, fmt.Errorf("Expected *v0_1Config, received %#v", c)
			},
		},
	})

	config := new(Config)
	err = p.Parse(in, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
