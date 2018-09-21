package room

type RoomServer interface {
	Start() // when start
	OnConn(string)
	OnMsg(string, []byte)
	OnClose(string)
}
