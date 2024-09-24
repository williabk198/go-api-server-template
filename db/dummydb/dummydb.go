package dummydb

import (
	"github.com/google/uuid"
	"github.com/williabk198/go-api-server-template/db"
)

type dummyDB struct{}

// Person implements db.Database.
func (d dummyDB) Person() db.Datastore[db.Person, uuid.UUID] {
	return personDatastore{}
}

func NewSession() db.Database {
	return dummyDB{}
}
