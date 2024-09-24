package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/williabk198/go-api-server-template/db"
)

/* Common Testing Types and Functions */

type args struct {
	w http.ResponseWriter
	r *http.Request
}

type wantResp[T any] struct {
	statusCode int
	data       T
}

func assertResponse[T any](t *testing.T, wantResp wantResp[T], gotRespWriter http.ResponseWriter) {

	gotResp := gotRespWriter.(*httptest.ResponseRecorder).Result()

	// Ensure the returned status code is what we expect
	if gotResp.StatusCode != wantResp.statusCode {
		t.Errorf("recieved status code %d, wanted %d", gotResp.StatusCode, wantResp.statusCode)
		return
	}

	// Ensure that the returned response is in JSON format
	var gotData T
	if err := json.NewDecoder(gotResp.Body).Decode(&gotData); err != nil {
		t.Errorf("failed to decode response body %v", err)
		return
	}

	// Ensure that the data retured is what is expected
	assert.Equal(t, wantResp.data, gotData)
}

func encodeJSONBody(t *testing.T, data any) io.Reader {
	t.Helper()

	rawData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal request body for test: %v", err)
	}

	return bytes.NewBuffer(rawData)
}

/* Mock Data Section */

type mockDatastore[T db.Entity, U db.Identifier] struct {
	mock.Mock
}

func (md *mockDatastore[T, U]) Get(ctx context.Context, id U) (*T, error) {
	args := md.Called(ctx, id)
	return args.Get(0).(*T), args.Error(1)
}

func (md *mockDatastore[T, U]) Insert(ctx context.Context, data *T) error {
	args := md.Called(ctx, data)
	return args.Error(0)
}

func (md *mockDatastore[T, U]) Remove(ctx context.Context, id U) (*T, error) {
	args := md.Called(ctx, id)
	return args.Get(0).(*T), args.Error(1)
}

func (md *mockDatastore[T, U]) Update(ctx context.Context, filter *T) error {
	args := md.Called(ctx, filter)
	return args.Error(0)
}
