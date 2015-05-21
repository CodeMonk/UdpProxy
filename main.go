package main;

import (
    "fmt"
    "flag"

    "github.com/CodeMonky/UdpProxy/udp"
)

var (
    listenPort = flag.Int("port", 8042, "Port to listen on")
)

func init() {
    flag.Parse()
}


func main() {

    // Read in configuration & flags

    // Start client -> server proxy
    clientChannel = make (chan udp.UdpMessage, 10)
    serverChannel = make (chan udp.UdpMessage, 10)

    // After we receive our first message, we can start our 
    firstMessageReceived := false

    go udp.receiver(listenPort, clientChannel)
    // server -> client proxy
    
    // Start client -> server proxy
}
