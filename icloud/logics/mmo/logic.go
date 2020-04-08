package mmo

import (
	"errors"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-toolkit/toolkit/gen_mmo"

	"github.com/itfantasy/gonode-icloud/icloud/gunpeer"
	"github.com/itfantasy/gonode-icloud/icloud/opcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/errorcode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/evncode"
	"github.com/itfantasy/gonode-icloud/icloud/opcode/paramcode"
)

type MmoHandler struct {
}

func (m *MmoHandler) HandleConn(id string) {
	gen_mmo.AddPeer(gen_mmo.NewMmoPeer(id, m))
}

func (m *MmoHandler) HandleClose(id string) {
	gen_mmo.RemovePeer(id)
}

func (m *MmoHandler) HandleMsg(id string, msg []byte) {
	gonode.Debug(msg)
	opCode, datas, err := gunpeer.ParseMsg(msg)
	if err != nil {
		gonode.LogError(err)
		return
	}
	peer, ok := gen_mmo.GetPeer(id)
	if !ok {
		gonode.LogError(errors.New("peer missing..." + id))
		return
	}
	switch opCode {
	case opcode.CreateWorld:
		m.handleCreateWorld(peer, opCode, datas)
		break
	case opcode.EnterWorld:
		m.handleEnterWorld(peer, opCode, datas)
		break
	case opcode.ExitWorld:
		m.handleExitWorld(peer, opCode, datas)
		break
	case opcode.Move:
		m.handleMove(peer, opCode, datas)
		break
	case opcode.RaiseGenericEvent:
		m.handleRaiseGenericEvent(peer, opCode, datas)
		break
	case opcode.SetItemProperties:
		m.handleSetProperties(peer, opCode, datas)
		break
	case opcode.SpawnItem:
		m.handleSpawnItem(peer, opCode, datas)
		break
	case opcode.DestroyItem:
		m.handleDestroyItem(peer, opCode, datas)
		break
	case opcode.SubscribeItem:
		m.handleSubscribeItem(peer, opCode, datas)
		break
	case opcode.UnsubscribeItem:
		m.handleUnsubscribeItem(peer, opCode, datas)
		break
	case opcode.SetViewDistance:
		m.handleSetViewDistance(peer, opCode, datas)
		break
	case opcode.AttachInterestArea:
		m.handleAttachInterestArea(peer, opCode, datas)
		break
	case opcode.DetachInterestArea:
		m.handleDetachInterestArea(peer, opCode, datas)
		break
	case opcode.AddInterestArea:
		m.handleAddInterestArea(peer, opCode, datas)
		break
	case opcode.RemoveInterestArea:
		m.handleRemoveInterestArea(peer, opCode, datas)
		break
	case opcode.GetItemProperties:
		m.handleGetProperties(peer, opCode, datas)
		break
	case opcode.MoveInterestArea:
		m.handleMoveInterestArea(peer, opCode, datas)
		break
	case opcode.RadarSubscribe:
		m.handleRadarSubscribe(peer, opCode, datas)
		break
	default:
		break
	}
}

func GetVector(p *gunpeer.PeerDatas, key byte) (*gen_mmo.Vector, bool) {
	val, ok := p.Get(key)
	if !ok {
		return nil, false
	}
	result, ok := val.(*gen_mmo.Vector)
	if !ok {
		return nil, false
	}
	return result, true
}

func GetBoundingBox(p *gunpeer.PeerDatas, key byte) (*gen_mmo.BoundingBox, bool) {
	val, ok := p.Get(key)
	if !ok {
		return nil, false
	}
	result, ok := val.(*gen_mmo.BoundingBox)
	if !ok {
		return nil, false
	}
	return result, true
}

func (m *MmoHandler) handleError(peer *gen_mmo.MmoPeer, opCode byte, err error) {
	errcode, _ := gonode.ErrorInfo(err)
	if errcode != 0 {
		gunpeer.SendResponse(peer.PeerId(), int16(errcode), opCode, map[byte]interface{}{})
	}
	gonode.LogError(err.Error(), gonode.LogSource(1))
}

