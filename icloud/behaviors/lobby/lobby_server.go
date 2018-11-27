package lobby

type LobbyServer interface {
	Start() // when start
	OnMsg(string, []byte)
	OnServerMsg(string, []byte)
}
