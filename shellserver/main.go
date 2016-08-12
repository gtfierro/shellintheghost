package main

import (
	"github.com/gtfierro/shellintheghost/server"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: ./shellserver <base svc uri> <path to shell>")
	}
	client := bw2.ConnectOrExit("")
	vk := client.SetEntityFromEnvironOrExit()
	client.OverrideAutoChainTo(true)

	server := server.NewServerService(client, vk, os.Args[1], os.Args[2])
	err := server.AddTerminal("0")
	if err != nil {
		log.Fatal(err)
	}
	x := make(chan bool)
	<-x
	log.Println("EXITING SERVEr")
}
