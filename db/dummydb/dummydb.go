package dummydb

import "github.com/williabk198/go-api-server-template/db"

type dummyDB struct{}

// Person implements db.Database.
func (d dummyDB) Person() db.Datastore[db.Person] {
	return personDatastore{}
}

func NewSession() db.Database {
	return dummyDB{}
}
