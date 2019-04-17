package main

import (
	"fmt"
	"strings"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-icloud/icloud/logics/game"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
)

type RoomServer struct {
}

func (this *RoomServer) Setup() *gen_server.NodeInfo {
	conf, err := ini.Load(io.CurDir() + "conf.ini")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	nodeInfo := new(gen_server.NodeInfo)

	nodeInfo.Id = conf.Get("node", "id")
	nodeInfo.Url = conf.Get("node", "url")
	nodeInfo.Pub = conf.GetInt("node", "pub", 0) > 0
	nodeInfo.BackEnds = conf.Get("node", "backends")

	nodeInfo.LogLevel = conf.Get("log", "loglevel")
	nodeInfo.LogComp = conf.Get("log", "logcomp")

	nodeInfo.RegComp = conf.Get("reg", "regcomp")

	return nodeInfo
}
func (this *RoomServer) Start() {

}
func (this *RoomServer) OnConn(id string) {
	game.HandleConn(id)
}
func (this *RoomServer) OnMsg(id string, msg []byte) {
	if strings.Contains(id, "lobby") {
		// native logic for lobbyserver
	} else {
		game.HandleMsg(id, msg)
	}
}
func (this *RoomServer) OnClose(id string) {
	game.HandleClose(id)
}

func main() {
	server := new(RoomServer)
	gonode.Node().Initialize(server)
}
