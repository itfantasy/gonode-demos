package main

import (
	"errors"
	"fmt"
	//	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/gnbuffers"
	//	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/stl"
	"github.com/itfantasy/gonode_demo/icloud/opcode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/actorparam"
	"github.com/itfantasy/gonode_demo/icloud/opcode/cacheop"
	"github.com/itfantasy/gonode_demo/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/evncode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/gameparam"
	"github.com/itfantasy/gonode_demo/icloud/opcode/paramcode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/recvgroup"

	"github.com/itfantasy/gonode_demo/icloud/lib/room"
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
		case opcode.RaiseEvent:
			handleRaiseEvent(id, opCode, parser)
			break
		case opcode.SetProperties:
			handleSetProperties(id, opCode, parser)
		default:
			gonode.Send(id, msg)
			break
		}
	}
}

var actorNr int32 = 1
var _insRoom *room.Room = nil

func insRoom() *room.Room {
	if _insRoom == nil {
		_insRoom = room.NewRoom(3001)
	}
	return _insRoom
}

func handleErrors(id string, opCode byte, err error) {
	//gonode.Error(err.Error())
	fmt.Print("[ERROR]::")
	fmt.Println(err)
}

func handleSetProperties(id string, opCode byte, parser *gnbuffers.GnParser) {
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
		buf.PushObject(hash.KeyValuePairs())

		suc := false
		actor, err := insRoom().ActorsManager().AddNewActor(id, actorNr)
		if err != nil {
			handleErrors(id, opCode, err)
			return
		} else {
			suc = true
		}
		actorNrs := insRoom().ActorsManager().GetAllActorNrs()
		buf.PushByte(paramcode.Actors)
		buf.PushObject(actorNrs)

		gonode.Send(id, buf.Bytes())
		if suc {
			pubJoinEvent(id, opCode, parser, actor.ActorNr)
		}
		actorNr += 1
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
		buf.PushObject(actorNr)

		//buf.PushByte(paramcode.ActorProperties)
		//hash := stl.NewHashTable()
		//buf.PushObject(hash.KeyValuePairs())

		buf.PushByte(paramcode.GameProperties)
		hash := stl.NewHashTable()
		hash.Set(gameparam.LobbyProperties, true)
		hash.Set(gameparam.CleanupCacheOnLeave, true)
		hash.Set(gameparam.MaxPlayers, 4)
		hash.Set(gameparam.IsVisible, true)
		hash.Set(gameparam.IsOpen, true)
		hash.Set(gameparam.MasterClientId, actorNr)
		buf.PushObject(hash.KeyValuePairs())

		suc := false
		actor, err := insRoom().ActorsManager().AddNewActor(id, actorNr)
		if err != nil {
			handleErrors(id, opCode, err)
			return
		} else {
			suc = true
		}
		actorNrs := insRoom().ActorsManager().GetAllActorNrs()
		buf.PushByte(paramcode.Actors)
		buf.PushObject(actorNrs)

		gonode.Send(id, buf.Bytes())
		if suc {
			pubJoinEvent(id, opCode, parser, actor.ActorNr)
		}
		actorNr += 1
	}
}

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
		actor, exist := insRoom().ActorsManager().GetActorByPeerId(id)
		if !exist {
			handleErrors(id, opCode, errors.New("cannot find the actor by the peeerid:"+id))
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
				handleErrors(id, opCode, err)
				return
			}

			//fmt.Print("paramCode:")
			//fmt.Println(paramCode)
			if paramCode == paramcode.ReceiverGroup { // ReceiverGroup
				if oRecvGroup, err := parser.Object(); err != nil {
					handleErrors(id, opCode, err)
					return
				} else {
					recvGroup = oRecvGroup.(byte)
				}
				//parser.Byte() // Data
			} else if paramCode == paramcode.Cache { // the event which will be cached
				if oCacheOp, err := parser.Object(); err != nil {
					handleErrors(id, opCode, err)
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
			if actor, exist := insRoom().ActorsManager().GetActorByPeerId(id); exist {
				evnCache().AddEvent(actor.ActorNr, eventCode, evn.Bytes())
			}
		}

		// handle the recvGroup
		ids := gonode.Node().NetWorker().GetAllConnIds() // get the ids in the same room
		if recvGroup == recvgroup.MasterClient {
			gonode.Send(id, evn.Bytes())
		} else if recvGroup == recvgroup.All {
			for _, item := range ids {
				//fmt.Println(evn.Bytes())
				gonode.Send(item, evn.Bytes())
			}
		} else {
			for _, item := range ids {
				if item != id {
					//fmt.Println(evn.Bytes())
					gonode.Send(item, evn.Bytes())
				}
			}
		}

	}
}

func evnCache() *room.RoomEventCache {
	return insRoom().EventCache()
}

//HiveGame.PublishEventCache (line 1410)
func pubEventCache(id string) {
	for _, item := range evnCache().Events() {
		event := item.(*room.CustomEvent)
		fmt.Print("pub the cache event:")
		fmt.Print(" actorNr:")
		fmt.Print(event.ActorNr)
		fmt.Print(" Code:")
		fmt.Println(event.Code)
		fmt.Print("Data:")
		fmt.Println(event.Data)
		gonode.Send(id, event.Data)
	}
}

func pubJoinEvent(id string, opCode byte, parser *gnbuffers.GnParser, actorNr int32) {
	if evn, err := gnbuffers.BuildBuffer(1024); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {

		// handle the cache events publist
		pubEventCache(id)

		evn.PushByte(1) // event
		evn.PushByte(evncode.Join)

		hashTable := stl.NewHashTable()
		hashTable.Set(actorparam.Nickname, "")
		evn.PushByte(paramcode.ActorProperties)
		evn.PushObject(hashTable.KeyValuePairs())

		evn.PushByte(paramcode.ActorNr)
		evn.PushObject(actorNr)

		actorNrs := insRoom().ActorsManager().GetAllActorNrs()
		evn.PushByte(paramcode.Actors)
		evn.PushObject(actorNrs)

		ids := gonode.Node().NetWorker().GetAllConnIds()
		for _, item := range ids {
			//			fmt.Println(evn.Bytes())
			gonode.Send(item, evn.Bytes())
		}

		fmt.Print("cache the join event:")
		fmt.Print(" actorNr:")
		fmt.Println(actorNr)
		evnCache().AddEvent(actorNr, evncode.Join, evn.Bytes())
	}
}
