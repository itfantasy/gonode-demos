package opcode

// ----------------- lobby&room -------------------

const (
	AuthenticateOnce byte = 231
	Authenticate     byte = 230
	JoinLobby        byte = 229
	LeaveLobby       byte = 228
	CreateGame       byte = 227
	JoinGame         byte = 226
	JoinRandomGame   byte = 225
	Leave            byte = 254
	RaiseEvent       byte = 253
	SetProperties    byte = 252
	GetProperties    byte = 251
	ChangeGroups     byte = 250
	FindFriends      byte = 222
	GetLobbyStats    byte = 221
	GetRegions       byte = 220
	WebRpc           byte = 219
	ServerSettings   byte = 218
	GetGameList      byte = 217
)
