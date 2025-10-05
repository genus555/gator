package main

import(
	"github.com/genus555/gator/internal/database"
	"context"
)

func LoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg_ptr.CurrentUserName)
		if err != nil {return err}

		return handler(s, cmd, user)
	}
}
