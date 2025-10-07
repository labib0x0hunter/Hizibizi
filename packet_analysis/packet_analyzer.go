package main

// Collected .cap files from ...
// http.cap : A simple HTTP request and response.
// https://wiki.wireshark.org/SampleCaptures?#hypertext-transport-protocol-http

import (
	"encoding/binary"
	"fmt"
	"strings"
	"syscall"
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
	HeaderLength   uint8   // Header legnth
	DSCP           uint8   //
	ECN            uint8   //
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
	HeaderSize uint8  // Header size
	Window     uint16 // Window size
	Checksum   uint16 // TCP checksum
	UrgentPtr  uint16 // Urgent pointer
}

type udpHeader struct {
	SrcPort  uint16 // Source port
	DstPort  uint16 // Destination port
	Length   uint16 // udp header + data
	Checksum uint16 // checksum
}

type Packet struct {
	pcktHdr    packetHeader   // packet header
	ethrHdr    ethernetHeader // ether header
	ipvHdr     ipv4Header     // ipv header
	protHdr    interface{}    // protocol header
	payload    []byte         // payload
	packetNo   int            // packet number
	payloadLen int            // payload len
}

const (
	pcapHeaderSize     = unsafe.Sizeof(pcapHeader{})
	packetHeaderSize   = unsafe.Sizeof(packetHeader{})
	ethernetHeaderSize = unsafe.Sizeof(ethernetHeader{})
	// ipv4HeaderSize     = unsafe.Sizeof(ipv4Header{})
	// tcpHeaderSize      = unsafe.Sizeof(tcpHeader{})
	// udpHeaderSize      = unsafe.Sizeof(udpHeader{})
)

func Mmap(filename string, length int, prot int, flags int) (data []byte, err error) {
	fd, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer syscall.Close(fd)

	if length == -1 {
		var info syscall.Stat_t
		if err = syscall.Fstat(fd, &info); err != nil {
			return
		}
		length = int(info.Size)
	}

	if length == 0 {
		err = fmt.Errorf("cannot map zero size")
		return
	}

	data, err = syscall.Mmap(
		fd,
		0,
		length,
		prot,
		flags,
	)
	return
}

func Munmap(data []byte) (err error) {
	return syscall.Munmap(data)
}

type File struct {
	data       []byte
	readOffset int64
	res        strings.Builder
	pcktNum    int
}

func (f *File) Open(filename string) (err error) {
	f.data, err = Mmap(filename, -1, syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return
	}
	f.readOffset = 0
	f.pcktNum = 1
	return
}

func (f *File) Close() {
	Munmap(f.data)
}

func (f *File) ReadPcapHeader() pcapHeader {
	if f.readOffset+int64(pcapHeaderSize) >= int64(len(f.data)) {
		// out of index -->>>><<<<-- \\
	}
	hdr := *(*pcapHeader)(unsafe.Pointer(&f.data[f.readOffset]))
	f.readOffset += int64(pcapHeaderSize)
	if hdr.MagicByte != 2712847316 {
		// Not pcap file
	}
	return hdr
}

func (f *File) ReadPacketHeader() packetHeader {
	if f.readOffset+int64(packetHeaderSize) > int64(len(f.data)) {
		// out of index -->>><<<---
	}
	hdr := *(*packetHeader)(unsafe.Pointer(&f.data[f.readOffset]))
	f.readOffset += int64(packetHeaderSize)
	return hdr
}

func (f *File) Analyze() {
	fmt.Println(f.res.String())
}

func (f *File) ReadEthernetHeader() ethernetHeader {
	hdr := *(*ethernetHeader)(unsafe.Pointer(&f.data[f.readOffset]))
	hdr.EtherType = binary.BigEndian.Uint16(f.data[f.readOffset+12 : f.readOffset+14])
	f.readOffset += int64(ethernetHeaderSize)
	return hdr
}

