package users

import "fmt"

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

func AddUser(user *User) (*UserResponse, error) {
	fmt.Println("User received ", user.FirstName)
	resp := &UserResponse{Id: "test"}
	return resp, nil
}
