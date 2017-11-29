package main

import (
	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/gnbuffers"
	"github.com/itfantasy/gonode_demo/lobbyandroom/opcode"
	"github.com/itfantasy/gonode_demo/lobbyandroom/opcode/errorcode"
	"github.com/itfantasy/gonode_demo/lobbyandroom/opcode/paramcode"
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
		case opcode.JoinRandomGame:
			handleJoinRandomGame(id, opCode, parser)
			break
		case opcode.JoinGame:
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
	} else {
		buf.PushByte(0)      // resp
		buf.PushByte(opCode) // opcode
		buf.PushShort(0)     // retcode
		gonode.Send(id, buf.Bytes())
	}
}

func handleJoinRandomGame(id string, opCode byte, parser *gnbuffers.GnParser) {
	if buf, err := gnbuffers.BuildBuffer(256); err != nil {
		handleErrors(id, opCode, err)
	} else {
		buf.PushByte(0)
		buf.PushByte(opCode)
		buf.PushShort(errorcode.Ok)
		buf.PushByte(paramcode.GameId)
		buf.PushObject("game1123")
		buf.PushByte(paramcode.Address)
		buf.PushObject("192.168.10.94:5055")
		gonode.Send(id, buf.Bytes())
	}
}

func handleJoinGame(id string, opCode byte, parser *gnbuffers.GnParser) {
	if buf, err := gnbuffers.BuildBuffer(256); err != nil {
		handleErrors(id, opCode, err)
	} else {

	}
}