func (f *File) ReadIpv4Header() ipv4Header {
	// prevOffset := f.readOffset
	var hdr ipv4Header
	versionIHL := *(*uint8)(unsafe.Pointer(&f.data[f.readOffset]))
	f.readOffset += int64(unsafe.Sizeof(versionIHL))
	hdr.Version = versionIHL >> 4
	hdr.HeaderLength = 4 * (versionIHL & ((1 << 4) - 1))

	dscpECN := *(*uint8)(unsafe.Pointer(&f.data[f.readOffset]))
	f.readOffset += int64(unsafe.Sizeof(dscpECN))
	hdr.DSCP = dscpECN >> 6
	hdr.ECN = dscpECN & ((1 << 2) - 1)

	hdr.TotalLen = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += int64(unsafe.Sizeof(hdr.TotalLen))

	hdr.Identification = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += int64(unsafe.Sizeof(hdr.Identification))

	hdr.FlagsFragOff = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += int64(unsafe.Sizeof(hdr.FlagsFragOff))

	hdr.TTL = *(*uint8)(unsafe.Pointer(&f.data[f.readOffset]))
	f.readOffset += int64(unsafe.Sizeof(hdr.TTL))

	hdr.Protocol = *(*uint8)(unsafe.Pointer(&f.data[f.readOffset]))
	f.readOffset += int64(unsafe.Sizeof(hdr.Protocol))

	hdr.Checksum = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += int64(unsafe.Sizeof(hdr.Checksum))

	ipSize := unsafe.Sizeof(hdr.SrcIp)
	copy(hdr.SrcIp[:], f.data[f.readOffset:f.readOffset+int64(ipSize)])
	f.readOffset += int64(ipSize)

	copy(hdr.DstIp[:], f.data[f.readOffset:f.readOffset+int64(ipSize)])
	f.readOffset += int64(ipSize)

	// Read option
	optionLen := hdr.HeaderLength - 20
	f.readOffset += int64(optionLen)

	return hdr
}

func (f *File) ReadTCPHeader() tcpHeader {
	// prevOffset := f.readOffset
	hdr := tcpHeader{}
	hdr.SrcPort = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += int64(unsafe.Sizeof(hdr.SrcPort))
	hdr.DstPort = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += int64(unsafe.Sizeof(hdr.DstPort))
	hdr.SeqNum = binary.BigEndian.Uint32(f.data[f.readOffset : f.readOffset+4])
	f.readOffset += int64(unsafe.Sizeof(hdr.SeqNum))
	hdr.AckNum = binary.BigEndian.Uint32(f.data[f.readOffset : f.readOffset+4])
	f.readOffset += int64(unsafe.Sizeof(hdr.AckNum))
	DataOffsetFlag := binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += int64(unsafe.Sizeof(DataOffsetFlag))
	hdr.DataOffset = uint8(DataOffsetFlag >> 12)
	hdr.Flag = DataOffsetFlag & ((1 << 9) - 1)
	hdr.HeaderSize = hdr.DataOffset * 4
	hdr.Window = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += 2
	hdr.Checksum = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += 2
	hdr.UrgentPtr = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += 2

	// Read option
	optionLen := hdr.HeaderSize - 20
	f.readOffset += int64(optionLen)
	return hdr
}

func (f *File) ReadUDPHeader() udpHeader {
	var p udpHeader
	p.SrcPort = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += 2
	p.DstPort = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += 2
	p.Length = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += 2
	p.Checksum = binary.BigEndian.Uint16(f.data[f.readOffset : f.readOffset+2])
	f.readOffset += 2
	return p
}

func (f *File) ReadPayload(n int) (p []byte) {
	// Read payload
	// p = f.data[f.readOffset:f.readOffset+int64(n)]
	f.readOffset += int64(n)
	// f.res.WriteString(fmt.Sprintf("Payload Length: %d\n", n))
	// f.res.WriteString("\n\n")
	return p
}

func (f *File) ReadPacket() Packet {
	var p Packet
	var payloadLen int64
	p.pcktHdr = f.ReadPacketHeader()
	p.ethrHdr = f.ReadEthernetHeader()
	p.ipvHdr = f.ReadIpv4Header()
	// p.tcpHdr = f.ReadTCPHeader()

	payloadLen = int64(p.ipvHdr.TotalLen) - int64(p.ipvHdr.HeaderLength)

	switch p.ipvHdr.Protocol {
	case 1:
		panic("ipvHdrProt")
	case 2:
		panic("ipvHdrProt")
	case 6:
		p.protHdr = f.ReadTCPHeader()
		if v, ok := p.protHdr.(tcpHeader); ok {
			payloadLen -= int64(v.HeaderSize)
		}
	case 17:
		p.protHdr = f.ReadUDPHeader()
		payloadLen -= int64(8)
	default:
	}

	p.payload = f.ReadPayload(int(payloadLen))
	p.packetNo = f.pcktNum
	f.pcktNum++
	p.payloadLen = int(payloadLen)
	return p
}

