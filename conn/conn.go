// The type Conn implements the net.Conn interface so that normal go services can be delivered
// over bosswave as though it were a network connection
// "Listening" is done by reading bosswave "binary" messages off of a channel
// "Writing" is done by publishing bosswave "binary" messages on a URI dialed by the client given to the connection
package conn

import (
	//"fmt"
	bw2 "gopkg.in/immesys/bw2bind.v5"
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
	}
}

type Conn struct {
	client *bw2.BW2Client
	vk     string
	write  string
	read   chan *bw2.SimpleMessage
}

func (c *Conn) Read(p []byte) (n int, err error) {
	//fmt.Println("reading")
	msg := <-c.read
	//msg.Dump()
	// unpack message using Blob 1.0.0.0
	// right now just grabs the first 1.0.0.0 PO
	po := msg.GetOnePODF(bw2.PODFBlob)
	// get byte contents of PO
	if po == nil {
		return 0, nil
	}
	contents := po.GetContents()
	// copy into return slice
	copy(p, contents)
	//fmt.Printf("READ>", string(p))
	return len(contents), nil

}

func (c *Conn) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	po, err := bw2.LoadPayloadObject(bw2.FromDotForm(bw2.PODFBlob), p)
	if err != nil {
		return 0, err
	}
	//fmt.Printf("WRITE>", string(p))
	//fmt.Println("Writing", string(p), "|")
	err = c.client.Publish(&bw2.PublishParams{
		URI:            c.write,
		AutoChain:      true,
		PayloadObjects: []bw2.PayloadObject{po},
	})

	return len(p), err
}

func (c *Conn) Close() error {
	close(c.read)
	return nil
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
