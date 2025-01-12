package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
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
	//run probe with TTL 1 to 10
	for i := 1; i <= 1; i++ {
		err = probe(addr, i)
		//print separator
		if err != nil {
			fmt.Println("Failed to probe:", err)
			//return
		}
		fmt.Println("\n--------------------------------------------------")

	}

}

func probe(addr *net.IPAddr, TTL int) error {
	sendSyn("127.0.0.1", addr.String(), 80, TTL)
	icmpconn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("Failed to open ICMP socket:", err)

	}
	defer icmpconn.Close()
	icmpconn.SetDeadline(time.Now().Add(5 * time.Second))

	err = readICMPResponse(icmpconn)

	return err
}

func sendSyn(laddr, raddr string, port uint16, ttl int) {
	conn, err := net.Dial("ip4:tcp", raddr)
	if err != nil {
		log.Fatalf("Dial: %s\n", err)
	}
	defer conn.Close()

	packet := TCPHeader{
		Source:      0xaa47 + uint16(ttl), // Random ephemeral port
		Destination: port,
		SeqNum:      rand.Uint32(),
		AckNum:      0,
		DataOffset:  5,      // 4 bits
		Reserved:    0,      // 3 bits
		ECN:         0,      // 3 bits
		Ctrl:        2,      // 6 bits (000010, SYN bit set)
		Window:      0xaaaa, // The amount of data that it is able to accept in bytes
		Checksum:    0,      // Kernel will set this if it's 0
		Urgent:      0,
		Options:     []TCPOption{},
	}

	data := packet.Marshal()

	packet.Checksum = Csum(data, to4byte(GetOutboundIP().String()), to4byte(raddr))
	//packet.Checksum = Csum(data, to4byte("10.183.158.238"), to4byte(raddr))

	data = packet.Marshal()

	//print the local address
	//fmt.Printf("Local address: %s\n", conn.LocalAddr())
	//set TTL on conn , convert to packet conn
	//print status
	fmt.Printf("Sending SYN to %s:%d with TTL %d\n", raddr, port, ttl)
	rawConn, err := conn.(net.PacketConn).(syscall.Conn).SyscallConn()
	if err != nil {
		log.Fatalf("Failed to get raw connection: %s\n", err)
	}
	rawConn.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TTL, 20)
		if err != nil {
			log.Fatalf("Failed to set socket option: %s\n", err)
		}
	})

	//close the fd in rawconn
	defer rawConn.Control(func(fd uintptr) {
		syscall.Close(int(fd))
	})

	numWrote, err := conn.Write(data)
	if err != nil {
		log.Fatalf("Write: %s\n", err)
	}
	if numWrote != len(data) {
		log.Fatalf("Short write. Wrote %d/%d bytes\n", numWrote, len(data))
	}

}

func readICMPResponse(conn *icmp.PacketConn) error {
	// Read the reply
	reply := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		fmt.Println("Failed to receive ICMP reply:", err)
		return err
	}

	// Parse the reply
	recvMsg, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		fmt.Println("Failed to parse ICMP reply:", err)
		return err
	}

	fmt.Printf("ICMP From Peer %v\n", peer)
	//unmarshal the message body
	switch body := recvMsg.Body.(type) {
	case *icmp.Echo:
		fmt.Printf("Echo reply ID: %d, Seq: %d\n", body.ID, body.Seq)
	case *icmp.TimeExceeded:
		fmt.Printf("ICMP: Time exceeded")
	}

	return nil

}

func to4byte(addr string) [4]byte {
	parts := strings.Split(addr, ".")
	b0, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("to4byte: %s (latency works with IPv4 addresses only, but not IPv6!)\n", err)
	}
	b1, _ := strconv.Atoi(parts[1])
	b2, _ := strconv.Atoi(parts[2])
	b3, _ := strconv.Atoi(parts[3])
	return [4]byte{byte(b0), byte(b1), byte(b2), byte(b3)}
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

const (
	FIN = 1  // 00 0001
	SYN = 2  // 00 0010
	RST = 4  // 00 0100
	PSH = 8  // 00 1000
	ACK = 16 // 01 0000
	URG = 32 // 10 0000
)

