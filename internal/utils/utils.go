package utils

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bhlox/ecom/internal/db"
	"github.com/go-playground/validator/v10"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

func ParseJSONReq(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func GetAuthHeader(r *http.Request, prefix string) string {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return ""
	}
	tokenHeader := strings.TrimPrefix(authorization, fmt.Sprintf("%v ", prefix))
	return tokenHeader
}

func HashString(text string) (string, error) {
	hashedText, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("something went wrong with hashing text")
	}
	hashedPasswordStr := base64.StdEncoding.EncodeToString(hashedText)
	return hashedPasswordStr, nil
}

func CreateDbTestContainer(t testing.TB, ctx context.Context, insertStatement ...string) (*postgres.PostgresContainer, *sql.DB) {
	t.Helper()

	username := "postgres"
	password := "postgres"
	dbName := "test-db"

	schemaDir := "../../../sql/schema/"

	files, err := os.ReadDir(schemaDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	// Need to sync the sql schemas of the database to the testcontainer
	// 1. Filter out non-SQL files and read their contents
	var sqlScripts []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			content, err := os.ReadFile(filepath.Join(schemaDir, file.Name()))
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			// Extract the SQL queries between -- +goose StatementBegin and -- +goose StatementEnd
			upStart := strings.Index(string(content), "-- +goose StatementBegin")
			downStart := strings.Index(string(content), "-- +goose StatementEnd")
			if upStart != -1 && downStart != -1 {
				upStart += len("-- +goose StatementBegin\n") // Skip the marker itself
				sqlQuery := string(content)[upStart:downStart]
				sqlScripts = append(sqlScripts, sqlQuery)
			}
		}
	}

	// 2. Join all SQL scripts into a single string
	allScripts := strings.Join(sqlScripts, "\n")

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		testcontainers.WithEnv(map[string]string{"TESTCONTAINERS_RYUK_DISABLED": "true"}),
		// postgres.WithInitScripts(filepath.Join("../../..", "testdata", "table_users.sql")),
		// postgres.WithInitScripts("../../../testdata/user/user.sql"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// 3. execute all commands
	_, _, err = pgContainer.Exec(ctx, []string{"psql", "-U", username, "-d", dbName, "-c", allScripts})
	if err != nil {
		t.Fatal(err)
	}
	for _, queries := range insertStatement {
		_, _, err = pgContainer.Exec(ctx, []string{"psql", "-U", username, "-d", dbName, "-c", queries})
		if err != nil {
			t.Fatal(err)
		}
	}

	err = pgContainer.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		t.Fatal(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	database, err := db.InitDB(connStr)
	if err != nil {
		t.Fatal(err)
	}
	return pgContainer, database
}
