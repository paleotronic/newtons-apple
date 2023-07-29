package main

import (
	"errors"
	"fmt"
	"log"
	"newtonsapple/proto"
	"time"

	"github.com/jakecoffman/cp"
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

var PhysicsService = newPhysicsService()

func (s *internalPhysicsService) reportData(deltas [][2]int) []*proto.ProtocolMessage {
	var allData = []byte{
		byte(len(deltas)),
	}
	const maxPayload = 128
	payloads := []*proto.ProtocolMessage{}
	var count = 0
	for _, d := range deltas {
		if len(allData) >= maxPayload {
			allData[0] = byte(count)
			payloads = append(
				payloads,
				&proto.ProtocolMessage{
					Type: proto.MsgUpdateMemMorePending,
					Body: allData,
				},
			)
			count = 0
			allData = []byte{0x00}
		}
		data, err := s.serialize(
			[]proto.Argument{
				{Type: proto.ArgTypeWord, Value: d[0]},
				{Type: proto.ArgTypeByte, Value: d[1]},
			},
		)
		if err != nil {
			return nil
		}
		allData = append(allData, data...)
		count++
	}
	allData[0] = byte(count)
	payloads = append(
		payloads,
		&proto.ProtocolMessage{
			Type: proto.MsgUpdateMem,
			Body: allData,
		},
	)
	return payloads
}

func newPhysicsService() *internalPhysicsService {
	ips := &internalPhysicsService{
		buffer: nil,
		pe:     NewPhysicsEngine(0, 0, 40, 48, 20*time.Millisecond),
	}
	// ips.pe.space.SetGravity(cp.Vector{0, 2})
	return ips
}

type internalPhysicsService struct {
	// stuff
	buffer []byte
	pe     *PhysicsEngine
	w      telnet.Writer
}

func (s *internalPhysicsService) sendWelcome(w telnet.Writer) {
	s.w = w
	log.Printf("Set writer to %+v", s.w)
	s.sendMessage(
		w,
		&proto.ProtocolMessage{
			Type: proto.MsgGreeting,
			Body: append([]byte("HELLO\r")),
		},
	)
}

func (s internalPhysicsService) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	log.Printf("PhysicsService: new connection")

	s.sendWelcome(w)

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

func (s *internalPhysicsService) sendMessage(w telnet.Writer, resp *proto.ProtocolMessage) {
	oi.LongWrite(w, []byte{byte(resp.Type)})
	oi.LongWrite(w, []byte{
		byte(len(resp.Body) & 0xff),
		byte(len(resp.Body) / 256),
	})
	oi.LongWrite(w, resp.Body)
	log.Printf("Sending message: Type = $%x, Payload Length = %d bytes [%+v]", byte(resp.Type), len(resp.Body), resp.Body)
}

func (s *internalPhysicsService) checkForMessage(w telnet.Writer) {
	if len(s.buffer) >= 3 {
		t := proto.MessageType(s.buffer[0])
		size := int(s.buffer[1]) + int(s.buffer[2])*256
		if len(s.buffer) >= size+3 {
			parcel := s.buffer[3 : size+3]
			s.buffer = s.buffer[size+3:]
			resp, err := s.handleMessage(&proto.ProtocolMessage{Type: t, Size: size, Body: parcel}, w)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				return
			}
			if resp != nil {
				s.sendMessage(w, resp)
			}
		}
	}
}

