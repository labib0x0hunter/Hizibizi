package main

// Collected .cap files from ...
// http.cap : A simple HTTP request and response.
// https://wiki.wireshark.org/SampleCaptures?#hypertext-transport-protocol-http

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

type pcapHeader struct {
	MagicByte    uint32 // Magic-byte = 0xa1b2c3d4 = 2712847316
	VersionMajor uint16 // File format version
	VersionMinor uint16 // File format version
	TimeZone     int32  // TimeZone offset, normally 0
	TimeStamp    uint32 // TimStamp accuracy, normally 0
	SnapLen      uint32 // Max capture length per packet
	LinkType     uint32 // Link-Layer type
}

type packetHeader struct {
	TimeStampSecond      uint32 // Timestamp (second)
	TimeStampMicroSecond uint32 // Timestamp (microsecond)
	InclLen              uint32 // Number of bytes in the file
	OrigLen              uint32 // Original length of the packet
}

type ethernetHeader struct {
	DstMac    [6]byte // Destination MAC-Address
	SrcMac    [6]byte // Source MAC-Address
	EtherType uint16  // Ether type , Big-Endian encoding must be
}

type ipv4Header struct {
	Version        uint8   // Version -
	DSCP_ECN       uint8   //
	TotalLen       uint16  // Total packet length = header + payload, must be big-endian
	Identification uint16  // Unique id for fragmentation, must be big-endian
	FlagsFragOff   uint16  // Fragmentation Flags and offset, must be big-endian
	TTL            uint8   // Pocket lifetime
	Protocol       uint8   // Next layer protocal (TCP/UDP/ICMP)
	Checksum       uint16  // Header error checking, must be big-endian
	SrcIp          [4]byte // Source ip address 4 chunks
	DstIp          [4]byte // Destination ip address 4 chunks
}

type tcpHeader struct {
	SrcPort    uint16 // Source port
	DstPort    uint16 // Destination port
	SeqNum     uint32 // Sequence number
	AckNum     uint32 // Acknowledge number
	DataOffset uint8  // Dataoffset +
	Flag       uint16 //  Flag bits
	Window     uint16 // Window size
	Checksum   uint16 // TCP checksum
	UrgentPtr  uint16 // Urgent pointer
	option     []byte
}

type Packet struct {
	pcktHdr packetHeader
	ethrHdr ethernetHeader
	ipvHdr  ipv4Header
	tcpHdr  tcpHeader
	payload []byte
}

const (
	pcapHeaderSize     = unsafe.Sizeof(pcapHeader{})
	packetHeaderSize   = unsafe.Sizeof(packetHeader{})
	ethernetHeaderSize = unsafe.Sizeof(ethernetHeader{})
	ipv4HeaderSize     = unsafe.Sizeof(ipv4Header{})
	tcpHeaderSize      = unsafe.Sizeof(tcpHeader{})
)

type File struct {
	file *os.File
}

func (f *File) Open(filename string) {
	fl, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	f.file = fl
}

func (f *File) Close() {
	f.file.Close()
}

func (f *File) ReadPacket() Packet {
	var p Packet
	p.pcktHdr = f.ReadPacketHeader()
	p.ethrHdr = f.ReadEthernetHeader()
	p.ipvHdr = f.ReadIpv4Header()
	p.tcpHdr = f.ReadTCPHeader()
	payloadLen := p.ipvHdr.TotalLen - uint16(ipv4HeaderSize) - uint16(p.tcpHdr.DataOffset*4)
	p.payload = make([]byte, payloadLen)
	f.file.Read(p.payload)
	return p
}

func (f *File) Analyze(p *Packet) {
	AnalyzePacketHeader(&p.pcktHdr)
	AnalyzeEthernetHeader(&p.ethrHdr)
	AnalyzeIpv4Header(&p.ipvHdr)
	AnalyzeTCPHeader(&p.tcpHdr)
	// fmt.Println("Payload: " + string(p.payload))
	fmt.Println("Payload len: ", len(p.payload))
	fmt.Printf("\n\n\n")
}

func (f *File) ReadPcapHeader() pcapHeader {
	buf := make([]byte, pcapHeaderSize)
	f.file.Read(buf)
	pcapHdr := *(*pcapHeader)(unsafe.Pointer(&buf[0]))
	return pcapHdr
}

