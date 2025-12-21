package main

import (
	"fmt"
	cfg "github.com/esanchezverges/gator/internal/config"
	"os"
)

var st state
var cmds commands = commands{commands: make(map[string]func(*state, command) error)}

func main() {
	if err := setConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	cmds.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Too few arguments")
		os.Exit(1)
		return
	}
	if err := cmds.run(&st, command{name: args[1], args: args[2:]}); err != nil {
		fmt.Printf("Error running command %v: %v", args[1], err)
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

func handlerLogin(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("Nil reference on state")
	}
	if len(cmd.args) != 1 {
		return fmt.Errorf("Unknown arguments")
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("The user %v has been logged in succesfully\n", cmd.args[0])
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
	config *cfg.Config
}
