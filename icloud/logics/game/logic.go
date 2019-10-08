package game

import (
	"errors"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/core/binbuf"

	"github.com/itfantasy/gonode-icloud/icloud/opcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/actorparam"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/cacheop"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/evncode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/gameparam"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/paramcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/recvgroup"
	"github.com/itfantasy/gonode/utils/stl"

	"github.com/itfantasy/gonode-toolkit/toolkit/gen_room"
)

func HandleConn(id string) {
	gen_room.AddPeer(gen_room.NewRoomPeer(id))
}

func HandleMsg(id string, msg []byte) {
	parser := binbuf.BuildParser(msg, 0)
	if opCode := parser.Byte(); parser.Error() != nil {
		gonode.LogError(parser.Error())
		return
	} else {
		peer, ok := gen_room.GetPeer(id)
		if !ok {
			return
		}
		switch opCode {
		case opcode.Authenticate:
			handleAuthenticate(peer, opCode, parser)
			break
		case opcode.CreateGame:
			handleCreateGame(peer, opCode, parser)
			break
		case opcode.JoinGame:
			handleJoinGame(peer, opCode, parser)
			break
		case opcode.RaiseEvent:
			handleRaiseEvent(peer, opCode, parser)
			break
		case opcode.SetProperties:
			handleSetProperties(peer, opCode, parser)
		default:
			break
		}
	}
}

func HandleClose(id string) {
	peer, ok := gen_room.GetPeer(id)
	if !ok {
		return
	}
	room, actor, err := gen_room.GetActorInRoom(peer.PeerId(), peer.RoomId())
	if err != nil {
		gonode.LogError(err)
		return
	}
	pubLeaveEvent(peer, actor, room)
	_, _, err2 := gen_room.LeaveRoom(peer.PeerId(), room.RoomId())
	if err2 != nil {
		gonode.LogError(err2)
		return
	}
	if room.IsEmpty() {
		gen_room.DisposeRoom(room.RoomId(), id)
	}
	gen_room.RemovePeer(id)
}

func handleSetProperties(peer *gen_room.RoomPeer, opCode byte, parser *binbuf.BinParser) {
	datas, _ := binbuf.BuildBuffer(256).
		PushByte(0).             // resp
		PushShort(0).            // retcode
		PushByte(opCode).Bytes() // opcode
	gonode.Send(peer.PeerId(), datas)
}

func handleAuthenticate(peer *gen_room.RoomPeer, opCode byte, parser *binbuf.BinParser) {
	datas, _ := binbuf.BuildBuffer(256).
		PushByte(0).             // resp
		PushShort(0).            // retcode
		PushByte(opCode).Bytes() // opcode
	gonode.Send(peer.PeerId(), datas)
}

func handleCreateGame(peer *gen_room.RoomPeer, opCode byte, parser *binbuf.BinParser) {
	parser.Byte()
	_roomId := parser.Object()
	if parser.Error() != nil {
		gonode.LogError(parser.Error())
		return
	}
	roomId, ok := _roomId.(string)
	if !ok {
		gonode.LogError(errors.New("cannot get the roomId!"))
		return
	}

	room, actor, err := gen_room.CreateRoom(peer.PeerId(), roomId)
	if err != nil {
		gonode.LogError(err)
	}
	peer.SetRoomId(room.RoomId())

	buf := binbuf.BuildBuffer(256)
	buf.PushByte(0)
	buf.PushShort(errorcode.Ok)
	buf.PushByte(opCode)
	buf.PushByte(paramcode.ActorNr)
	buf.PushObject(actor.ActorNr())
	buf.PushByte(paramcode.GameProperties)
	hash := stl.NewHashTable()
	list2 := stl.NewList(0)
	hash.Set(gameparam.LobbyProperties, list2.Values())
	hash.Set(gameparam.CleanupCacheOnLeave, true)
	hash.Set(gameparam.MaxPlayers, byte(4))
	hash.Set(gameparam.IsVisible, true)
	hash.Set(gameparam.IsOpen, true)
	hash.Set(gameparam.MasterClientId, room.MasterId())
	hash.Set(gameparam.CleanupCacheOnLeave, true)
	buf.PushObject(hash.KeyValuePairs())
	buf.PushByte(paramcode.Actors)
	buf.PushObject(room.ActorsManager().GetAllActorNrs())
	datas, _ := buf.Bytes()
	gonode.Send(peer.PeerId(), datas)

	pubJoinEvent(peer, actor, room)
}

