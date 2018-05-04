package main

import (
	//"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/gnbuffers"
	//	"github.com/itfantasy/gonode/utils/json"
	//"github.com/itfantasy/gonode/utils/stl"
	"github.com/itfantasy/gonode_demo/icloud/opcode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/paramcode"
)

func HandleMsg(id string, msg []byte) {
	parser := gnbuffers.BuildParser(msg, 0)
	if opCode, err := parser.Byte(); err != nil {
		gonode.Node().Logger().Error(err.Error())
		return
	} else {
		switch opCode {
		case opcode.Authenticate:
			handleAuthenticate(id, opCode, parser)
			break
		case opcode.CreateGame:
			handleCreateGame(id, opCode, parser)
			break
		case opcode.JoinGame:
			handleJoinGame(id, opCode, parser)
			break
		case opcode.JoinRandomGame:
			handleJoinRandomGame(id, opCode, parser)
			break
		default:
			gonode.Send(id, msg)
			break
		}
	}
}

func handleErrors(id string, opCode byte, err error) {
	gonode.Error(err.Error())
}

func handleAuthenticate(id string, opCode byte, parser *gnbuffers.GnParser) {
	if buf, err := gnbuffers.BuildBuffer(256); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {
		buf.PushByte(0)      // resp
		buf.PushShort(0)     // retcode
		buf.PushByte(opCode) // opcode
		gonode.Send(id, buf.Bytes())
	}
}

func handleCreateGame(id string, opCode byte, parser *gnbuffers.GnParser) {
	if buf, err := gnbuffers.BuildBuffer(256); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {
		buf.PushByte(0)
		buf.PushShort(errorcode.Ok)
		buf.PushByte(opCode)
		buf.PushByte(paramcode.GameId)
		buf.PushObject("game1123")
		buf.PushByte(paramcode.Address)
		buf.PushObject("192.168.10.85:5056")
		gonode.Send(id, buf.Bytes())
	}
}

func handleJoinGame(id string, opCode byte, parser *gnbuffers.GnParser) {
	if buf, err := gnbuffers.BuildBuffer(256); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {
		buf.PushByte(0)
		buf.PushShort(errorcode.Ok)
		buf.PushByte(opCode)
		buf.PushByte(paramcode.GameId)
		buf.PushObject("game1123")
		buf.PushByte(paramcode.Address)
		buf.PushObject("192.168.10.85:5056")
		gonode.Send(id, buf.Bytes())
	}
}

func handleJoinRandomGame(id string, opCode byte, parser *gnbuffers.GnParser) {
	if buf, err := gnbuffers.BuildBuffer(256); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {
		buf.PushByte(0)
		buf.PushShort(errorcode.NoMatchFound)
		buf.PushByte(opCode)
		gonode.Send(id, buf.Bytes())
	}
}
