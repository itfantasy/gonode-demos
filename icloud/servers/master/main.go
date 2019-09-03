package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"

	"github.com/itfantasy/gonode-icloud/icloud/logics/master"
	"github.com/itfantasy/gonode-toolkit/toolkit/gen_lobby"
)

type MasterServer struct {
}

func (this *MasterServer) Setup() *gen_server.NodeInfo {
	conf, err := ini.Load(io.CurDir() + "conf.ini")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	info := new(gen_lobby.LobbyServerInfo)

	info.Id = conf.Get("node", "id")
	info.Url = conf.Get("node", "url")

	info.LogLevel = conf.Get("log", "loglevel")
	info.LogComp = conf.Get("log", "logcomp")

	info.RegComp = conf.Get("reg", "regcomp")

	defaultRoomUrl := conf.Get("room", "defaulturl")
	master.SetDefaultRoomUrl(defaultRoomUrl)

	return info.ExpandToNodeInfo()
}

func (this *MasterServer) Start() {

}
func (this *MasterServer) OnConn(id string) {
	fmt.Println("new conn !! " + id)
}
func (this *MasterServer) OnMsg(id string, msg []byte) {
	if gonode.Label(id) == "room" {
		master.HandleServerMsg(id, msg)
	} else if gonode.IsCntId(id) {
		master.HandleMsg(id, msg)
	}
}
func (this *MasterServer) OnClose(id string, reason error) {
	fmt.Println("conn closed !! " + id)
}

func main() {
	gonode.Bind(new(MasterServer))
	gonode.Launch()
}
