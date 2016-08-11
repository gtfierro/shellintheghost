package main

import (
	"bufio"
	"fmt"
	"github.com/gtfierro/shellintheghost/client"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"log"
	"os"
)

func printme(s string) {
	//fmt.Println("GOT")
	fmt.Print(s)
}

func main() {
	cl := bw2.ConnectOrExit("")
	vk := cl.SetEntityFromEnvironOrExit()
	cl.OverrideAutoChainTo(true)

	terminal, err := client.NewClient(cl, vk, "gabe.pantry/terminals/s.shell/_/i.term/slot/0", printme)
	if err != nil {
		log.Fatal(err)
	}
	in := bufio.NewReader(os.Stdin)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if err := terminal.Write(append([]byte(line), '\n')); err != nil {
			log.Fatal(err)
		}
	}
}
