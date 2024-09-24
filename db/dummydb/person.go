package dummydb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/williabk198/go-api-server-template/db"
)

type personDatastore struct{}

// Get implements db.Datastore.
func (p personDatastore) Get(ctx context.Context, id uuid.UUID) (*db.Person, error) {

	result := &db.Person{
		ID:          id,
		FirstName:   "Testy",
		LastName:    "McTesterson",
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Removed:     db.NewBool(false),
	}

	return result, nil
}

// Insert implements db.Datastore.
func (p personDatastore) Insert(ctx context.Context, item *db.Person) error {
	item.ID = uuid.New()
	return nil
}

// Remove implements db.Datastore.
func (p personDatastore) Remove(ctx context.Context, id uuid.UUID) (*db.Person, error) {
	result := &db.Person{
		ID:          id,
		FirstName:   "Testy",
		LastName:    "McTesterson",
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Removed:     db.NewBool(true),
	}

	return result, nil
}

// Update implements db.Datastore.
func (p personDatastore) Update(ctx context.Context, item *db.Person) error {
	return nil
}
