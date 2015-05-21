package udp

import (
	"fmt"
	"net"
	"strconv"
)

type UdpMessage struct {
	Address *net.UDPAddr
	Payload []byte
}

type UdpProxy struct {
	listenClient *net.UDPAddr
	sendClient   *net.UDPAddr

	clientConn *net.UDPConn
	serverConn *net.UDPConn

	fromClient chan *UdpMessage
	fromServer chan *UdpMessage
	toClient   chan *UdpMessage
	toServer   chan *UdpMessage
}

func dieErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func dumpConn(desc string, conn *net.UDPConn) {
	fmt.Printf("%s: %v (%s->%s)\n", desc, conn,
		conn.LocalAddr(), conn.RemoteAddr())
}
func (u *UdpProxy) dumpConnectionInformation() {
	dumpConn("Client", u.clientConn)
	dumpConn("Server", u.serverConn)
}

func (u *UdpProxy) sender(conn *net.UDPConn, c chan *UdpMessage) {
	for {
		msg := <-c

		var n int
		var err error
		if conn.RemoteAddr() == nil {
			// We don't have an address to respond to, so, send
			// it on to the address in the payload
			fmt.Printf("Sending: FAKED REMOTE:  %d bytes to %s->%s\n", len(msg.Payload),
				conn.LocalAddr(), u.sendClient)
			n, err = conn.WriteToUDP(msg.Payload, u.sendClient)
		} else {
			fmt.Printf("Sending  %d bytes to %s->%s\n", len(msg.Payload),
				conn.LocalAddr(), conn.RemoteAddr())
			n, err = conn.Write(msg.Payload)
		}
		dieErr(err)

		if n != len(msg.Payload) {
			fmt.Printf("Error:  Wrote %d instead of %d bytes!\n", n, len(msg.Payload))
		}
	}
}

func (u *UdpProxy) receiver(con *net.UDPConn, c chan *UdpMessage) {
	buf := make([]byte, 10000)

	for {
		n, addr, err := con.ReadFromUDP(buf)
		dieErr(err)
		if u.sendClient == nil {
			u.sendClient = addr
		}
		c <- &UdpMessage{addr, buf[0:n]}
	}
}

func (u *UdpProxy) startSendersAndReceivers() {

	// Start senders
	go u.sender(u.serverConn, u.toServer)
	go u.sender(u.clientConn, u.toClient)

	// Start receivers
	go u.receiver(u.serverConn, u.fromServer)
	go u.receiver(u.clientConn, u.fromClient)

}

func (u *UdpProxy) startClientReceiver() {
	// Wait for message, then send it on the channel
	var err error

	u.clientConn, err = net.ListenUDP("udp", u.listenClient)
	dieErr(err)

	// now that we have a connection, start senders and receivers
	u.startSendersAndReceivers()

	// There is a real receiver now, so we can exit.
}

func (u *UdpProxy) Initialize(server string, port int) (chan *UdpMessage,
	chan *UdpMessage, chan *UdpMessage, chan *UdpMessage) {

	var err error

	// Setup Listening port -- connection comes from first packet
	u.listenClient, err = net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	dieErr(err)

	// Set up server connection
	sAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", server, port))
	dieErr(err)

	u.serverConn, err = net.DialUDP("udp", nil, sAddr)
	dieErr(err)

	u.fromClient = make(chan *UdpMessage, 10)
	u.fromServer = make(chan *UdpMessage, 10)
	u.toClient = make(chan *UdpMessage, 10)
	u.toServer = make(chan *UdpMessage, 10)

	go u.startClientReceiver()

	// Return the channels
	return u.fromClient, u.fromServer, u.toClient, u.toServer
}

func (u *UdpProxy) Proxy(msg *UdpMessage, c chan *UdpMessage) {
	// Just pop the body into the channel
	c <- msg
}
