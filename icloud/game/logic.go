package main

import (
	"fmt"
	//	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/gnbuffers"
	//	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/stl"
	"github.com/itfantasy/gonode_demo/icloud/opcode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/actorparam"
	"github.com/itfantasy/gonode_demo/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/evncode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/gameparam"
	"github.com/itfantasy/gonode_demo/icloud/opcode/paramcode"
	"github.com/itfantasy/gonode_demo/icloud/opcode/recvgroup"
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

func handleErrors(id string, opCode byte, err error) {
	gonode.Error(err.Error())
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

		buf.PushByte(paramcode.Actors)
		list := stl.NewListInt(10)
		list.Add(actorNr)
		buf.PushObject(list.Values())

		gonode.Send(id, buf.Bytes())
		pubJoinEvent(id, opCode, parser)
		//actorNr += 1
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

		buf.PushByte(paramcode.Actors)
		list := stl.NewListInt(10)
		list.Add(actorNr)
		buf.PushObject(list.Values())

		gonode.Send(id, buf.Bytes())
		pubJoinEvent(id, opCode, parser)
		//actorNr += 1
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
		parser.Byte()                 // ParameterCode.Code
		parser.Byte()                 // gntypes.Byte
		eventCode, _ := parser.Byte() // eventCode

		evn.PushByte(1) // event
		evn.PushByte(eventCode)

		evn.PushByte(paramcode.ActorNr)
		evn.PushObject(actorNr)
		//actorNr += 1

		evn.PushByte(paramcode.Code)
		evn.PushObject(byte(eventCode))

		parser.Byte() // ReceiverGroup
		var recvGroup byte = recvgroup.Others
		if oRecvGroup, err := parser.Object(); err != nil {
			handleErrors(id, opCode, err)
			return
		} else {
			recvGroup = oRecvGroup.(byte)
		}

		evn.PushByte(paramcode.Data)
		data := parser.Bytes()
		evn.PushBytes(data[5:])

		ids := gonode.Node().NetWorker().GetAllConnIds() // 获得同房间的id列表

		if recvGroup == recvgroup.MasterClient {
			gonode.Send(id, evn.Bytes())
		} else if recvGroup == recvgroup.All {
			for _, item := range ids {
				fmt.Println(evn.Bytes())
				gonode.Send(item, evn.Bytes())
			}
		} else {
			for _, item := range ids {
				if item != id {
					fmt.Println(evn.Bytes())
					gonode.Send(item, evn.Bytes())
				}
			}
		}
	}
}

func pubJoinEvent(id string, opCode byte, parser *gnbuffers.GnParser) {
	if evn, err := gnbuffers.BuildBuffer(1024); err != nil {
		handleErrors(id, opCode, err)
		return
	} else {

		evn.PushByte(1) // event
		evn.PushByte(evncode.Join)

		hashTable := stl.NewHashTable()
		hashTable.Set(actorparam.Nickname, "")
		evn.PushByte(paramcode.ActorProperties)
		evn.PushObject(hashTable.KeyValuePairs())

		evn.PushByte(paramcode.ActorNr)
		evn.PushObject(actorNr)

		intArray := stl.NewListInt(10)
		intArray.Add(actorNr)
		evn.PushByte(paramcode.Actors)
		evn.PushObject(intArray.Values())

		ids := gonode.Node().NetWorker().GetAllConnIds()
		for _, item := range ids {
			fmt.Println(evn.Bytes())
			gonode.Send(item, evn.Bytes())
		}
	}
}
