package game

import (
	"errors"
	"fmt"
	//	"strconv"
	//	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/core/binbuf"
	//	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode-icloud/icloud/opcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/actorparam"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/cacheop"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/evncode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/gameparam"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/paramcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/recvgroup"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/servereventcode"
	"github.com/itfantasy/gonode/utils/stl"
	"github.com/itfantasy/gonode/utils/strs"

	"github.com/itfantasy/gonode-icloud/icloud/behaviors/room"
	"github.com/itfantasy/gonode-icloud/icloud/peers"
)

func HandleConn(id string) {
	if strs.StartsWith(id, "cnt") {
		insPeerManager().AddPeer(peers.NewClientPeer(id))
	}
}

func HandleMsg(id string, msg []byte) {
	parser := binbuf.BuildParser(msg, 0)
	if opCode, err := parser.Byte(); err != nil {
		gonode.LogError(err)
		return
	} else {
		peer, ok := insPeerManager().GetClientPeer(id)
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
			gonode.Send(id, msg)
			break
		}
	}
}

func HandleClose(id string) {
	peer, ok := insPeerManager().GetClientPeer(id)
	if !ok {
		return
	}
	fmt.Println("the conn has been closed:" + peer.PeerId())
	room := insRoom(peer.RoomId())
	if actor, exist := room.ActorsManager().GetActorByPeerId(peer.PeerId()); exist {
		fmt.Print("has found the target actor:")
		fmt.Println(actor.ActorNr)
		room.EventCache().RemoveEventsByActor(actor.ActorNr) // remove the events of the actor from the eventcache
		room.ActorsManager().RemoveActorByNr(actor.ActorNr)  // remove the actor from the actormanager
		if room.MasterClientId == actor.ActorNr {
			if newactor, exist := insRoom(peer.RoomId()).ActorsManager().GetActorByIndex(0); exist {
				room.MasterClientId = newactor.ActorNr // testcode:use the first actor of the left room actors as the masterclient
			} else {
				room.MasterClientId = 0 // testcode:when the room has no actors
			}
		}

		fmt.Print("try to pub the disconnect event :")
		fmt.Println(actor.ActorNr)
		//pubDisconnectEvent(id, actor.ActorNr)
		pubLeaveEvent(peer, actor.ActorNr)

		// check the room is empty, and send a remove roomstate event to lobby
		if room.IsEmpty() {
			disposeRoom(room.Name)
			sendRemoveRoomState(room.Name)
		}
	}
}

var _insPeerManager *peers.PeerManager = nil

func insPeerManager() *peers.PeerManager {
	if _insPeerManager == nil {
		_insPeerManager = peers.NewPeerManager()
	}
	return _insPeerManager
}

var actorNr int32 = 1
var _insRoomManager *room.RoomManager = nil

func insRoom(roomId string) *room.Room {
	if _insRoomManager == nil {
		_insRoomManager = room.NewRoomManager()
	}
	return _insRoomManager.FetchRoom(roomId)
}

func disposeRoom(roomId string) {
	if _insRoomManager == nil {
		_insRoomManager = room.NewRoomManager()
	}
	_insRoomManager.DisposeRoom(roomId)
}

func handleErrors(id string, opCode byte, err error) {
	//gonode.Error(err.Error())
	fmt.Print("[ERROR]::")
	fmt.Println(err)
}

func handleSetProperties(peer *peers.ClientPeer, opCode byte, parser *binbuf.BinParser) {
	if buf, err := binbuf.BuildBuffer(256); err != nil {
		handleErrors(peer.PeerId(), opCode, err)
		return
	} else {
		buf.PushByte(0)      // resp
		buf.PushShort(0)     // retcode
		buf.PushByte(opCode) // opcode
		gonode.Send(peer.PeerId(), buf.Bytes())
	}
}

func handleAuthenticate(peer *peers.ClientPeer, opCode byte, parser *binbuf.BinParser) {
	if buf, err := binbuf.BuildBuffer(256); err != nil {
		handleErrors(peer.PeerId(), opCode, err)
		return
	} else {
		buf.PushByte(0)      // resp
		buf.PushShort(0)     // retcode
		buf.PushByte(opCode) // opcode
		gonode.Send(peer.PeerId(), buf.Bytes())
	}
}

