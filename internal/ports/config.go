package ports

// ConfigPort defines the interface for configuration adapters (driven port)
type ConfigPort interface {
	ApplyConfig(config map[string]interface{}) error
	GetConfig() map[string]interface{}
}
