package config

import (
	"fmt"
	"os"
)

type conf struct {
	getEnv func(string) string
}

func newConfig() *conf {
	return &conf{
		getEnv: func(key string) string { return os.Getenv(key) },
	}
}

func (c *conf) getConf() Config {
	switch env := c.getRuntimeEnvironment(); env {
	case envNameProd:
		return c.getProdConfig()
	case envNameDev:
		return c.getDevConfig()
	default:
		panic(fmt.Sprintf("Unknown runtime enviroment: %s", env))
	}
}

func (c *conf) getRuntimeEnvironment() string {
	environments := map[string]bool{envNameDev: true, envNameProd: true}
	environment := c.getEnv(envVarEnvironment)
	if _, ok := environments[environment]; ok {
		return environment
	}

	return envNameDev
}

var loadedConfig Config = newConfig().getConf()

// GetConfiguration returns the cofiguration values required at runtime
func GetConfiguration() Config {
	return loadedConfig
}
