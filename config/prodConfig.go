package config

func (c *conf) getProdConfig() Config {
	return Config{
		Endpoint: "",
		TableNames: TableNames{
			Items: c.getEnv(envVarTableNameItems),
			Lists: c.getEnv(envVarTableNameLists),
		},
	}
}
