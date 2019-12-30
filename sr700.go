package main

import "io"

type SR700 struct {
	port io.ReadWriteCloser
}

func (s *SR700) Connect() error {
	// AA 55 61 74 63 00 00 00 00 00 00 00 AA FA
	/*
		pkt := Packet{
			Header:      Init,
			Flag:        ControllerSent,
			Control:     Read,
			Fan:         0x0,
			Timer:       0x0,
			Heat:        0x0,
			Temperature: 0x0,
		}
	*/

	return nil
}

type State struct {
	Fan  Speed
	Heat Heat
}
