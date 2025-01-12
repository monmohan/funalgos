package packets

import (
	"flag"
	"fmt"
	"testing"
)

// create the flag for command line parameters
var iface string
var filter string

func init() {
	flag.StringVar(&iface, "iface", "en0", "Interface to capture packets")
	flag.StringVar(&filter, "filter", "icmp or udp", "BPF filter for packet capture")

}

// test listDevices
func TestListDevices(t *testing.T) {
	listDevices()
	// Output:
	// Devices found:
	// Name: en0
	// Description: Ethernet
	// Name: lo0
	// Description: Loopback
}

// write test for readIncomingPacket
func TestCaptureTCPPackets(t *testing.T) {
	//kill this after 1 mins

	captureTCPPackets("10.58.182.116", "10.58.182.116", "utun4")
	// Output:
	// <nil>
}

// test printRouteToIP
func TestPrintRouteToIP(t *testing.T) {
	printRouteToIP("18.155.68.105")
	// Output:
	// <nil>
}

// test synWithDial
func TestSynWithDial(t *testing.T) {
	synWithDial("10.58.182.116")

	// Output:
	// <nil>
}

// test capturePackets
func TestCapturePackets(t *testing.T) {
	//capturePackets("tcp and dst port 8181")
	//capture ICMP packets
	//read interface name and filter from command line parameters
	flag.Parse()
	//print iface and filter
	fmt.Printf("Interface: %s, Filter: %s\n", iface, filter)
	// filter for udp and src host

	capturePackets(iface, filter)
	// Output:
	// <nil>
}
