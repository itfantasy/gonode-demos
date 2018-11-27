package main

import (
	"fmt"

	//	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-icloud/icloud/behaviors/room"
)

type RoomServer struct {
}

func (this *RoomServer) Start() {

}
func (this *RoomServer) OnConn(id string) {
	HandleConn(id)
}
func (this *RoomServer) OnMsg(id string, msg []byte) {
	// receive the msg from client
	fmt.Println(msg)
	HandleMsg(id, msg)
}
func (this *RoomServer) OnClose(id string) {
	HandleClose(id)
}
func main() {
	server := new(RoomServer)
	boot := new(room.RoomBoot)
	boot.Initialize(server)
}
