package room

import (
	"errors"

	"github.com/itfantasy/gonode/utils/stl"
)

type ActorsManager struct {
	allActors *stl.List // List<Actor>
}

func NewActorsManager() *ActorsManager {
	actorsManager := new(ActorsManager)
	actorsManager.allActors = stl.NewList(10)
	return actorsManager
}

func (this *ActorsManager) AddNewActor(peerId string, actorNr int32) (*Actor, error) {
	if _, exist := this.GetActorByNr(actorNr); exist {
		return nil, errors.New("there has been an actor with the same actorNr!")
	}
	if _, exist := this.GetActorByPeerId(peerId); exist {
		return nil, errors.New("there has been an actor with the same peerId!")
	}

	actor := NewActor(peerId, actorNr)
	this.allActors.Add(actor)
	return actor, nil
}

func (this *ActorsManager) GetActorByNr(actorNr int32) (*Actor, bool) {
	for _, item := range this.allActors.Values() {
		actor := item.(*Actor)
		if actor.ActorNr == actorNr {
			return actor, true
		}
	}
	return nil, false
}

func (this *ActorsManager) GetActorByPeerId(peerId string) (*Actor, bool) {
	for _, item := range this.allActors.Values() {
		actor := item.(*Actor)
		if actor.PeerId == peerId {
			return actor, true
		}
	}
	return nil, false
}

func (this *ActorsManager) GetActorByIndex(index int) (*Actor, bool) {
	actor, err := this.allActors.Get(index)
	if err != nil {
		return nil, false
	}
	return actor.(*Actor), true
}

func (this *ActorsManager) RemoveActorByNr(actorNr int32) bool {
	if actor, exist := this.GetActorByNr(actorNr); exist {
		this.allActors.Remove(actor)
		return true
	}
	return false
}

func (this *ActorsManager) RemoveActorByPeer(peerId string) bool {
	if actor, exist := this.GetActorByPeerId(peerId); exist {
		this.allActors.Remove(actor)
		return true
	}
	return false
}

func (this *ActorsManager) GetAllActorNrs() []int32 {
	list := make([]int32, 0, this.allActors.Count())
	for _, item := range this.allActors.Values() {
		actor := item.(*Actor)
		list = append(list, actor.ActorNr)
	}
	return list
}

func (this *ActorsManager) GetAllPeerIds() []string {
	list := make([]string, 0, this.allActors.Count())
	for _, item := range this.allActors.Values() {
		actor := item.(*Actor)
		list = append(list, actor.PeerId)
	}
	return list
}

func (this *ActorsManager) ActorsCount() int {
	return this.allActors.Count()
}

func (this *ActorsManager) ClearAll() {
	this.allActors.Clear()
}
