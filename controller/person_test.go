package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/go-api-server-template/db"
)

func Test_personDataHandler_Add(t *testing.T) {
	testUUID, _ := uuid.NewRandom()
	testLogger := slog.Default()

	mockPersonStore := &mockDatastore[db.Person, uuid.UUID]{}

	// Setup the expected parameters to be passes to the `Insert` when it is called.
	// `mock.Anything` is used as the first parameter to `Insert` since the value of the context being
	// passed in doesn't really matter(unless if you are testing a context that has a timeout/expiration).
	mockPersonStore.On("Insert", mock.Anything, &db.Person{
		FirstName:   "Testy",
		LastName:    "McTesterson",
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Removed:     db.NewBool(false),
	}).Run(func(args mock.Arguments) {
		// Using `Run` here since the `Insert` function mutates the passed in db.Person data
		// So, we need to replicate that here.
		person := args.Get(1).(*db.Person)
		person.ID = testUUID
	}).Return(error(nil)) // Set what the `Insert` function should return if the `On` conditions were met.

	mockPersonStore.On("Insert", mock.Anything, &db.Person{
		FirstName:   "Corrupted Value",
		LastName:    "McTesterson",
		DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Removed:     db.NewBool(false),
	}).Return(fmt.Errorf("mock error"))

	tests := []struct {
		name     string
		pdh      personDataHandler
		args     args
		wantResp wantResp[dataResponse[person]]
	}{
		{
			name: "Success",
			pdh: personDataHandler{
				personDatastore: mockPersonStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/person", encodeJSONBody(t, person{
					FirstName:   "Testy",
					LastName:    "McTesterson",
					DateOfBirth: "1970-01-01",
				})),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusOK,
				data: dataResponse[person]{
					baseResponse: baseResponse{Success: true},
					Data: person{
						ID:          testUUID.String(),
						FirstName:   "Testy",
						LastName:    "McTesterson",
						DateOfBirth: "1/1/1970",
					},
				},
			},
		},
		{
			name: "Bad Request Format",
			pdh: personDataHandler{
				personDatastore: mockPersonStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/person", strings.NewReader("malformed data")),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusBadRequest,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "failed to read request"},
				},
			},
		},
		{
			name: "Bad Request Data",
			pdh: personDataHandler{
				personDatastore: mockPersonStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/person", encodeJSONBody(t, person{
					FirstName:   "Testy",
					LastName:    "McTesterson",
					DateOfBirth: "1/1/1970",
				})),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusUnprocessableEntity,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "malformed request data"},
				},
			},
		},
		{
			name: "Database Error",
			pdh: personDataHandler{
				personDatastore: mockPersonStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/person", encodeJSONBody(t, person{
					FirstName:   "Corrupted Value",
					LastName:    "McTesterson",
					DateOfBirth: "1970-01-01",
				})),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusInternalServerError,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "server encountered an error processing the request"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pdh.Add(tt.args.w, tt.args.r)
			assertResponse(t, tt.wantResp, tt.args.w)
		})
	}
}

