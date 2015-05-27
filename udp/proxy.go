package udp

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type ProxyHandler interface {
	MessageLogger(bool, *net.UDPAddr, []byte)
}

type UdpProxy struct {
	serverAddr *net.UDPAddr
	handlers []ProxyHandler
}


func dieErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func getUdpListener(port int) (*net.UDPConn, error) {

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

func (u *UdpProxy) AddHandler(handler ProxyHandler) {
	u.handlers = append(u.handlers, handler)
}

func (u *UdpProxy) RemoveHandler(handler ProxyHandler) error {
	if len(u.handlers) < 1 {
		return fmt.Errorf("No handlers installed.")
	}
	// grow new array, and copy all but removed into it.
	newHandlers := make([]ProxyHandler, len(u.handlers))
	count := 0

	for _, item := range u.handlers {
		if item != handler {
			newHandlers[count] = item
			count ++
		}
	}

	if count == len(u.handlers) {
		// didn't remove anything
		return fmt.Errorf("Could not find handler %v", handler)
	}

	u.handlers = newHandlers[0:count]

	return nil
}


func (u *UdpProxy) Run(listenPort int, destServer string) {

	// Wait for connections.  For each connection, spawn
	// a routine to send request over to server, and sender
	// server's response to client.

	var err error
	u.serverAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		destServer, listenPort))
	dieErr(err)

	conn, err := getUdpListener(listenPort)
	dieErr(err)

	for {
		buf := make([]byte, 10000)
		n, addr, err := conn.ReadFromUDP(buf)
		dieErr(err)

		// We have a connection!
		go u.doProxy(conn, addr, buf[0:n])
	}

}

func (u *UdpProxy) callHandlers(server bool, addr *net.UDPAddr, data []byte) {
	for _, handler := range(u.handlers) {
		go handler.MessageLogger(server, addr, data)
	}
}

func (u *UdpProxy) doProxy(clientConn *net.UDPConn, src *net.UDPAddr, buf []byte) {

	// Send/Receive message from server

	u.callHandlers(false, src, buf)

	//fmt.Printf("client -> server: %s\n", buf)
	response, err := sendRecv(u.serverAddr, buf, 60)
	if err == nil {
		u.callHandlers(true, u.serverAddr, response)
		_, err = clientConn.WriteToUDP(response, src)
		dieErr(err)
	}
}

func (u *UdpProxy) respondToClient(conn *net.UDPConn,
	addr *net.UDPAddr, msg []byte) {

	_, err := conn.WriteToUDP(msg, addr)
	dieErr(err)
}

func sendRecv(saddr *net.UDPAddr, msg []byte, timeout int) ([]byte, error) {
	// send msg to server, and wait for a response

	conn, err := net.DialUDP("udp", nil, saddr)
	dieErr(err)

	// This should work for the read and the write
	err = conn.SetDeadline(time.Now().Add(time.Second *
		time.Duration(timeout)))
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
