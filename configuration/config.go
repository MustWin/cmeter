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

type Driver map[string]Parameters

func (driver Driver) Type() string {
	var driverType []string

	for k := range driver {
		driverType = append(driverType, k)
	}

	if len(driverType) > 1 {
		panic("multiple drivers specified in the configuration or environment: %s" + strings.Join(driverType, ", "))
	}

	if len(driverType) == 1 {
		return driverType[0]
	}

	return ""
}

func (driver Driver) Parameters() Parameters {
	return driver[driver.Type()]
}

func (driver Driver) setParameter(key string, value interface{}) {
	driver[driver.Type()][key] = value
}

func (driver *Driver) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var driverMap map[string]Parameters
	err := unmarshal(&driverMap)
	if err == nil && len(driverMap) > 0 {
		types := make([]string, 0, len(driverMap))
		for k := range driverMap {
			types = append(types, k)
		}

		if len(types) > 1 {
			return fmt.Errorf("Must provide exactly one driver type. provided: %v", types)
		}

		*driver = driverMap
		return nil
	}

	var driverType string
	if err = unmarshal(&driverType); err != nil {
		return err
	}

	*driver = Driver{
		driverType: Parameters{},
	}

	return nil
}

func (driver Driver) MarshalYAML() (interface{}, error) {
	if driver.Parameters() == nil {
		return driver.Type(), nil
	}

	return map[string]Parameters(driver), nil
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
	Rate int64 `yaml:"rate"`
}

type Marker struct {
	Env   string `yaml:"env,omitempty"`
	Label string `yaml:"label,omitempty"`
}

type TrackerConfig struct {
	Marker Marker `yaml:"marker,omitempty"`
}

type Config struct {
	Log        LogConfig       `yaml:"log"`
	Containers Driver          `yaml:"containers"`
	Reporting  Driver          `yaml:"reporting"`
	Collector  CollectorConfig `yaml:"collector"`
	Tracking   TrackerConfig   `yaml:"tracking"`
}

type v1_0Config Config

func newConfig() *Config {
	config := &Config{
		Log: LogConfig{
			Level:     "debug",
			Formatter: "text",
			Fields:    make(map[string]interface{}),
		},

		Containers: make(Driver),

		Reporting: make(Driver),

		Tracking: TrackerConfig{
			Marker: Marker{
				Env:   "CMETER_TRACKING",
				Label: "cmeter.tracking",
			},
		},

		Collector: CollectorConfig{
			Rate: 10000,
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
			Version: MajorMinorVersion(1, 0),
			ParseAs: reflect.TypeOf(v1_0Config{}),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				if v1_0, ok := c.(*v1_0Config); ok {
					if v1_0.Containers.Type() == "" {
						return nil, fmt.Errorf("no containers configuration provided")
					}

					return (*Config)(v1_0), nil
				}

				return nil, fmt.Errorf("Expected *v1_0Config, received %#v", c)
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
