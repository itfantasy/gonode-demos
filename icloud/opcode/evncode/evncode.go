package evncode

const (
	/// <summary>
	///   Specifies that no event code is set.
	/// </summary>
	NoCodeSet byte = 0

	/// <summary>
	///   The event code for the <see cref="JoinEvent"/>.
	/// </summary>
	Join byte = 255

	/// <summary>
	///   The event code for the <see cref="LeaveEvent"/>.
	/// </summary>
	Leave byte = 254

	/// <summary>
	///   The event code for the <see cref="PropertiesChangedEvent"/>.
	/// </summary>
	PropertiesChanged byte = 253

	/// <summary>
	/// The event code for the <see cref="DisconnectEvent"/>.
	/// </summary>
	Disconnect byte = 252

	/// <summary>
	/// The event code for the <see cref="ErrorInfoEvent"/>.
	/// </summary>
	ErrorInfo byte = 251

	CacheSliceChanged byte = 250

	EventCacheSlicePurged byte = 249

	GameList          byte = 230
	GameListUpdate    byte = 229
	QueueState        byte = 228
	AppStats          byte = 226
	GameServerOffline byte = 225
	LobbyStats        byte = 224
)
