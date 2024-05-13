package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bhlox/ecom/internal/db"
	"github.com/bhlox/ecom/internal/types"
	"github.com/bhlox/ecom/internal/utils"
	"github.com/stretchr/testify/assert"
)

type userTestCase struct {
	method       string
	testTitle    string
	payload      any
	route        string
	expectedCode int
	expectedMsg  string
}

func TestUserRoutes(t *testing.T) {
	ctx := context.Background()

	userTestDataJson, err := os.ReadFile("../../../testdata/user/user.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	var userTestData []db.User
	err = json.Unmarshal(userTestDataJson, &userTestData)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	var insertStatementsSlice []string
	for _, user := range userTestData {
		insertStatement := fmt.Sprintf("INSERT INTO users (id, first_name,last_name,email,password) VALUES ('%v', '%v', '%v', '%v', '%v');", user.ID, user.FirstName, user.LastName, user.Email, user.Password)
		insertStatementsSlice = append(insertStatementsSlice, insertStatement)
	}

	insertUsersStatement := strings.Join(insertStatementsSlice, " ")

	_, database := utils.CreateDbTestContainer(t, ctx, insertUsersStatement)

	router := http.NewServeMux()
	v1Router := http.NewServeMux()
	v1Router.Handle("/v1/", http.StripPrefix("/v1", router))

	userHandler := NewHandler(db.New(database))
	userHandler.RegisterRoutes(v1Router)

	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	testCases := []userTestCase{
		{
			method:    "POST",
			testTitle: "register should pass",
			payload: types.RegisterPayload{
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
			},
			route:        "/v1/register",
			expectedCode: http.StatusCreated,
			expectedMsg:  fmt.Sprintf("should return a status of code of %v", http.StatusCreated),
		},
		{
			method:    "POST",
			testTitle: "register should fail",
			payload: types.RegisterPayload{
				Firstname: "test",
				Lastname:  "test",
				Email:     "testing@test.com",
				Password:  "Password123",
			},
			route:        "/v1/register",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  fmt.Sprintf("should return a status of code of %v", http.StatusOK),
		},
		{
			method:    "POST",
			testTitle: "login should pass",
			payload: types.LoginPayload{
				Email:    "testing@test.com",
				Password: "Password123",
			},
			route:        "/v1/login",
			expectedCode: http.StatusOK,
			expectedMsg:  fmt.Sprintf("should return a status of code of %v", http.StatusOK),
		},
		{
			method:    "POST",
			testTitle: "login should fail",
			payload: types.LoginPayload{
				Email:    "john.doe@example.com",
				Password: "password12",
			},
			route:        "/v1/login",
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  fmt.Sprintf("should return a status of code of %v", http.StatusUnauthorized),
		},
	}

	for _, c := range testCases {
		t.Run(c.testTitle, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(c.payload)
			req, err := http.NewRequest(c.method, c.route, bytes.NewBuffer(payloadBytes))
			if err != nil {
				t.Errorf("Could not create request: %v", err)
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, c.expectedCode, rr.Code, c.expectedMsg)

		})
	}

	t.Run("get all users", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/users", nil)
		if err != nil {
			t.Errorf("Could not create request: %v", err)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
	})
	// t.Cleanup(func() {
	// 	if err := pgContainer.Terminate(ctx); err != nil {
	// 		t.Fatalf("failed to terminate container: %s", err)
	// 	}
	// })
}
