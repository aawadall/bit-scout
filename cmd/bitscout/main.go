package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aawadall/bit-scout/internal/index"
	"github.com/aawadall/bit-scout/internal/loaders"
	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting bitscout")

	// Initialize loader registry
	registry := loaders.NewLoaderRegistry()
	registry.Register("filesystem", loaders.NewFilesystemLoader("."))

	// Load documents
	documents, err := registry.LoadAll()
	if err != nil {
		log.Error().Msgf("Error loading documents: %s", err)
		return
	}

	log.Info().Msgf("Loaded %d documents", len(documents))

	// Initialize and configure index
	idx := index.NewSimpleIndex()
	config := map[string]interface{}{
		"max_results": 10,
		"dimensions":  []string{"fileSize", "lastModified", "fileExtension"},
	}
	if err := idx.Configure(config); err != nil {
		log.Error().Msgf("Error configuring index: %s", err)
		return
	}

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

	// Start interactive search
	startInteractiveSearch(idx)
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