func (m *MmoHandler) handleCreateWorld(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	worldName, _ := datas.GetString(paramcode.WorldName)
	boundingBox, _ := GetBoundingBox(datas, paramcode.BoundingBox)
	tileDimensions, _ := GetVector(datas, paramcode.TileDimensions)

	_, err := gen_mmo.CreateWorld(peer.PeerId(), worldName, boundingBox, tileDimensions)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.WorldName: worldName,
		})
	}
}

func (m *MmoHandler) handleEnterWorld(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	usrName, _ := datas.GetString(paramcode.Username)
	worldName, _ := datas.GetString(paramcode.WorldName)
	interestAreaId, _ := datas.GetByte(paramcode.InterestAreaId)
	position, _ := GetVector(datas, paramcode.NewPosition)
	rotation, _ := GetVector(datas, paramcode.Rotation)
	viewDistanceEnter, _ := GetVector(datas, paramcode.ViewDistanceEnter)
	viewDistanceExit, _ := GetVector(datas, paramcode.ViewDistanceExit)
	properties, _ := datas.GetHash(paramcode.Properties)

	world, _, err := gen_mmo.EnterWorld(peer.PeerId(), usrName, worldName, interestAreaId, position, rotation, viewDistanceEnter, viewDistanceExit, properties)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.WorldName:      world.Name(),
			paramcode.BoundingBox:    world.Area(),
			paramcode.TileDimensions: world.TileDimensions(),
		})
	}
}

func (m *MmoHandler) handleExitWorld(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	if err := gen_mmo.ExitWorld(peer.PeerId()); err != nil {
		m.handleError(peer, opCode, err)
	}
}

func (m *MmoHandler) handleMove(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	itemId, _ := datas.GetString(paramcode.ItemId)
	position, _ := GetVector(datas, paramcode.NewPosition)
	rotation, _ := GetVector(datas, paramcode.Rotation)

	item, err := gen_mmo.Move(peer.PeerId(), itemId, position, rotation)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.ItemId: item.Id(),
		})
	}
}

func (m *MmoHandler) handleRaiseGenericEvent(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	itemId, _ := datas.GetString(paramcode.ItemId)
	evnCode, _ := datas.GetByte(paramcode.CustomEventCode)
	eventReceiver, _ := datas.GetByte(paramcode.EventReceiver)
	var evnData []byte = nil // fetch the left buff from the datas

	item, err := gen_mmo.RaiseGenericEvent(peer.PeerId(), itemId, evnCode, evnData, eventReceiver)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.ItemId: item.Id(),
		})
	}
}

func (m *MmoHandler) handleSetProperties(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	itemId, _ := datas.GetString(paramcode.ItemId)
	propertiesSet, _ := datas.GetHash(paramcode.PropertiesSet)
	propertiesUnset, _ := datas.GetArray(paramcode.PropertiesUnset)

	item, err := gen_mmo.SetProperties(peer.PeerId(), itemId, propertiesSet, propertiesUnset)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.ItemId: item.Id(),
		})
	}
}

func (m *MmoHandler) handleSpawnItem(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	itemId, _ := datas.GetString(paramcode.ItemId)
	itemType, _ := datas.GetByte(paramcode.ItemType)
	position, _ := GetVector(datas, paramcode.NewPosition)
	rotation, _ := GetVector(datas, paramcode.Rotation)
	properties, _ := datas.GetHash(paramcode.Properties)

	item, err := gen_mmo.SpawnItem(peer.PeerId(), itemId, itemType, position, rotation, properties)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.ItemId: item.Id(),
		})
	}
}

func (m *MmoHandler) handleDestroyItem(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	itemId, _ := datas.GetString(paramcode.ItemId)

	err := gen_mmo.DestroyItem(peer.PeerId(), itemId)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.ItemId: itemId,
		})
	}
}

