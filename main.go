package main

import (
	"flag"

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

	proxy.Run(*listenPort, *server)

}
