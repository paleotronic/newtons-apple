package main

import (
	"context"
	"log"
	"time"

	"go.bug.st/serial"
)

type SerialService interface {
	ServeSerial(ctx context.Context, r Reader, w Writer)
}

func ListenAndServeSerial(ctx context.Context, s SerialService, device string, baudRate int, dataBits int, stopBits int, parity string) error {
	p := serial.NoParity
	switch parity {
	case "N":
		p = serial.NoParity
	case "E":
		p = serial.NoParity
	case "O":
		p = serial.OddParity
	case "M":
		p = serial.MarkParity
	case "S":
		p = serial.SpaceParity
	}
	port, err := serial.Open(device, &serial.Mode{
		BaudRate: baudRate,
		DataBits: dataBits,
		Parity:   p,
		StopBits: serial.StopBits(stopBits),
	})
	if err != nil {
		return err
	}

	s.ServeSerial(ctx, port, port)

	return nil
}

func (s internalPhysicsService) ServeSerial(ctx context.Context, r Reader, w Writer) {
	log.Printf("PhysicsService: new SERIAL connection")

	// s.sendWelcome(w)

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]

	var running = true

	var bytesIn = make(chan byte, 4096)
	go func(ch chan byte) {
		for {
			n, err := r.Read(p)

			if n > 0 {
				ch <- p[0]
			}

			if nil != err {
				running = false
				break
			}
		}
		s.pe.Stop()
	}(bytesIn)

	for running {
		select {
		case b := <-bytesIn:
			log.Printf("PhysicsService: received byte ($%x)", b)
			// oi.LongWrite(w, p[:n])
			s.buffer = append(s.buffer, b)
			s.checkForMessage(w)
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}
