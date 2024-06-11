package controller

import (
	"log/slog"
	"net/http"

	"github.com/williabk198/go-api-server-template/db"
)

const (
	requestDateFormat      string = "2006-01-02"
	requestTimeFormat      string = "15:04Z07:00"
	requestTimestampFormat string = requestDateFormat + "T" + requestTimeFormat
)

// DataHandler defines simple HTTP handlers that interact with database data.
type DataHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	GetSpecific(w http.ResponseWriter, r *http.Request)
	Remove(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

// Controller defines the different parts of the controller.
type Controller interface {
	Person() DataHandler
}

type controller struct {
	database db.Database
	logger   *slog.Logger
}

func (c controller) Person() DataHandler {
	return personDataHandler{
		personDatastore: c.database.Person(),
		logger:          c.logger,
	}
}

func NewController(logger *slog.Logger, database db.Database) Controller {
	return controller{
		database: database,
		logger:   logger,
	}
}