type TCPHeader struct {
	Source      uint16
	Destination uint16
	SeqNum      uint32
	AckNum      uint32
	DataOffset  uint8 // 4 bits
	Reserved    uint8 // 3 bits
	ECN         uint8 // 3 bits
	Ctrl        uint8 // 6 bits
	Window      uint16
	Checksum    uint16 // Kernel will set this if it's 0
	Urgent      uint16
	Options     []TCPOption
}

type TCPOption struct {
	Kind   uint8
	Length uint8
	Data   []byte
}

// Parse packet into TCPHeader structure
func NewTCPHeader(data []byte) *TCPHeader {
	var tcp TCPHeader
	r := bytes.NewReader(data)
	binary.Read(r, binary.BigEndian, &tcp.Source)
	binary.Read(r, binary.BigEndian, &tcp.Destination)
	binary.Read(r, binary.BigEndian, &tcp.SeqNum)
	binary.Read(r, binary.BigEndian, &tcp.AckNum)

	var mix uint16
	binary.Read(r, binary.BigEndian, &mix)
	tcp.DataOffset = byte(mix >> 12)  // top 4 bits
	tcp.Reserved = byte(mix >> 9 & 7) // 3 bits
	tcp.ECN = byte(mix >> 6 & 7)      // 3 bits
	tcp.Ctrl = byte(mix & 0x3f)       // bottom 6 bits

	binary.Read(r, binary.BigEndian, &tcp.Window)
	binary.Read(r, binary.BigEndian, &tcp.Checksum)
	binary.Read(r, binary.BigEndian, &tcp.Urgent)

	return &tcp
}

func (tcp *TCPHeader) HasFlag(flagBit byte) bool {
	return tcp.Ctrl&flagBit != 0
}

func (tcp *TCPHeader) Marshal() []byte {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, tcp.Source)
	binary.Write(buf, binary.BigEndian, tcp.Destination)
	binary.Write(buf, binary.BigEndian, tcp.SeqNum)
	binary.Write(buf, binary.BigEndian, tcp.AckNum)

	var mix uint16
	mix = uint16(tcp.DataOffset)<<12 | // top 4 bits
		uint16(tcp.Reserved)<<9 | // 3 bits
		uint16(tcp.ECN)<<6 | // 3 bits
		uint16(tcp.Ctrl) // bottom 6 bits
	binary.Write(buf, binary.BigEndian, mix)

	binary.Write(buf, binary.BigEndian, tcp.Window)
	binary.Write(buf, binary.BigEndian, tcp.Checksum)
	binary.Write(buf, binary.BigEndian, tcp.Urgent)

	for _, option := range tcp.Options {
		binary.Write(buf, binary.BigEndian, option.Kind)
		if option.Length > 1 {
			binary.Write(buf, binary.BigEndian, option.Length)
			binary.Write(buf, binary.BigEndian, option.Data)
		}
	}

	out := buf.Bytes()

	// Pad to min tcp header size, which is 20 bytes (5 32-bit words)
	pad := 20 - len(out)
	for i := 0; i < pad; i++ {
		out = append(out, 0)
	}

	return out
}

// TCP Checksum
func Csum(data []byte, srcip, dstip [4]byte) uint16 {

	pseudoHeader := []byte{
		srcip[0], srcip[1], srcip[2], srcip[3],
		dstip[0], dstip[1], dstip[2], dstip[3],
		0,                  // zero
		6,                  // protocol number (6 == TCP)
		0, byte(len(data)), // TCP length (16 bits), not inc pseudo header
	}

	sumThis := make([]byte, 0, len(pseudoHeader)+len(data))
	sumThis = append(sumThis, pseudoHeader...)
	sumThis = append(sumThis, data...)
	//fmt.Printf("% x\n", sumThis)

	lenSumThis := len(sumThis)
	var nextWord uint16
	var sum uint32
	for i := 0; i+1 < lenSumThis; i += 2 {
		nextWord = uint16(sumThis[i])<<8 | uint16(sumThis[i+1])
		sum += uint32(nextWord)
	}
	if lenSumThis%2 != 0 {
		//fmt.Println("Odd byte")
		sum += uint32(sumThis[len(sumThis)-1])
	}

	// Add back any carry, and any carry from adding the carry
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)

	// Bitwise complement
	return uint16(^sum)
}
