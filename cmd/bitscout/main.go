package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aawadall/bit-scout/internal/api"
	"github.com/aawadall/bit-scout/internal/engine"
	"github.com/aawadall/bit-scout/internal/index"
	"github.com/aawadall/bit-scout/internal/loaders"
	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
)

// Adapter for index.SimpleIndex to ports.IndexPort
// Only implements required methods
// (AddDocument, Search, Count, Close)
type simpleIndexAdapter struct {
	idx *index.SimpleIndex
}

func (a *simpleIndexAdapter) AddDocument(doc interface{}) error {
	d, ok := doc.(models.Document)
	if !ok {
		return fmt.Errorf("expected models.Document, got %T", doc)
	}
	return a.idx.AddDocument(d)
}

func (a *simpleIndexAdapter) Search(query string) ([]interface{}, error) {
	results, err := a.idx.Search(query)
	if err != nil {
		return nil, err
	}
	out := make([]interface{}, len(results))
	for i, d := range results {
		out[i] = d
	}
	return out, nil
}

func (a *simpleIndexAdapter) Count() (int, error) {
	return a.idx.Count()
}

func (a *simpleIndexAdapter) Close() error {
	return a.idx.Close()
}

// Adapter for loaders.FilesystemLoader to ports.LoaderPort
// Only implements required method (Load)
type filesystemLoaderAdapter struct {
	loader *loaders.FilesystemLoader
}

func (a *filesystemLoaderAdapter) Load(source string) ([]interface{}, error) {
	docs, err := a.loader.Load()
	if err != nil {
		return nil, err
	}
	out := make([]interface{}, len(docs))
	for i, d := range docs {
		out[i] = d
	}
	return out, nil
}

// LoaderConfig represents a loader configuration from the starter config
// Example: { "name": "filesystem", "type": "FilesystemLoader", "config": { "root": "." } }
type LoaderConfig struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// APIConfig represents an API configuration from the starter config
// Example: { "name": "graphql", "type": "GraphQL", "config": { "listen": ":8080" } }
type APIConfig struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// StarterConfig holds the structure for the starter JSON config
// Only index config is used for now, but features can be extended
// as needed.
type StarterConfig struct {
	Index   map[string]interface{} `json:"indexes"`
	Loaders []LoaderConfig         `json:"loaders"`
	Apis    []APIConfig            `json:"apis"`
	// Features map[string]features.ExtractorConfig `json:"features"` // Uncomment if you want to support feature config
}

func loadStarterConfig(path string) (*StarterConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg StarterConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	log.Info().Msg("Starting bitscout")

	// Parse flags
	daemon := flag.Bool("daemon", false, "Run as a background daemon (no interactive search)")
	configPath := flag.String("config", "config/starter_config.json", "Path to starter config JSON file")
	flag.Parse()

	// Initialize EngineCore
	core := engine.NewEngineCore()

	// Initialize loader registry and register loader
	registry := loaders.NewLoaderRegistry()
	filesystemLoader := loaders.NewFilesystemLoader(".")
	registry.Register("filesystem", filesystemLoader)
	// Register loader with core using adapter
	core.RegisterLoader("filesystem", &filesystemLoaderAdapter{loader: filesystemLoader})

	// Load documents
	documents, err := registry.LoadAll()
	if err != nil {
		log.Error().Msgf("Error loading documents: %s", err)
		return
	}

	log.Info().Msgf("Loaded %d documents", len(documents))

	// Load starter config
	cfg, err := loadStarterConfig(*configPath)
	if err != nil {
		log.Warn().Msgf("Could not load config file %s: %s. Using default config.", *configPath, err)
	}

	// Initialize and configure index
	idx := index.NewSimpleIndex()
	if cfg != nil && cfg.Index != nil {
		if err := idx.Configure(cfg.Index); err != nil {
			log.Error().Msgf("Error configuring index from config file: %s", err)
			return
		}
	} else {
		config := map[string]interface{}{
			"max_results": 10,
			"dimensions":  []string{"fileSize", "lastModified", "fileExtension"},
		}
		if err := idx.Configure(config); err != nil {
			log.Error().Msgf("Error configuring index: %s", err)
			return
		}
	}
	// Register index with core using adapter
	core.RegisterIndex("simple", &simpleIndexAdapter{idx: idx})

	// Add documents to index
	if err := idx.AddDocuments(documents); err != nil {
		log.Error().Msgf("Error adding documents to index: %s", err)
		return
	}

	// Get index statistics
	count, err := idx.Count()
	if err != nil {
		log.Error().Msgf("Error getting index count: %s", err)
	} else {
		log.Info().Msgf("Index contains %d documents", count)
	}

	size, err := idx.Size()
	if err != nil {
		log.Error().Msgf("Error getting index size: %s", err)
	} else {
		log.Info().Msgf("Index size: %d bytes", size)
	}

	if *daemon {
		log.Info().Msgf("Running in daemon mode. No interactive search. PID: %d", os.Getpid())
		// Keep the process alive
		select {}
	} else {
		// Create your API implementation (inject dependencies as needed)
		gqlAPI := &api.GraphQLAPI{}
		if err := gqlAPI.Start(); err != nil {
			log.Error().Msgf("Failed to start GraphQL server: %s", err)
		}
	}
}
