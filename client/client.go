package client

import (
	"fmt"
	"github.com/gtfierro/shellintheghost/conn"
	"golang.org/x/crypto/ssh/terminal"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"io"
	"strings"
)

type Client struct {
	client *bw2.BW2Client
	conn   *conn.Conn
	term   *terminal.Terminal
	output io.Writer
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
func NewClient(client *bw2.BW2Client, vk, termURI string, output io.Writer) (*Client, error) {
	fmt.Println("Client write to", termURI)
	c := &Client{
		client: client,
		output: output,
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
	go io.Copy(c.output, c.conn)
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
