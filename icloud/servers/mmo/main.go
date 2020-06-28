package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/io"
	"github.com/itfantasy/gonode/utils/yaml"

	"github.com/itfantasy/gonode-icloud/icloud/logics/mmo"
	"github.com/itfantasy/gonode-toolkit/toolkit/gen_mmo"
)

type MmoServer struct {
	handler *mmo.MmoHandler
}

func (m *MmoServer) Setup() *gen_server.NodeInfo {
	conf, err := io.LoadFile(io.CurrentDir() + "conf.yaml")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	info := new(gen_mmo.MmoServerInfo)
	if err := yaml.Unmarshal(conf, info); err != nil {
		fmt.Println(err)
		return nil
	}
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
