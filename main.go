package main

import (
	"flag"

	"github.com/CodeMonk/UdpProxy/udp"
)

var (
	listenPort  = flag.Int("port", 8042, "Port to listen on")
	server      = flag.String("server", "192.168.2.3", "Destination to proxy to")
	fakeServer  = flag.Bool("fakeServer", false, "Start fake udp server, instead of proxy")
	fakeClient  = flag.Bool("fakeClient", false, "Start fake udp server, instead of proxy")
	fakeReverse = flag.Bool("fakeReverse", false, "Reverse fake udp server, responses?")
)

func init() {
	flag.Parse()
}

func main() {

	// Read in configuration & flags

	// If we're doing a fake server, just do it and exit
	switch {
	case *fakeServer:
		udp.EchoServer(*listenPort, *fakeReverse)
		return
	case *fakeClient:
		udp.FloodServer(*server, *listenPort, 100, true)
		return
	default:
		// Start client -> server proxy
		proxy := &udp.UdpProxy{}
		proxy.Run(*listenPort, *server)
	}

}
