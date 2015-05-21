package udp

import (
	"fmt"
	"net"
)

func EchoServer(listenPort int, reverseData bool) {
	// Start a fake echo server that will return the received data backward
	// as a response
	conn, err := getUdpListener(listenPort)
	dieErr(err)

	buf := make([]byte, 10000)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		dieErr(err)

		// We have a connection!
		go respond(conn, addr, buf[0:n], reverseData)
	}
}

func respond(conn *net.UDPConn, addr *net.UDPAddr, msg []byte, reverseData bool) {
	// Reverse our string
	if reverseData {
		msg = reverse(msg)
	}
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Printf("Error writing: %s\n", err.Error())
	}
}

func reverse(src []byte) []byte {
	length := len(src)
	dest := make([]byte, length)

	for i := range src {
		length--
		dest[length] = src[i]
	}

	return dest
}
