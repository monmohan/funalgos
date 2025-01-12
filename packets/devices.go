package packets

import (
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/routing"
)

func listDevices() {
	//list all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// Print device information
	fmt.Println("Devices found:")
	for _, device := range devices {
		//print name and IP addresses
		fmt.Println("Name:", device.Name)
		for _, address := range device.Addresses {
			fmt.Println("IP:", address.IP)

		}
	}

}

func captureTCPPackets(srcIp string, destIp string, iface string) {
	//list all devices
	handle, err := pcap.OpenLive(iface, 6000, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// Set BPF filter for tcp and source ip is srcIp and destination IP as 18.155.68.94
	err = handle.SetBPFFilter("tcp and src host " + srcIp + " or dst host " + destIp)

	if err != nil {
		log.Fatal(err)
	}
	// Packet source
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {

		//Get TCP SYN
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		//if tcpLayer != nil && tcpLayer.(*layers.TCP).SYN {
		if tcpLayer != nil && tcpLayer.(*layers.TCP).SYN {
			//Get ethernet layer
			ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
			if ethernetLayer != nil {
				//print DST MAC address
				fmt.Println("Destination MAC address:", ethernetLayer.(*layers.Ethernet).DstMAC)
			}
			//Get IP layer
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer != nil {
				//print IP address
				fmt.Println("Destination IP address:", ipLayer.(*layers.IPv4).DstIP)

			}
			//print packet
			fmt.Println(packet)
			/*buffer := gopacket.NewSerializeBuffer()
			gopacket.SerializePacket(buffer, gopacket.SerializeOptions{}, packet)
			synWithDial(destIp, buffer.Bytes())*/
		}

	}

}

/*func synWithDial(destIp string, data []byte) {
	conn, err := net.Dial("ip4:tcp", destIp)
	if err != nil {
		log.Fatalf("Dial: %s\n", err)
	}
	defer conn.Close()
	numWrote, err := conn.Write(data)
	if err != nil {
		log.Fatalf("Write: %s\n", err)
	}
	if numWrote != len(data) {
		log.Fatalf("Short write. Wrote %d/%d bytes\n", numWrote, len(data))
	}

}*/

func synWithDial(destIp string) {
	conn, err := net.Dial("ip4:tcp", destIp)
	if err != nil {
		log.Fatalf("Dial: %s\n", err)
	}
	defer conn.Close()
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(65123),
		DstPort: layers.TCPPort(8181),
		Seq:     rand.Uint32(),
		Ack:     0,
		Window:  0xaaaa,
		SYN:     true,
	}
	buffer := gopacket.NewSerializeBuffer()
	tcpLayer.SetNetworkLayerForChecksum(&layers.IPv4{
		SrcIP: net.ParseIP("10.183.158.238"),
		DstIP: net.ParseIP(destIp),
	})
	gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}, tcpLayer)
	//ensure checksum is calculated

	data := buffer.Bytes()

	numWrote, err := conn.Write(data)
	if err != nil {
		log.Fatalf("Write: %s\n", err)
	}
	if numWrote != len(data) {
		log.Fatalf("Short write. Wrote %d/%d bytes\n", numWrote, len(data))
	}

}

/**
* Does not work on MacOS
 */
func printRouteToIP(ip string) {
	router, err := routing.New()
	if err != nil {
		log.Fatal(err)
	}
	//print all routes to 18.155.68.94, convert string to net.IP
	iface, gw, src, err := router.Route(net.ParseIP(ip))

	//print interface, gateway and source IP
	fmt.Println("Interface:", iface)
	fmt.Println("Gateway:", gw)
	fmt.Println("Source IP:", src)

}

func capturePackets(dev string, filter string) {
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
