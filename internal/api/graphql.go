package api

// GraphQL Implementation to API port

import (
	"errors"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/aawadall/bit-scout/internal/ports"
)

// GraphQLAPI is a minimal implementation of the APIPort interface for GraphQL.
type GraphQLAPI struct{}

func (g *GraphQLAPI) Name() string {
	return "GraphQL"
}

func (g *GraphQLAPI) Start() error {
	// TODO: Implement GraphQL server startup
	return errors.New("GraphQL Start not implemented")
}

func (g *GraphQLAPI) Stop() error {
	// TODO: Implement GraphQL server shutdown
	return errors.New("GraphQL Stop not implemented")
}

func (g *GraphQLAPI) Search(query ports.SearchQuery) (ports.SearchResults, error) {
	// TODO: Implement GraphQL search
	return ports.SearchResults{}, errors.New("GraphQL Search not implemented")
}

func (g *GraphQLAPI) Stats() (ports.Stats, error) {
	// TODO: Implement GraphQL stats
	return ports.Stats{}, errors.New("GraphQL Stats not implemented")
}

func (g *GraphQLAPI) Index(doc models.Document) error {
	// TODO: Implement GraphQL index
	return errors.New("GraphQL Index not implemented")
}
