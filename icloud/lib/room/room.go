package room

// the room struct
type Room struct {
	RoomId        int32
	actorsManager *ActorsManager
	eventCache    *RoomEventCache
}

func NewRoom(RoomId int32) *Room {
	room := new(Room)
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
