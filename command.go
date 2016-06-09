package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

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
		// when commands are combined, split on `;`
		commands := regexp.MustCompile(`[;]`).Split(cmd.Action, -1)

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
			cmd := exec.Command(command, args...)

			// puts the terminal connected into raw mode, data is given as-is
			// to the program (pty), and the system does not interpret any of
			// the special characters
			// - https://en.wikipedia.org/wiki/Cooked_mode
			// - http://stackoverflow.com/a/13104579/1346257
			fd := os.Stdin.Fd()
			oldState, err := terminal.MakeRaw(int(fd))
			if err != nil {
				log.Fatal(err)
			}
			defer terminal.Restore(int(fd), oldState)

			// execute the command and use a pseudo-terminal
			tty, err := pty.Start(cmd)
			if err != nil {
				log.Println("EHLLO")
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
				// restore terminal to oldState when interrupted
				terminal.Restore(int(fd), oldState)
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
