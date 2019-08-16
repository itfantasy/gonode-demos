package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"

	"github.com/itfantasy/gonode-icloud/icloud/behaviors/lobby"
	"github.com/itfantasy/gonode-icloud/icloud/logics/master"
)

type MasterServer struct {
}

func (this *MasterServer) Setup() *gen_server.NodeInfo {
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

	redisConf := conf.Get("comps", "redis")
	if err := lobby.RegisterCoreRedis(redisConf); err != nil {
		fmt.Println(err)
		return nil
	}

	defaultRoomUrl := conf.Get("room", "defaulturl")
	master.SetDefaultRoomUrl(defaultRoomUrl)

	return nodeInfo
}

func (this *MasterServer) Start() {

}
func (this *MasterServer) OnConn(id string) {
	fmt.Println("new conn !! " + id)
}
func (this *MasterServer) OnMsg(id string, msg []byte) {
	if gonode.Label(id) == "room" {
		master.HandleServerMsg(id, msg)
	} else {
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
