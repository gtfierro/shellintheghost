package client

import (
	"bufio"
	"fmt"
	"github.com/gtfierro/shellintheghost/conn"
	"golang.org/x/crypto/ssh/terminal"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"strings"
)

type Client struct {
	client *bw2.BW2Client
	conn   *conn.Conn
	term   *terminal.Terminal
}

// given a URI of form "/a/b/c/key", returns key
func getURIKey(uri string) string {
	li := strings.LastIndex(uri, "/")
	if li > 0 {
		return uri[li+1:]
	}
	return uri
}

// termURI is the slot URI
func NewClient(client *bw2.BW2Client, vk, termURI string, handleOutput func(string)) (*Client, error) {
	fmt.Println("Client write to", termURI)
	c := &Client{
		client: client,
	}
	termName := getURIKey(termURI)
	base := strings.TrimSuffix(termURI, "slot/"+termName)
	fmt.Println("Client subscribe to", base+"signal/"+termName)
	sub, err := client.Subscribe(&bw2.SubscribeParams{
		URI: base + "signal/" + termName,
	})
	if err != nil {
		return nil, err
	}
	c.conn = conn.NewConn(client, vk, termURI, sub)
	go func() {
		r := bufio.NewReader(c.conn)
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				fmt.Println("ERR", err)
				break
			}
			//fmt.Println("client read", line)
			handleOutput(line)
		}
	}()
	//t := &term{c.conn}
	//c.term = terminal.NewTerminal(t, "")
	//go func() {
	//	defer c.conn.Close()
	//	for {
	//		line, err := c.term.ReadLine()
	//		if err != nil {
	//			fmt.Println(err)
	//			break
	//		}
	//		fmt.Println("client read", line)
	//		handleOutput(line)
	//	}
	//}()
	return c, nil
}

func (c *Client) Write(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

type term struct {
	c *conn.Conn
}

func (t *term) Read(b []byte) (int, error) {
	fmt.Println("Read input %s", b)
	return t.c.Write(b)
}

func (t *term) Write(b []byte) (int, error) {
	fmt.Println("Write output %s", b)
	return t.c.Read(b)
}
