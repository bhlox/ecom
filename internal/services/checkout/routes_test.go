package checkout

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

func TestCheckout(t *testing.T) {
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

	var insertUsersStatements []string
	for _, user := range userTestData {
		insertStatement := fmt.Sprintf("INSERT INTO users (id, first_name,last_name,email,password) VALUES ('%v', '%v', '%v', '%v', '%v');", user.ID, user.FirstName, user.LastName, user.Email, user.Password)
		insertUsersStatements = append(insertUsersStatements, insertStatement)
	}
	insertUsersStatement := strings.Join(insertUsersStatements, " ")

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

	_, database := utils.CreateDbTestContainer(t, ctx, insertUsersStatement, insertProductStatements)

	router := http.NewServeMux()
	v1Router := http.NewServeMux()
	v1Router.Handle("/v1/", http.StripPrefix("/v1", router))

	userHandler := user.NewHandler(db.New(database))
	userHandler.RegisterRoutes(v1Router)

	checkoutHandler := NewHandler(db.New(database))
	checkoutHandler.RegisterRoutes(v1Router)

	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	testCases := []struct {
		payload    any
		shouldFail bool
		token      any
	}{
		{
			payload: types.OrderPayload{
				OrderItems: []types.OrderItemPayload{
					{
						ProductID: 3,
						Quantity:  2,
					},
					{
						ProductID: 4,
						Quantity:  1,
					},
				},
				Address: "testing avenue",
			},
			shouldFail: false,
		},
		{
			payload:    "checkers",
			shouldFail: true,
		},
		{
			payload: map[string]any{
				"orderItems": map[string]any{
					"ProductID": 7,
					"QUantity":  15,
				},
				"address": "legit address",
			},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run("checking out", func(t *testing.T) {

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

			assert.NotNilf(t, tokenResp.Token, "token should exist")

			payloadBytes, _ = json.Marshal(tc.payload)
			req, err = http.NewRequest("POST", "/v1/checkout", bytes.NewBuffer(payloadBytes))
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
				t.Errorf(err.Error())
			}
			// fmt.Println(rr.Body.String())
			if tc.shouldFail {
				assert.Contains(t, rr.Body.String(), "error")
			} else {
				fmt.Println(response.CreatedOrderItems)
				assert.NotNilf(t, response.CreatedOrderItems, "should contain the product details")
			}
		})
	}
}
