package udp

import (
	"fmt"
	"net"
	"strconv"
)

type UdpMessage struct {
	From    *net.UDPAddr
	To      *net.UDPAddr
	Payload []byte
}

type UdpProxy struct {
	client       *net.UDPAddr
	server       *net.UDPAddr
	listenClient *net.UDPAddr
	listenServer *net.UDPAddr

	clientConn *net.UDPConn
	serverConn *net.UDPConn

	fromClient chan *UdpMessage
	fromServer chan *UdpMessage

	toClient chan []byte
	toServer chan []byte
}

func dieErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (u *UdpProxy) dumpConnectionInformation() {
	fmt.Printf("Client (from / to): %s / %s\n",
		u.listenClient, u.client)
	fmt.Printf("Server (from / to): %s / %s\n",
		u.listenServer, u.server)
}

func (u *UdpProxy) clientReceiver() {

	buf := make([]byte, 10000)

	// Read our first packet to get connection information
	n, addr, err := u.clientConn.ReadFromUDP(buf)
	dieErr(err)

	// Finish our initialization which will create our client
	u.finishInitialization(addr)

	// And send our response upward
	u.fromClient <- &UdpMessage{addr, u.listenClient, buf[0:n]}

	// There is a real receiver now, so we can exit.
}
func (u *UdpProxy) StartClientReceiver() {
	// Wait for message, then send it on the channel
	var err error

	u.clientConn, err = net.ListenUDP("udp", u.listenClient)
	dieErr(err)

	go u.clientReceiver()
}

func sender(conn *net.UDPConn, c chan []byte) {
	for {
		payload := <-c
		fmt.Printf("About to send %d bytes to %s->%s", len(payload),
			conn.LocalAddr(), conn.RemoteAddr())
		n, err := conn.Write(payload)
		dieErr(err)
		if n != len(payload) {
			fmt.Printf("Error:  Wrote %d instead of %d bytes!\n", n, len(payload))
		}
	}
}

func receiver(con *net.UDPConn, c chan *UdpMessage) {
	buf := make([]byte, 10000)

	for {
		n, err := con.Read(buf)
		dieErr(err)
		c <- &UdpMessage{nil, nil, buf[0:n]}
	}
}

func (u *UdpProxy) StartSendersAndReceivers() {

	// Start senders
	go sender(u.serverConn, u.toServer)
	go sender(u.clientConn, u.toClient)

	// Start receivers
	go receiver(u.serverConn, u.fromServer)
	go receiver(u.clientConn, u.fromClient)

}

func (u *UdpProxy) Initialize(server string, port int) (chan *UdpMessage,
	chan *UdpMessage, chan []byte, chan []byte) {
	var err error
	u.listenClient, err = net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	dieErr(err)

	u.server, err = net.ResolveUDPAddr("udp",
		fmt.Sprintf("%s:%d", server, port))
	dieErr(err)

	u.fromClient = make(chan *UdpMessage, 10)
	u.fromServer = make(chan *UdpMessage, 10)
	u.toClient = make(chan []byte, 10)
	u.toServer = make(chan []byte, 10)

	u.StartClientReceiver()

	// Return the channels?
	return u.fromClient, u.fromServer, u.toClient, u.toServer
}

func (u *UdpProxy) finishInitialization(fromClient *net.UDPAddr) {
	// Build the rest of our stuff, and start server receiver
	var err error
	u.client = fromClient
	u.listenServer, err = net.ResolveUDPAddr("udp",
		fmt.Sprintf(":%d", fromClient.Port))
	dieErr(err)

	// Create our real connections
	u.listenServer.Port++
	fmt.Printf("Dialing ServerCon Local: %s, remote: %s\n", u.listenServer, u.server)
	u.serverConn, err = net.DialUDP("udp", u.listenServer, u.server)
	dieErr(err)

	u.StartSendersAndReceivers()
}

func (u *UdpProxy) Proxy(msg *UdpMessage, c chan []byte) {
	// Just pop the body into the channel
	c <- msg.Payload
}
