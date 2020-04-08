package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"

	"github.com/itfantasy/gonode-icloud/icloud/logics/master"
	"github.com/itfantasy/gonode-toolkit/toolkit"
	"github.com/itfantasy/gonode-toolkit/toolkit/gen_lobby"
)

type MasterServer struct {
}

func (m *MasterServer) Setup() *gen_server.NodeInfo {
	conf, err := ini.Load(io.CurrentDir() + "conf.ini")
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
	if err := gen_lobby.InitGameDB(conf.Get("gamedb", "comp")); err != nil {
		return nil
	}
	return info.ExpandToNodeInfo()
}

func (m *MasterServer) Start() {

}
func (m *MasterServer) OnConn(id string) {
	fmt.Println("new conn !! " + id)
	if gonode.IsPeer(id) {
		master.HandleConn(id)
	}
}
func (m *MasterServer) OnMsg(id string, msg []byte) {
	if gonode.IsPeer(id) {
		master.HandleMsg(id, msg)
	} else if gonode.Label(id) == toolkit.LABEL_ROOM {
		master.HandleServerMsg(id, msg)
	}
}
func (m *MasterServer) OnClose(id string, reason error) {
	fmt.Println("conn closed !! " + id + " -- reason:" + reason.Error())
	if gonode.IsPeer(id) {
		master.HandleClose(id)
	}
}

func main() {
	gonode.Bind(new(MasterServer))
	gonode.Launch()
}
