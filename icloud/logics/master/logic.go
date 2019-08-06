package master

import (
	"errors"
	"fmt"
	//"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/core/binbuf"
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
	parser := binbuf.BuildParser(msg, 0)
	if opCode := parser.Byte(); parser.Error() != nil {
		gonode.LogError(parser.Error())
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

func handleUpdateGameState(id string, opCode byte, parser *binbuf.BinParser) {

}

func handleRemoveGameState(id string, opCode byte, parser *binbuf.BinParser) {
	gameId := parser.String()
	if parser.Error() != nil {
		handleErrors(id, opCode, parser.Error())
		return
	}
	insLobby().RemoveRoomState(gameId)
}

func HandleMsg(id string, msg []byte) {
	parser := binbuf.BuildParser(msg, 0)
	if opCode := parser.Byte(); parser.Error() != nil {
		gonode.LogError(parser.Error())
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
			//gonode.Send(id, msg)
			break
		}
	}
}

func handleErrors(id string, opCode byte, err error) {
	gonode.LogError(err)
}

var _insLobby *lobby.Lobby = nil

func insLobby() *lobby.Lobby {
	if _insLobby == nil {
		_insLobby = lobby.NewLobby("default")
	}
	return _insLobby
}

func handleAuthenticate(id string, opCode byte, parser *binbuf.BinParser) {
	datas, _ := binbuf.BuildBuffer(256).
		PushByte(0).             // resp
		PushShort(0).            // retcode
		PushByte(opCode).Bytes() // opcode
	gonode.Send(id, datas)
}

func handleCreateGame(id string, opCode byte, parser *binbuf.BinParser) {
	gameId := lobby.GenerateRoomId()
	if ret, err := insLobby().CreateRoomState(gameId, ""); err != nil {
		handleErrors(id, opCode, err)
		return
	} else if !ret {
		handleErrors(id, opCode, errors.New("cannot create a roomstate:"+gameId))
		return
	} else {
		datas, _ := binbuf.BuildBuffer(256).
			PushByte(0).
			PushShort(errorcode.Ok).
			PushByte(opCode).
			PushByte(paramcode.GameId).
			PushObject(gameId).
			PushByte(paramcode.Address).
			PushObject(tempRoomUrl).Bytes()
		gonode.Send(id, datas)
	}

}

func handleJoinGame(id string, opCode byte, parser *binbuf.BinParser) {
	datas, _ := binbuf.BuildBuffer(256).
		PushByte(0).
		PushShort(errorcode.Ok).
		PushByte(opCode).
		PushByte(paramcode.GameId).
		PushObject("game1123").
		PushByte(paramcode.Address).
		PushObject(tempRoomUrl).Bytes()
	gonode.Send(id, datas)

}

func handleJoinRandomGame(id string, opCode byte, parser *binbuf.BinParser) {
	buf := binbuf.BuildBuffer(256)
	gameId, exist := insLobby().RandomRoomStateId()
	if !exist {
		datas, _ := buf.PushByte(0).
			PushShort(errorcode.NoMatchFound).
			PushByte(opCode).Bytes()
		gonode.Send(id, datas)
	} else {
		datas, _ := buf.PushByte(0).
			PushShort(errorcode.Ok).
			PushByte(opCode).
			PushByte(paramcode.GameId).
			PushObject(gameId).
			PushByte(paramcode.Address).
			PushObject(tempRoomUrl).Bytes()
		gonode.Send(id, datas)
	}

}
