package main

import (
	"errors"
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

	switch msg.Type {
	case proto.MsgInitSystem:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "mode", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return err
		}
		log.Printf("Arguments: %+v", params)
		// TODO: other types here
	}

	return nil
}

func (s *internalPhysicsService) deserialize(data []byte, args []proto.Argument) (map[string]any, int, error) {
	ptr := 0
	argIndex := 0
	out := map[string]any{}
	for argIndex < len(args) && ptr < len(data) {
		var arg = args[argIndex]
		switch arg.Type {
		case proto.ArgTypeByte:
			if len(data)-ptr >= 1 {
				out[arg.Name] = data[ptr]
				ptr++
			} else {
				return out, -1, errors.New("packet truncated expecting byte")
			}
		case proto.ArgTypeWord:
			if len(data)-ptr >= 2 {
				out[arg.Name] = int(data[argIndex]) + int(data[argIndex+1])*256
				ptr += 2
			} else {
				return out, -1, errors.New("packet truncated expecting word")
			}
		}
	}
	return out, ptr, nil
}
