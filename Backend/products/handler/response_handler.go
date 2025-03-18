package handler

import (
	"encoding/json"
	"log"
	"net/http"
	customerrors "web-service/errors"
	"web-service/model"
)

func HandleError(w http.ResponseWriter, err error, message string) {
	statusCode := http.StatusInternalServerError
	errorMessage := message
	var errorDetails string

	if customErr, ok := err.(customerrors.CustomErrorInterface); ok {
		statusCode = customErr.GetStatusCode()
		errorMessage = customErr.GetMessage()
		errorDetails = customErr.Error()
	} else if err != nil {
		errorDetails = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.ErrorResponse{
		Error:   errorMessage,
		Details: errorDetails,
	}
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		log.Printf("Error encoding JSON response: %v\n", encodeErr)
	}

	log.Printf("HTTP Error: Status=%d, Message=\"%s\", Error=\"%v\"", statusCode, errorMessage, err)
}

func HandleSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		HandleError(w, err, "Error encoding response")
		return
	}
	log.Printf("HTTP Success: Status=%d, Message=\"Response sent successfully\"", statusCode)
}
