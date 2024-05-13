package order

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/bhlox/ecom/internal/db"
	"github.com/bhlox/ecom/internal/services/user"
	"github.com/bhlox/ecom/internal/types"
	"github.com/bhlox/ecom/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
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

	var insertUsersStatementsSlice []string
	for _, user := range userTestData {
		insertStatement := fmt.Sprintf("INSERT INTO users (id, first_name,last_name,email,password) VALUES ('%v', '%v', '%v', '%v', '%v');", user.ID, user.FirstName, user.LastName, user.Email, user.Password)
		insertUsersStatementsSlice = append(insertUsersStatementsSlice, insertStatement)
	}
	insertUsersStatement := strings.Join(insertUsersStatementsSlice, " ")

	csvFile, err := os.Open("../../../testdata/golang-ecom-product.csv")
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	var insertStatementsSlice []string
	for i, record := range records {
		// skipping because this just contains the column (e.g. name, description, etc)
		if i == 0 {
			continue
		}
		insertStatement := fmt.Sprintf("INSERT INTO products (id, name, description, image, price, quantity) VALUES ('%s', '%s','%s', '%s','%s', '%s');", record[0], record[1], record[2], record[3], record[4], record[5])
		insertStatementsSlice = append(insertStatementsSlice, insertStatement)
	}
	insertProductStatements := strings.Join(insertStatementsSlice, " ")

	purchasesSql, err := os.ReadFile("../../../testdata/dummy-purchases.sql")
	if err != nil {
		t.Fatalf("Error reading SQL file: %v", err)
	}
	purchasesSqlStatement := string(purchasesSql)

	_, database := utils.CreateDbTestContainer(t, ctx, insertUsersStatement, insertProductStatements, purchasesSqlStatement)

	router := http.NewServeMux()
	v1Router := http.NewServeMux()
	v1Router.Handle("/v1/", http.StripPrefix("/v1", router))

	userHandler := user.NewHandler(db.New(database))
	userHandler.RegisterRoutes(v1Router)

	orderHandler := NewHandler(db.New(database))
	orderHandler.RegisterRoutes(v1Router)

	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	t.Run("testing getting own order", func(t *testing.T) {
		payload := types.LoginPayload{
			Email:    "testing@test.com",
			Password: "Password123",
		}
		payloadBytes, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(payloadBytes))
		if err != nil {
			t.Errorf("Could not create request: %v", err)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		tokenResp := struct {
			Token string `json:"token"`
		}{}
		err = json.Unmarshal(rr.Body.Bytes(), &tokenResp)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		if tokenResp.Token == "" {
			t.Errorf("no token detected")
		}

		req, err = http.NewRequest("GET", "/v1/order/1", nil)
		if err != nil {
			t.Errorf("Could not create request: %v", err)
		}
		rr = httptest.NewRecorder()
		
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenResp.Token))
		router.ServeHTTP(rr, req)
		var response struct {
			Totals            float64        `json:"totals"`
			CreatedOrderItems []db.OrderItem `json:"createdOrderItems"`
		}

		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("something happend when marshalling")
		}

		assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
	})

	t.Run("testing getting other's order", func(t *testing.T) {
		payload := types.LoginPayload{
			Email:    "testing@test.com",
			Password: "Password123",
		}
		payloadBytes, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(payloadBytes))
		if err != nil {
			t.Errorf("Could not create request: %v", err)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		tokenResp := struct {
			Token string `json:"token"`
		}{}
		err = json.Unmarshal(rr.Body.Bytes(), &tokenResp)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		if tokenResp.Token == "" {
			t.Errorf("no token detected")
		}

		req, err = http.NewRequest("GET", "/v1/order/2", nil)
		if err != nil {
			t.Errorf("Could not create request: %v", err)
		}
		rr = httptest.NewRecorder()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenResp.Token))
		router.ServeHTTP(rr, req)
		assert.Contains(t, rr.Body.String(), "error")
	})

}
