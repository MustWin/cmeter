package configuration

type Config struct {
}

func Resolve(args []string) (*Config, error) {
	return &Config{}, nil
}
