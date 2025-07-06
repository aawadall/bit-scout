package ports

// PersistencePort defines the interface for persistence adapters (driven port)
type PersistencePort interface {
	Save(data interface{}) error
	Load() (interface{}, error)
	Close() error
}