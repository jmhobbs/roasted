package sr700

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

type Roaster struct {
	port  io.ReadWriteCloser
	id    ID
	buf   []byte
	debug bool
}

func New(port io.ReadWriteCloser) *Roaster {
	return &Roaster{
		port:  port,
		id:    DefaultID,
		buf:   make([]byte, PacketLength),
		debug: false,
	}
}

func (s *Roaster) SetDebug(debug bool) {
	s.debug = debug
}

func (s *Roaster) Write(pkt Packet) error {
	if s.debug {
		log.Println("-> ", BytesToHexString(pkt.Bytes()))
	}

	// todo: validation

	n, err := s.port.Write(pkt.Bytes())
	if err != nil {
		return err
	}
	if n != PacketLength {
		return fmt.Errorf("incomplete write, wrote %d, expected %d", n, PacketLength)
	}
	return nil
}

func (s *Roaster) Read() (Packet, error) {
	n, err := s.port.Read(s.buf)

	if err != nil {
		return Packet{}, err
	}

	if n != PacketLength {
		return Packet{}, fmt.Errorf("incomplete read, read %d, expected %d", n, PacketLength)
	}

	if s.debug {
		log.Println("<- ", BytesToHexString(s.buf))
	}

	return ParsePacket(s.buf)
}

func (s *Roaster) Connect() (*Program, error) {
	err := s.Write(Packet{
		Header:      Init,
		Flag:        ControllerSent,
		Control:     Read,
		ID:          s.id,
		Fan:         0x0,
		Timer:       0x0,
		Heat:        0x0,
		Temperature: 0x0,
	})
	if err != nil {
		return nil, err
	}

	// Read current settings
	pkt, err := s.Read()
	if err != nil {
		return nil, err
	}

	// Grab our device "ID" for future communication
	s.id = pkt.ID

	// Read current program
	pgm := Program{}
	for {
		pkt, err = s.Read()
		if err != nil {
			return nil, err
		}

		pgm = append(pgm, State{
			Fan:   pkt.Fan,
			Timer: pkt.Timer,
			Heat:  pkt.Heat,
		})

		if pkt.Flag == TerminalSequenceLine {
			break
		}
	}

	return &pgm, nil
}

func (s *Roaster) Roast(fan Speed, heat Heat, duration time.Duration) (Temperature, error) {
	timer := NewTimerFromDuration(duration)

	err := s.Write(Packet{
		Header:  Normal,
		Flag:    ControllerSent,
		Control: Roast,
		ID:      s.id,
		Fan:     fan,
		Timer:   timer,
		Heat:    heat,
	})

	if err != nil {
		return 0, err
	}

	pkt, err := s.Read()
	if err != nil {
		return 0, err
	}

	if pkt.Fan != fan || pkt.Heat != heat || pkt.Timer != timer {
		return 0, errors.New("roaster responded with incorrect settings")
	}

	return pkt.Temperature, nil
}

func (s *Roaster) Cool(fan Speed, duration time.Duration) (Temperature, error) {
	timer := NewTimerFromDuration(duration)

	err := s.Write(Packet{
		Header:  Normal,
		Flag:    ControllerSent,
		Control: Cooling,
		ID:      s.id,
		Fan:     fan,
		Timer:   timer,
		Heat:    Cool,
	})

	if err != nil {
		return 0, err
	}

	pkt, err := s.Read()
	if err != nil {
		return 0, err
	}

	if pkt.Fan != fan || pkt.Timer != timer {
		return 0, errors.New("roaster responded with incorrect settings")
	}

	return pkt.Temperature, nil
}

func (s *Roaster) Stop() (Temperature, error) {
	err := s.Write(Packet{
		Header:  Normal,
		Flag:    ControllerSent,
		Control: Stop,
		ID:      s.id,
		Fan:     0x01,
		Timer:   0x00,
		Heat:    Cool,
	})

	if err != nil {
		return 0, err
	}

	pkt, err := s.Read()
	if err != nil {
		return 0, err
	}

	return pkt.Temperature, nil
}
