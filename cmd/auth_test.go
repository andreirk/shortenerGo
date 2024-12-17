package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go/adv-demo/config"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initDB() *gorm.DB {
	err := godotenv.Load("test.env")
	if err != nil {
		panic("Error loading .env file")
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database")
	}
	return db
}

func initData(db *gorm.DB, param string) {
	passHash := "$2a$10$J5YzQCdVbIP1i1purfjYEOEEtlsr4UgX7VCOSkDkWoJIldoJ69Vaa"
	id := db.Create(&user.User{
		Email:    param,
		Password: passHash,
		Name:     "Bob Ankle",
	})
	fmt.Println("created id", id)
}
func cleanData(db *gorm.DB, param any) {
	db.Unscoped().
		Where("email = ?", param).
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	// Prepare
	db := initDB()
	testEmail := "test@test.ru"
	initData(db, testEmail)
	defer cleanData(db, testEmail)

	ts := httptest.NewServer(App(config.LoadConfig("test")))
	defer ts.Close()
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    testEmail,
		Password: "123456",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.AccessToken == "" {
		t.Fatal("access token empty")
	}
}

func TestLoginFailed(t *testing.T) {
	// Prepare
	db := initDB()
	testEmail := "test@test.ru"
	initData(db, testEmail)
	defer cleanData(db, testEmail)

	// Test
	ts := httptest.NewServer(App(config.LoadConfig("test")))
	defer ts.Close()
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.com",
		Password: "11111",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 401 {
		t.Fatalf("Expected %d got %d", 401, res.StatusCode)
	}
}
