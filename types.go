package main

import "time"

type Speed uint8

type Timer uint8

func (t Timer) Value() time.Duration {
	return time.Second * time.Duration(float64(t)/10.0*60.0)
}

type Temperature uint16

func (t Temperature) Valid() bool {
	return 0xFF00 != t
}

type Header uint16

const (
	Normal Header = 0xAAAA
	Init   Header = 0xAA55
)

type Flag uint8

const (
	ControllerSent          Flag = 0x63
	RoasterSent             Flag = 0x00
	NonTerminalSequenceLine Flag = 0xAA
	TerminalSequenceLine    Flag = 0xAF
	CurrentSettings         Flag = 0xA0
)

type Control uint16

const (
	Read    Control = 0x0000
	Idle    Control = 0x0201
	Roast   Control = 0x0402
	Cooling Control = 0x0404
	Stop    Control = 0x0801
)

type Heat uint8

const (
	Cool   Heat = 0x00
	Low    Heat = 0x01
	Medium Heat = 0x02
	High   Heat = 0x03
)