func handleJoinGame(peer *gen_room.RoomPeer, opCode byte, parser *binbuf.BinParser) {
	parser.Byte()
	_roomId := parser.Object()
	if parser.Error() != nil {
		gonode.LogError(parser.Error())
		return
	}
	roomId, ok := _roomId.(string)
	if !ok {
		gonode.LogError(errors.New("cannot get the roomId!"))
		return
	}

	room, actor, err := gen_room.JoinRoom(peer.PeerId(), roomId)
	if err != nil {
		gonode.LogError(err)
	}
	peer.SetRoomId(room.RoomId())

	buf := binbuf.BuildBuffer(1024)
	buf.PushByte(0)
	buf.PushShort(errorcode.Ok)
	buf.PushByte(opCode)
	buf.PushByte(paramcode.ActorNr)
	buf.PushObject(actor.ActorNr())
	buf.PushByte(paramcode.GameProperties)
	hash := stl.NewHashTable()
	hash.Set(gameparam.LobbyProperties, true)
	hash.Set(gameparam.CleanupCacheOnLeave, true)
	hash.Set(gameparam.MaxPlayers, byte(4))
	hash.Set(gameparam.IsVisible, true)
	hash.Set(gameparam.IsOpen, true)
	hash.Set(gameparam.MasterClientId, room.MasterId())
	hash.Set(gameparam.CleanupCacheOnLeave, true)
	buf.PushObject(hash.KeyValuePairs())
	buf.PushByte(paramcode.Actors)
	buf.PushObject(room.ActorsManager().GetAllActorNrs())
	datas, _ := buf.Bytes()
	gonode.Send(peer.PeerId(), datas)

	pubJoinEvent(peer, actor, room)
}

func handleRaiseEvent(peer *gen_room.RoomPeer, opCode byte, parser *binbuf.BinParser) {
	// send self resp
	datas, _ := binbuf.BuildBuffer(1024).
		PushByte(0).
		PushShort(errorcode.Ok).
		PushByte(opCode).Bytes()
	gonode.Send(peer.PeerId(), datas)

	// pub the event to others
	parser.Byte()              // ParameterCode.Code
	parser.Byte()              // gntypes.Byte
	eventCode := parser.Byte() // eventCode
	var recvGroup byte = recvgroup.Others
	var cacheOp byte = cacheop.DoNotCache
	for {
		paramCode := parser.Byte()
		if parser.Error() != nil {
			gonode.LogError(parser.Error())
			return
		}
		if paramCode == paramcode.ReceiverGroup { // ReceiverGroup
			if oRecvGroup := parser.Object(); parser.Error() != nil {
				gonode.LogError(parser.Error())
				return
			} else {
				recvGroup = oRecvGroup.(byte)
			}
		} else if paramCode == paramcode.Cache { // the event which will be cached
			if oCacheOp := parser.Object(); parser.Error() != nil {
				gonode.LogError(parser.Error())
				return
			} else {
				cacheOp = oCacheOp.(byte)
			}
		} else if paramCode == paramcode.Data { // when get the data
			break
		}
	}
	_, actor, err := gen_room.GetActorInRoom(peer.PeerId(), peer.RoomId())
	if err != nil {
		gonode.LogError(err)
	}

	evn := binbuf.BuildBuffer(1024)
	evn.PushByte(1) // event
	evn.PushByte(eventCode)
	evn.PushByte(paramcode.ActorNr)
	evn.PushObject(actor.ActorNr())
	evn.PushByte(paramcode.Code)
	evn.PushObject(byte(eventCode))
	evn.PushByte(paramcode.Data)
	data := parser.Bytes()
	evn.PushBytes(data[5:])

	datas2, _ := evn.Bytes()
	addToRoomCache := (cacheOp == cacheop.AddToRoomCache || cacheOp == cacheop.AddToRoomCacheGlobal)
	gen_room.RaiseEvent(peer.PeerId(), peer.RoomId(), datas2, recvGroup, addToRoomCache)
}

func pubJoinEvent(peer *gen_room.RoomPeer, actor *gen_room.Actor, room *gen_room.RoomEntity) {
	gen_room.RcvCacheEvent(peer.PeerId(), room.RoomId())

	evn := binbuf.BuildBuffer(1024)
	evn.PushByte(1)
	evn.PushByte(evncode.Join)
	hashTable := stl.NewHashTable()
	hashTable.Set(actorparam.Nickname, "")
	evn.PushByte(paramcode.ActorProperties)
	evn.PushObject(hashTable.KeyValuePairs())
	evn.PushByte(paramcode.ActorNr)
	evn.PushObject(actor.ActorNr())
	evn.PushByte(paramcode.Actors)
	evn.PushObject(room.ActorsManager().GetAllActorNrs())
	datas, _ := evn.Bytes()
	gen_room.RaiseEvent(peer.PeerId(), room.RoomId(), datas, gen_room.RcvGroup_All, true)
}

func pubLeaveEvent(peer *gen_room.RoomPeer, actor *gen_room.Actor, room *gen_room.RoomEntity) {
	evn := binbuf.BuildBuffer(1024)
	evn.PushByte(1)
	evn.PushByte(evncode.Leave)
	evn.PushByte(paramcode.ActorNr)
	evn.PushObject(actor.ActorNr())
	evn.PushByte(paramcode.Actors)
	evn.PushObject(room.ActorsManager().GetAllActorNrs())
	evn.PushByte(paramcode.IsInactive)
	evn.PushObject(false)
	datas, _ := evn.Bytes()
	gen_room.RaiseEvent(peer.PeerId(), room.RoomId(), datas, gen_room.RcvGroup_Others, false)
}
