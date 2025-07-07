package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
		// Start interactive search (using the registered index)
		startInteractiveSearch(idx)
	}
}

func startInteractiveSearch(idx index.Index) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n=== Bit-Scout Search Interface ===")
	fmt.Println("Type 'quit' to exit, 'help' for commands")
	fmt.Println("Enter search query:")

	for {
		fmt.Print("> ")
		query, err := reader.ReadString('\n')
		if err != nil {
			log.Error().Msgf("Error reading input: %s", err)
			break
		}

		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		if query == "quit" || query == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if query == "help" {
			printHelp()
			continue
		}

		if query == "config" {
			showIndexConfig(idx)
			continue
		}

		// Complex Queries
		if query == "complex" {
			interactiveComplexQuery(idx)
			continue
		}
		// Perform search
		results, err := idx.Search(query)
		if err != nil {
			log.Error().Msgf("Search error: %s", err)
			fmt.Printf("Error performing search: %s\n", err)
			continue
		}

		// Display results
		displaySearchResults(results, query)
	}
}

func printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  <query>     - Search for documents containing the query")
	fmt.Println("  complex     - Interactive complex query builder")
	fmt.Println("  config      - Show current index configuration")
	fmt.Println("  help        - Show this help message")
	fmt.Println("  quit/exit   - Exit the application")
	fmt.Println("\nSearch types:")
	fmt.Println("  Simple:     <text> - searches in content, metadata, and paths")
	fmt.Println("  Advanced:   <dimension><operator><value> - boolean queries")
	fmt.Println("\nAdvanced query operators:")
	fmt.Println("  =           - equals (e.g., fileExtension=go)")
	fmt.Println("  !=          - not equals (e.g., fileExtension!=md)")
	fmt.Println("  <, <=       - less than, less than or equal")
	fmt.Println("  >, >=       - greater than, greater than or equal")
	fmt.Println("  contains    - contains text (e.g., filename contains main)")
	fmt.Println("\nExamples:")
	fmt.Println("  fileExtension=go")
	fmt.Println("  fileSize<1000")
	fmt.Println("  filename contains README")
	fmt.Println("  fileExtension=go and fileSize<1000")
	fmt.Println()
}

func showIndexConfig(idx index.Index) {
	config, err := idx.ShowConfig()
	if err != nil {
		log.Error().Msgf("Error getting index config: %s", err)
		fmt.Printf("Error getting index configuration: %s\n", err)
		return
	}

	fmt.Println("\n=== Index Configuration ===")
	if len(config) == 0 {
		fmt.Println("No configuration set")
	} else {
		for key, value := range config {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
	fmt.Println()
}

func displaySearchResults(results []models.Document, query string) {
	if len(results) == 0 {
		fmt.Printf("No results found for '%s'\n", query)
		return
	}

	fmt.Printf("\nFound %d result(s) for '%s':\n", len(results), query)
	fmt.Println(strings.Repeat("-", 50))

	for i, doc := range results {
		filename := doc.Meta["filename"]
		if filename == "" {
			filename = "Unknown"
		}

		fmt.Printf("%d. %s\n", i+1, filename)
		fmt.Printf("   Path: %s\n", doc.Source)
		fmt.Printf("   ID: %s\n", doc.ID)
		fmt.Printf("   Size: %s bytes\n", doc.Meta["fileSize"])
		fmt.Printf("   Modified: %s\n", doc.Meta["lastModified"])
		fmt.Printf("   Extension: %s\n", doc.Meta["extension"])
		fmt.Println()
	}
}

// Interactive Complex Query
// REPL for complex queries
// Takes arguments until special token is sent
func interactiveComplexQuery(idx index.Index) {
	reader := bufio.NewReader(os.Stdin)

	// Display intro and quick help
	fmt.Println("\n=== Complex Query Interface ===")
	fmt.Println("Given the following dimensions, operators, and logic operators, build a query to search the index.")
	fmt.Println("When done, type ## and enter to submit your query.")
	fmt.Println("Enter 'help' to display available dimensions, operators, and logic operators")
	fmt.Println("Enter 'quit' to exit")

	// Display available dimensions
	dimensions, err := idx.ShowConfig()
	if err != nil {
		log.Error().Msgf("Error getting index config: %s", err)
		fmt.Printf("Error getting index configuration: %s\n", err)
		return
	}

	// Display available dimensions
	fmt.Println("\nAvailable dimensions:")
	for _, dimension := range dimensions {
		fmt.Printf("  %s\n", dimension)
	}

	// display available operators
	fmt.Println("\nAvailable operators:")
	fmt.Println("  =           - equals (e.g., fileExtension=go)")
	fmt.Println("  !=          - not equals (e.g., fileExtension!=md)")
	fmt.Println("  <, <=       - less than, less than or equal")
	fmt.Println("  >, >=       - greater than, greater than or equal")
	fmt.Println("  contains    - contains text (e.g., filename contains main)")

	// display available logic operators
	fmt.Println("\nAvailable logic operators:")
	fmt.Println("  and         - logical AND")
	fmt.Println("  or          - logical OR")
	fmt.Println("  not         - logical NOT")

	// display available query types
	fmt.Println("\nAvailable query types:")
	fmt.Println("  simple      - simple text search")

	// Initialize query state
	query := ""
	queryParts := []string{}

	// REPL loop
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Error().Msgf("Error reading input: %s", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == "quit" || input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "help" {
			printHelp()
			continue
		}

		if input == "##" {
			// submit query
			results, err := idx.Search(query)
			if err != nil {
				log.Error().Msgf("Error searching: %s", err)
				fmt.Printf("Error searching: %s\n", err)
				continue
			}
			displaySearchResults(results, query)
			break
		}

		// Add input to query parts
		queryParts = append(queryParts, input)
		query = strings.Join(queryParts, " ")
		fmt.Printf("Current query: %s\n", query)
	}

}
