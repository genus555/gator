REQUIRMENTS:
    - Postgres
    - Go

Also requires Goose:
    - run "go install github.com/pressly/goose/v3/cmd/goose@latest" in terminal to install Goose
    - run "goose -version" in terminal to check if installation worked correctly

Setting Up Program:
    - Postgres needs a password. Run "sudo passwd {your_password}"
    - In the .gatorconfig.json file you'll see a "db_url"
    - Replace existing URL with your own in this format:
        "postgres://{your_username}:{your_password}@localhost:5432/gator?sslmode=disable"

Commands (gator {command}):
    *login {user_name}: login existing user
	*register {user_name}: register a user
	*reset: completely reset the database
	*users: list all registered users
    *agg {minutes}: scrape through feed/{minutes} and save all posts into database 
	*addfeed {feed_name} {url}: add feed to the database
	*feeds: lists all feeds
	*feed {url}: find feed by url
	*follow {url}: add feed to user follow list
	*following: list all feeds user is following
	*unfollow {url}: unfollow feed by url
	*browse {num_of_entries}: show {num_of_entries} posts sorted by publish date
        ** num_of_entries optional. default 2
    *browsebyfeed {feed_url} {num_of_entries}: show {num_of_entries} posts in feed
        ** num_of_entries optional. default 2