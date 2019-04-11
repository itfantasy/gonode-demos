package lobby

import (
	"errors"

	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/redis"
	"github.com/itfantasy/gonode/utils/crypt"
	"github.com/itfantasy/gonode/utils/rand"
)

var coreRedis *redis.Redis

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
	return coreRedis.HSet("Lobby-"+this.guid, roomid, info)
}

func (this *Lobby) RemoveRoomState(roomid string) (bool, error) {
	return coreRedis.HDel("Lobby-"+this.guid, roomid)
}

func (this *Lobby) GetRoomState(roomid string) (string, error) {
	return coreRedis.HGet("Lobby-"+this.guid, roomid)
}

func (this *Lobby) RandomRoomStateId() (string, bool) {
	keys, err := coreRedis.HKeys("Lobby-" + this.guid)
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

func RegisterCoreRedis(redisConf string) error {
	comp, err := components.NewComponent(redisConf)
	if err != nil {
		return err
	}

	red, ok := comp.(*redis.Redis)
	if !ok {
		return errors.New("redis comp init faild!")
	}

	coreRedis = red
	return nil
}
