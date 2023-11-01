package controllertests

import (
	"bytes"
	"dsi/api/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {

	err := refreshUserTable(server.DB)
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		fullname     string
		email        string
		age          int
		mobileNumb   string
		errorMessage string
	}{
		{
			inputJSON:    `{"fullname":"Pet", "email": "pet@gmail.com", "password": "password", "age": 10, "mobile_number": "7764124555126"}`,
			statusCode:   201,
			fullname:     "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"fullname":"Frank", "email": "pet@gmail.com", "password": "password", "age": 10, "mobile_number": "7764124555126"}`,
			statusCode:   500,
			errorMessage: "email Already Taken",
		},
		{
			inputJSON:    `{"fullname":"Kan", "email": "kangmail.com", "password": "password", "age": 10, "mobile_number": "7764124555126"}`,
			statusCode:   422,
			errorMessage: "invalid Email",
		},
		{
			inputJSON:    `{"fullname": "", "email": "kan@gmail.com", "password": "password", "age": 10, "mobile_number": "7764124555126"}`,
			statusCode:   422,
			errorMessage: "required fullname",
		},
		{
			inputJSON:    `{"fullname": "Kan", "email": "", "password": "password", "age": 10, "mobile_number": "7764124555126"}`,
			statusCode:   422,
			errorMessage: "required Email",
		},
		{
			inputJSON:    `{"fullname": "Kan", "email": "kan@gmail.com", "password": "", "age": 10, "mobile_number": "7764124555126"}`,
			statusCode:   422,
			errorMessage: "required Password",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseBytes := rr.Body.Bytes()
		responseMap := make(map[string]interface{})
		err = json.Unmarshal(responseBytes, &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["fullname"], v.fullname)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable(server.DB)
	if err != nil {
		log.Fatal(err)
	}

	// Seed users and get the first user's details
	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}

	// Get the first user's details
	var user models.User
	for _, u := range users {
		if u.ID == 1 { // Adjust the condition to match the user you want
			user = u
			break
		}
	}

	AuthID = user.ID
	AuthEmail = user.Email
	AuthPassword = "password" // Note the password in the database is already hashed, we want unhashed

	//Login the user and get the authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When no token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// User 2 trying to use User 1 token
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range userSample {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &responseMap) // Use rr.Body.Bytes() instead of []byte(rr.Body.String())
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
