package main

import (
	"github.com/genus555/gator/internal/config"
	"fmt"
)

type state struct {
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
	f, ok := c.commandRegistry[cmd.name];
	if !ok {return fmt.Errorf("Unknown command: %s", cmd.name);}
	return f(s, cmd);
}
func (c *commands) register(name string, f func(*state, command) error) {
	c.commandRegistry[name] = f;
}

func createCommand(sysArgs []string) (command, error) {
	if len(sysArgs) == 0 {return command{}, fmt.Errorf("Not enough arguments");}

	new_command := command{
		name: sysArgs[0],
		args: sysArgs[1:],
	};

	return new_command, nil;
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username detected.");
	}

	username := cmd.args[0];
	err := config.SetUser(*s.cfg_ptr, username);
	if err != nil {return err;}

	fmt.Printf("Username \"%s\" has been set.\n", username);

	return nil;
}