package auth_test

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"go/adv-demo/config"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/user"
	"go/adv-demo/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "a@a.ru",
	}, nil
}
func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const testEmail = "a@a.ru"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err := authService.Register(testEmail, "Vasia", "12")
	if err != nil {
		t.Fatal(err)
	}
	if email != testEmail {
		t.Errorf("got email %s, expected %s", email, testEmail)
	}
}

func bootstrap() (*auth.Handler, sqlmock.Sqlmock, error) {
	dataBase, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dataBase,
	}))
	if err != nil {

		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := auth.Handler{
		Config: &config.Config{
			Auth: config.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil

}

func TestLoginHandlerFailed(t *testing.T) {
	const testEmail = "a@a.ru"
	correctPassHash := "$2a$10$J5YzQCdVbIP1i1purfjYEOEEtlsr4UgX7VCOSkDkWoJIldoJ69Vaa"
	wrongPassword := "123456"
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow(testEmail, correctPassHash)
	mock.ExpectQuery("Select").WillReturnRows(rows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub Gorm database connection", err)
		return
	}
	data, err := json.Marshal(&auth.LoginRequest{
		Email:    testEmail,
		Password: wrongPassword,
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)
	if w.Code == http.StatusOK {
		t.Errorf("got %d, expected %d ", w.Code, http.StatusOK)
	}

}

func TestRegisterHandlerSuccess(t *testing.T) {
	const testEmail = "a@a.ru"
	password := "123456"
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("Select").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub Gorm database connection", err)
		return
	}
	data, err := json.Marshal(&auth.RegisterRequest{
		Email:    testEmail,
		Password: password,
		Name:     "Vasia",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("got %d, expected %d ", w.Code, http.StatusCreated)
	}

}
