package users

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	DateOfBirth string `json:"DateOfBirth"`
	Email       string `json:"Email"`
	PhoneNumber string `json:"PhoneNumber"`
}

type UserResponse struct {
	Id string `json:"UserId"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserService struct {
	store *UserStore
}

func CreateUserService(store *UserStore) (*UserService, error) {
	svc := &UserService{store: store}
	return svc, nil
}

// AddRequest accepts a new request for adding a user to our system.
func (svc *UserService) AddRequest(w http.ResponseWriter, r *http.Request) {
	var u User

	// Decode the incoming http request of content-type JSON.
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		encodeAndReturnError(w, 400, err)
	}

	// Handle user request.
	resp, err := svc.addUser(&u)
	if err != nil {
		encodeAndReturnError(w, 500, err)
	}

	// Encode and send the response back as JSON.
	encodeAndReturn(w, resp)
	return
}

func (svc *UserService) addUser(user *User) (*UserResponse, error) {
	id := generateUserId(user)
	err := svc.store.UpsertUser(id, user)
	if err != nil {
		return nil, err
	}

	resp := &UserResponse{Id: id}
	return resp, nil
}

// helpers
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

func generateUserId(user *User) string {
	u := fmt.Sprintf("%+v", user)

	hash := sha256.New()
	hash.Write([]byte(u))
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return user.FirstName + "_" + user.LastName + "_" + sha
}