func (f *File) ReadPacketHeader() packetHeader {
	buf1 := make([]byte, packetHeaderSize)
	f.file.Read(buf1)
	hdr1 := *(*packetHeader)(unsafe.Pointer(&buf1[0]))
	return hdr1
}

func (f *File) ReadEthernetHeader() ethernetHeader {
	buf2 := make([]byte, ethernetHeaderSize)
	f.file.Read(buf2)
	hdr2 := *(*ethernetHeader)(unsafe.Pointer(&buf2[0]))
	hdr2.EtherType = binary.BigEndian.Uint16(buf2[12:14])
	return hdr2
}

func (f *File) ReadIpv4Header() ipv4Header {
	buf3 := make([]byte, ipv4HeaderSize)
	f.file.Read(buf3)
	hd3 := *(*ipv4Header)(unsafe.Pointer(&buf3[0]))
	hd3.TotalLen = binary.BigEndian.Uint16(buf3[2:4])
	hd3.Identification = binary.BigEndian.Uint16(buf3[4:6])
	hd3.FlagsFragOff = binary.BigEndian.Uint16(buf3[6:8])
	hd3.Checksum = binary.BigEndian.Uint16(buf3[10:12])
	return hd3
}

func (f *File) ReadTCPHeader() tcpHeader {
	buf4 := make([]byte, tcpHeaderSize)
	f.file.Read(buf4)
	hd4 := tcpHeader{}
	hd4.SrcPort = binary.BigEndian.Uint16(buf4[0:2])
	hd4.DstPort = binary.BigEndian.Uint16(buf4[2:4])
	hd4.SeqNum = binary.BigEndian.Uint32(buf4[4:8])
	hd4.AckNum = binary.BigEndian.Uint32(buf4[8:12])
	DataOffsetFlag := binary.BigEndian.Uint16(buf4[12:14])
	hd4.DataOffset = uint8(DataOffsetFlag >> 12)
	hd4.Flag = DataOffsetFlag & ((1 << 9) - 1)
	hd4.Window = binary.BigEndian.Uint16(buf4[14:16])
	hd4.Checksum = binary.BigEndian.Uint16(buf4[16:18])
	hd4.UrgentPtr = binary.BigEndian.Uint16(buf4[18:20])

	if hd4.DataOffset*4 > 20 {
		hd4.option = make([]byte, hd4.DataOffset*4-20)
		f.file.Read(hd4.option)
	}

	return hd4
}

func AnalyzePcapHeader(hdr *pcapHeader) {
	var res strings.Builder
	res.WriteString("::: pcap header :::\n")
	res.WriteString(fmt.Sprintf("Magic Byte: 0x%x\n", hdr.MagicByte))
	res.WriteString(fmt.Sprintf("Version   : %d.%d\n", hdr.VersionMajor, hdr.VersionMinor))
	res.WriteString(fmt.Sprintf("Timezone. : %d\n", hdr.TimeZone))
	res.WriteString(fmt.Sprintf("Timestamp : %d\n", hdr.TimeStamp))
	res.WriteString(fmt.Sprintf("Max capture len : %d\n", hdr.SnapLen))
	linktype := ""
	switch hdr.LinkType {
	case 1:
		linktype = "1 (Ethernet)"
	default:
		linktype = "unknown"
	}
	res.WriteString("Link type : " + linktype + "\n")
	fmt.Println(res.String() + "\n")
}

func AnalyzePacketHeader(hdr *packetHeader) {
	var res strings.Builder
	res.WriteString("::: packet header :::\n")
	res.WriteString(fmt.Sprintf("Timestamp: %d (second)\n", hdr.TimeStampSecond))
	res.WriteString(fmt.Sprintf("Timestamp: %d (microsecond)\n", hdr.TimeStampMicroSecond))
	res.WriteString(fmt.Sprintf("File length: %d bytes\n", hdr.InclLen))
	res.WriteString(fmt.Sprintf("Packet length: %d bytes\n", hdr.OrigLen))
	fmt.Println(res.String() + "\n")
}

