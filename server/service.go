package main

import (
	"log"
	"newtonsapple/proto"

	"github.com/reiver/go-telnet"
)

var PhysicsService = newPhysicsService()

func newPhysicsService() *internalPhysicsService {
	return &internalPhysicsService{
		buffer: nil,
	}
}

type internalPhysicsService struct {
	// stuff
	buffer []byte
}

func (s internalPhysicsService) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]

	for {
		n, err := r.Read(p)

		if n > 0 {
			log.Printf("PhysicsService: received %d bytes (%s)", n, string(p[:n]))
			// oi.LongWrite(w, p[:n])
			s.buffer = append(s.buffer, p[:n]...)
			s.checkForMessage()
		}

		if nil != err {
			break
		}
	}
}

func (s *internalPhysicsService) checkForMessage() {
	if len(s.buffer) >= 3 {
		t := proto.MessageType(s.buffer[0])
		size := int(s.buffer[1]) + int(s.buffer[2])*256
		if len(s.buffer) >= size+3 {
			parcel := s.buffer[3 : size+3]
			s.buffer = s.buffer[size+3:]
			s.handleMessage(&proto.ProtocolMessage{Type: t, Size: size, Body: parcel})
		}
	}
}

func (s *internalPhysicsService) handleMessage(msg *proto.ProtocolMessage) error {
	log.Printf("handleMessage: received message '%s' (%d bytes)", msg.Type, len(msg.Body))
	return nil
}
