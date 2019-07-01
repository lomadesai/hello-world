package users

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// AddRequest accepts a new request for adding a user to our system.
func AddRequest(w http.ResponseWriter, r *http.Request) {
	var u User

	// Decode the incoming http request of content-type JSON.
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		encodeAndReturnError(w, 400, err)
	}

	// Handle user request.
	resp, err := AddUser(&u)
	if err != nil {
		encodeAndReturnError(w, 500, err)
	}

	// Encode and send the response back as JSON.
	encodeAndReturn(w, resp)
	return
}

func encodeAndReturn(w http.ResponseWriter, v interface{}) int {
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")

	err := encoder.Encode(v)
	if err != nil {
		_ = encodeAndReturnError(w, 500, err)
	}
	return 200
}

func encodeAndReturnError(w http.ResponseWriter, statusCode int, err error) int {
	w.WriteHeader(statusCode)

	var e ErrorResponse
	e.Error = err.Error()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(e)
	return statusCode
}
