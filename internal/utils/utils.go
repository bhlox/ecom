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
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
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

func CreateDbTestContainer(t testing.TB, ctx context.Context, insertStatement ...string) *sql.DB {
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

	// 3. initiate docker test container config
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%v", password),
			fmt.Sprintf("POSTGRES_USER=%v", username),
			fmt.Sprintf("POSTGRES_DB=%v", dbName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://%v:%v@%s/%v?sslmode=disable", username, password, hostAndPort, dbName)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	pool.MaxWait = 120 * time.Second

	var database *sql.DB

	fmt.Println("initiaing connect to DB")
	if err = pool.Retry(func() error {
		database, err = db.InitDB(databaseUrl)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	// 4. execute all commands
	_, err = database.Exec(allScripts)
	if err != nil {
		t.Fatalf("error executing query")
	}
	for _, queries := range insertStatement {
		_, err = database.Exec(queries)
		if err != nil {
			t.Fatal("error executing query")
		}
	}

	return database
}
