// The type Conn implements the net.Conn interface so that normal go services can be delivered
// over bosswave as though it were a network connection
// "Listening" is done by reading bosswave "binary" messages off of a channel
// "Writing" is done by publishing bosswave "binary" messages on a URI dialed by the client given to the connection
package conn

import (
	"github.com/gtfierro/shellintheghost/ponum"
	bw2 "gopkg.in/immesys/bw2bind.v5"
	"io"
	"time"
)

type Addr struct {
	VK string
}

func (a Addr) Network() string {
	return "bosswave"
}

func (a Addr) String() string {
	return a.VK
}

func NewConn(client *bw2.BW2Client, vk, write string, read chan *bw2.SimpleMessage) *Conn {
	return &Conn{
		client: client,
		write:  write,
		read:   read,
		vk:     vk,
		closed: false,
	}
}

type Conn struct {
	client *bw2.BW2Client
	vk     string
	write  string
	read   chan *bw2.SimpleMessage
	closed bool
}

func (c *Conn) isLogoutMessage(msg *bw2.SimpleMessage) bool {
	return msg.GetOnePODF(ponum.PODFShellLogout) != nil
}

func (c *Conn) Read(p []byte) (n int, err error) {
	if c.closed || c.read == nil {
		return 0, io.EOF
	}
	msg := <-c.read
	if msg == nil {
		return 0, io.EOF
	}
	// unpack message using Blob 1.0.0.0
	// right now just grabs the first 1.0.0.0 PO
	po := msg.GetOnePODF(ponum.PODFShellRaw)
	// get byte contents of PO
	if po == nil {
		if c.isLogoutMessage(msg) {
			return 0, c.Close()
		}
		return 0, nil
	}
	contents := po.GetContents()
	// copy into return slice
	copy(p, contents)
	return len(contents), nil

}

func (c *Conn) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	po, err := bw2.LoadPayloadObject(bw2.FromDotForm(ponum.PODFShellRaw), p)
	if err != nil {
		return 0, err
	}
	err = c.client.Publish(&bw2.PublishParams{
		URI:            c.write,
		AutoChain:      true,
		PayloadObjects: []bw2.PayloadObject{po},
	})

	return len(p), err
}

func (c *Conn) Close() error {
	c.closed = true
	return io.EOF
}

func (c *Conn) Leave() error {
	// send the logout signal
	po, err := bw2.LoadPayloadObject(bw2.FromDotForm(ponum.PODFShellLogout), []byte("bye!"))
	if err != nil {
		return err
	}
	return c.client.Publish(&bw2.PublishParams{
		URI:            c.write,
		AutoChain:      true,
		PayloadObjects: []bw2.PayloadObject{po},
	})
}

func (c *Conn) LocalAddr() Addr {
	return Addr{VK: c.vk}
}

func (c *Conn) RemoteAddr() Addr {
	return c.LocalAddr()
}

// need these for net.conn interface
func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}
func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}
