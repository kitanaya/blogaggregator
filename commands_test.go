package main

import (
	"fmt"
	"testing"

	"github.com/kitanaya/blogaggregator/internal/config"
)

func TestCommandsRunLogin(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg := config.Config{DBURL: "postgres://example"}
	s := state{cfg: &cfg}
	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	if err := cmds.run(&s, command{name: "login", args: []string{"alice"}}); err != nil {
		t.Fatal(err)
	}
	stored, err := config.Read()
	if err != nil {
		t.Fatal(err)
	}
	if stored.CurrentUserName != "alice" || stored.DBURL != "postgres://example" || cfg.CurrentUserName != "alice" {
		t.Fatalf("login did not preserve configuration: stored=%+v memory=%+v", stored, cfg)
	}
}

func TestLoginRejectsInvalidArguments(t *testing.T) {
	for _, args := range [][]string{nil, {"alice", "bob"}} {
		t.Run(fmt.Sprint(args), func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())
			cfg := config.Config{DBURL: "postgres://example", CurrentUserName: "seb"}
			defer func() {
				if failure := recover(); failure != nil {
					t.Fatalf("login panicked instead of returning an error: %v", failure)
				}
			}()
			if err := handlerLogin(&state{cfg: &cfg}, command{name: "login", args: args}); err == nil {
				t.Fatal("expected an argument error")
			}
			if cfg.CurrentUserName != "seb" {
				t.Fatal("invalid login changed current user")
			}
		})
	}
}
