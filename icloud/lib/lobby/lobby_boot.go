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

func (this *LobbyBoot) SelfNodeInfo() (*gen_server.NodeInfo, error) {
	conf, err := ini.Load(io.CurDir() + "conf.ini")
	if err != nil {
		return nil, err
	}
	nodeInfo := new(gen_server.NodeInfo)
	nodeInfo.Tag = conf.Get("node", "tag")
	nodeInfo.Id = conf.Get("node", "id")
	nodeInfo.Url = conf.Get("node", "url")
	nodeInfo.RedUrl = conf.Get("redis", "url")
	nodeInfo.RedPool = conf.GetInt("redis", "pool", 0)
	nodeInfo.RedDB = conf.GetInt("redis", "db", 0)
	nodeInfo.RedAuth = conf.Get("redis", "auth")
	nodeInfo.AutoDetect = conf.GetInt("net", "autodetect", 0) > 0
	nodeInfo.Net = conf.Get("net", "net")
	return nodeInfo, nil
}
func (this *LobbyBoot) IsInterestedIn(string) bool {
	return false
}
func (this *LobbyBoot) Start() {
	fmt.Println("node starting...")
	this.server.Start()
}
func (this *LobbyBoot) Update() {

}
func (this *LobbyBoot) OnConn(id string) {
	fmt.Println("new conn !! " + id)
}
func (this *LobbyBoot) OnMsg(id string, msg []byte) {
	if strings.Contains(id, "room") {
		// native logic for roomserver
		// update the roomstate
		// delete the roomstate
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
func (this *LobbyBoot) OnReload(tag string) error {
	return nil
}
func (this *LobbyBoot) CreateConnId() string {
	return "cnt" + strconv.Itoa(rand.Intn(100000))
}
func (this *LobbyBoot) Initialize(server LobbyServer) {
	this.server = server
	gonode.Node().Initialize(this)
}
