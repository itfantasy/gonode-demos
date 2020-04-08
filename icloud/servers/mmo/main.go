package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"

	"github.com/itfantasy/gonode-icloud/icloud/logics/mmo"
	"github.com/itfantasy/gonode-toolkit/toolkit/gen_mmo"
)

type MmoServer struct {
	handler *mmo.MmoHandler
}

func (m *MmoServer) Setup() *gen_server.NodeInfo {
	conf, err := ini.Load(io.CurrentDir() + "conf.ini")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	info := new(gen_mmo.MmoServerInfo)
	info.Id = conf.Get("node", "id")
	info.Url = conf.Get("node", "url")
	info.LogLevel = conf.Get("log", "loglevel")
	info.LogComp = conf.Get("log", "logcomp")
	info.RegComp = conf.Get("reg", "regcomp")

	m.handler = new(mmo.MmoHandler)
	return info.ExpandToNodeInfo()
}

func (m *MmoServer) Start() {

}
func (m *MmoServer) OnConn(id string) {
	if gonode.IsPeer(id) {
		m.handler.HandleConn(id)
	}
}
func (m *MmoServer) OnMsg(id string, msg []byte) {
	if gonode.IsPeer(id) {
		m.handler.HandleMsg(id, msg)
	}
}
func (m *MmoServer) OnClose(id string, reason error) {
	fmt.Println("conn closed !! " + id + " -- reason:" + reason.Error())
	if gonode.IsPeer(id) {
		m.handler.HandleClose(id)
	}
}

func main() {
	gonode.Bind(new(MmoServer))
	gonode.Launch()
}
