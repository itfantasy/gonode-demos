package room

type Actor struct {
	PeerId  string
	ActorNr int32
}

func NewActor(peerId string, actorNr int32) *Actor {
	actor := new(Actor)
	actor.PeerId = peerId
	actor.ActorNr = actorNr
	return actor
}
