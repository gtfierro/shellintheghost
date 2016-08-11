package main

import (
	"github.com/gtfierro/shellintheghost/server"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"log"
)

func main() {
	client := bw2.ConnectOrExit("")
	vk := client.SetEntityFromEnvironOrExit()
	client.OverrideAutoChainTo(true)

	server := server.NewServerService(client, vk, "gabe.pantry/terminals")
	err := server.AddTerminal("0")
	if err != nil {
		log.Fatal(err)
	}
	x := make(chan bool)
	<-x
}
