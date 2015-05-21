
package udp;

import (
    "net"
"strconv"
"fmt"
)

type UdpMessage struct {
    From *UDPAddr
    To   *UDPAddr
    Payload      []bytes
}

func dieErr(err error) {
    if err != nil {
        panic(err.Error())
    }
}
func receiver(port int, c chan UdpMessage) {
    // Wait for message, then send it on the channel
    srv,err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
    dieErr(err)
    con, err := net.ListenUDP("udp", srv)
    dieErr(err)
        
    var buf [10000]byte
    for {
        n, addr, err := con.ReadFromUDP(buf)
        fmt.Printf("Read %d bytes from %s\n", n, addr.String)
        
    }
}

