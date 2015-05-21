package udp

import (
	"bytes"
	"fmt"
	"net"
)

type Result struct {
	success bool
	error   string
}

// Create count separate connections to the server, validating the data
func FloodServer(server string, port int, count int, reverseData bool) {

	c := make(chan *Result)

	// make server addr
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		server, port))
	dieErr(err)

	for i := 0; i < count; i++ {
		go testServer(i, saddr, reverseData, c)
	}

	for i := 0; i < count; i++ {
		result := <-c
		// And check our result
		if !result.success {
			fmt.Printf("Error: %s", result.error)
		}
	}

}
func testServer(itteration int, saddr *net.UDPAddr,
	reverseData bool, c chan *Result) {

	message := []byte(fmt.Sprintf("Message Number %d", itteration))
	response, err := sendRecv(saddr, message, 60)

	// and build our Result
	result := &Result{}
	if err != nil {
		result.success = false
		result.error = err.Error()
	} else {
		var dataValid bool
		if reverseData {
			dataValid = bytes.Compare(response, reverse(message)) == 0
		} else {
			dataValid = bytes.Compare(response, message) == 0
		}
		if dataValid {
			result.success = true
			result.error = ""
		} else {
			result.success = false
			result.error = fmt.Sprintf("%d: Data not valid.  Sent <%s> received <%s> reverse = %v",
				itteration, message, response, reverseData)
		}
	}
	c <- result

	return
}
