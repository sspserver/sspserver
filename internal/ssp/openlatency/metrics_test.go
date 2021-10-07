package openlatency

import (
	"fmt"
	"net"
	"testing"

	"github.com/aeden/traceroute"
)

func Test_TR(t *testing.T) {
	options := traceroute.TracerouteOptions{}
	options.SetRetries(1)
	options.SetMaxHops(traceroute.DEFAULT_MAX_HOPS + 1)
	options.SetFirstHop(traceroute.DEFAULT_FIRST_HOP)

	host := "google.com"
	ipAddr, err := net.ResolveIPAddr("ip", host)
	fmt.Println(">", host, ipAddr, err)
	if err != nil {
		return
	}

	c := make(chan traceroute.TracerouteHop)
	go func() {
		for {
			hop, ok := <-c
			if !ok {
				fmt.Println()
				return
			}
			fmt.Println(hop)
		}
	}()

	_, err = traceroute.Traceroute(host, &options, c)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
