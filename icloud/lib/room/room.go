package room

import (
	"errors"

	"github.com/itfantasy/gonode/utils/stl"
)

// the room struct
type Room struct {
	Name           string
	MasterClientId int32
	actorsManager  *ActorsManager
	eventCache     *RoomEventCache
}

func NewRoom(name string) *Room {
	room := new(Room)
	room.Name = name
	room.MasterClientId = 0
	room.actorsManager = NewActorsManager()
	room.eventCache = NewRoomEventCache()
	return room
}

func (this *Room) ActorsManager() *ActorsManager {
	return this.actorsManager
}

func (this *Room) EventCache() *RoomEventCache {
	return this.eventCache
}

func (this *Room) IsEmpty() bool {
	return this.actorsManager.ActorsCount() <= 0
}

type RoomManager struct {
	dict *stl.Dictionary
}

func NewRoomManager() *RoomManager {
	roomManager := new(RoomManager)
	roomManager.dict = stl.NewDictionary()
	return roomManager
}

func (this *RoomManager) CreateRoom(name string) (*Room, error) {
	if this.dict.ContainsKey(name) {
		return nil, errors.New("the roommanager has contained a room with the same name:" + name)
	}
	room := NewRoom(name)
	this.dict.Set(name, room)
	return room, nil
}

func (this *RoomManager) FetchRoom(name string) *Room {
	item, exist := this.dict.Get(name)
	if exist {
		return item.(*Room)
	} else {
		room := NewRoom(name)
		this.dict.Set(name, room)
		return room
	}
}

func (this *RoomManager) DisposeRoom(name string) {
	// need dispose the actorsManager and the eventCache
	for _, val := range this.dict.KeyValuePairs() {
		room := val.(*Room)
		room.actorsManager.ClearAll()
		room.eventCache.ClearCache()
	}
}
