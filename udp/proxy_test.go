package udp

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"testing"
)

func TestReverse(t *testing.T) {
	one := []byte("One")
	reversed := reverse(one)

	if bytes.Compare(reversed, []byte("enO")) != 0 {
		t.Errorf("Reverse did not work. %s -> %s\n", one, reversed)
	}

	back := reverse(reversed)
	if bytes.Compare(back, one) != 0 {
		t.Errorf("Reverse did not work. %s -> %s\n", reversed, back)
	}
}

type Message struct {
	Server bool
	Addr *net.UDPAddr
	Data []byte
}
type Handler struct {
	Messages []*Message
}

func (h *Handler) MessageLogger(server bool, addr *net.UDPAddr, data []byte) {
	m := &Message{server, addr, data}
	h.Messages = append(h.Messages, m)
	fmt.Printf("One: Server:%v, data: %s\n", server, data)
}

func TestAddHandler(t *testing.T) {

	proxy := &UdpProxy{}

	h1 := &Handler{}
	h2 := &Handler{}
	
	proxy.AddHandler(h1)

	if len(proxy.handlers) != 1 {
		t.Errorf("Handler count should be 1, not %d\n", len(proxy.handlers))
	}
	proxy.AddHandler(h2)

	if len(proxy.handlers) != 2 {
		t.Errorf("Handler count should be 2, not %d\n", len(proxy.handlers))
	}
}

func TestRemoveHandler(t *testing.T) {

	proxy := &UdpProxy{}

	h1 := &Handler{}
	h2 := &Handler{}
	
	proxy.AddHandler(h1)
	proxy.AddHandler(h2)

	if len(proxy.handlers) != 2 {
		t.Errorf("Handler count should be 2, not %d\n", len(proxy.handlers))
	}

	// remove one
	err := proxy.RemoveHandler(h1)
	if err != nil {
		t.Errorf("Error removing proxy 1: %s", err.Error())
	}
	if len(proxy.handlers) != 1 {
		t.Errorf("Should have one handler left, not %d\n", len(proxy.handlers))
	}

	// remove one again
	err = proxy.RemoveHandler(h1)
	if err == nil {
		t.Error("Should have gotten error trying to re-remove proxy 1\n")
	}

	// remove two
	err = proxy.RemoveHandler(h2)
	if err != nil {
		t.Errorf("Error removing proxy 2: %s", err.Error())
	}
	if len(proxy.handlers) != 0 {
		t.Errorf("Should have no handlers left, not %d\n", len(proxy.handlers))
	}

	// remove something again
	err = proxy.RemoveHandler(h1)
	if err == nil {
		t.Error("Should have returned error removing with no handlers\n")
	}
}

func TestHandler(t *testing.T) {

	proxy := &UdpProxy{}

	h1 := &Handler{}
	h2 := &Handler{}

	proxy.AddHandler(h1)
	proxy.AddHandler(h2)

	clientAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		"1.2.3.4", 567))
	if err != nil {
		t.Error(err.Error())
	}
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		"4.3.2.1", 765))
	if err != nil {
		t.Error(err.Error())
	}

	clientMessage := Message{false, clientAddr, []byte("From Client")}
	serverMessage := Message{true, serverAddr, []byte("From Server")}

	// Add client and server message
	proxy.callHandlers(clientMessage.Server, clientMessage.Addr,
		clientMessage.Data)
	proxy.callHandlers(serverMessage.Server, serverMessage.Addr,
		serverMessage.Data)

	// need to wait for handlers to be called
	time.Sleep(time.Millisecond * 1)

	// Now check the mssages
	if len(h1.Messages) != len(h2.Messages) || len(h2.Messages) != 2 {
		t.Errorf("Should be 2 (not %d:%d) messages!\n", len(h1.Messages),
			len(h2.Messages))
	}
	// if h1.Messages[0] != h2.Messages[0] || h1.Messages[1] != h2.Messages[1] {
	// 	t.Errorf("h1 != h2 messages")
	// }

	// if *h1.Messages[0] != clientMessage {
	// 	t.Errorf("Client message does not match! %v\n", h1.Messages[0])
	// }

	// if *h2.Message[1] != serverMessage {
	// 	t.Errorf("Server message does not match! %v\n", h2.Message[1])
	// }


}

