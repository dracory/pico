package migrations

import (
	"testing"

	"project/internal/app"
	"project/internal/config"

	_ "modernc.org/sqlite"
)

func newTestApp(t *testing.T) app.AppInterface {
	t.Helper()
	cfg := config.New()
	cfg.SetAppEnv("testing")
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseName(":memory:")

	a, err := app.New(cfg)
	if err != nil {
		t.Fatalf("app.New failed: %v", err)
	}

	return a
}

func TestMigrateAll_NilApp(t *testing.T) {
	err := MigrateAll(nil)
	if err == nil {
		t.Fatal("expected error for nil app")
	}
}

func TestMigrateAll_Success(t *testing.T) {
	a := newTestApp(t)
	defer a.Close()

	err := MigrateAll(a)
	if err != nil {
		t.Fatalf("MigrateAll failed: %v", err)
	}
}

func TestMigrateAll_Idempotent(t *testing.T) {
	a := newTestApp(t)
	defer a.Close()

	if err := MigrateAll(a); err != nil {
		t.Fatalf("first MigrateAll failed: %v", err)
	}
	if err := MigrateAll(a); err != nil {
		t.Fatalf("second MigrateAll should be idempotent: %v", err)
	}
}

func TestCreateUsersTable_Signature(t *testing.T) {
	m := &CreateUsersTable{}
	if m.Signature() != "2026_07_09_0001_create_users_table" {
		t.Errorf("expected signature '2026_07_09_0001_create_users_table', got %q", m.Signature())
	}
}

func TestCreateUsersTable_Description(t *testing.T) {
	m := &CreateUsersTable{}
	if m.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestCreateUsersTable_UpWithoutSchema(t *testing.T) {
	m := &CreateUsersTable{}
	err := m.Up()
	if err == nil {
		t.Fatal("expected error when schema is nil")
	}
}

func TestCreateUsersTable_DownWithoutSchema(t *testing.T) {
	m := &CreateUsersTable{}
	err := m.Down()
	if err == nil {
		t.Fatal("expected error when schema is nil")
	}
}
