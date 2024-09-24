package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/williabk198/go-api-server-template/db"
)

type personDataHandler struct {
	personDatastore db.Datastore[db.Person, uuid.UUID]
	logger          *slog.Logger
}

func (pdh personDataHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jsonEncoder := json.NewEncoder(w)

	var person person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		pdh.logger.Error("failed to parse JSON request", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, jsonEncoder)
		return
	}

	dbPerson, err := person.asDatabaseModel()
	if err != nil {
		pdh.logger.Error("failed to read request data", "error", err)
		sendErrorResponse(w, http.StatusUnprocessableEntity, jsonEncoder)
		return
	}

	err = pdh.personDatastore.Insert(ctx, dbPerson)
	if err != nil {
		pdh.logger.Error("failed to insert user into database", "error", err)
		sendErrorResponse(w, http.StatusInternalServerError, jsonEncoder)
		return
	}

	respData := pdh.personFromDatabaseModel(dbPerson)
	sendDataResponse(respData, jsonEncoder)
}

func (pdh personDataHandler) GetSpecific(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jsonEncoder := json.NewEncoder(w)
	personID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		pdh.logger.Error("failed to parse UUID from URL parameter", "error", err)
		sendErrorResponse(w, http.StatusNotFound, jsonEncoder)
		return
	}

	dbPerson, err := pdh.personDatastore.Get(ctx, personID)
	if err != nil {
		if errors.Is(err, db.ErrNoResultsFound) {
			sendErrorResponse(w, http.StatusNotFound, jsonEncoder)
			return
		}
		pdh.logger.Error("failed to get user from database", "error", err)
		sendErrorResponse(w, http.StatusInternalServerError, jsonEncoder)
		return
	}

	respData := pdh.personFromDatabaseModel(dbPerson)
	sendDataResponse(respData, jsonEncoder)
}

func (pdh personDataHandler) Remove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jsonEncoder := json.NewEncoder(w)
	personID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		pdh.logger.Error("failed to parse UUID from URL parameter", "error", err)
		sendErrorResponse(w, http.StatusNotFound, jsonEncoder)
		return
	}

	_, err = pdh.personDatastore.Remove(ctx, personID)
	if err != nil {
		if errors.Is(err, db.ErrNoResultsFound) {
			sendErrorResponse(w, http.StatusNotFound, jsonEncoder)
			return
		}
		pdh.logger.Error("failed to get user from database", "error", err)
		sendErrorResponse(w, http.StatusInternalServerError, jsonEncoder)
		return
	}

	jsonEncoder.Encode(baseResponse{Success: true})
}

func (pdh personDataHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jsonEncoder := json.NewEncoder(w)
	var person person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		pdh.logger.Error("failed to parse JSON request", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, jsonEncoder)
		return
	}

	dbPerson, err := person.asDatabaseModel()
	if err != nil {
		pdh.logger.Error("failed to read request data", "error", err)
		sendErrorResponse(w, http.StatusUnprocessableEntity, jsonEncoder)
		return
	}

	err = pdh.personDatastore.Update(ctx, dbPerson)
	if err != nil {
		if errors.Is(err, db.ErrNoResultsFound) {
			sendErrorResponse(w, http.StatusNotFound, jsonEncoder)
			return
		}
		pdh.logger.Error("failed to get user from database", "error", err)
		sendErrorResponse(w, http.StatusInternalServerError, jsonEncoder)
		return
	}

	respData := pdh.personFromDatabaseModel(dbPerson)
	sendDataResponse(respData, jsonEncoder)
}

func (pdh personDataHandler) personFromDatabaseModel(dbUser *db.Person) person {
	return person{
		ID:          dbUser.ID.String(),
		FirstName:   dbUser.FirstName,
		LastName:    dbUser.LastName,
		DateOfBirth: dbUser.DateOfBirth.Format("1/2/2006"),
		Removed:     *dbUser.Removed,
	}
}

type person struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dob"`
	Removed     bool   `json:"removed"`
}

func (p person) asDatabaseModel() (*db.Person, error) {
	var dateOfBirth time.Time
	// var removed bool
	var err error

	id := uuid.Nil
	if p.ID != "" {
		id, err = uuid.Parse(p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field 'id': %w", err)
		}
	}

	if p.DateOfBirth != "" {
		dateOfBirth, err = time.Parse(requestDateFormat, p.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field 'dob': %w", err)
		}
	}

	return &db.Person{
		ID:          id,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		DateOfBirth: dateOfBirth,
		Removed:     db.NewBool(p.Removed),
	}, nil
}
