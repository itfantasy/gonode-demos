package room

import (
	"fmt"
	"strings"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
)

type RoomBoot struct {
	server RoomServer
}

func (this *RoomBoot) Setup() *gen_server.NodeInfo {
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
func (this *RoomBoot) Start() {

}
func (this *RoomBoot) OnConn(id string) {
	this.server.OnConn(id)
}
func (this *RoomBoot) OnMsg(id string, msg []byte) {
	if strings.Contains(id, "lobby") {
		// native logic for lobbyserver
	} else {
		this.server.OnMsg(id, msg)
	}
}
func (this *RoomBoot) OnClose(id string) {
	this.server.OnClose(id)
}

func (this *RoomBoot) Initialize(server RoomServer) {
	this.server = server
	gonode.Node().Initialize(this)
}
