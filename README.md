Engage
======

Engage is a project that allows you to create your own command line
application where the commands are defined in a config file. This project was
created to allow you to bundle common command line tasks in a config file and
create your own named application.

It changes this:
```json
{
    "name": "engage",
    "usage": "make it so",
    "authors": [
        {
            "name": "Jean Luc Picard",
            "email": "j.l.picard@starfleet.org"
        }
    ],
    "version": "NCC-1701-D",
    "commands": [
        {
            "name": "photon",
            "action": "./photon-torpedoes.sh",
            "usage": "Arm the photon torpedoes"
        },
        {
            "name": "hail",
            "action": "echo",
            "usage": "Hail ship"
        },
        {
            "name": "warp",
            "action": "./engines.sh && ./warp.sh",
            "usage": "Use to jump to warp speed"
        }
    ]
}
```

Into this:
```bash
$ engage

NAME:
   engage - make it so

USAGE:
   main [global options] command [command options] [arguments...]
   
VERSION:
   NCC-1701-D
   
AUTHOR(S):
   Jean Luc Picard <j.l.picard@starfleet.org> 
   
COMMANDS:
    photon      Arm the photon torpedoes
    hail        Hail ship
    warp        Use to jump to warp speed

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
```


Using Engage
------------

Instead of creating your own named command line application, you can use engage
itself.

```bash
$ go get github.com/erroneousboat/engage
$ cd $GOPATH/src/github.com/erroneousboat/engage

# This will install the binary in the bin folder of your $GOPATH, if you want
# to run engage from the command line make sure this folder is on your system
# PATH: export PATH=$GOPATH/bin:$PATH
$ go install
```

Now create your `engage.json` file and put it in the `$GOPATH/bin` folder.
See the *Usage* section on how to create the `engage.json` file.

Create own cli app
------------------

First clone this project, and make sure you have
[The Go Programming Language](https://golang.org/dl/) installed. Then think
about an incredible name you want to use for your application,
and then use it in the following command:

```bash
$ make APP_NAME=[name-of-app]
```

This will place your new application in the `bin/` folder. If you like you can
place this application together with the `engage.json` file somewhere you like.
And if you want to be able to reference from anywhere on your system, then
create a symlink to the application in one of the following locations:

* `$HOME/bin`           (yourself only)
* `/usr/local/bin`      (you and other local users)
* `/usr/local/sbin`     (root only)

Create symlink:

```bash
$ ln -s ~/location/of/app $HOME/bin
```

Usage
-----

Now we need to create our config file `engage.json`, an example:

```json
{
    "name": "engage",
    "usage": "make it so",
    "authors": [
        {
            "name": "Jean Luc Picard",
            "email": "j.l.picard@starfleet.org"
        }
    ],
    "version": "NCC-1701-D",
    "commands": [
        {
            "name": "photon",
            "action": "./photon-torpedoes.sh",
            "usage": "Arm the photon torpedoes"
        },
        {
            "name": "hail",
            "action": "echo",
            "usage": "Hail ship"
        },
        {
            "name": "warp",
            "action": "./engines.sh && ./warp.sh",
            "usage": "Use to jump to warp speed"
        }
    ]
}
```

Save the config file along-side your newly created app, and you'l be able to
issue the following commands:

```bash
$ engage

NAME:
   engage - make it so

USAGE:
   main [global options] command [command options] [arguments...]
   
VERSION:
   NCC-1701-D
   
AUTHOR(S):
   Jean Luc Picard <j.l.picard@starfleet.org> 
   
COMMANDS:
    photon      Arm the photon torpedoes
    hail        Hail ship
    warp        Use to jump to warp speed

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version

# Executing a single command
$ engage photon
Firing photons topedoes...

# With single commands you'll be able to add arguments
$ engage hail klingon ship
klingon ship

# This combines two script calls
$ engage warp
engines operational
warp speed 9, engage ...
```

As you can see in the example file you can combine commands by using
`&&` and `;`.

Todo
----

- [ ] command alias
- [ ] flags
- [ ] subcommands
- [ ] arguments with combined commands
