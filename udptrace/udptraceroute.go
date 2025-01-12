package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/icmp"
)

func main() {
	// Get the IP address to trace
	if len(os.Args) != 2 {
		fmt.Println("Usage: gotrace [ip_address]")
		return
	}
	ipAddr := os.Args[1]

	// Resolve the IP address
	addr, err := net.ResolveIPAddr("ip4", ipAddr)
	if err != nil {
		fmt.Println("Failed to resolve IP address:", err)
		return
	}
	fmt.Println("Resolved IP address:", addr)
	/*icmpconn, err := setUpICMPListener()
	if err != nil {
		fmt.Println("Failed to set up ICMP listener:", err)
		return
	}*/
	go setUpICMPListenerViaCapture("utun4", GetOutboundIP())

	//run probe with TTL 1 to 10
	for i := 1; i <= 15; i++ {
		//print separator
		fmt.Println("\n--------------------------------------------------")
		peer, err := probe(addr, i)
		if err != nil {
			fmt.Println("Failed to probe:", err)

		}
		if peer != nil && peer.String() == addr.String() {
			fmt.Println("Reached destination, hops needed to reach destination:", i)
			break
		}
		//print separator

	}

}

func probe(addr *net.IPAddr, TTL int) (net.Addr, error) {
	laddr := net.UDPAddr{
		IP:   GetOutboundIP(),
		Port: 0,
	}
	fmt.Println("Sending UDP packet to", addr, "from", laddr, "with TTL", TTL)

	conn, err := net.DialUDP("udp", &laddr, &net.UDPAddr{IP: addr.IP, Port: 33434})
	if err != nil {
		fmt.Println("Failed to open UDP connection:", err)
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	rawConn, err := conn.SyscallConn()
	if err != nil {
		fmt.Println("Failed to get raw connection:", err)
		return nil, err
	}
	rawConn.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TTL, TTL)
		if err != nil {
			fmt.Println("Failed to set socket option:", err)

		}
	})

	_, err = conn.Write([]byte("Go Ping"))
	if err != nil {
		fmt.Println("Failed to send UDP message:", err)
		return nil, err
	}
	//fmt.Scanf("Press enter to continue")

	//readUDPResponse(conn)

	/*peer, err := readICMPResponse(icmpconn)
	if err != nil {
		return peer, err
	}*/
	return nil, nil

}

func readUDPResponse(conn *net.UDPConn) (net.Addr, error) {
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		return nil, err
	}

	response := buffer[:n]
	fmt.Println("Received response from server:", string(response))
	return nil, nil

}

func setUpICMPListener() (*icmp.PacketConn, error) {

	// Open a raw socket for ICMP messages
	conn, err := icmp.ListenPacket("ip4:icmp", GetOutboundIP().String())
	if err != nil {
		fmt.Println("Failed to open ICMP socket:", err)
		return nil, err
	}
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	return conn, nil
}

func readICMPResponse(conn *icmp.PacketConn) (net.Addr, error) {
	reply := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return nil, fmt.Errorf("no ICMP reply: %v, %v", err, peer)

	}

	packet := gopacket.NewPacket(reply[:n], layers.LayerTypeICMPv4, gopacket.Default)
	icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
	if icmpLayer == nil {
		return nil, fmt.Errorf("failed to parse ICMP reply")

	}
	icmpPacket, _ := icmpLayer.(*layers.ICMPv4)
	fmt.Printf("Reply from %v: ", peer)

	fmt.Printf("Type: %v, Code: %v, ID: %v, Seq: %v, ", icmpPacket.TypeCode.Type(), icmpPacket.TypeCode.Code(), icmpPacket.Id, icmpPacket.Seq)

	switch icmpPacket.TypeCode.Type() {
	case layers.ICMPv4TypeEchoReply:
		fmt.Println("Echo reply")
	case layers.ICMPv4TypeDestinationUnreachable:
		fmt.Println("Destination unreachable")
	case layers.ICMPv4TypeTimeExceeded:
		fmt.Println("Time exceeded")
	default:
		fmt.Println("Unknown")
	}
	return peer, nil

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

func setUpICMPListenerViaCapture(iface string, destIp net.IP) {

	ifc, e := findInterface(GetOutboundIP())
	if e != nil {
		fmt.Println("Failed to find interface:", e)
		return
	}
	fmt.Println("Found interface:", ifc.Name)
	handle, err := pcap.OpenLive(ifc.Name, 1600, false, pcap.BlockForever)
	// print what is captured
	fmt.Println("Capturing on interface", iface)

	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// Set BPF filter
	err = handle.SetBPFFilter("icmp and dst host " + destIp.String())
	fmt.Println("Filter set for ICMP packets to destination IP", destIp)

	if err != nil {
		log.Fatal(err)
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
	}

}

func findInterface(ip net.IP) (*net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.Equal(ip) {
					return &iface, nil
				}
			case *net.IPAddr:
				if v.IP.Equal(ip) {
					return &iface, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no interface found for IP %s", ip)
}