func handleCreateGame(peer *peers.ClientPeer, opCode byte, parser *binbuf.BinParser) {
	if buf, err := binbuf.BuildBuffer(256); err != nil {
		handleErrors(peer.PeerId(), opCode, err)
		return
	} else {
		parser.Byte() // 255
		_roomId, err := parser.Object()
		if err != nil {
			handleErrors(peer.PeerId(), opCode, err)
			return
		}
		roomId, ok := _roomId.(string)
		if !ok {
			handleErrors(peer.PeerId(), opCode, errors.New("cannot get the roomId!"))
			return
		}

		buf.PushByte(0)
		buf.PushShort(errorcode.Ok)
		buf.PushByte(opCode)

		buf.PushByte(paramcode.ActorNr)
		buf.PushObject(actorNr)

		//buf.PushByte(paramcode.ActorProperties)
		//hash := stl.NewHashTable()
		//buf.PushObject(hash.KeyValuePairs())

		buf.PushByte(paramcode.GameProperties)
		hash := stl.NewHashTable()
		list2 := stl.NewList(0)
		hash.Set(gameparam.LobbyProperties, list2.Values())
		hash.Set(gameparam.CleanupCacheOnLeave, true)
		hash.Set(gameparam.MaxPlayers, byte(4))
		hash.Set(gameparam.IsVisible, true)
		hash.Set(gameparam.IsOpen, true)
		hash.Set(gameparam.MasterClientId, actorNr)
		hash.Set(gameparam.CleanupCacheOnLeave, true)
		buf.PushObject(hash.KeyValuePairs())

		suc := false
		actor, err := insRoom(roomId).ActorsManager().AddNewActor(peer.PeerId(), actorNr)
		if err != nil {
			handleErrors(peer.PeerId(), opCode, err)
			return
		} else {
			suc = true
			if insRoom(peer.RoomId()).MasterClientId == 0 {
				insRoom(peer.RoomId()).MasterClientId = actor.ActorNr // testcode: use the first actor as the masterclient of the room
			}
			peer.SetRoomId(roomId)
		}
		actorNrs := insRoom(peer.RoomId()).ActorsManager().GetAllActorNrs()
		buf.PushByte(paramcode.Actors)
		buf.PushObject(actorNrs)

		gonode.Send(peer.PeerId(), buf.Bytes())
		if suc {
			pubJoinEvent(peer, opCode, parser, actor.ActorNr)
		}
		actorNr += 1
	}
}

func handleJoinGame(peer *peers.ClientPeer, opCode byte, parser *binbuf.BinParser) {
	if buf, err := binbuf.BuildBuffer(1024); err != nil {
		handleErrors(peer.PeerId(), opCode, err)
		return
	} else {

		parser.Byte() // 255
		_roomId, err := parser.Object()
		if err != nil {
			handleErrors(peer.PeerId(), opCode, err)
			return
		}
		roomId, ok := _roomId.(string)
		if !ok {
			handleErrors(peer.PeerId(), opCode, errors.New("cannot get the roomId!"))
			return
		}

		buf.PushByte(0)
		buf.PushShort(errorcode.Ok)
		buf.PushByte(opCode)

		buf.PushByte(paramcode.ActorNr)
		buf.PushObject(actorNr)

		//buf.PushByte(paramcode.ActorProperties)
		//hash := stl.NewHashTable()
		//buf.PushObject(hash.KeyValuePairs())

		buf.PushByte(paramcode.GameProperties)
		hash := stl.NewHashTable()
		hash.Set(gameparam.LobbyProperties, true)
		hash.Set(gameparam.CleanupCacheOnLeave, true)
		hash.Set(gameparam.MaxPlayers, byte(4))
		hash.Set(gameparam.IsVisible, true)
		hash.Set(gameparam.IsOpen, true)
		hash.Set(gameparam.MasterClientId, insRoom(peer.RoomId()).MasterClientId)
		hash.Set(gameparam.CleanupCacheOnLeave, true)
		buf.PushObject(hash.KeyValuePairs())

		suc := false
		actor, err := insRoom(roomId).ActorsManager().AddNewActor(peer.PeerId(), actorNr)
		if err != nil {
			handleErrors(peer.PeerId(), opCode, err)
			return
		} else {
			suc = true
			peer.SetRoomId(roomId)
		}
		actorNrs := insRoom(peer.RoomId()).ActorsManager().GetAllActorNrs()
		buf.PushByte(paramcode.Actors)
		buf.PushObject(actorNrs)

		gonode.Send(peer.PeerId(), buf.Bytes())
		if suc {
			pubJoinEvent(peer, opCode, parser, actor.ActorNr)
		}
		actorNr += 1
	}
}

