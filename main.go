package main

import (
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/kardianos/osext"
)

func main() {
	// Get location of the binary
	binPath, err := osext.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// Change the current working directory to the path of the binary
	os.Chdir(path.Dir(binPath))

	// Initialize new app
	newApp, cliCommands, _ := NewApp()

	app := cli.NewApp()
	app.Name = newApp.Name
	app.Usage = newApp.Usage
	app.Authors = newApp.Authors
	app.Version = newApp.Version
	app.Commands = cliCommands

	app.Run(os.Args)
}
