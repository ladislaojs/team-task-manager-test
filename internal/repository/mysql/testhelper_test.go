package mysqlrepo_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	container, err := mysql.Run(ctx,
		"mysql:8.4",
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
	)
	if err != nil {
		t.Fatalf("start mysql container: %v", err)
	}
	t.Cleanup(func() { container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "parseTime=true&loc=UTC")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := runMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	return db
}

func runMigrations(db *sql.DB) error {
	_, file, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(file), "..", "..", "..", "migrations")

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	var upFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, filepath.Join(migrationsDir, e.Name()))
		}
	}
	sort.Strings(upFiles)

	for _, path := range upFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, stmt := range strings.Split(string(content), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				return err
			}
		}
	}

	return nil
}
