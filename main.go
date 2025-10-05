package main
import _ "github.com/lib/pq"
import (
	"fmt";
	"github.com/genus555/gator/internal/config"
	"github.com/genus555/gator/internal/database"
	"os"
	"database/sql"
)

var c commands
var dbURL string = "postgres://genus555:postgres@localhost:5432/gator?sslmode=disable"

func init() {
	c.commandRegistry = make(map[string]func(*state, command) error)

	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAggregate)
	c.register("addfeed", LoggedIn(handlerAddFeed))
	c.register("feeds", handlerFeeds)
	c.register("feed", handlerFeed)
	c.register("follow", LoggedIn(handlerFollow))
	c.register("following", LoggedIn(handlerFollowing))
	c.register("unfollow", LoggedIn(handlerUnfollow))
}

func main() {
	cfg, err := config.ReadJson()
	if err != nil {fmt.Println(err)}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {fmt.Println(err)}

	dbQueries := database.New(db)

	current_state := state{db: dbQueries, cfg_ptr: &cfg}

	sysArgs := os.Args[1:]
	current_command, err := createCommand(sysArgs)
	if err != nil {fmt.Println(err)}

	if err := c.run(&current_state, current_command); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}