package main

import (
	"github.com/aawadall/bit-scout/internal/corpus"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting bitscout")

	registry := corpus.NewLoaderRegistry()
	registry.Register("filesystem", corpus.NewFilesystemLoader("."))

	documents, err := registry.LoadAll()
	if err != nil {
		log.Error().Msgf("Error loading documents: %s", err)
		return
	}

	log.Info().Msgf("Loaded %d documents", len(documents))

	log.Info().Msgf("Documents: %d", len(documents))
	for _, document := range documents {
		document.Print()
	}
	
}