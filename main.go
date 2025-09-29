package main
import (
	"fmt";
	"github.com/genus555/gator/internal/config"
	"os"
)

var c commands;

func init() {
	c.commandRegistry = make(map[string]func(*state, command) error)

	c.register("login", handlerLogin);
}

func main() {
	cfg, err := config.ReadJson();
	if err != nil {fmt.Println(err);}

	current_state := state{cfg_ptr: &cfg};
	fmt.Printf("Current state: %v\n", current_state);

	sysArgs := os.Args[1:];
	current_command, err := createCommand(sysArgs);
	if err != nil {fmt.Println(err);}

	if err := c.run(&current_state, current_command); err != nil {
		fmt.Println(err);
		os.Exit(1);
	}
}