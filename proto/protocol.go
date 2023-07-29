package proto

import "fmt"

type MessageType uint8

const (
	MsgInitSystem              MessageType = 0x01
	MsgDefineObject            MessageType = 0x02
	MsgSetObjectMass           MessageType = 0x04
	MsgSetObjectVelocity       MessageType = 0x05
	MsgSetObjectPosition       MessageType = 0x06
	MsgSetObjectHeading        MessageType = 0x07
	MsgSetObjectElasticity     MessageType = 0x08
	MsgDefineGlobalForce       MessageType = 0x09
	MsgRemoveObject            MessageType = 0x0a
	MsgAddBlockingRegionRect   MessageType = 0x0b
	MsgAddBlockingRegionCircle MessageType = 0x0c
	MsgDefineObjectShapeRect   MessageType = 0x0d
	MsgDefineObjectShapeCircle MessageType = 0x0e
	MsgSetObjectColor          MessageType = 0x0f
	MsgRequestDeltas           MessageType = 0x10
	MsgSetObjectType           MessageType = 0x11
	MsgStartPhysics            MessageType = 0x12
	MsgStopPhysics             MessageType = 0x13
	MsgGetPosition             MessageType = 0x14
	MsgGetColor                MessageType = 0x15
	MsgGetOOB                  MessageType = 0x16
	MsgGetCollision            MessageType = 0x17
	//
	MsgGreeting             MessageType = 0x7f
	MsgClearMem             MessageType = 0x80
	MsgUpdateMem            MessageType = 0x81
	MsgGetPositionResponse  MessageType = 0x82
	MsgGetColorResponse     MessageType = 0x83
	MsgGetOOBResponse       MessageType = 0x84
	MsgUpdateMemMorePending MessageType = 0x85
	MsgGetCollisionResponse MessageType = 0x86
	//
	MsgOk    MessageType = 0xf0
	MsgError MessageType = 0xf1
)

func (t MessageType) String() string {
	switch t {
	case MsgInitSystem:
		return "init-system"
	case MsgDefineObject:
		return "define-object"
	case MsgSetObjectMass:
		return "set-object-mass"
	case MsgSetObjectVelocity:
		return "set-object-velocity"
	case MsgSetObjectPosition:
		return "set-object-position"
	case MsgSetObjectHeading:
		return "set-object-heading"
	case MsgSetObjectElasticity:
		return "set-object-elasticity"
	case MsgDefineGlobalForce:
		return "define-global-force"
	case MsgRemoveObject:
		return "remove-object"
	case MsgAddBlockingRegionRect:
		return "add-blocking-rect"
	case MsgAddBlockingRegionCircle:
		return "add-blocking-circle"
	case MsgDefineObjectShapeRect:
		return "define-object-rect"
	case MsgDefineObjectShapeCircle:
		return "define-object-circle"
	case MsgSetObjectColor:
		return "set-object-color"
	case MsgSetObjectType:
		return "set-object-type"
	case MsgRequestDeltas:
		return "request-video-deltas"
	case MsgStartPhysics:
		return "start-physics"
	case MsgStopPhysics:
		return "stop-physics"
	case MsgGetPosition:
		return "get-object-position"
	case MsgGetColor:
		return "get-object-color"
	case MsgGetOOB:
		return "get-object-oob-state"
	case MsgGetCollision:
		return "get-collision-state"
	case MsgGreeting:
		return "greeting"
	case MsgClearMem:
		return "clear-memory"
	default:
		return fmt.Sprintf("unknown-message ($%.2x)", int(t))
	}
}

type ArgType int

const (
	ArgTypeByte       ArgType = 0x00
	ArgTypeWord       ArgType = 0x01
	ArgTypeSignedByte ArgType = 0x02
)

type ProtocolMessage struct {
	Type MessageType
	Size int
	Body []byte
}

type Argument struct {
	Name  string
	Type  ArgType
	Value any
}
