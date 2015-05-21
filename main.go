package main

import (
	"flag"
	"fmt"

	"github.com/CodeMonk/UdpProxy/udp"
)

var (
	listenPort = flag.Int("port", 8042, "Port to listen on")
	server     = flag.String("server", "192.168.2.3", "Destination to proxy to")
)

func init() {
	flag.Parse()
}

func main() {

	// Read in configuration & flags

	// Start client -> server proxy
	proxy := &udp.UdpProxy{}

	fromClient, fromServer, toClient, toServer := proxy.Initialize(*server,
		*listenPort)

	for {
		select {
		case msg := <-fromClient:
			_ = msg
			fmt.Printf("Received Client Message! from %s Payload:%s\n",
				msg.From, string(msg.Payload))

			// And, send it on
			proxy.Proxy(msg, toServer)
		case msg := <-fromServer:
			_ = msg
			fmt.Printf("Received Server Message! from %s Payload:%s\n",
				msg.From, string(msg.Payload))

			// And, send it on
			proxy.Proxy(msg, toClient)
		}
	}

	// server -> client proxy
	// Start client -> server proxy
}
