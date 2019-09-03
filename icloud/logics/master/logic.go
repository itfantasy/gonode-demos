package master

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/paramcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/servereventcode"
	"github.com/itfantasy/gonode/core/binbuf"
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
			break
		}
	}
}

func handleErrors(id string, opCode byte, err error) {
	gonode.LogError(err)
}

func handleAuthenticate(id string, opCode byte, parser *binbuf.BinParser) {
	datas, _ := binbuf.BuildBuffer(256).
		PushByte(0).             // resp
		PushShort(0).            // retcode
		PushByte(opCode).Bytes() // opcode
	gonode.Send(id, datas)
}

func handleCreateGame(id string, opCode byte, parser *binbuf.BinParser) {
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
	datas, _ := buf.PushByte(0).
		PushShort(errorcode.Ok).
		PushByte(opCode).
		PushByte(paramcode.GameId).
		PushObject("game1123").
		PushByte(paramcode.Address).
		PushObject(tempRoomUrl).Bytes()
	gonode.Send(id, datas)
}
