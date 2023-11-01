package modeltests

import (
	"dsi/api/controllers"
	"dsi/api/models"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	server       = controllers.Server{}
	userInstance = models.User{}
)

func TestMain(m *testing.M) {

	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "postgres" {
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		dialector := postgres.New(postgres.Config{
			DSN: dsn,
		})
		server.DB, err = gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable(db *gorm.DB) error {
	migrator := db.Migrator()

	if migrator.HasTable(&models.User{}) {
		if err := migrator.DropTable(&models.User{}); err != nil {
			return err
		}
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {
	err := refreshUserTable(server.DB)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		FullName:     "USER Tst",
		Email:        "initest@gmail.com",
		Password:     "password",
		Age:          15,
		MobileNumber: "0813364571898",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		{
			FullName:     "Usertst Pertama",
			Email:        "usersatutst@gmail.com",
			Password:     "$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Age:          40,
			MobileNumber: "087812311111",
		},
		{
			FullName:     "Usertst Kedua",
			Email:        "userduatst@gmail.com",
			Password:     "$2y$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Age:          28,
			MobileNumber: "123123123123",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