func (s *internalPhysicsService) handleMessage(msg *proto.ProtocolMessage, w telnet.Writer) (*proto.ProtocolMessage, error) {
	log.Printf("handleMessage: received message '%s' (%d bytes)", msg.Type, len(msg.Body))

	switch msg.Type {
	case proto.MsgRequestDeltas:
		deltas := s.pe.screen.GetDeltasWithBase(1024)
		if len(deltas) > 0 {
			payloads := s.reportData(deltas)
			for _, m := range payloads[:len(payloads)-1] {
				s.sendMessage(w, m)
			}
			return payloads[len(payloads)-1], nil
		} else {
			return &proto.ProtocolMessage{
				Type: proto.MsgOk,
				Body: []byte{1},
			}, nil
		}
	case proto.MsgStopPhysics:
		s.pe.Stop()
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgStartPhysics:
		s.pe.Start()
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgGetCollision:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		collided, with := s.pe.GetCollidedWith(int(id))
		if collided {
			return &proto.ProtocolMessage{
				Type: proto.MsgGetCollisionResponse,
				Body: []byte{1, byte(with)},
			}, nil
		} else {
			return &proto.ProtocolMessage{
				Type: proto.MsgGetCollisionResponse,
				Body: []byte{0, 0},
			}, nil
		}
	case proto.MsgGetPosition:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		x, y := s.pe.GetObjectPos(int(id))
		return &proto.ProtocolMessage{
			Type: proto.MsgGetPositionResponse,
			Body: []byte{byte(x), byte(y)},
		}, nil
	case proto.MsgGetOOB:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		oob := s.pe.GetObjectOOB(int(id))
		return &proto.ProtocolMessage{
			Type: proto.MsgGetOOBResponse,
			Body: []byte{byte(oob)},
		}, nil
	case proto.MsgGetColor:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		c := s.pe.GetObjectColor(int(id))
		return &proto.ProtocolMessage{
			Type: proto.MsgGetColorResponse,
			Body: []byte{byte(c)},
		}, nil
	case proto.MsgAddBlockingRegionRect:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
				{Name: "x", Type: proto.ArgTypeByte},
				{Name: "y", Type: proto.ArgTypeByte},
				{Name: "w", Type: proto.ArgTypeByte},
				{Name: "h", Type: proto.ArgTypeByte},
				{Name: "color", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		w := int(params["w"].(byte))
		h := int(params["h"].(byte))
		x := int(params["x"].(byte)) + w/2
		y := int(params["y"].(byte)) + h/2
		c := int(params["color"].(byte))
		s.pe.addRect(
			int(id), float64(w), float64(h), 1000, cp.Vector{X: float64(x) + 0.5, Y: float64(y)}, cp.Vector{X: 0, Y: 0},
			c, cp.BODY_STATIC,
			false,
		)
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgDefineObjectShapeRect:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
				{Name: "w", Type: proto.ArgTypeByte},
				{Name: "h", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		w := int(params["w"].(byte))
		h := int(params["h"].(byte))
		s.pe.SetObjectRect(int(id), w, h)
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgSetObjectPosition:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
				{Name: "x", Type: proto.ArgTypeByte},
				{Name: "y", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		x := int(params["x"].(byte))
		y := int(params["y"].(byte))
		s.pe.SetObjectPos(int(id), x, y)
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgSetObjectType:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
				{Name: "type", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		kind := int(params["type"].(byte))
		s.pe.SetObjectType(int(id), kind)
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgSetObjectMass:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
				{Name: "mass", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		mass := int(params["mass"].(byte))
		s.pe.SetObjectMass(int(id), mass)
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgSetObjectColor:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
				{Name: "color", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		col := int(params["color"].(byte) & 0x0f)
		s.pe.SetObjectColor(int(id), col)
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgSetObjectVelocity:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
				{Name: "velX", Type: proto.ArgTypeSignedByte},
				{Name: "velY", Type: proto.ArgTypeSignedByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		velX := float64(params["velX"].(int8))
		velY := float64(params["velY"].(int8))
		s.pe.SetObjectVelocity(int(id), velX, velY)
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil
	case proto.MsgDefineObject:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "objectId", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		log.Printf("Arguments: %+v", params)
		id := params["objectId"].(byte)
		s.pe.RemoveObject(int(id))
		s.pe.addCircle(
			int(id), 1, 10, cp.Vector{X: 20, Y: 24}, cp.Vector{X: 0, Y: 0},
			15, cp.BODY_DYNAMIC,
			false,
		)
		log.Printf("Created object with Id: $%.2x", id)
		return &proto.ProtocolMessage{
			Type: proto.MsgOk,
			Body: []byte{1},
		}, nil

	case proto.MsgInitSystem:
		params, _, err := s.deserialize(
			msg.Body,
			[]proto.Argument{
				{Name: "mode", Type: proto.ArgTypeByte},
			},
		)
		if err != nil {
			return nil, err
		}
		s.pe.Stop()
		s.pe = NewPhysicsEngine(0, 0, 40, 48, 20*time.Millisecond)
		s.pe.screen.Clear()
		log.Printf("Arguments: %+v", params)
		data, err := s.serialize(
			[]proto.Argument{
				{Type: proto.ArgTypeWord, Value: 0x400},
				{Type: proto.ArgTypeWord, Value: 0x3ff},
				{Type: proto.ArgTypeByte, Value: 0x00},
			},
		)
		if err != nil {
			return nil, err
		}
		return &proto.ProtocolMessage{
			Type: proto.MsgClearMem,
			Body: data,
		}, nil
	}

	return nil, nil
}

func (s *internalPhysicsService) serialize(args []proto.Argument) ([]byte, error) {
	var data = []byte{}
	for _, arg := range args {
		switch arg.Type {
		case proto.ArgTypeByte:
			if b, ok := arg.Value.(byte); ok {
				data = append(data, b)
			} else if b, ok := arg.Value.(int); ok {
				data = append(data, byte(b))
			} else {
				return data, errors.New("serialize expected byte value")
			}
		case proto.ArgTypeWord:
			if b, ok := arg.Value.(uint16); ok {
				data = append(data, byte(b&0xff), byte(b>>8))
			} else if b, ok := arg.Value.(int); ok {
				data = append(data, byte(b&0xff), byte(b>>8))
			} else {
				return data, errors.New("serialize expected word value")
			}
		default:
			return data, fmt.Errorf("unsupported type: %d", arg.Type)
		}
	}
	return data, nil
}

func byteToInt8(b byte) int8 {
	if b&0x80 != 0 {
		return -int8(128 - int(b&0x7f))
	} else {
		return int8(b)
	}
}

func (s *internalPhysicsService) deserialize(data []byte, args []proto.Argument) (map[string]any, int, error) {
	ptr := 0
	argIndex := 0
	out := map[string]any{}
	for argIndex < len(args) && ptr < len(data) {
		var arg = args[argIndex]
		log.Printf("De-ser: arg = %s", arg.Name)
		switch arg.Type {
		case proto.ArgTypeSignedByte:
			if len(data)-ptr >= 1 {
				out[arg.Name] = byteToInt8(data[ptr])
				ptr++
			} else {
				return out, -1, errors.New("packet truncated expecting byte")
			}
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
		argIndex++
	}
	return out, ptr, nil
}
