package proto

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
	default:
		return "unknown-message"
	}
}

type ArgType int

const (
	ArgTypeByte ArgType = 0x00
	ArgTypeWord ArgType = 0x01
)

type ProtocolMessage struct {
	Type MessageType
	Size int
	Body []byte
}

type Argument struct {
	Name string
	Type ArgType
}