func AnalyzeEthernetHeader(hdr *ethernetHeader) {
	var res strings.Builder
	res.WriteString("::: ethernet header :::\n")
	res.WriteString(fmt.Sprintf("Source Mac : %02x:%02x:%02x:%02x:%02x:%02x\n", hdr.SrcMac[0], hdr.SrcMac[1], hdr.SrcMac[2], hdr.SrcMac[3], hdr.SrcMac[4], hdr.SrcMac[5]))
	res.WriteString(fmt.Sprintf("Destination Mac : %02x:%02x:%02x:%02x:%02x:%02x\n", hdr.DstMac[0], hdr.DstMac[1], hdr.DstMac[2], hdr.DstMac[3], hdr.DstMac[4], hdr.DstMac[5]))
	etherType := ""
	switch hdr.EtherType {
	case 2048:
		etherType = "2048 (Ipv4)"
	default:
		etherType = "unknown ether type"
	}
	res.WriteString("Ether Type: " + etherType + "\n")
	fmt.Println(res.String() + "\n")
}

func AnalyzeIpv4Header(hdr *ipv4Header) {
	var res strings.Builder
	res.WriteString("::: ipv4 header :::\n")
	// v := hdr.Version >> 4              // left-shift 4bytes (first 4bytes)
	hl := hdr.Version & ((1 << 4) - 1) // (last 4bytes)
	res.WriteString(fmt.Sprintf("Header length: %d bytes\n", hl * 4))
	res.WriteString(fmt.Sprintf("Total packet length: %d bytes\n", hdr.TotalLen))
	protocol := ""
	switch hdr.Protocol {
	case 6:
		protocol = "6 (tcp)"
	default:
		protocol = "unknown protocol"
	}
	res.WriteString("Protocol : " + protocol + "\n")
	res.WriteString(fmt.Sprintf("Source Ip: %d:%d:%d:%d\n", hdr.SrcIp[0], hdr.SrcIp[1], hdr.SrcIp[2], hdr.SrcIp[3]))
	res.WriteString(fmt.Sprintf("Destination Ip: %d:%d:%d:%d\n", hdr.DstIp[0], hdr.DstIp[1], hdr.DstIp[2], hdr.DstIp[3]))
	fmt.Println(res.String() + "\n")
}

func AnalyzeTCPHeader(hdr *tcpHeader) {
	var res strings.Builder
	res.WriteString("::: TCP header :::\n")
	res.WriteString(fmt.Sprintf("Destination port: %d\n", hdr.DstPort))
	res.WriteString(fmt.Sprintf("Source port: %d\n", hdr.SrcPort))
	res.WriteString(fmt.Sprintf("Sequence number: %d\n", hdr.SeqNum))
	res.WriteString(fmt.Sprintf("Window size: %d\n", hdr.Window))
	res.WriteString(fmt.Sprintf("Dataoffset: %d\n", hdr.DataOffset))
	res.WriteString(fmt.Sprintf("Header length: %d\n", hdr.DataOffset*4))
	fl := ""
	switch hdr.Flag {
	case 2:
		fl = "2 (SYN)"
	case 16:
		fl = "16 (ACK)"
	default:
		fl = "unknown flag"
	}
	res.WriteString("Flag: " + fl + "\n")
	if hdr.DataOffset*4 > 20 {
		// fmt.Printf("TCP Options: % X\n", hdr.option)
	}
	fmt.Println(res.String() + "\n")
}

func main() {

	filename := "http.cap"
	var fi File
	fi.Open(filename)
	defer fi.Close()

	/* pcap header */
	pcapHdr := fi.ReadPcapHeader()
	AnalyzePcapHeader(&pcapHdr)

	pckt := fi.ReadPacket()
	fi.Analyze(&pckt)

	pckt = fi.ReadPacket()
	fi.Analyze(&pckt)

	// /* packet header */
	// pcktHdr := fi.ReadPacketHeader()
	// AnalyzePacketHeader(&pcktHdr)

	// /* ethernet header */
	// ethrHdr := fi.ReadEthernetHeader()
	// AnalyzeEthernetHeader(&ethrHdr)

	// /* Ipv4 header */
	// ipv4Hdr := fi.ReadIpv4Header()
	// AnalyzeIpv4Header(&ipv4Hdr)

	// // /* tcp header */
	// tcpHdr := fi.ReadTCPHeader()
	// AnalyzeTCPHeader(&tcpHdr)

}
