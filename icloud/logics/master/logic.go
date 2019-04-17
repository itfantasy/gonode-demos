package master

import (
	"errors"
	"fmt"
	//"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/gnbuffers"
	//	"github.com/itfantasy/gonode/utils/json"
	//"github.com/itfantasy/gonode/utils/stl"
	"github.com/itfantasy/gonode-icloud/icloud/behaviors/lobby"
	"github.com/itfantasy/gonode-icloud/icloud/opcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/paramcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/servereventcode"
	"github.com/itfantasy/gonode/utils/strs"
)

var tempRoomUrl string

func SetDefaultRoomUrl(url string) {
	tempRoomUrl = url
}

func HandleConn(id string) {
	if strs.StartsWith(id, "room") {

	}
}

func HandleClose(id string) {
	if strs.StartsWith(id, "room") {

	}
}

func HandleServerMsg(id string, msg []byte) {
	parser := gnbuffers.BuildParser(msg, 0)
	if opCode, err := parser.Byte(); err != nil {
		gonode.Node().Logger().Error(err.Error())
		return
	} else {
		switch opCode {
		case servereventcode.UpdateGameState:
			handleUpdateGameState(id, opCode, parser)
			break
		case servereventcode.RemoveGameState:
			fmt.Println("Receive the RemoveGameState Event!")
			fmt.Println(msg)
			handleRemoveGameState(id, opCode, parser)
			break
		}
	}
}

func handleUpdateGameState(id string, opCode byte, parser *gnbuffers.GnParser) {

}

func handleRemoveGameState(id string, opCode byte, parser *gnbuffers.GnParser) {
	gameId, err := parser.String()
	if err != nil {
		handleErrors(id, opCode, err)
		return
	}
	insLobby().RemoveRoomState(gameId)
}

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

var _insLobby *lobby.Lobby = nil

func insLobby() *lobby.Lobby {
	if _insLobby == nil {
		_insLobby = lobby.NewLobby("default")
	}
	return _insLobby
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
		gameId := lobby.GenerateRoomId()
		if ret, err := insLobby().CreateRoomState(gameId, ""); err != nil {
			handleErrors(id, opCode, err)
			return
		} else if !ret {
			handleErrors(id, opCode, errors.New("cannot create a roomstate:"+gameId))
			return
		} else {
			buf.PushByte(0)
			buf.PushShort(errorcode.Ok)
			buf.PushByte(opCode)
			buf.PushByte(paramcode.GameId)
			buf.PushObject(gameId)
			buf.PushByte(paramcode.Address)
			buf.PushObject(tempRoomUrl)
			gonode.Send(id, buf.Bytes())
		}
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
		buf.PushObject(tempRoomUrl)
		gonode.Send(id, buf.Bytes())
	}
}

func handleJoinRandomGame(id string, opCode byte, parser *gnbuffers.GnParser) {
	if buf, err := gnbuffers.BuildBuffer(256); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {
		gameId, exist := insLobby().RandomRoomStateId()
		if !exist {
			buf.PushByte(0)
			buf.PushShort(errorcode.NoMatchFound)
			buf.PushByte(opCode)
			gonode.Send(id, buf.Bytes())
		} else {
			buf.PushByte(0)
			buf.PushShort(errorcode.Ok)
			buf.PushByte(opCode)
			buf.PushByte(paramcode.GameId)
			buf.PushObject(gameId)
			buf.PushByte(paramcode.Address)
			buf.PushObject(tempRoomUrl)
			gonode.Send(id, buf.Bytes())
		}
	}
}
