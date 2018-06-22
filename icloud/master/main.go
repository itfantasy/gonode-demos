package main

import (
	"fmt"

	//	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-demos/icloud/lib/lobby"
)

type MasterServer struct {
}

func (this *MasterServer) Start() {

}
func (this *MasterServer) Update() {

}
func (this *MasterServer) OnMsg(id string, msg []byte) {
	// receive the msg from client
	fmt.Println(msg)
	HandleMsg(id, msg)
}
func (this *MasterServer) OnReload(string) error {
	return nil
}
func main() {
	server := new(MasterServer)
	boot := new(lobby.LobbyBoot)
	boot.Initialize(server)
}
