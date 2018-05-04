package gameparam

const (
	MaxPlayers          byte = 255
	IsVisible                = 254
	IsOpen                   = 253
	PlayerCount              = 252 // used in gamestate reproted to master
	Removed                  = 251 // used in gamestate reproted to master
	LobbyProperties          = 250
	CleanupCacheOnLeave      = 249 // TODO: add reading of this property to GameParameterReader and converting from flash and json
	MasterClientId           = 248
	ExpectedUsers            = 247

	MinValue = 235
)
