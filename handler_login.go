package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login requires exactly one username")
	}
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	_, err := fmt.Printf("User has been set to %s\n", cmd.args[0])
	return err
}