func handleRaiseEvent(peer *peers.ClientPeer, opCode byte, parser *binbuf.BinParser) {
	// send self resp
	if buf, err := binbuf.BuildBuffer(1024); err != nil {
		handleErrors(peer.PeerId(), opCode, err)
		return
	} else {
		buf.PushByte(0)
		buf.PushShort(errorcode.Ok)
		buf.PushByte(opCode)
		gonode.Send(peer.PeerId(), buf.Bytes())
	}

	// pub the event to others
	if evn, err := binbuf.BuildBuffer(1024); err != nil {
		handleErrors(peer.PeerId(), opCode, err)
		return
	} else {
		actor, exist := insRoom(peer.RoomId()).ActorsManager().GetActorByPeerId(peer.PeerId())
		if !exist {
			handleErrors(peer.PeerId(), opCode, errors.New("cannot find the actor by the peeerid:"+peer.PeerId()))
			return
		}

		parser.Byte()                 // ParameterCode.Code
		parser.Byte()                 // gntypes.Byte
		eventCode, _ := parser.Byte() // eventCode

		evn.PushByte(1) // event
		evn.PushByte(eventCode)

		evn.PushByte(paramcode.ActorNr)
		evn.PushObject(actor.ActorNr)
		//actorNr += 1

		evn.PushByte(paramcode.Code)
		evn.PushObject(byte(eventCode))

		var recvGroup byte = recvgroup.Others
		var cacheOp byte = cacheop.DoNotCache

		for {
			paramCode, err := parser.Byte()
			if err != nil {
				handleErrors(peer.PeerId(), opCode, err)
				return
			}
			//fmt.Print("paramCode:")
			//fmt.Println(paramCode)
			if paramCode == paramcode.ReceiverGroup { // ReceiverGroup
				if oRecvGroup, err := parser.Object(); err != nil {
					handleErrors(peer.PeerId(), opCode, err)
					return
				} else {
					recvGroup = oRecvGroup.(byte)
				}
				//parser.Byte() // Data
			} else if paramCode == paramcode.Cache { // the event which will be cached
				if oCacheOp, err := parser.Object(); err != nil {
					handleErrors(peer.PeerId(), opCode, err)
					return
				} else {
					cacheOp = oCacheOp.(byte)
				}
			} else if paramCode == paramcode.Data { // when get the data
				break
			}
		}

		evn.PushByte(paramcode.Data)
		data := parser.Bytes()
		evn.PushBytes(data[5:])

		// handle the cacheOp
		if cacheOp == cacheop.AddToRoomCache || cacheOp == cacheop.AddToRoomCacheGlobal {
			if actor, exist := insRoom(peer.RoomId()).ActorsManager().GetActorByPeerId(peer.PeerId()); exist {
				evnCache(peer).AddEvent(actor.ActorNr, eventCode, evn.Bytes())
			}
		}

		// handle the recvGroup
		ids := gonode.AllConnIds() // get the ids in the same room
		if recvGroup == recvgroup.MasterClient {
			gonode.Send(peer.PeerId(), evn.Bytes())
		} else if recvGroup == recvgroup.All {
			for _, item := range ids {
				//fmt.Println(evn.Bytes())
				gonode.Send(item, evn.Bytes())
			}
		} else {
			for _, item := range ids {
				if item != peer.PeerId() {
					//fmt.Println(evn.Bytes())
					gonode.Send(item, evn.Bytes())
				}
			}
		}

	}
}

func evnCache(peer *peers.ClientPeer) *room.RoomEventCache {
	return insRoom(peer.RoomId()).EventCache()
}

