package udp

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type UdpMessage struct {
	Address *net.UDPAddr
	Payload []byte
}

type UdpProxy struct {
	serverAddr *net.UDPAddr
}

func dieErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func (u *UdpProxy) FakeServer(listenPort int) {
	// Start a fake echo server that will return the received data backward
	// as a response
	conn, err := u.getUdpListener(listenPort)
	dieErr(err)

	buf := make([]byte, 10000)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		dieErr(err)

		// We have a connection!
		go u.fakeRespond(conn, addr, buf[0:n])
	}
}

func reverse(src []byte) []byte {
	length := len(src)
	dest := make([]byte, length)

	for i := range src {
		dest[length-i] = src[i]
	}

	return dest
}

func (u *UdpProxy) fakeRespond(conn *net.UDPConn, addr *net.UDPAddr, msg []byte) {
	// Reverse our string
	reversed := reverse(msg)
	_, err := conn.WriteToUDP(reversed, addr)
	if err != nil {
		fmt.Printf("Error writing: %s\n", err.Error())
	}
}

func (u *UdpProxy) getUdpListener(port int) (*net.UDPConn, error) {

	// Setup Listening connection
	caddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", caddr)
	if err != nil {
		return nil, err
	}

	return conn, err

}

func (u *UdpProxy) Run(listenPort int, destServer string) {

	// Wait for connections.  For each connection, spawn
	// a routine to send request over to server, and sender
	// server's response to client.

	var err error
	u.serverAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		destServer, listenPort))
	dieErr(err)

	conn, err := u.getUdpListener(listenPort)
	dieErr(err)

	for {
		buf := make([]byte, 10000)
		n, addr, err := conn.ReadFromUDP(buf)
		dieErr(err)

		// We have a connection!
		go u.doProxy(conn, addr, buf[0:n])
	}

}

func (u *UdpProxy) doProxy(clientConn *net.UDPConn, src *net.UDPAddr, buf []byte) {

	// Send/Receive message from server

	// LogHook (client -> server == buf)
	fmt.Printf("client -> server: %s\n", buf)
	response, err := u.sendToServer(buf)
	if err == nil {
		// LogHook (server -> client == response)
		fmt.Printf("server -> client: %s\n", response)
		_, err = clientConn.WriteToUDP(response, src)
		dieErr(err)
	}
}

func (u *UdpProxy) respondToClient(conn *net.UDPConn,
	addr *net.UDPAddr, msg []byte) {

	_, err := conn.WriteToUDP(msg, addr)
	dieErr(err)
}

func (u *UdpProxy) sendToServer(msg []byte) ([]byte, error) {
	// send msg to server, and wait for a response

	conn, err := net.DialUDP("udp", nil, u.serverAddr)
	dieErr(err)

	// This should work for the read and the write
	err = conn.SetDeadline(time.Now().Add(time.Second * 60))
	dieErr(err)

	_, err = conn.Write(msg)
	if err != nil {
		fmt.Printf("Error writing: %s\n", err.Error())
		return nil, err
	}

	// And, wait for a response
	buf := make([]byte, 10000)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading: %s\n", err.Error())
		return nil, err
	}
	conn.Close()

	// return our data
	return buf[0:n], nil
}
