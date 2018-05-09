// event cache
package room

import (
	"strconv"

	"github.com/itfantasy/gonode/utils/stl"
)

type CustomEvent struct {
	actorNr int32
	Code    byte
	Data    []byte
}

func NewCustomEvent(actor int32, eventCode byte, data []byte) *CustomEvent {
	event := new(CustomEvent)
	event.actorNr = actor
	event.Code = eventCode
	event.Data = make([]byte, 0, len(data))
	copy(event.Data, data)
	return event
}

type RoomEventCache struct {
	list *stl.List // *stl.List
}

func NewRoomEventCache() *RoomEventCache {
	this := new(RoomEventCache)
	this.list = stl.NewList(50)
	return this
}

func (this *RoomEventCache) AddEvent(actor int32, eventCode byte, data []byte) {
	this.list.Add(NewCustomEvent(actor, eventCode, data))
}

func (this *RoomEventCache) RemoveEventsByActor(actor int32) int {
	dirtyList := stl.NewList(10)
	for _, item := range this.list.Values() {
		customeEvent := item.(CustomEvent)
		if customeEvent.actorNr == actor {
			dirtyList.Add(customeEvent)
		}
	}
	for _, item := range dirtyList.Values() {
		this.list.Remove(item)
	}
	return dirtyList.Count()
}

func (this *RoomEventCache) Events() []interface{} {
	return this.list.Values()
}

type RoomEventCacheManager struct {
	_map *stl.Dictionary // roomid=>RoomEventCache
}

var _roomEventCacheManager *RoomEventCacheManager

func InsRoomEventCacheManager() *RoomEventCacheManager {
	if _roomEventCacheManager == nil {
		_roomEventCacheManager = new(RoomEventCacheManager)
		_roomEventCacheManager._map = stl.NewDictionary()
	}
	return _roomEventCacheManager
}

func (this *RoomEventCacheManager) FetchRoomCache(roomId int32) *RoomEventCache {
	key := strconv.Itoa(int(roomId))
	val, exist := this._map.Get(key)
	if exist {
		return val.(*RoomEventCache)
	} else {
		var cache *RoomEventCache = NewRoomEventCache()
		this._map.Set(key, cache)
		return cache
	}
}
