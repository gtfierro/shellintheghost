package main

import (
	"github.com/gtfierro/shellintheghost/client"
	"golang.org/x/crypto/ssh/terminal"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"io"
	"log"
	"os"
)

func main() {
	cl := bw2.ConnectOrExit("")
	vk := cl.SetEntityFromEnvironOrExit()
	cl.OverrideAutoChainTo(true)

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	term, err := client.NewClient(cl, vk, "gabe.pantry/terminals/s.shell/_/i.term/slot/0", os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	shell := terminal.NewTerminal(term, "")

	go io.Copy(shell, os.Stdin)
	<-term.Closed
}
