package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	path := filepath.Join(os.Getenv("HOME"), ".gatorconfig.json")
	if err := os.WriteFile(path, []byte(`{"db_url":"postgres://example","current_user_name":"seb"}`), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Read()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBURL != "postgres://example" || cfg.CurrentUserName != "seb" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestSetUserPersistsAndPreservesDBURL(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg := Config{DBURL: "postgres://example", CurrentUserName: "old"}
	if err := cfg.SetUser("alice"); err != nil {
		t.Fatal(err)
	}
	stored, err := Read()
	if err != nil {
		t.Fatal(err)
	}
	if stored.DBURL != "postgres://example" || stored.CurrentUserName != "alice" || cfg.CurrentUserName != "alice" {
		t.Fatalf("config not preserved: stored=%+v in-memory=%+v", stored, cfg)
	}
}
