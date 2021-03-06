package sr700

import (
	"reflect"
	"testing"
)

func TestPacketToBytes(t *testing.T) {
	tests := []struct {
		Packet Packet
		Bytes  []byte
	}{
		// Samples from https://roastero.readthedocs.io/en/stable/data/SR700%20Reversing/README/
		{
			Packet{
				Header:      Init,
				ID:          DefaultID,
				Flag:        ControllerSent,
				Control:     Read,
				Fan:         0x00,
				Timer:       0x00,
				Heat:        0x00,
				Temperature: 0x00,
			},
			[]byte{0xAA, 0x55, 0x61, 0x74, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xAA, 0xFA},
		},
		{
			Packet{
				Header:      Normal,
				ID:          DefaultID,
				Flag:        CurrentSettings,
				Control:     Read,
				Fan:         0x09,
				Timer:       0x3B,
				Heat:        0x02,
				Temperature: 0x00,
			},
			[]byte{0xAA, 0xAA, 0x61, 0x74, 0xA0, 0x00, 0x00, 0x09, 0x3B, 0x02, 0x00, 0x00, 0xAA, 0xFA},
		},
		{
			Packet{
				Header:      Normal,
				ID:          DefaultID,
				Flag:        NonTerminalSequenceLine,
				Control:     Read,
				Fan:         0x09,
				Timer:       0x03,
				Heat:        0x03,
				Temperature: 0x00,
			},
			[]byte{0xAA, 0xAA, 0x61, 0x74, 0xAA, 0x00, 0x00, 0x09, 0x03, 0x03, 0x00, 0x00, 0xAA, 0xFA},
		},
		{
			Packet{
				Header:      Normal,
				ID:          DefaultID,
				Flag:        NonTerminalSequenceLine,
				Control:     Read,
				Fan:         0x09,
				Timer:       0x01,
				Heat:        0x02,
				Temperature: 0x00,
			},
			[]byte{0xAA, 0xAA, 0x61, 0x74, 0xAA, 0x00, 0x00, 0x09, 0x01, 0x02, 0x00, 0x00, 0xAA, 0xFA},
		},
		{
			Packet{
				Header:      Normal,
				ID:          DefaultID,
				Flag:        TerminalSequenceLine,
				Control:     Read,
				Fan:         0x09,
				Timer:       0x1C,
				Heat:        0x00,
				Temperature: 0x00,
			},
			[]byte{0xAA, 0xAA, 0x61, 0x74, 0xAF, 0x00, 0x00, 0x09, 0x1C, 0x00, 0x00, 0x00, 0xAA, 0xFA},
		},
		{
			Packet{
				Header:      Normal,
				ID:          DefaultID,
				Flag:        ControllerSent,
				Control:     Idle,
				Fan:         0x01,
				Timer:       0x3B,
				Heat:        0x01,
				Temperature: 0x00,
			},
			[]byte{0xAA, 0xAA, 0x61, 0x74, 0x63, 0x02, 0x01, 0x01, 0x3B, 0x01, 0x00, 0x00, 0xAA, 0xFA},
		},
	}

	for _, test := range tests {
		actual := test.Packet.Bytes()
		if !reflect.DeepEqual(actual, test.Bytes) {
			t.Errorf("Packet.Bytes() does not match.\nexpected: %X\n  actual: %X\n\n%v", test.Bytes, actual, test.Packet)
		}
	}
}

func Test_ParsePacket(t *testing.T) {
	tests := []struct {
		Bytes  []byte
		Packet Packet
	}{
		{
			[]byte{0xAA, 0xAA, 0x61, 0x74, 0xA0, 0x00, 0x00, 0x09, 0x14, 0x01, 0x00, 0x00, 0xAA, 0xFA},
			Packet{
				Header:      Normal,
				ID:          DefaultID,
				Flag:        CurrentSettings,
				Control:     Read,
				Fan:         0x09,
				Timer:       0x14,
				Heat:        0x01,
				Temperature: 0x00,
			},
		},
		{
			[]byte{0xAA, 0xAA, 0x61, 0x74, 0xA0, 0x00, 0x00, 0x09, 0x14, 0x01, 0x01, 0x02, 0xAA, 0xFA},
			Packet{
				Header:      Normal,
				ID:          DefaultID,
				Flag:        CurrentSettings,
				Control:     Read,
				Fan:         0x09,
				Timer:       0x14,
				Heat:        0x01,
				Temperature: 0x0102,
			},
		},
	}

	for _, test := range tests {
		actual, err := ParsePacket(test.Bytes)
		if err != nil {
			t.Errorf("unexpected parse error for %v: %v", test.Bytes, err)
			continue
		}

		if actual != test.Packet {
			t.Errorf("error parsing packet\nexpected: %v\n  actual: %v", test.Packet, actual)
		}
	}
}

func TestBytesToHexString(t *testing.T) {
	actual := BytesToHexString([]byte{0x00, 0x10, 0x0F, 0x44})
	expected := "0x00 0x10 0x0F 0x44"
	if actual != expected {
		t.Errorf("output did not match\nexpected: %s\n  actual: %s", expected, actual)
	}
}
