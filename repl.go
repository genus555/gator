package main
import _ "github.com/lib/pq"
import (
	"github.com/genus555/gator/internal/config"
	"github.com/genus555/gator/internal/database"
	"fmt"
	"context"
	"time"
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
	user, err := s.db.CreateUser(context.Background(), username)
	if err != nil {return err}

	err = config.SetUser(*s.cfg_ptr, user.Name)
	if err != nil {return err}

	fmt.Printf("User \"%s\" has been created.\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {return err}

	fmt.Println("Successfully reset database")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {return err}

	for _, user := range users {
		if user.Name == s.cfg_ptr.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
		fmt.Printf("* %s\n", user.Name)
	}
	}
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Insert how many minutes between requests")
	}

	time_between_reqs := cmd.args[0]
	tbq, err := time.ParseDuration(time_between_reqs+"m")
	if err != nil {return err}

	fmt.Printf("Collecting feeds every %v minutes\n", time_between_reqs)

	ticker := time.NewTicker(tbq)
	for ; ; <-ticker.C {scrapeFeeds(s)}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Missing URL or Name of feed")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:		cmd.args[0],
		Url:		cmd.args[1],
		UserID:	user.ID,
	})
	if err != nil {return err}

	newFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {return err}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		FeedID:		feed.ID,
		UserID:		user.ID,
	})

	fmt.Println(newFeed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {return err}

	for _, feed := range feeds {
		user, err := s.db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {return err}

		fmt.Printf("%s:\n   URL: \"%s\"\n   Created By: \"%s\"\n",
		feed.Name, feed.Url, user.Name)
	}

	return nil
}

func handlerFeed(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No  URL detected")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {return err}

	fmt.Printf("Feed: %s\nURL: %s\n",feed.Name, feed.Url)
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No feed detected")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {return err}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		FeedID:		feed.ID,
		UserID:		user.ID,
	})
	if err != nil {return err}

	fmt.Printf("Feed: %s\nCurrent User: %s\n", feed.Name, user.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds_followed, err := s.db.GetFeedFollowerForUser(context.Background(), user.ID)
	if err != nil {return err}

	fmt.Printf("Current User: %s\nFollowing:\n", user.Name)
	for _, feed := range feeds_followed {
		f, err := s.db.GetFeedFromFeedID(context.Background(), feed.FeedID)
		if err != nil {return err}

		fmt.Printf("   -\"%s\"\n", f.Name)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No URL detected")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {return err}

	err = s.db.Unfollow(context.Background(), database.UnfollowParams{
		FeedID:		feed.ID,
		UserID:		user.ID,
	})
	if err != nil {return err}

	fmt.Printf("%s has unfollowed \"%s\"\n", user.Name, feed.Name)
	return nil
}