package main

import (
	"github.com/gtfierro/shellintheghost/server"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./shellserver <base svc uri>")
	}
	client := bw2.ConnectOrExit("")
	vk := client.SetEntityFromEnvironOrExit()
	client.OverrideAutoChainTo(true)

	server := server.NewServerService(client, vk, os.Args[1], "/bin/bash")
	err := server.AddTerminal("0")
	if err != nil {
		log.Fatal(err)
	}
	x := make(chan bool)
	<-x
	log.Println("EXITING SERVEr")
}
