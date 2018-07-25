package room

type RoomServer interface {
	Start()  // when start
	Update() // timer update
	OnConn(string)
	OnMsg(string, []byte)
	OnReload(string) error
	OnClose(string)
}
