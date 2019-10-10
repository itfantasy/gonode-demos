package master

import (
	"errors"
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-icloud/icloud/gunpeer"
	"github.com/itfantasy/gonode-icloud/icloud/opcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/paramcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/servereventcode"
	"github.com/itfantasy/gonode/core/binbuf"

	"github.com/itfantasy/gonode-toolkit/toolkit"
	"github.com/itfantasy/gonode-toolkit/toolkit/gen_lobby"
	"github.com/itfantasy/gonode/utils/snowflake"
)

func HandleConn(id string) {
	gen_lobby.AddPeer(gen_lobby.NewLobbyPeer(id))
}

func HandleClose(id string) {
	gen_lobby.RemovePeer(id)
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
	opCode, datas, err := gunpeer.ParseMsg(msg)
	if err != nil {
		gonode.LogError(err)
		return
	}
	peer, ok := gen_lobby.GetPeer(id)
	if !ok {
		gonode.LogError(errors.New("peer missing..." + id))
		return
	}
	switch opCode {
	case opcode.Authenticate:
		handleAuthenticate(peer, opCode, datas)
		break
	case opcode.CreateGame:
		handleCreateGame(peer, opCode, datas)
		break
	case opcode.JoinGame:
		handleJoinGame(peer, opCode, datas)
		break
	case opcode.JoinRandomGame:
		handleJoinRandomGame(peer, opCode, datas)
		break
	default:
		break
	}
}

func handleError(peer *gen_lobby.LobbyPeer, opCode byte, err error) {
	gonode.LogError(err)
}

func handleAuthenticate(peer *gen_lobby.LobbyPeer, opCode byte, datas *gunpeer.PeerDatas) {
	gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, nil)
}

func handleCreateGame(peer *gen_lobby.LobbyPeer, opCode byte, datas *gunpeer.PeerDatas) {
	room, err := gen_lobby.CreateRoom(peer.PeerId(), snowflake.Generate())
	if err != nil {
		handleError(peer, opCode, err)
		return
	}
	info, err := gonode.GetNodeInfo(room.NodeId)
	if err != nil {
		handleError(peer, opCode, err)
		return
	}
	pub, ok := info.UsrDatas[toolkit.USRDATA_PUBDOMAIN]
	if !ok {
		pub = info.Url
	}
	gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
		paramcode.GameId:  room.RoomId,
		paramcode.Address: pub,
	})
}

func handleJoinGame(peer *gen_lobby.LobbyPeer, opCode byte, datas *gunpeer.PeerDatas) {
	roomId, _ := datas.GetString(paramcode.GameId)
	if datas.Err() != nil {
		handleError(peer, opCode, datas.Err())
		return
	}
	room, err := gen_lobby.JoinRoom(peer.PeerId(), roomId)
	if err != nil {
		handleError(peer, opCode, err)
		return
	}
	info, err := gonode.GetNodeInfo(room.NodeId)
	if err != nil {
		handleError(peer, opCode, err)
		return
	}
	pub, ok := info.UsrDatas[toolkit.USRDATA_PUBDOMAIN]
	if !ok {
		pub = info.Url
	}
	gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
		paramcode.GameId:  room.RoomId,
		paramcode.Address: pub,
	})
}

func handleJoinRandomGame(peer *gen_lobby.LobbyPeer, opCode byte, datas *gunpeer.PeerDatas) {
	room, err := gen_lobby.JoinRandomRoom(peer.PeerId())
	if err != nil {
		handleError(peer, opCode, err)
		return
	}
	info, err := gonode.GetNodeInfo(room.NodeId)
	if err != nil {
		handleError(peer, opCode, err)
		return
	}
	pub, ok := info.UsrDatas[toolkit.USRDATA_PUBDOMAIN]
	if !ok {
		pub = info.Url
	}
	gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
		paramcode.GameId:  room.RoomId,
		paramcode.Address: pub,
	})
}
