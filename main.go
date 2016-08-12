package main

import (
	"fmt"
	"github.com/gtfierro/shellintheghost/client"
	"github.com/gtfierro/shellintheghost/server"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"io"
	"log"
	"os"
	"strings"
)

func startServer(c *cli.Context) error {
	baseURI := c.String("uri")
	if baseURI == "" {
		return fmt.Errorf("Requires base uri")
	}
	terminals := c.StringSlice("terminal")
	if terminals == nil {
		return fmt.Errorf("Requires at least one terminal endpoint")
	}
	shell := c.String("shell")

	client := bw2.ConnectOrExit("")
	vk := client.SetEntityFileOrExit(c.String("entity"))
	client.OverrideAutoChainTo(true)
	server := server.NewServerService(client, vk, baseURI, shell)

	for _, term := range terminals {
		termName := term
		go func() {
			fmt.Printf("Creating terminal: %s\n", termName)
			err := server.AddTerminal(termName)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	x := make(chan bool)
	<-x
	log.Println("EXITING SERVER")
	return nil
}

func startClient(c *cli.Context) error {
	baseURI := c.String("uri")
	if baseURI == "" {
		return fmt.Errorf("Requires base uri")
	}
	termName := c.String("terminal")
	if termName == "" {
		return fmt.Errorf("Requires terminal endpoint")
	}
	cl := bw2.ConnectOrExit("")
	vk := cl.SetEntityFileOrExit(c.String("entity"))
	cl.OverrideAutoChainTo(true)
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)

	shellURI := strings.TrimSuffix(baseURI, "/") + "/s.shell/_/i.term/slot/" + termName
	term, err := client.NewClient(cl, vk, shellURI, os.Stdout)
	if err != nil {
		return err
	}

	shell := terminal.NewTerminal(term, "")

	go io.Copy(shell, os.Stdin)
	<-term.Closed
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "Shell in the Ghost"
	app.Version = "0.2"

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Start SITG server and listen for connections",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uri,u",
					Usage: "Base Service URI",
				},
				cli.StringSliceFlag{
					Name:  "terminal, t",
					Usage: "Name of terminal. Typically one per user.",
				},
				cli.StringFlag{
					Name:   "entity,e",
					EnvVar: "BW2_DEFAULT_ENTITY",
					Usage:  "The entity to use",
				},
				cli.StringFlag{
					Name:  "shell,s",
					Value: "/bin/bash",
					Usage: "Shell to expose on server",
				},
			},
			Action: startServer,
		},
		{
			Name:  "client",
			Usage: "Connect to SITG server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uri,u",
					Usage: "Base Service URI",
				},
				cli.StringFlag{
					Name:  "terminal, t",
					Usage: "Name of terminal",
				},
				cli.StringFlag{
					Name:   "entity,e",
					EnvVar: "BW2_DEFAULT_ENTITY",
					Usage:  "The entity to use",
				},
			},
			Action: startClient,
		},
	}
	app.Run(os.Args)
}
