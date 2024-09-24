package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ErrNoResultsFound is a database agnostic error that indicates that no results were found
var ErrNoResultsFound = fmt.Errorf("db query yeilded no results")

// Database defines the interactions with the database
type Database interface {
	Person() Datastore[Person, uuid.UUID]
}

// Datastore defines the basic interactions for a database entry
type Datastore[T Entity, U Identifier] interface {
	// Get retrieves an item from the database that matches the passed in item.
	Get(ctx context.Context, id U) (*T, error)
	// Insert puts the given item into the database.
	Insert(ctx context.Context, item *T) error
	// Remove marks the given item as removed in the database. This should NOT actually remove the item from the database.
	Remove(ctx context.Context, id U) (*T, error)
	// Update changes an item in the database to the given value.
	// The key of the database item to be updated must be provided in the passed in item.
	Update(ctx context.Context, item *T) error
}

// Entity is a type constraint which represents the items that are stored in the database.
type Entity interface {
	Person
}

type Identifier interface {
	int | uuid.UUID
}

// NullBool is a database agnostic type to represent a nullable boolean value.
type NullBool *bool

func NewBool(b bool) NullBool {
	return &b
}
