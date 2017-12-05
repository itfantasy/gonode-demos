package main

import (
	"fmt"
	//	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/gnbuffers"
	//	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/stl"
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
			handleJoinGame(id, opCode, parser)
			break
		case opcode.RaiseEvent:
			handleRaiseEvent(id, opCode, parser)
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

func handleJoinRandomGame(id string, opCode byte, parser *gnbuffers.GnParser) {
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
		buf.PushObject("192.168.10.94:5055")
		gonode.Send(id, buf.Bytes())
	}
}

func handleJoinGame(id string, opCode byte, parser *gnbuffers.GnParser) {
	if buf, err := gnbuffers.BuildBuffer(1024); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {
		buf.PushByte(0)
		buf.PushShort(errorcode.Ok)
		buf.PushByte(opCode)
		buf.PushByte(paramcode.ActorNr)
		buf.PushObject(0)
		buf.PushByte(paramcode.ActorProperties)
		hash := stl.NewHashTable()
		buf.PushObject(hash.KeyValuePairs())
		buf.PushByte(paramcode.GameProperties)
		hash2 := stl.NewHashTable()
		buf.PushObject(hash2.KeyValuePairs())
		buf.PushByte(paramcode.Actors)
		intArray := make([]int32, 0, 0)
		buf.PushObject(intArray)
		gonode.Send(id, buf.Bytes())
	}
}

var actorNr = 7

func handleRaiseEvent(id string, opCode byte, parser *gnbuffers.GnParser) {
	// send self resp
	if buf, err := gnbuffers.BuildBuffer(1024); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {
		buf.PushByte(0)
		buf.PushShort(errorcode.Ok)
		buf.PushByte(opCode)
		gonode.Send(id, buf.Bytes())
	}

	// pub the event to others

	if evn, err := gnbuffers.BuildBuffer(1024); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {
		parser.Byte()                 // ParameterCode.Code
		parser.Byte()                 // gntypes.Byte
		eventCode, _ := parser.Byte() // eventCode

		evn.PushByte(1) // event
		evn.PushByte(eventCode)

		evn.PushByte(paramcode.ActorNr)
		evn.PushObject(actorNr)
		actorNr += 1

		evn.PushByte(paramcode.Code)
		evn.PushObject(byte(eventCode))

		//parser.Byte() // ParameterCode.Data

		evn.PushByte(paramcode.Data) // hashtable

		data := parser.Bytes()
		evn.PushBytes(data[5:])

		ids := gonode.Node().NetWorker().GetAllConnIds()
		for _, item := range ids {
			//if item != id {
			fmt.Println(evn.Bytes())
			gonode.Send(item, evn.Bytes())
			//}
		}
	}
}
