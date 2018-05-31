package room

type RoomServer interface {
	Start()  // when start
	Update() // timer update
	OnMsg(string, []byte)
	OnReload(string) error
	OnClose(string)
}