func Test_personDataHandler_GetSpecific(t *testing.T) {
	testUUID, _ := uuid.NewRandom()
	errorUUID, _ := uuid.NewRandom()
	dneUUID, _ := uuid.NewRandom()

	testLogger := slog.Default()
	mockPersonDatastore := &mockDatastore[db.Person, uuid.UUID]{}
	mockPersonDatastore.On("Get", mock.Anything, testUUID).Return(
		&db.Person{
			ID:          testUUID,
			FirstName:   "Some",
			LastName:    "Tester",
			DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			Removed:     db.NewBool(false),
		},
		error(nil),
	)
	mockPersonDatastore.On("Get", mock.Anything, errorUUID).Return(
		(*db.Person)(nil),
		fmt.Errorf("mockError"),
	)
	mockPersonDatastore.On("Get", mock.Anything, dneUUID).Return(
		(*db.Person)(nil),
		db.ErrNoResultsFound,
	)

	tests := []struct {
		name      string
		pdh       personDataHandler
		args      args
		urlParams map[string]string
		wantResp  wantResp[dataResponse[person]]
	}{
		{
			name: "Success",
			pdh: personDataHandler{
				personDatastore: mockPersonDatastore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/person/{id}", nil),
			},
			urlParams: map[string]string{"id": testUUID.String()},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusOK,
				data: dataResponse[person]{
					baseResponse: baseResponse{Success: true},
					Data: person{
						ID:          testUUID.String(),
						FirstName:   "Some",
						LastName:    "Tester",
						DateOfBirth: "1/1/1970",
					},
				},
			},
		},
		{
			name: "Bad UUID",
			pdh: personDataHandler{
				personDatastore: mockPersonDatastore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/person/{id}", nil),
			},
			urlParams: map[string]string{"id": "badUUID"},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusNotFound,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "not found"},
				},
			},
		},
		{
			name: "Database Error",
			pdh: personDataHandler{
				personDatastore: mockPersonDatastore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/person/{id}", nil),
			},
			urlParams: map[string]string{"id": errorUUID.String()},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusInternalServerError,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "server encountered an error processing the request"},
				},
			},
		},
		{
			name: "ID not in Database",
			pdh: personDataHandler{
				personDatastore: mockPersonDatastore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/person/{id}", nil),
			},
			urlParams: map[string]string{"id": dneUUID.String()},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusNotFound,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "not found"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*Setup request context and headers*/
			chiContext := chi.NewRouteContext()
			for k, v := range tt.urlParams {
				chiContext.URLParams.Add(k, v)
			}
			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, chiContext))

			tt.pdh.GetSpecific(tt.args.w, tt.args.r)
			assertResponse(t, tt.wantResp, tt.args.w)
		})
	}
}

func Test_personDataHandler_Remove(t *testing.T) {

	testUUID, _ := uuid.NewRandom()
	errorUUID, _ := uuid.NewRandom()
	dneUUID, _ := uuid.NewRandom()

	testLogger := slog.Default()
	mockPersonDatastore := &mockDatastore[db.Person, uuid.UUID]{}
	mockPersonDatastore.On("Remove", mock.Anything, testUUID).Return(
		&db.Person{
			ID:          testUUID,
			FirstName:   "Some",
			LastName:    "Tester",
			DateOfBirth: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			Removed:     db.NewBool(false),
		},
		error(nil),
	)
	mockPersonDatastore.On("Remove", mock.Anything, errorUUID).Return(
		(*db.Person)(nil),
		fmt.Errorf("mockError"),
	)
	mockPersonDatastore.On("Remove", mock.Anything, dneUUID).Return(
		(*db.Person)(nil),
		db.ErrNoResultsFound,
	)

	tests := []struct {
		name      string
		pdh       personDataHandler
		args      args
		urlParams map[string]string
		wantResp  wantResp[dataResponse[person]]
	}{
		{
			name: "Success",
			pdh: personDataHandler{
				personDatastore: mockPersonDatastore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/person/{id}", nil),
			},
			urlParams: map[string]string{"id": testUUID.String()},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusOK,
				data: dataResponse[person]{
					baseResponse: baseResponse{Success: true},
				},
			},
		},
		{
			name: "Bad UUID",
			pdh: personDataHandler{
				personDatastore: mockPersonDatastore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/person/{id}", nil),
			},
			urlParams: map[string]string{"id": "badUUID"},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusNotFound,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "not found"},
				},
			},
		},
		{
			name: "Database Error",
			pdh: personDataHandler{
				personDatastore: mockPersonDatastore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/person/{id}", nil),
			},
			urlParams: map[string]string{"id": errorUUID.String()},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusInternalServerError,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "server encountered an error processing the request"},
				},
			},
		},
		{
			name: "ID not in Database",
			pdh: personDataHandler{
				personDatastore: mockPersonDatastore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/person/{id}", nil),
			},
			urlParams: map[string]string{"id": dneUUID.String()},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusNotFound,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "not found"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*Setup request context and headers*/
			chiContext := chi.NewRouteContext()
			for k, v := range tt.urlParams {
				chiContext.URLParams.Add(k, v)
			}
			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, chiContext))

			tt.pdh.Remove(tt.args.w, tt.args.r)
			assertResponse(t, tt.wantResp, tt.args.w)
		})
	}
}

