package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func main() {
	// Get location of the binary
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// When used, resolve symlinks
	binPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		log.Fatal(err)
	}

	// Change the current working directory to the path of the binary
	os.Chdir(filepath.Dir(binPath))

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
