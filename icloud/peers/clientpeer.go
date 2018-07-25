package peers

import (
	"github.com/itfantasy/gonode/utils/stl"
)

type Peer interface {
	PeerId() string
}

type PeerManager struct {
	dict *stl.Dictionary
}

func NewPeerManager() *PeerManager {
	this := new(PeerManager)
	this.dict = stl.NewDictionary()
	return this
}

func (this *PeerManager) AddPeer(peer Peer) error {
	return this.dict.Add(peer.PeerId(), peer)
}

func (this *PeerManager) RemovePeer(peer Peer) error {
	return this.dict.Remove(peer.PeerId())
}

func (this *PeerManager) GetPeer(peerId string) (Peer, bool) {
	ret, exist := this.dict.Get(peerId)
	if !exist {
		return nil, false
	}
	peer, ok := ret.(Peer)
	if !ok {
		return nil, false
	}
	return peer, true
}

func (this *PeerManager) GetClientPeer(peerId string) (*ClientPeer, bool) {
	peer, ok := this.GetPeer(peerId)
	if !ok {
		return nil, false
	}
	cntpeer, ok := peer.(*ClientPeer)
	if !ok {
		return nil, false
	}
	return cntpeer, true
}

type ClientPeer struct {
	peerId string
	roomId string // mark the gameId of this player
}

func NewClientPeer(peerId string) *ClientPeer {
	this := new(ClientPeer)
	this.peerId = peerId
	return this
}

func (this *ClientPeer) PeerId() string {
	return this.peerId
}

func (this *ClientPeer) RoomId() string {
	return this.roomId
}

func (this *ClientPeer) SetRoomId(roomId string) {
	this.roomId = roomId
}
