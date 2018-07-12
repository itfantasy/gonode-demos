package gameparam

const (
	MaxPlayers          byte = 255
	IsVisible           byte = 254
	IsOpen              byte = 253
	PlayerCount         byte = 252 // used in gamestate reproted to master
	Removed             byte = 251 // used in gamestate reproted to master
	LobbyProperties     byte = 250
	CleanupCacheOnLeave byte = 249 // TODO: add reading of this property to GameParameterReader and converting from flash and json
	MasterClientId      byte = 248
	ExpectedUsers       byte = 247

	MinValue byte = 235
)
