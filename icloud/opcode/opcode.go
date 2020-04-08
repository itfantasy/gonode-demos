package opcode

// ----------------- lobby&room -------------------

const (
	AuthenticateOnce byte = 231
	Authenticate          = 230
	JoinLobby             = 229
	LeaveLobby            = 228
	CreateGame            = 227
	JoinGame              = 226
	JoinRandomGame        = 225
	Leave                 = 254
	RaiseEvent            = 253
	SetProperties         = 252
	GetProperties         = 251
	ChangeGroups          = 250
	FindFriends           = 222
	GetLobbyStats         = 221
	GetRegions            = 220
	WebRpc                = 219
	ServerSettings        = 218
	GetGameList           = 217

	// ----- MMO
	Nil                = 0
	CreateWorld        = 90
	EnterWorld         = 91
	ExitWorld          = 92
	Move               = 93
	RaiseGenericEvent  = 94
	SetItemProperties  = 95 // mmo
	SpawnItem          = 96
	DestroyItem        = 97
	SubscribeItem      = 98 //// Manually subscribes item (does not affect interest area updates).
	UnsubscribeItem    = 99
	SetViewDistance    = 100
	AttachInterestArea = 101
	DetachInterestArea = 102
	AddInterestArea    = 103
	RemoveInterestArea = 104
	GetItemProperties  = 105 // mmo
	MoveInterestArea   = 106
	RadarSubscribe     = 107
	UnsubscribeCounter = 108
	SubscribeCounter   = 109
)
