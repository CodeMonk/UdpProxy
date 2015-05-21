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

func (u *UdpProxy) clientReceiver(con *net.UDPConn) {

	buf := make([]byte, 10000)

	// Read our first packet to get connection information
	n, addr, err := con.ReadFromUDP(buf)
	dieErr(err)
	con.Close()
	u.finishInitialization(addr)

	// And send our response onward
	u.fromClient <- &UdpMessage{addr, u.listenClient, buf[0:n]}

	// Finally, start our real loop
	for {
		// Receive Payload

		// Send it on channel
	}

}
func (u *UdpProxy) StartClientReceiver() {
	// Wait for message, then send it on the channel
	con, err := net.ListenUDP("udp", u.listenClient)
	dieErr(err)

	go u.clientReceiver(con)
}

//func StartSenders(server string, port int, clientMessage *UdpMessage,
//	toClient, toServer chan []byte) {
//
//	serverAddr, err := net.ResolveUDPAddr("udp",
//		fmt.Sprintf("%s:%d", server, port))
//	dieErr(err)
//
//}

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

	u.dumpConnectionInformation()
	panic("foo")

	// Start Senders!
	//u.clientConn, err = net.DialUDP("udp", laddr, raddr)
	//dieErr(err)
}
