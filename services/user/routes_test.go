package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/panjiasmoroart/gopher-ecom/types"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	// t.Run("should fail if the user payload is empty", func(t *testing.T) {
	// 	// Empty payload scenario
	// 	req, err := http.NewRequest(http.MethodPost, "/register", nil)
	// 	if err != nil {
	// 		t.Fatalf("failed to create request: %v", err)
	// 	}

	// 	rr := httptest.NewRecorder()
	// 	router := mux.NewRouter()
	// 	router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
	// 	router.ServeHTTP(rr, req)

	// 	if rr.Code != http.StatusBadRequest {
	// 		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	// 	}

	// 	// Check for the error message in response
	// 	expectedError := `{"error":"invalid payload"}`
	// 	if rr.Body.String() != expectedError {
	// 		t.Errorf("expected response body %s, got %s", expectedError, rr.Body.String())
	// 	}
	// })

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "panji",
			LastName:  "asmoro",
			Email:     "invalid",
			Password:  "12345",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "salman",
			LastName:  "alfarisi",
			Email:     "valid@gmail.com",
			Password:  "12345",
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(u types.User) error {
	return nil
}
