package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/kr/pty"
)

// App is the struct that holds the application specification
type App struct {
	Name     string       `json:"name"`
	Usage    string       `json:"usage"`
	Authors  []cli.Author `json:"authors"`
	Version  string       `json:"version"`
	Commands []Command    `json:"commands"`
}

// Command struct consists out of individual command line commands
type Command struct {
	Name   string `json:"name"`
	Usage  string `json:"usage"`
	Action string `json:"action"`
}

// createAction will return a closure that will be used as Action for a
// cli.Command
func (cmd Command) createAction() func(c *cli.Context) {
	action := func(c *cli.Context) {
		// when commands are combined, split on && and ;
		commands := regexp.MustCompile("[&&;]").Split(cmd.Action, -1)

		for _, commandStr := range commands {
			// remove trailing spaces when && or ; is used
			commandStr := strings.TrimSpace(commandStr)

			// split on space in order to get the command and the arguments
			commandArr := strings.Split(commandStr, " ")

			// this allows us to use additional arguments on the commandline
			// arg is defined in config file and cli_arg on the commandline
			// Example: `command arg1 arg2 cli_arg1 cli_arg2`
			command := commandArr[0]
			args := append(commandArr[1:], c.Args()...)

			// execute the command and use a pseudo-terminal
			cmd := exec.Command(command, args...)
			tty, err := pty.Start(cmd)
			if err != nil {
				log.Fatal(err)
			}
			defer tty.Close()

			// pipe tty session to stdout
			go func() {
				io.Copy(os.Stdout, tty)
			}()

			// pipe stdin to tty session
			go func() {
				io.Copy(tty, os.Stdin)
			}()

			err = cmd.Wait()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return action
}

// NewApp is the constructor for the App struct. It will read the config file
// and based on the configuration will construct a cli application with the
// commands specified in the config file.
func NewApp() (App, []cli.Command, error) {
	var app App
	var cliCommands []cli.Command

	file, err := os.Open("engage.json")
	if err != nil {
		return app, cliCommands, err
	}

	if err := json.NewDecoder(file).Decode(&app); err != nil {
		return app, cliCommands, err
	}

	for _, command := range app.Commands {

		cliCommand := cli.Command{
			Name:   command.Name,
			Usage:  command.Usage,
			Action: command.createAction(),
		}

		cliCommands = append(cliCommands, cliCommand)
	}

	return app, cliCommands, nil
}
