package gunpeer

import (
	"errors"
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-toolkit/toolkit/gen_room"
	"github.com/itfantasy/gonode/core/binbuf"

	"github.com/itfantasy/gonode-icloud/icloud/opcode/gameparam"
)

const _errconst string = "type matching failed.. "

type PeerDatas struct {
	_raw    []byte
	_map    map[byte]interface{}
	_errMsg string
	_err    error
}

func NewPeerDatas(raw []byte) *PeerDatas {
	p := new(PeerDatas)
	p._raw = raw
	p._map = make(map[byte]interface{})
	p._errMsg = ""
	p._err = nil
	return p
}

func (p *PeerDatas) Set(key byte, val interface{}) {
	p._map[key] = val
}

func (p *PeerDatas) Get(key byte) (interface{}, bool) {
	ret, exist := p._map[key]
	if !exist {
		p._errMsg = "key not exist.. " + fmt.Sprint(key)
	}
	return ret, true
}

func (p *PeerDatas) RawBytes() []byte {
	return p._raw
}

func (p *PeerDatas) Err() error {
	if p._err == nil && p._errMsg != "" {
		p._err = errors.New(p._errMsg)
	}
	return p._err
}

func (p *PeerDatas) GetBool(key byte) (bool, bool) {
	val, ok := p.Get(key)
	if !ok {
		return false, false
	}
	boolVal, ok := val.(bool)
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "(bool)"
		return false, false
	}
	return boolVal, true
}

func (p *PeerDatas) GetByte(key byte) (byte, bool) {
	val, ok := p.Get(key)
	if !ok {
		return 0, false
	}
	byteVal, ok := val.(byte)
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "(byte)"
		return 0, false
	}
	return byteVal, true
}

func (p *PeerDatas) GetShort(key byte) (int16, bool) {
	val, ok := p.Get(key)
	if !ok {
		return 0, false
	}
	int16Val, ok := val.(int16)
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "(int16)"
		return 0, false
	}
	return int16Val, true
}

func (p *PeerDatas) GetInt(key byte) (int32, bool) {
	val, ok := p.Get(key)
	if !ok {
		return 0, false
	}
	int32Val, ok := val.(int32)
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "(int32)"
		return 0, false
	}
	return int32Val, true
}

func (p *PeerDatas) GetLong(key byte) (int64, bool) {
	val, ok := p.Get(key)
	if !ok {
		return 0, false
	}
	int64Val, ok := val.(int64)
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "(int64)"
		return 0, false
	}
	return int64Val, true
}

func (p *PeerDatas) GetString(key byte) (string, bool) {
	val, ok := p.Get(key)
	if !ok {
		return "", false
	}
	stringVal, ok := val.(string)
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "(string)"
		return "", false
	}
	return stringVal, true
}

func (p *PeerDatas) GetFloat(key byte) (float32, bool) {
	val, ok := p.Get(key)
	if !ok {
		return 0, false
	}
	float32Val, ok := val.(float32)
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "(float32)"
		return 0, false
	}
	return float32Val, true
}

func (p *PeerDatas) GetInts(key byte) ([]int32, bool) {
	val, ok := p.Get(key)
	if !ok {
		return nil, false
	}
	intsVal, ok := val.([]int32)
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "([]int32)"
		return nil, false
	}
	return intsVal, true
}

func (p *PeerDatas) GetArray(key byte) ([]interface{}, bool) {
	val, ok := p.Get(key)
	if !ok {
		return nil, false
	}
	arrayVal, ok := val.([]interface{})
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "([]interface{})"
		return nil, false
	}
	return arrayVal, true
}

func (p *PeerDatas) GetHash(key byte) (map[interface{}]interface{}, bool) {
	val, ok := p.Get(key)
	if !ok {
		return nil, false
	}
	hashVal, ok := val.(map[interface{}]interface{})
	if !ok {
		p._errMsg = _errconst + fmt.Sprint(key) + "(map[interface{}]interface{})"
		return nil, false
	}
	return hashVal, true
}

func ParseMsg(msg []byte) (byte, *PeerDatas, error) {
	datas := NewPeerDatas(msg)
	parser := binbuf.BuildParser(msg, 0)
	opCode := parser.Byte()
	for !parser.OverFlow() {
		key := parser.Byte()
		val := parser.Object()
		datas.Set(key, val)
	}
	return opCode, datas, parser.Error()
}

func SendResponse(peerId string, retCode int16, opCode byte, datas map[byte]interface{}) {
	buf := binbuf.BuildBuffer(1024)
	buf.PushByte(0)
	buf.PushShort(retCode)
	buf.PushByte(opCode)
	if datas != nil {
		for k, v := range datas {
			buf.PushByte(k)
			buf.PushObject(v)
		}
	}
	bytes, err := buf.Bytes()
	if err != nil {
		gonode.LogError(err)
		return
	}
	if err := gonode.Send(peerId, bytes); err != nil {
		gonode.LogError(err)
	}
}

func EventDatas(evnCode byte, datas map[byte]interface{}) ([]byte, error) {
	buf := binbuf.BuildBuffer(1024)
	buf.PushByte(1)
	buf.PushByte(evnCode)
	if datas != nil {
		for k, v := range datas {
			buf.PushByte(k)
			buf.PushObject(v)
		}
	}
	bytes, err := buf.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func RoomToHash(room *gen_room.RoomEntity) map[interface{}]interface{} {
	hash := make(map[interface{}]interface{})
	list := make([]interface{}, 0, 0)
	hash[gameparam.LobbyProperties] = list
	hash[gameparam.CleanupCacheOnLeave] = true
	hash[gameparam.MaxPlayers] = room.MaxPeers()
	hash[gameparam.IsVisible] = true
	hash[gameparam.IsOpen] = true
	hash[gameparam.MasterClientId] = room.MasterId()
	return hash
}
