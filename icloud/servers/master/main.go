package main

import (
	"fmt"

	//	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-icloud/icloud/behaviors/lobby"
	"github.com/itfantasy/gonode-icloud/icloud/logics/master"
)

type MasterServer struct {
}

func (this *MasterServer) Start() {

}
func (this *MasterServer) OnMsg(id string, msg []byte) {
	// receive the msg from client
	fmt.Println(msg)
	master.HandleMsg(id, msg)
}
func (this *MasterServer) OnServerMsg(id string, msg []byte) {
	master.HandleServerMsg(id, msg)
}
func main() {
	server := new(MasterServer)
	boot := new(lobby.LobbyBoot)
	boot.Initialize(server)
}
