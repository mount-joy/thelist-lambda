package config

import "os"

type conf struct {
	getEnv func(string) string
}

func newConfig() *conf {
	return &conf{
		getEnv: func(key string) string { return os.Getenv(key) },
	}
}

func (c *conf) getConf() Config {
	switch c.getRuntimeEnvironment() {
	case envNameProd:
		return c.getProdConfig()
	default:
		return c.getDevConfig()
	}
}

var environments map[string]bool = map[string]bool{envNameDev: true, envNameProd: true}

func (c *conf) getRuntimeEnvironment() string {
	value := c.getEnv(envVarEnvironment)
	if environments[value] {
		return value
	}

	return envNameDev
}

var loadedConfig Config = newConfig().getConf()

// GetConfiguration returns the cofiguration values required at runtime
func GetConfiguration() Config {
	return loadedConfig
}
