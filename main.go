package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	api "github.com/esanchezverges/gator/internal/api"
	cfg "github.com/esanchezverges/gator/internal/config"
	"github.com/esanchezverges/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"html"
	"os"
	"time"
)

var st state
var cmds commands = commands{commands: make(map[string]func(*state, command) error)}

func main() {
	if err := setConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	if err := setDb(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerUsers)
	cmds.register("reset", handlerReset)
	cmds.register("agg", handlerAggregate)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Too few arguments")
		os.Exit(1)
		return
	}

	if err := cmds.run(&st, command{name: args[1], args: args[2:]}); err != nil {
		fmt.Printf("Error running command %v:\n %v\n", args[1], err)
		os.Exit(1)
		return
	}
}

func setConfig() error {
	c, err := cfg.Read()
	if err != nil {
		return fmt.Errorf("Error reading the config: %v", err)
	}
	st.config = &c
	return nil
}

func setDb() error {
	db, err := sql.Open("postgres", st.config.Dburl)
	if err != nil {
		return fmt.Errorf("There was an error opening the db connection: %v", err)
	}
	st.db = database.New(db)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("Nil reference on state")
	}
	if len(cmd.args) != 1 {
		return fmt.Errorf("Unknown arguments")
	}

	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)

	if err != nil {
		return err
	}
	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("The user %v has been logged in succesfully\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("Nil reference on state")
	}
	if len(cmd.args) != 1 {
		return fmt.Errorf("Unknown arguments")
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	newUser, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	fmt.Println("Created user: ", newUser.Name)
	s.config.SetUser(newUser.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("Nil reference on state")
	}
	if len(cmd.args) > 0 {
		return fmt.Errorf("Unknown arguments: %v\n", cmd.args)
	}
	err := s.db.DeleteAll(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Deleted all users from the database")
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("Nil reference on state\n")
	}
	//if len(cmd.args) != 1 {
	//return fmt.Errorf("Unknown arguments: %v\n", cmd.args)
	//}
	rssFeed, err := api.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	if rssFeed == nil {
		return fmt.Errorf("Nil reference on rssFeed\n")
	}
	fmt.Println("Found at:", "https://www.wagslane.dev/index.xml")
	bytes, err := xml.MarshalIndent(*rssFeed, "", "    ")
	if err != nil {
		return fmt.Errorf("Error marshalling the rssfeed: %v\n", err)
	}
	fmt.Println(html.UnescapeString(string(bytes)))
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("Nil reference on state")
	}
	if len(cmd.args) > 0 {
		return fmt.Errorf("Unknown arguments: %v\n", cmd.args)
	}
	usrs, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range usrs {
		if u.Name == s.config.Currentusername {
			fmt.Println("*", u.Name, "(current)")
		} else {
			fmt.Println("*", u.Name)
		}
	}
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	command := c.commands[cmd.name]
	return command(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

type commands struct {
	commands map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

type state struct {
	db     *database.Queries
	config *cfg.Config
}
