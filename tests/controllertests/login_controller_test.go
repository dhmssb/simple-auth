package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {

	err := refreshUserTable(server.DB)
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("this is the error %v\n", err)
	}

	samples := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     "password", //Note the password has to be this, not the hashed one from the database
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "Wrong email",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {
		token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {
	refreshUserTable(server.DB)

	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("this is the error %v\n", err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "initest@gmail.com", "password": "password"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "initest@gmail.com", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: "incorrect Password",
		},
		{
			inputJSON:    `{"email": "frank@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "incorrect Details",
		},
		{
			inputJSON:    `{"email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "invalid Email",
		},
		{
			inputJSON:    `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required Email",
		},
		{
			inputJSON:    `{"email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "required Password",
		},
		{
			inputJSON:    `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required Email",
		},
	}

	for _, v := range samples {
		t.Run(v.inputJSON, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
			if err != nil {
				t.Errorf("Error creating the request: %v", err)
				return // Return early if there's an error
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.Login)
			handler.ServeHTTP(rr, req)

			if rr.Code != v.statusCode {
				t.Errorf("Expected status code %d, but got %d", v.statusCode, rr.Code)
			}

			if v.statusCode == 200 {
				if rr.Body.String() == "" {
					t.Error("Expected non-empty response body, but got an empty one")
				}
			}

			if v.statusCode == 422 && v.errorMessage != "" {
				responseBytes := rr.Body.Bytes()
				responseMap := make(map[string]interface{})
				if err := json.Unmarshal(responseBytes, &responseMap); err != nil {
					t.Errorf("Error converting response to JSON: %v", err)
					return // Return early if there's an error
				}
				if responseMap["error"] != v.errorMessage {
					t.Errorf("Expected error message '%s', but got '%s'", v.errorMessage, responseMap["error"])
				}
			}

		})
	}
}
