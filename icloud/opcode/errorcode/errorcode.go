package errorcode

const (
	InvalidRequestParameters int16 = -6
	ArgumentOutOfRange       int16 = -4

	OperationDenied     int16 = -3
	OperationInvalid    int16 = -2
	InternalServerError int16 = -1

	Ok int16 = 0

	InvalidAuthentication int16 = 32767 // 0x7FFF // codes start at short.MaxValue
	GameIdAlreadyExists   int16 = 32766 // 0x7FFF - 1
	GameFull              int16 = 32765 // 0x7FFF - 2
	GameClosed            int16 = 32764 // 0x7FFF - 3
	AlreadyMatched        int16 = 32763 // 0x7FFF - 4
	ServerFull            int16 = 32762 // 0x7FFF - 5
	UserBlocked           int16 = 32761 // 0x7FFF - 6
	NoMatchFound          int16 = 32760 // 0x7FFF - 7
	RedirectRepeat        int16 = 32759 // 0x7FFF - 8
	GameIdNotExists       int16 = 32758 // 0x7FFF - 9

	// for authenticate requests. Indicates that the max ccu limit has been reached
	MaxCcuReached int16 = 32757 // 0x7FFF - 10

	// for authenticate requests. Indicates that the application is not subscribed to this region / private cloud.
	InvalidRegion int16 = 32756 // 0x7FFF - 11

	// for authenticate requests. Indicates that the call to the external authentication service failed.
	CustomAuthenticationFailed int16 = 32755 // 0x7FFF - 12

	AuthenticationTokenExpired int16 = 32753 // 0x7FFF - 14
	// for authenticate requests. Indicates that the call to the external authentication service failed.

	PluginReportedError int16 = 32752 //0x7FFF - 15
	PluginMismatch      int16 = 32751 // 0x7FFF - 16

	JoinFailedPeerAlreadyJoined    int16 = 32750 // 0x7FFF - 17
	JoinFailedFoundInactiveJoiner  int16 = 32749 // 0x7FFF - 18
	JoinFailedWithRejoinerNotFound int16 = 32748 // 0x7FFF - 19
	JoinFailedFoundExcludedUserId  int16 = 32747 // 0x7FFF - 20
	JoinFailedFoundActiveJoiner    int16 = 32746 // 0x7FFF - 21

	HttpLimitReached       int16 = 32745 // 0x7FFF - 22
	ExternalHttpCallFailed int16 = 32744 // 0x7FFF - 23
)
