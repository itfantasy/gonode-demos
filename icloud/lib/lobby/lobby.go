package lobby

import (
	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/utils/crypt"
	"github.com/itfantasy/gonode/utils/rand"
)

// the lobby struct
type Lobby struct {
	guid string
}

func NewLobby(guid string) *Lobby {
	lobby := new(Lobby)
	lobby.guid = guid
	return lobby
}

func GenerateRoomId() string {
	return "game" + crypt.C32()
}

func (this *Lobby) CreateRoomState(roomid string, info string) (bool, error) {
	return gonode.Node().CoreRedis().HSet("Lobby-"+this.guid, roomid, info)
}

func (this *Lobby) RemoveRoomState(roomid string) (bool, error) {
	return gonode.Node().CoreRedis().HDel("Lobby-"+this.guid, roomid)
}

func (this *Lobby) GetRoomState(roomid string) (string, error) {
	return gonode.Node().CoreRedis().HGet("Lobby-"+this.guid, roomid)
}

func (this *Lobby) RandomRoomStateId() (string, bool) {
	keys, err := gonode.Node().CoreRedis().HKeys("Lobby-" + this.guid)
	if err != nil {
		return "", false
	}
	l := len(keys)
	if l <= 0 {
		return "", false
	}
	i := rand.Random(1, l)
	return keys[i-1], true
}
