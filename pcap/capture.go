package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

//convert dev and filter to comman line arguments
//go run capture.go en0 "tcp and port 80"

func main() {
	//make them command line flags
	//read interface name and filter from command
	if len(os.Args) != 3 {
		fmt.Println("Usage: capture [interface] [filter]")
		return
	}

	dev := os.Args[1]
	filter := os.Args[2]

	capture(dev, filter)
}

func capture(dev string, filter string) {

	//handle, err := pcap.OpenLive("lo0", 1600, false, pcap.BlockForever)
	handle, err := pcap.OpenLive(dev, 1600, false, pcap.BlockForever)
	// print what is captured
	fmt.Println("Capturing packets on interface en0")

	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// Set BPF filter
	err = handle.SetBPFFilter(filter)
	fmt.Println("Filter set to", filter)
	if err != nil {
		log.Fatal(err)
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
	}

}
