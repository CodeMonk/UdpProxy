package main

import (
	"flag"

	"github.com/CodeMonk/UdpProxy/udp"
)

var (
	listenPort = flag.Int("port", 8042, "Port to listen on")
	server     = flag.String("server", "192.168.2.3", "Destination to proxy to")
	fake       = flag.Bool("fakeServer", false, "Start fake udp server, instead of proxy")
)

func init() {
	flag.Parse()
}

func main() {

	// Read in configuration & flags

	// Start client -> server proxy
	proxy := &udp.UdpProxy{}

	if *fake {
		proxy.FakeServer(*listenPort)
	} else {
		proxy.Run(*listenPort, *server)
	}

}
