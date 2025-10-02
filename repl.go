package main
import _ "github.com/lib/pq"
import (
	"github.com/genus555/gator/internal/config"
	"github.com/genus555/gator/internal/database"
	"github.com/google/uuid"
	"fmt"
	"time"
	"context"
)

type state struct {
	db			*database.Queries
	cfg_ptr		*config.Config
}

type command struct {
	name		string
	args		[]string
}

type commands struct {
	commandRegistry		map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commandRegistry[cmd.name]
	if !ok {return fmt.Errorf("Unknown command: %s", cmd.name)}
	return f(s, cmd)
}
func (c *commands) register(name string, f func(*state, command) error) {
	c.commandRegistry[name] = f
}

func createCommand(sysArgs []string) (command, error) {
	if len(sysArgs) == 0 {return command{}, fmt.Errorf("Not enough arguments")}

	new_command := command{
		name: sysArgs[0],
		args: sysArgs[1:],
	}

	return new_command, nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username detected.")
	}

	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username); if err != nil {
		return err
	}
	err = config.SetUser(*s.cfg_ptr, user.Name)
	if err != nil {return err}

	fmt.Printf("Username \"%s\" has been set.\n", user.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len (cmd.args) == 0 {
		return fmt.Errorf("No name detected")
	}

	username := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		username,
	})
	if err != nil {return err}

	err = config.SetUser(*s.cfg_ptr, user.Name)
	if err != nil {return err}

	fmt.Printf("User \"%s\" has been created.\n", user.Name)

	return nil
}