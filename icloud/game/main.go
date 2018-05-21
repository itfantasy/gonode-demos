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
func (this *RoomServer) OnReload(string) error {
	return nil
}
func main() {
	server := new(RoomServer)
	boot := new(room.RoomBoot)
	boot.Initialize(server)
}
