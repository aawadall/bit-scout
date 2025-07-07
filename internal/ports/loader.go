package ports

// LoaderPort defines the interface for loader adapters (driven port)
type LoaderPort interface {
	Load(source string) ([]interface{}, error)
}
