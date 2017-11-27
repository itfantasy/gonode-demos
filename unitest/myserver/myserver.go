package main

import (
	"fmt"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	//"github.com/itfantasy/gonode/gnbuffers"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
	//	"github.com/itfantasy/gonode/utils/timer"
)

type MyServer struct {
}

func (this *MyServer) SelfNodeInfo() (*gen_server.NodeInfo, error) {
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
func (this *MyServer) IsInterestedIn(string) bool {
	return true
}
func (this *MyServer) Start() {
	fmt.Println("node starting...")
}
func (this *MyServer) Update() {

}
func (this *MyServer) OnConn(id string) {
	fmt.Println("[" + id + "] has connected to this node!")
}
func (this *MyServer) OnMsg(id string, msg []byte) {
	/*
		buffer := gnbuffers.BuildBuffer(1024)
		buffer.PushInt(666)
		buffer.PushString("万能的刚哥!!")
		fmt.Println(buffer.Bytes())
	*/
	fmt.Println("[" + id + "]'s msg has been received!")
	fmt.Println(len(msg))
	fmt.Println(msg)
	/*
		parser := gnbuffers.BuildParser(msg, 0)
		if val, err := parser.Int(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}
		if val, err := parser.String(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}

		buffer := gnbuffers.BuildBuffer(4096)
		buffer.PushInt(999)
		buffer.PushString("犀利的王大侠!!!")
		gonode.Node().NetWorker().Send(id, buffer.Bytes())
	*/
}
func (this *MyServer) OnClose(id string) {

}
func (this *MyServer) OnShell(id string, msg string) {

}
func (this *MyServer) OnReload(tag string) error {
	return nil
}
func (this *MyServer) CreateConnId() string {
	return "cnt"
}
func main() {
	myserver := new(MyServer)
	gonode.Node().Initialize(myserver)
}
