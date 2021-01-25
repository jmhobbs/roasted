package sr700

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const PacketLength int = 14

type Packet struct {
	Header      Header
	Flag        Flag
	ID          ID
	Control     Control
	Fan         Speed
	Timer       Timer
	Heat        Heat
	Temperature Temperature
}

func (p Packet) String() string {
	var b strings.Builder
	fmt.Fprint(&b, "Packet{\n")
	fmt.Fprintf(&b, "  % 12s: %04X\n", "Header", p.Header)
	fmt.Fprintf(&b, "  % 12s: %04X\n", "ID", p.ID)
	fmt.Fprintf(&b, "  % 12s:   %02X\n", "Flag", p.Flag)
	fmt.Fprintf(&b, "  % 12s: %04X\n", "Control", p.Control)
	fmt.Fprintf(&b, "  % 12s:   %02X\n", "Fan", p.Fan)
	fmt.Fprintf(&b, "  % 12s:   %02X\n", "Timer", p.Timer)
	fmt.Fprintf(&b, "  % 12s:   %02X\n", "Heat", p.Heat)
	fmt.Fprintf(&b, "  % 12s: %04X\n", "Temperature", p.Temperature)
	fmt.Fprintf(&b, "  % 12s: %04X\n", "Footer", 0xAAFA)
	fmt.Fprint(&b, "}")
	return b.String()
}

func (p Packet) Bytes() []byte {
	buf := make([]byte, 14)

	binary.BigEndian.PutUint16(buf[0:], uint16(p.Header))
	binary.BigEndian.PutUint16(buf[2:], uint16(p.ID))
	binary.BigEndian.PutUint16(buf[4:], uint16(p.Flag)<<8)
	binary.BigEndian.PutUint16(buf[5:], uint16(p.Control))
	binary.BigEndian.PutUint16(buf[7:], uint16(p.Fan)<<8|uint16(p.Timer))
	binary.BigEndian.PutUint16(buf[9:], uint16(p.Heat)<<8|uint16(p.Temperature))
	binary.BigEndian.PutUint16(buf[12:], 0xAAFA)

	return buf
}

func ParsePacket(buf []byte) (Packet, error) {
	if len(buf) != 14 {
		return Packet{}, fmt.Errorf("invalid packet length %d, must be 14", len(buf))
	}

	var p Packet

	// todo: check these values are in range?
	p.Header = Header(binary.BigEndian.Uint16(buf[0:2]))
	p.ID = ID(binary.BigEndian.Uint16(buf[2:]))
	p.Flag = Flag(buf[4])
	p.Control = Control(binary.BigEndian.Uint16(buf[5:]))
	p.Fan = Speed(buf[7])
	p.Timer = Timer(buf[8])
	p.Heat = Heat(buf[9])
	p.Temperature = Temperature(binary.BigEndian.Uint16(buf[10:]))

	return p, nil
}

func BytesToHexString(bytes []byte) string {
	s := []string{}
	for _, b := range bytes {
		s = append(s, fmt.Sprintf("0x%02X", b))
	}
	return strings.Join(s, " ")
}
