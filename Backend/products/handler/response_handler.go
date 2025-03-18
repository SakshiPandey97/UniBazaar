package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"web-service/errors"
	"web-service/model"
)

func HandleError(w http.ResponseWriter, err error, message string) {
	switch err := err.(type) {
	case *errors.DatabaseError:
		handleErrorResponse(w, err.Message, err, err.StatusCode)

	case *errors.NotFoundError:
		handleErrorResponse(w, err.Message, err, err.StatusCode)

	case *errors.S3Error:
		handleErrorResponse(w, err.Message, err, err.StatusCode)

	case *errors.BadRequestError:
		handleErrorResponse(w, err.Message, err, err.StatusCode)

	default:
		handleErrorResponse(w, message, err, http.StatusInternalServerError)
	}
}

func handleErrorResponse(w http.ResponseWriter, message string, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.ErrorResponse{
		Error:   message,
		Details: "",
	}

	if err != nil {
		response.Details = err.Error()
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		log.Printf("Error encoding JSON response: %v\n", encodeErr)
	}

	log.Printf("Error [%d]: %s - %v\n", statusCode, message, err)
}

func HandleSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		handleErrorResponse(w, "Error encoding response", err, http.StatusInternalServerError)
		return
	}
	log.Printf("Success [%d]: Response sent successfully\n", statusCode)
}
