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

type Containers map[string]Parameters

func (containers Containers) Type() string {
	var containersType []string

	for k := range containers {
		containersType = append(containersType, k)
	}

	if len(containersType) > 1 {
		panic("multiple containers drivers specified in the configuration or environment: " + strings.Join(containersType, ", "))
	}

	if len(containersType) == 1 {
		return containersType[0]
	}

	return ""
}

func (containers Containers) Parameters() Parameters {
	return containers[containers.Type()]
}

func (containers Containers) setParameter(key string, value interface{}) {
	containers[containers.Type()][key] = value
}

func (containers *Containers) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var containersMap map[string]Parameters
	err := unmarshal(&containersMap)
	if err == nil && len(containersMap) > 1 {
		types := make([]string, 0, len(containersMap))
		for k := range containersMap {
			types = append(types, k)
		}

		if len(types) > 1 {
			return fmt.Errorf("Must provide exactly one containers type. provided: %v", types)
		}

		*containers = containersMap
		return nil
	}

	var containersType string
	if err = unmarshal(&containersType); err != nil {
		return err
	}

	*containers = Containers{
		containersType: Parameters{},
	}

	return nil
}

func (containers Containers) MarshalYAML() (interface{}, error) {
	if containers.Parameters() == nil {
		return containers.Type(), nil
	}

	return map[string]Parameters(containers), nil
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

type TrackerConfig struct {
	ServiceKeyLabel string `yaml:"service_key_label,omitempty"`
	TrackingLabel   string `yaml:"tracking_label,omitempty"`
}

type Config struct {
	Log        LogConfig       `yaml:"log"`
	Containers Containers      `yaml:"containers"`
	MockApi    MockApiConfig   `yaml:"mockapi"`
	Collector  CollectorConfig `yaml:"collector"`
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

		Containers: make(Containers),

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
					if v0_1.Containers.Type() == "" {
						return nil, fmt.Errorf("no containers configuration provided")
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