func (f *File) WritePcapHeader(hdr *pcapHeader) {
	f.res.WriteString("::: pcap header ::: ")
	f.res.WriteString(fmt.Sprintf("Version (%d.%d)", hdr.VersionMajor, hdr.VersionMinor))
	linktype := ""
	switch hdr.LinkType {
	case 1:
		linktype = "Ethernet"
	default:
		linktype = "unknown"
	}
	f.res.WriteString(", " + linktype + "\n\n")
}

func (f *File) WritePacket(p *Packet) {
	// packet header
	f.res.WriteString(fmt.Sprintf("%d |", p.packetNo))
	f.res.WriteString(fmt.Sprintf(" %d |", p.pcktHdr.TimeStampMicroSecond))

	// ether header
	etherType := ""
	switch p.ethrHdr.EtherType {
	case 2048:
		etherType = "Ipv4"
	default:
		etherType = "unknown ether type"
	}
	f.res.WriteString(" " + etherType + "::")

	// ipv4 header
	protocol := ""
	switch p.ipvHdr.Protocol {
	case 1:
		protocol = "ICMP"
	case 2:
		protocol = "IGMP"
	case 6:
		protocol = "TCP"
	case 17:
		protocol = "UDP"
	default:
		protocol = fmt.Sprintf("unknown protocol: %d", p.ipvHdr.Protocol)
	}
	f.res.WriteString(protocol + " |")

	switch v := p.protHdr.(type) {
	case tcpHeader:
		f.res.WriteString(fmt.Sprintf(" %d:%d:%d:%d::%d =>", p.ipvHdr.SrcIp[0], p.ipvHdr.SrcIp[1], p.ipvHdr.SrcIp[2], p.ipvHdr.SrcIp[3], v.SrcPort))
		f.res.WriteString(fmt.Sprintf(" %d:%d:%d:%d::%d", p.ipvHdr.DstIp[0], p.ipvHdr.DstIp[1], p.ipvHdr.DstIp[2], p.ipvHdr.DstIp[3], v.DstPort))

		// Tcp header
		f.res.WriteString(fmt.Sprintf(" | Seq %d", v.SeqNum))
		if v.AckNum > 0 {
			f.res.WriteString(fmt.Sprintf(" | Ack %d", v.AckNum))
		}
		f.res.WriteString(fmt.Sprintf(" | Win %d", v.Window))

		fl := ""
		for i := 0; i <= 8; i++ {
			if ((v.Flag >> i) & 1) == 0 {
				continue
			}
			if len(fl) > 0 {
				fl += "-"
			}
			if len(fl) == 0 {
				fl += "("
			}
			switch i {
			case 0:
				fl += "FIN"
			case 1:
				fl += "SYN"
			case 2:
				fl += "RST"
			case 3:
				fl += "PSH"
			case 4:
				fl += "ACK"
			case 5:
				fl += "URG"
			case 6:
				fl += "ECE"
			case 7:
				fl += "CWR"
			case 8:
				fl += "NS"
			default:
				fl += "unknown flag"
			}
		}
		fl += ")"
		f.res.WriteString(" | " + fl)

	case udpHeader:
		f.res.WriteString(fmt.Sprintf(" %d:%d:%d:%d::%d =>", p.ipvHdr.SrcIp[0], p.ipvHdr.SrcIp[1], p.ipvHdr.SrcIp[2], p.ipvHdr.SrcIp[3], v.SrcPort))
		f.res.WriteString(fmt.Sprintf(" %d:%d:%d:%d::%d", p.ipvHdr.DstIp[0], p.ipvHdr.DstIp[1], p.ipvHdr.DstIp[2], p.ipvHdr.DstIp[3], v.DstPort))
	default:
		// panic("unknown protocol")
		fmt.Println(v)
	}

	f.res.WriteString(fmt.Sprintf(" | Len %d", p.payloadLen))

	f.res.WriteString("\n")
}

func ReadFile(fi *File) {
	hdr := fi.ReadPcapHeader()
	fi.WritePcapHeader(&hdr)
	for fi.readOffset < int64(len(fi.data)) {
		p := fi.ReadPacket()
		fi.WritePacket(&p)
	}
}

func main() {

	// filename := "http.cap"
	filename := "dns.cap"
	var fi File
	fi.Open(filename)
	defer fi.Close()

	ReadFile(&fi)
	fi.Analyze()
}