func (m *MmoHandler) handleSubscribeItem(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	itemId, _ := datas.GetString(paramcode.ItemId)
	propertiesRevision, _ := datas.GetInt(paramcode.PropertiesRevision)

	item, err := gen_mmo.SubscribeItem(peer.PeerId(), itemId, int(propertiesRevision))
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.ItemId: item.Id(),
		})
	}
}

func (m *MmoHandler) handleUnsubscribeItem(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	itemId, _ := datas.GetString(paramcode.ItemId)

	item, err := gen_mmo.UnsubscribeItem(peer.PeerId(), itemId)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.InterestAreaId: byte(0),
			paramcode.ItemId:         item.Id(),
		})
	}
}

func (m *MmoHandler) handleSetViewDistance(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	interestAreaId, _ := datas.GetByte(paramcode.InterestAreaId)
	viewDistanceEnter, _ := GetVector(datas, paramcode.ViewDistanceEnter)
	viewDistanceExit, _ := GetVector(datas, paramcode.ViewDistanceExit)

	interestArea, err := gen_mmo.SetViewDistance(peer.PeerId(), interestAreaId, viewDistanceEnter, viewDistanceExit)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.InterestAreaId: interestArea.Id(),
		})
	}
}

func (m *MmoHandler) handleAttachInterestArea(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	interestAreaId, _ := datas.GetByte(paramcode.InterestAreaId)
	itemId, _ := datas.GetString(paramcode.ItemId)

	interestArea, err := gen_mmo.AttachInterestArea(peer.PeerId(), interestAreaId, itemId)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.InterestAreaId: interestArea.Id(),
			paramcode.ItemId:         interestArea.AttachedItem().Id(),
		})
	}
}

func (m *MmoHandler) handleDetachInterestArea(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	interestAreaId, _ := datas.GetByte(paramcode.InterestAreaId)

	interestArea, err := gen_mmo.DetachInterestArea(peer.PeerId(), interestAreaId)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.InterestAreaId: interestArea.Id(),
		})
	}
}

func (m *MmoHandler) handleAddInterestArea(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	interestAreaId, _ := datas.GetByte(paramcode.InterestAreaId)
	itemId, _ := datas.GetString(paramcode.ItemId)
	position, _ := GetVector(datas, paramcode.NewPosition)
	viewDistanceEnter, _ := GetVector(datas, paramcode.ViewDistanceEnter)
	viewDistanceExit, _ := GetVector(datas, paramcode.ViewDistanceExit)

	interestArea, err := gen_mmo.AddInterestArea(peer.PeerId(), interestAreaId, itemId, position, viewDistanceEnter, viewDistanceExit)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.InterestAreaId: interestArea.Id(),
		})
	}
}

func (m *MmoHandler) handleRemoveInterestArea(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	interestAreaId, _ := datas.GetByte(paramcode.InterestAreaId)

	err := gen_mmo.RemoveInterestArea(peer.PeerId(), interestAreaId)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.InterestAreaId: interestAreaId,
		})
	}
}

func (m *MmoHandler) handleGetProperties(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	itemId, _ := datas.GetString(paramcode.ItemId)
	propertiesRevision, _ := datas.GetInt(paramcode.PropertiesRevision)

	properties, err := gen_mmo.GetProperties(peer.PeerId(), itemId, int(propertiesRevision))
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		if properties.Updated {
			gunpeer.SendEvent(peer.PeerId(), evncode.ItemProperties, map[byte]interface{}{
				paramcode.ItemId:             properties.ItemId,
				paramcode.PropertiesRevision: properties.PropertiesRevision,
				paramcode.PropertiesSet:      properties.PropertiesSet,
			})
		}
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.ItemId: itemId,
		})
	}
}

func (m *MmoHandler) handleMoveInterestArea(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	interestAreaId, _ := datas.GetByte(paramcode.InterestAreaId)
	position, _ := GetVector(datas, paramcode.NewPosition)

	interestArea, err := gen_mmo.MoveInterestArea(peer.PeerId(), interestAreaId, position)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.InterestAreaId: interestArea.Id(),
		})
	}
}