func Test_personDataHandler_Update(t *testing.T) {
	testUUID, _ := uuid.NewRandom()
	errorUUID, _ := uuid.NewRandom()
	dneUUID, _ := uuid.NewRandom()

	testLogger := slog.Default()
	mockUserStore := &mockDatastore[db.Person, uuid.UUID]{}
	mockUserStore.On("Update", mock.Anything, &db.Person{
		ID:          testUUID,
		FirstName:   "Another",
		LastName:    "Tester",
		DateOfBirth: time.Date(1992, 1, 27, 0, 0, 0, 0, time.UTC),
		Removed:     db.NewBool(false),
	}).Return(error(nil))
	mockUserStore.On("Update", mock.Anything, &db.Person{
		ID:          errorUUID,
		FirstName:   "Another",
		LastName:    "Tester",
		DateOfBirth: time.Date(1992, 1, 27, 0, 0, 0, 0, time.UTC),
		Removed:     db.NewBool(false),
	}).Return(fmt.Errorf("mock error"))
	mockUserStore.On("Update", mock.Anything, &db.Person{
		ID:          dneUUID,
		FirstName:   "Another",
		LastName:    "Tester",
		DateOfBirth: time.Date(1992, 1, 27, 0, 0, 0, 0, time.UTC),
		Removed:     db.NewBool(false),
	}).Return(db.ErrNoResultsFound)

	tests := []struct {
		name     string
		pdh      personDataHandler
		args     args
		wantResp wantResp[dataResponse[person]]
	}{
		{
			name: "Success",
			pdh: personDataHandler{
				personDatastore: mockUserStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPut, "/person/{id}", encodeJSONBody(t, person{
					ID:          testUUID.String(),
					FirstName:   "Another",
					LastName:    "Tester",
					DateOfBirth: "1992-01-27",
					Removed:     false,
				})),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusOK,
				data: dataResponse[person]{
					baseResponse: baseResponse{Success: true},
					Data: person{
						ID:          testUUID.String(),
						FirstName:   "Another",
						LastName:    "Tester",
						DateOfBirth: "1/27/1992",
						Removed:     false,
					},
				},
			},
		},
		{
			name: "Bad Request Format",
			pdh: personDataHandler{
				personDatastore: mockUserStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPut, "/person/{id}", strings.NewReader("malformed data")),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusBadRequest,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "failed to read request"},
				},
			},
		},
		{
			name: "Bad UUID in Request",
			pdh: personDataHandler{
				personDatastore: mockUserStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPut, "/person/{id}", encodeJSONBody(t, person{
					ID: "BadUUID",
				})),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusUnprocessableEntity,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "malformed request data"},
				},
			},
		},
		{
			name: "Database Error",
			pdh: personDataHandler{
				personDatastore: mockUserStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPut, "/person/{id}", encodeJSONBody(t, person{
					ID:          errorUUID.String(),
					FirstName:   "Another",
					LastName:    "Tester",
					DateOfBirth: "1992-01-27",
					Removed:     false,
				})),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusInternalServerError,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "server encountered an error processing the request"},
				},
			},
		},
		{
			name: "ID not in Database",
			pdh: personDataHandler{
				personDatastore: mockUserStore,
				logger:          testLogger,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPut, "/person/{id}", encodeJSONBody(t, person{
					ID:          dneUUID.String(),
					FirstName:   "Another",
					LastName:    "Tester",
					DateOfBirth: "1992-01-27",
					Removed:     false,
				})),
			},
			wantResp: wantResp[dataResponse[person]]{
				statusCode: http.StatusNotFound,
				data: dataResponse[person]{
					baseResponse: baseResponse{Message: "not found"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pdh.Update(tt.args.w, tt.args.r)
			assertResponse(t, tt.wantResp, tt.args.w)
		})
	}
}
