// This package is used by a driver or other bosswave process. It defines a bosswave
// service. Arbitrary terminals can be created; this results in advertising a new slot.
//
// Each slot corresponds to a single terminal and all input on the slot is directed to that
// terminal, so its probably best to not have multiple people writing here. Each connection
// to that slot (which is typically N=1) has "response" URI  on which responses are published.
// This response URI is the slot name appended with ":out
// For example:
//   Read input: /a/b/c/s.shell/0/i.term/slot/terminal01
//   Output published on: /a/b/c/s.shell/0/i.term/signal/terminal01:out
//
// A client sends bytes by publishing on the slot. This is fed into a terminal VT100 emulator
// (from golang std lib) that runs on the server side
package server

import (
	"fmt"
	"github.com/gtfierro/shellintheghost/conn"
	"github.com/kr/pty"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"io"
	"os/exec"
	"syscall"
)

type Server struct {
	client  *bw2.BW2Client
	vk      string
	svc     *bw2.Service
	iface   *bw2.Interface
	command string
}

func NewServerService(client *bw2.BW2Client, vk, uri, command string) *Server {
	svc := client.RegisterService(uri, "s.shell")
	iface := svc.RegisterInterface("_", "i.term")
	return &Server{
		client:  client,
		svc:     svc,
		iface:   iface,
		vk:      vk,
		command: command,
	}
}

func (s *Server) AddTerminal(slotname string) error {
	fmt.Println("Server listens on", s.iface.SlotURI(slotname))
	fmt.Println("Server writes on", s.iface.SignalURI(slotname))
	sub, err := s.client.Subscribe(&bw2.SubscribeParams{
		URI: s.iface.SlotURI(slotname),
	})
	if err != nil {
		return err
	}
	conn := conn.NewConn(s.client, s.vk, s.iface.SignalURI(slotname), sub)

	psty, tty, err := pty.Open()
	if err != nil {
		return err
	}

	go io.Copy(psty, conn)
	go io.Copy(conn, psty)

	c := exec.Command(s.command)
	c.Stdin = tty
	c.Stdout = tty
	c.Stderr = tty
	c.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
	go func() {
		err := c.Start()
		if err != nil {
			fmt.Println("ERR>", err)
			psty.Close()
		}
		err = c.Wait()
		if err != nil {
			//TODO: stop forwarding signals
			//TODO: send the halt signal to the client
			fmt.Println("ERR>", err)
			psty.Close()
		}
		fmt.Println("Exiting server shell")
	}()
	return nil
}
