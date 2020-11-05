package config

// TableNames contains the dynamodb table names
type TableNames struct {
	Items string
	Lists string
}

// Config contains the cofiguration values required at runtime
type Config struct {
	Endpoint   string
	TableNames TableNames
}