//HiveGame.PublishEventCache (line 1410)
func pubEventCache(peer *peers.ClientPeer) {
	for _, item := range evnCache(peer).Events() {
		event := item.(*room.CustomEvent)
		fmt.Print("pub the cache event:")
		fmt.Print(" actorNr:")
		fmt.Print(event.ActorNr)
		fmt.Print(" Code:")
		fmt.Println(event.Code)
		fmt.Print("Data:")
		fmt.Println(event.Data)
		gonode.Send(peer.PeerId(), event.Data)
	}
}

func pubLeaveEvent(peer *peers.ClientPeer, actorNr int32) {
	if evn, err := binbuf.BuildBuffer(1024); err != nil {
		handleErrors(peer.PeerId(), 0, err)
		return
	} else {
		evn.PushByte(1)
		evn.PushByte(evncode.Leave)

		evn.PushByte(paramcode.ActorNr)
		evn.PushObject(actorNr)

		actorNrs := insRoom(peer.RoomId()).ActorsManager().GetAllActorNrs()
		evn.PushByte(paramcode.Actors)
		evn.PushObject(actorNrs)

		evn.PushByte(paramcode.IsInactive)
		evn.PushObject(false)

		ids := insRoom(peer.RoomId()).ActorsManager().GetAllPeerIds()
		for _, item := range ids {
			if item != peer.PeerId() {
				fmt.Println(evn.Bytes())
				gonode.Send(item, evn.Bytes())
			}
		}
	}
}

func pubDisconnectEvent(peer *peers.ClientPeer, actorNr int32) {
	if evn, err := binbuf.BuildBuffer(1024); err != nil {
		handleErrors(peer.PeerId(), 0, err)
		return
	} else {
		evn.PushByte(1)
		evn.PushByte(evncode.Leave)

		evn.PushByte(paramcode.ActorNr)
		evn.PushObject(actorNr)

		actorNrs := make([]int32, 0, 0)
		evn.PushByte(paramcode.Actors)
		evn.PushObject(actorNrs)

		evn.PushByte(paramcode.IsInactive)
		evn.PushObject(true)

		ids := insRoom(peer.RoomId()).ActorsManager().GetAllPeerIds()
		for _, item := range ids {
			if item != peer.PeerId() {
				fmt.Println(evn.Bytes())
				gonode.Send(item, evn.Bytes())
			}
		}
	}
}

func pubJoinEvent(peer *peers.ClientPeer, opCode byte, parser *binbuf.BinParser, actorNr int32) {
	if evn, err := binbuf.BuildBuffer(1024); err != nil {
		handleErrors(peer.PeerId(), opCode, err)
		return
	} else {

		// handle the cache events publist
		pubEventCache(peer)

		evn.PushByte(1) // event
		evn.PushByte(evncode.Join)

		hashTable := stl.NewHashTable()
		hashTable.Set(actorparam.Nickname, "")
		evn.PushByte(paramcode.ActorProperties)
		evn.PushObject(hashTable.KeyValuePairs())

		evn.PushByte(paramcode.ActorNr)
		evn.PushObject(actorNr)

		actorNrs := insRoom(peer.RoomId()).ActorsManager().GetAllActorNrs()
		evn.PushByte(paramcode.Actors)
		evn.PushObject(actorNrs)

		ids := insRoom(peer.RoomId()).ActorsManager().GetAllPeerIds()
		for _, item := range ids {
			//			fmt.Println(evn.Bytes())
			gonode.Send(item, evn.Bytes())
		}

		fmt.Print("cache the join event:")
		fmt.Print(" actorNr:")
		fmt.Println(actorNr)
		evnCache(peer).AddEvent(actorNr, evncode.Join, evn.Bytes())
	}
}

func sendRemoveRoomState(roomId string) {
	if buf, err := binbuf.BuildBuffer(256); err != nil {
		handleErrors("lobby", servereventcode.RemoveGameState, err)
		return
	} else {
		buf.PushByte(servereventcode.RemoveGameState)
		buf.PushString(roomId)
		gonode.Send("lobby", buf.Bytes())
	}
}
