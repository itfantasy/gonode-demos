package main

import (
	"fmt"

	//	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/lobby"
)

type LobbySvcLogic struct {
}

func (this *LobbySvcLogic) Start() {

}
func (this *LobbySvcLogic) Update() {

}
func (this *LobbySvcLogic) OnMsg(id string, msg []byte) {
	// receive the msg from client
	fmt.Println(msg)
	HandleMsg(id, msg)
}
func (this *LobbySvcLogic) OnReload(string) error {
	return nil
}
func main() {
	logic := new(LobbySvcLogic)
	svc := new(lobby.Lobby)
	svc.Initialize(logic)
}
