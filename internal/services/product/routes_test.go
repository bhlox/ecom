package product

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
	"github.com/bhlox/ecom/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestProducts(t *testing.T) {
	ctx := context.Background()

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
	for _, record := range records {
		insertStatement := fmt.Sprintf("INSERT INTO products (id, name, description, image, price, quantity) VALUES ('%s', '%s','%s', '%s','%s', '%s');", record[0], record[1], record[2], record[3], record[4], record[5])
		insertStatementsSlice = append(insertStatementsSlice, insertStatement)
	}
	insertProductStatements := strings.Join(insertStatementsSlice, " ")

	database := utils.CreateDbTestContainer(t, ctx, insertProductStatements)

	router := http.NewServeMux()
	v1Router := http.NewServeMux()
	v1Router.Handle("/v1/", http.StripPrefix("/v1", router))

	productHandler := NewHandler(db.New(database))
	productHandler.RegisterRoutes(v1Router)

	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	testCasesInsertProduct := []struct {
		payload    any
		shouldFail bool
	}{
		{
			payload: db.CreateProductParams{
				Name:        "limited glass",
				Description: "limited",
				Price:       "45.99",
				Image:       "random image link",
				Quantity:    9,
			},
			shouldFail: false,
		},
		{
			payload:    "",
			shouldFail: true,
		},
		{
			payload: map[string]any{
				"name":         "sunshine",
				"dedscription": "bright",
				"price":        "9999.44",
				"image":        "a random one",
				"quantity":     5,
			},
			shouldFail: true,
		},
	}

	testCasesGetProduct := []struct {
		id         any
		shouldFail bool
	}{
		{
			id:         23,
			shouldFail: false,
		},
		{
			id:         "14",
			shouldFail: false,
		},
		{
			id:         "4fh555sd",
			shouldFail: true,
		},
		{
			id:         -23.42,
			shouldFail: true,
		},
		{
			id:         234234234,
			shouldFail: true,
		},
	}

	for _, c := range testCasesGetProduct {
		t.Run(fmt.Sprintf("getting product of %v", c.id), func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("/v1/products/%v", c.id), nil)
			if err != nil {
				t.Errorf("Could not create request: %v", err)
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			var response db.Product
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Errorf(err.Error())
			}
			if c.shouldFail {
				assert.Contains(t, rr.Body.String(), "error")
			} else {
				assert.NotNilf(t, rr.Body, "should contain the product details")
			}
		})
	}

	for i, d := range testCasesInsertProduct {
		t.Run(fmt.Sprintf("creating product case index: %v", i), func(t *testing.T) {
			payloadBytes, _ := json.Marshal(d.payload)
			req, err := http.NewRequest("POST", "/v1/products", bytes.NewBuffer(payloadBytes))
			if err != nil {
				t.Errorf("Could not create request: %v", err)
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			if d.shouldFail {
				assert.Contains(t, rr.Body.String(), "error")
			} else {
				assert.NotNilf(t, rr.Body, "should contain the product details")
			}
		})
	}
}
