package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinefinder/internal/model"
)

type mockUserService struct{}

func (m *mockUserService) Create(user model.User) (*model.User, error) {
	user.ID = 1
	return &user, nil
}

func (m *mockUserService) List() ([]model.User, error) {
	return []model.User{
		{
			ID:    1,
			Email: "test@example.com",
		},
	}, nil
}

func (m *mockUserService) GetByID(id int) (*model.User, error) {
	return &model.User{
		ID:    id,
		Email: "test@example.com",
	}, nil
}

func (m *mockUserService) ValidateUser(email string, password string) (*model.User, error) {
	return &model.User{
		ID:    1,
		Email: email,
	}, nil
}

func TestCreateUser_Success(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	body := model.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperado 201, veio %d", resp.StatusCode)
	}
}