package seeders

import (
	"testing"

	"project/database/migrations"
	"project/internal/app"
	"project/internal/config"

	_ "modernc.org/sqlite"
)

func newTestApp(t *testing.T) app.AppInterface {
	t.Helper()
	t.Setenv("APP_HOST", "localhost")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_ENV", "testing")
	t.Setenv("DB_DRIVER", "sqlite")
	t.Setenv("DB_DATABASE", ":memory:")

	cfg, err := config.NewFromEnv()
	if err != nil {
		t.Fatalf("config.NewFromEnv failed: %v", err)
	}

	a, err := app.New(cfg)
	if err != nil {
		t.Fatalf("app.New failed: %v", err)
	}

	return a
}

func TestSeedAll_NilApp(t *testing.T) {
	err := SeedAll(nil)
	if err == nil {
		t.Fatal("expected error for nil app")
	}
}

func TestSeedAll_Success(t *testing.T) {
	a := newTestApp(t)
	defer a.Close()

	if err := migrations.MigrateAll(a); err != nil {
		t.Fatalf("MigrateAll failed: %v", err)
	}

	if err := SeedAll(a); err != nil {
		t.Fatalf("SeedAll failed: %v", err)
	}
}

func TestUserSeeder_Signature(t *testing.T) {
	s := &UserSeeder{}
	if s.Signature() != "2026_07_09_0001_user_seeder" {
		t.Errorf("expected signature '2026_07_09_0001_user_seeder', got %q", s.Signature())
	}
}

func TestUserSeeder_RunNilApp(t *testing.T) {
	s := &UserSeeder{}
	err := s.Run()
	if err == nil {
		t.Fatal("expected error for nil app")
	}
}

func TestUserSeeder_RunSuccess(t *testing.T) {
	a := newTestApp(t)
	defer a.Close()

	if err := migrations.MigrateAll(a); err != nil {
		t.Fatalf("MigrateAll failed: %v", err)
	}

	s := NewUserSeeder(a)
	if err := s.Run(); err != nil {
		t.Fatalf("UserSeeder.Run failed: %v", err)
	}
}

func TestUserSeeder_RunIdempotent(t *testing.T) {
	a := newTestApp(t)
	defer a.Close()

	if err := migrations.MigrateAll(a); err != nil {
		t.Fatalf("MigrateAll failed: %v", err)
	}

	s := NewUserSeeder(a)
	if err := s.Run(); err != nil {
		t.Fatalf("first Run failed: %v", err)
	}
	if err := s.Run(); err != nil {
		t.Fatalf("second Run should be idempotent: %v", err)
	}
}
