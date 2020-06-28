package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/yaml"

	"github.com/itfantasy/gonode-icloud/icloud/logics/master"
	"github.com/itfantasy/gonode-toolkit/toolkit"
	"github.com/itfantasy/gonode-toolkit/toolkit/gen_lobby"
)

type MasterServer struct {
}

func (m *MasterServer) Setup() *gen_server.NodeInfo {
	conf, err := io.LoadFile(io.CurrentDir() + "conf.yaml")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	info := new(gen_lobby.LobbyServerInfo)
	if err := yaml.Unmarshal(conf, info); err != nil {
		fmt.Println(err)
		return nil
	}
	if err := gen_lobby.InitGameDB(info.GameDB); err != nil {
		fmt.Println(err)
		return nil
	}
	return info.ExpandToNodeInfo()
}

func (m *MasterServer) Start() {

}
func (m *MasterServer) OnConn(nodeid string) {
	fmt.Println("new conn !! " + nodeid)
	if gonode.IsPeer(nodeid) {
		master.HandleConn(nodeid)
	}
}
func (m *MasterServer) OnMsg(nodeid string, msg []byte) {
	if gonode.IsPeer(nodeid) {
		master.HandleMsg(nodeid, msg)
	} else if gonode.Label(nodeid) == toolkit.LABEL_ROOM {
		master.HandleServerMsg(nodeid, msg)
	}
}
func (m *MasterServer) OnClose(nodeid string, reason error) {
	fmt.Println("conn closed !! " + nodeid + " -- reason:" + reason.Error())
	if gonode.IsPeer(nodeid) {
		master.HandleClose(nodeid)
	}
}

func main() {
	gonode.Bind(new(MasterServer))
	gonode.Launch()
}
