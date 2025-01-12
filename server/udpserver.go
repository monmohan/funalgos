// udp server to respond to any udp packet
package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Create a new UDP server
	conn, err := net.ListenPacket("udp", GetOutboundIP().String()+":33434")
	if err != nil {
		fmt.Println("Failed to create UDP server:", err)
		return
	}
	//print the address of the server
	fmt.Println("Server address:", conn.LocalAddr())
	defer conn.Close()

	//create a listen loop
	for {
		handleClient(conn)
	}

}

func handleClient(conn net.PacketConn) {
	// Read from the UDP connection
	buf := make([]byte, 1500)
	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		fmt.Println("Failed to read from UDP server:", err)
		return
	}
	fmt.Printf("Received %d bytes from %v\n", n, addr)
	//print the data received
	fmt.Printf("Data: %s\n", buf[:n])

	// Send a response back to the client
	if _, err = conn.WriteTo([]byte("Hello from UDP server"), addr); err != nil {
		fmt.Println("Failed to write to Client:", err)
		return
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
