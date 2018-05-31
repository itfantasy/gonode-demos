package main

import (
	//	"fmt"

	//	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode_demo/icloud/lib/room"
)

type RoomServer struct {
}

func (this *RoomServer) Start() {

}
func (this *RoomServer) Update() {

}
func (this *RoomServer) OnMsg(id string, msg []byte) {
	// receive the msg from client
	//fmt.Println(msg)
	HandleMsg(id, msg)
}
func (this *RoomServer) OnReload(id string) error {
	return nil
}
func (this *RoomServer) OnClose(id string) {
	HandleClose(id)
}
func main() {
	server := new(RoomServer)
	boot := new(room.RoomBoot)
	boot.Initialize(server)
}
