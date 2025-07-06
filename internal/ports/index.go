package ports

// IndexPort defines the interface for index adapters (driven port)
type IndexPort interface {
	AddDocument(doc interface{}) error
	Search(query string) ([]interface{}, error)
	Count() (int, error)
	Close() error
}