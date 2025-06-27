package corpus

/*
Interface for corpus loader, and registry of loaders.
*/

// CorpusLoader defines the interface for loading documents from a source.
type CorpusLoader interface {
	// Load loads documents from the given source (e.g., directory, file, URI).
	Load(source string) ([]Document, error)
}
