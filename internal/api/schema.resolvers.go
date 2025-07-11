package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"fmt"
)

// Start is the resolver for the start field.
func (r *mutationResolver) Start(ctx context.Context) (*CommandResult, error) {
	panic(fmt.Errorf("not implemented: Start - start"))
}

// Stop is the resolver for the stop field.
func (r *mutationResolver) Stop(ctx context.Context) (*CommandResult, error) {
	panic(fmt.Errorf("not implemented: Stop - stop"))
}

// Index is the resolver for the index field.
func (r *mutationResolver) Index(ctx context.Context, document DocumentInput) (*CommandResult, error) {
	panic(fmt.Errorf("not implemented: Index - index"))
}

// Ping is the resolver for the ping field.
func (r *queryResolver) Ping(ctx context.Context) (*PingResult, error) {
	panic(fmt.Errorf("not implemented: Ping - ping"))
}

// Stats is the resolver for the stats field.
func (r *queryResolver) Stats(ctx context.Context) (*StatsResult, error) {
	panic(fmt.Errorf("not implemented: Stats - stats"))
}

// Search is the resolver for the search field.
func (r *queryResolver) Search(ctx context.Context, query QueryInput) (*SearchResult, error) {
	panic(fmt.Errorf("not implemented: Search - search"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
