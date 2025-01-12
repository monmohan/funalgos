package main

import (
	"fmt"
	"net"
)

func main() {
	protocol := "tcp"
	//netaddr, _ := net.ResolveIPAddr("ip4", "127.0.0.1")
	netaddr, _ := net.ResolveIPAddr("ip4", "10.58.182.116")
	conn, _ := net.ListenIP("ip4:"+protocol, netaddr)

	//loop until the port is 80
	for {
		buf := make([]byte, 1024)
		numRead, _, _ := conn.ReadFrom(buf)
		fmt.Println(buf[3])
		//check for byte 14
		if buf[3] == 0x50 {

			fmt.Printf("% X\n", buf[:numRead])
		}
	}

}
