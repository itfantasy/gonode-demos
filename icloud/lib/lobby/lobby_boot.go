package lobby

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
	//	"github.com/itfantasy/gonode/utils/timer"
)

type LobbyBoot struct {
	server LobbyServer
}

func (this *LobbyBoot) SelfInfo() (*gen_server.NodeInfo, error) {
	conf, err := ini.Load(io.CurDir() + "conf.ini")
	if err != nil {
		return nil, err
	}
	nodeInfo := new(gen_server.NodeInfo)

	nodeInfo.Id = conf.Get("node", "id")
	nodeInfo.Url = conf.Get("node", "url")
	nodeInfo.AutoDetect = conf.GetInt("node", "autodetect", 0) > 0
	nodeInfo.Public = conf.GetInt("node", "public", 0) > 0

	nodeInfo.RedUrl = conf.Get("redis", "url")
	nodeInfo.RedPool = conf.GetInt("redis", "pool", 0)
	nodeInfo.RedDB = conf.GetInt("redis", "db", 0)
	nodeInfo.RedAuth = conf.Get("redis", "auth")

	return nodeInfo, nil
}
func (this *LobbyBoot) OnDetect(string) bool {
	return false
}
func (this *LobbyBoot) Start() {
	fmt.Println("node starting...")
	this.server.Start()
}
func (this *LobbyBoot) OnConn(id string) {
	fmt.Println("new conn !! " + id)
}
func (this *LobbyBoot) OnMsg(id string, msg []byte) {
	if strings.Contains(id, "room") {
		this.server.OnServerMsg(id, msg)
	} else {
		this.server.OnMsg(id, msg)
	}
}
func (this *LobbyBoot) OnClose(id string) {
	fmt.Println("conn closed !! " + id)
}
func (this *LobbyBoot) OnShell(id string, msg string) {

}
func (this *LobbyBoot) OnRanId() string {
	return "cnt" + strconv.Itoa(rand.Intn(100000))
}
func (this *LobbyBoot) Initialize(server LobbyServer) {
	this.server = server
	gonode.Node().Initialize(this)
}
