package db

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Removed     NullBool
}
