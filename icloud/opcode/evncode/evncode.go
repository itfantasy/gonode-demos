package evncode

const (
	/// <summary>
	///   Specifies that no event code is set.
	/// </summary>
	NoCodeSet byte = 0

	/// <summary>
	///   The event code for the <see cref="JoinEvent"/>.
	/// </summary>
	Join = 255

	/// <summary>
	///   The event code for the <see cref="LeaveEvent"/>.
	/// </summary>
	Leave = 254

	/// <summary>
	///   The event code for the <see cref="PropertiesChangedEvent"/>.
	/// </summary>
	PropertiesChanged = 253

	/// <summary>
	/// The event code for the <see cref="DisconnectEvent"/>.
	/// </summary>
	Disconnect = 252

	/// <summary>
	/// The event code for the <see cref="ErrorInfoEvent"/>.
	/// </summary>
	ErrorInfo = 251

	CacheSliceChanged = 250

	EventCacheSlicePurged = 249
)
