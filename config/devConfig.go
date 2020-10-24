package config

func (c *conf) getDevConfig() Config {
	return Config{
		Endpoint:   "http://localhost:8000",
		TableNames: TableNames{Items: "items", Lists: "lists"},
	}
}
