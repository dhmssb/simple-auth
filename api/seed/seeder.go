package seed

import (
	"dsi/api/models"
	"log"

	"gorm.io/gorm"
)

var users = []models.User{
	{
		FullName:     "User Pertama",
		Email:        "usersatu@gmail.com",
		Password:     "password",
		Age:          40,
		MobileNumber: "087812311111",
	},
	{
		FullName:     "User Kedua",
		Email:        "userdua@gmail.com",
		Password:     "password",
		Age:          28,
		MobileNumber: "123123123123",
	},
}

func Load(db *gorm.DB) {
	// Drop the User table if it exists
	if err := db.Migrator().DropTable(&models.User{}); err != nil {
		log.Fatalf("Failed to drop User table: %v", err)
	}

	// Auto migrate the User table
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to auto migrate User table: %v", err)
	}

	// Seed the User table
	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Fatalf("Failed to seed User table: %v", err)
		}
	}
}
