package sr700

import (
	"fmt"
	"time"
)

type Speed uint8

func (s Speed) Valid() bool {
	return s >= 0x01 && s <= 0x09
}

func (s Speed) String() string {
	return fmt.Sprintf("Fan %d", s)
}

type Timer uint8

type ID uint16

const (
	DefaultID ID = 0x6174
)

func (t Timer) Value() time.Duration {
	return time.Second * time.Duration(float64(t)/10.0*60.0)
}

func NewTimerFromDuration(duration time.Duration) Timer {
	return Timer(duration.Seconds() / 60.0 * 10.0)
}

type Temperature uint16

const TemperatureBelow150F = 0xFF00

func (t Temperature) Valid() bool {
	return TemperatureBelow150F != t
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

var flagNames = map[Flag]string{
	ControllerSent:          "Controller Sent",
	RoasterSent:             "Roaster Sent",
	NonTerminalSequenceLine: "Program Line",
	TerminalSequenceLine:    "Final Program Line",
	CurrentSettings:         "Current Settings",
}

func (f Flag) String() string {
	s, ok := flagNames[f]
	if ok {
		return s
	}
	return "Invalid"
}

type Control uint16

const (
	Read    Control = 0x0000
	Idle    Control = 0x0201
	Roast   Control = 0x0402
	Cooling Control = 0x0404
	Stop    Control = 0x0801
)

var controlNames = map[Control]string{
	Read:    "Read",
	Idle:    "Idle",
	Roast:   "Roast",
	Cooling: "Cooling",
	Stop:    "Stop",
}

func (c Control) String() string {
	s, ok := controlNames[c]
	if ok {
		return s
	}
	return "Invalid"
}

type Heat uint8

const (
	Cool   Heat = 0x00
	Low    Heat = 0x01 // 390F
	Medium Heat = 0x02 // 455F
	High   Heat = 0x03 // 490F
)

var heatNames = map[Heat]string{
	Cool:   "Cool",
	Low:    "Low",
	Medium: "Medium",
	High:   "High",
}

func (h Heat) String() string {
	s, ok := heatNames[h]
	if ok {
		return s
	}
	return "Invalid"
}
