package controller

import (
	"encoding/json"
	"net/http"
)

// baseResponse is a type provides a basic response message back to the client
type baseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"msg,omitempty"`
}

// dataResponse is a type that will provide the requested data back to the client
type dataResponse[T any] struct {
	baseResponse
	Data T `json:"data"`
}

// sendDataResponse is a convenience function that sends a response back to the client with the requested data
func sendDataResponse[T any](respData T, jsonEncoder *json.Encoder) error {
	return jsonEncoder.Encode(
		dataResponse[T]{
			baseResponse: baseResponse{
				Success: true,
			},
			Data: respData,
		},
	)
}

// sendErrorResponse is a convenience to send an error back to the client
func sendErrorResponse(w http.ResponseWriter, statusCode int, jsonEncoder *json.Encoder) error {
	w.WriteHeader(statusCode)
	switch statusCode {
	case http.StatusBadRequest:
		return jsonEncoder.Encode(baseResponse{Message: "failed to read request"})
	case http.StatusNotFound:
		return jsonEncoder.Encode(baseResponse{Message: "not found"})
	case http.StatusUnprocessableEntity:
		return jsonEncoder.Encode(baseResponse{Message: "malformed request data"})
	case http.StatusInternalServerError:
		return jsonEncoder.Encode(baseResponse{Message: "server encountered an error processing the request"})
	}

	return nil
}