func (m *MmoHandler) handleRadarSubscribe(peer *gen_mmo.MmoPeer, opCode byte, datas *gunpeer.PeerDatas) {
	worldName, _ := datas.GetString(paramcode.WorldName)

	world, err := gen_mmo.RadarSubscribe(peer.PeerId(), worldName)
	if err != nil {
		m.handleError(peer, opCode, err)
	} else {
		gunpeer.SendResponse(peer.PeerId(), errorcode.Ok, opCode, map[byte]interface{}{
			paramcode.WorldName:      world.Name(),
			paramcode.BoundingBox:    world.Area(),
			paramcode.TileDimensions: world.TileDimensions(),
		})
	}
}

func (m *MmoHandler) OnItemGenericEvent(peer *gen_mmo.MmoPeer, event *gen_mmo.ItemGeneric) {
	gunpeer.SendEvent(peer.PeerId(), evncode.ItemGeneric, map[byte]interface{}{
		paramcode.EventCode: event.CustomEventCode,
		paramcode.EventData: event.EventData,
	})
}
func (m *MmoHandler) OnItemDestroyed(peer *gen_mmo.MmoPeer, itemId string) {
	gunpeer.SendEvent(peer.PeerId(), evncode.ItemDestroyed, map[byte]interface{}{
		paramcode.ItemId: itemId,
	})
}
func (m *MmoHandler) OnItemMoved(peer *gen_mmo.MmoPeer, event *gen_mmo.ItemMoved) {
	gunpeer.SendEvent(peer.PeerId(), evncode.ItemMoved, map[byte]interface{}{
		paramcode.ItemId:      event.ItemId,
		paramcode.OldPosition: event.OldPosition,
		paramcode.NewPosition: event.Position,
		paramcode.OldRotation: event.OldRotation,
		paramcode.Rotation:    event.Rotation,
	})
}
func (m *MmoHandler) OnItemPropertiesSet(peer *gen_mmo.MmoPeer, event *gen_mmo.ItemPropertiesSet) {
	gunpeer.SendEvent(peer.PeerId(), evncode.ItemPropertiesSet, map[byte]interface{}{
		paramcode.ItemId:             event.ItemId,
		paramcode.PropertiesRevision: event.PropertiesRevision,
		paramcode.PropertiesSet:      event.PropertiesSet,
		paramcode.PropertiesUnset:    event.PropertiesUnset,
	})
}
func (m *MmoHandler) OnWorldExited(peer *gen_mmo.MmoPeer, worldName string) {
	gunpeer.SendEvent(peer.PeerId(), evncode.WorldExited, map[byte]interface{}{
		paramcode.WorldName: worldName,
	})
}
func (m *MmoHandler) OnItemSubscribed(peer *gen_mmo.MmoPeer, event *gen_mmo.ItemSubscribed) {
	gunpeer.SendEvent(peer.PeerId(), evncode.ItemSubscribed, map[byte]interface{}{
		paramcode.InterestAreaId:     event.InterestAreaId,
		paramcode.ItemId:             event.ItemId,
		paramcode.ItemType:           event.ItemType,
		paramcode.NewPosition:        event.Position,
		paramcode.Rotation:           event.Rotation,
		paramcode.PropertiesRevision: event.PropertiesRevision,
	})
}
func (m *MmoHandler) OnItemUnsubscribed(peer *gen_mmo.MmoPeer, event *gen_mmo.ItemUnsubscribed) {
	gunpeer.SendEvent(peer.PeerId(), evncode.ItemUnsubscribed, map[byte]interface{}{
		paramcode.InterestAreaId: event.InterestAreaId,
		paramcode.ItemId:         event.ItemId,
	})
}
func (m *MmoHandler) OnRadarUpdate(peer *gen_mmo.MmoPeer, event *gen_mmo.RadarUpdate) {
	gunpeer.SendEvent(peer.PeerId(), evncode.RadarUpdate, map[byte]interface{}{
		paramcode.ItemId:      event.ItemId,
		paramcode.ItemType:    event.ItemType,
		paramcode.NewPosition: event.Position,
		paramcode.Remove:      event.Remove,
	})
}
