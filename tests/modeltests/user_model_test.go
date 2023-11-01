package modeltests

import (
	"dsi/api/models"
	"log"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestSaveUser(t *testing.T) {

	err := refreshUserTable(server.DB)
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		FullName:     "blablatest",
		Email:        "test123@gmail.com",
		Password:     "password",
		Age:          12,
		MobileNumber: "129312480",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.FullName, savedUser.FullName)
	assert.Equal(t, newUser.Age, savedUser.Age)
	assert.Equal(t, newUser.MobileNumber, savedUser.MobileNumber)

}

func TestGetUserByID(t *testing.T) {

	err := refreshUserTable(server.DB)
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	foundUser, err := userInstance.FindUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.FullName, user.FullName)
	assert.Equal(t, foundUser.Age, user.Age)
	assert.Equal(t, foundUser.MobileNumber, user.MobileNumber)
}

func TestDeleteAUser(t *testing.T) {

	err := refreshUserTable(server.DB)
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()

	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	isDeleted, err := userInstance.DeleteAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
