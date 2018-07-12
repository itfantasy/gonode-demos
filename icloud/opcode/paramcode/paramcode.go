package paramcode

const (
	/// <summary>
	///   The game id.
	/// </summary>
	GameId byte = 255

	/// <summary>
	///   The actor nr
	///   used as op-key and ev-key
	/// </summary>
	ActorNr byte = 254

	/// <summary>
	///   The target actor nr.
	/// </summary>
	TargetActorNr byte = 253

	/// <summary>
	///   The actors.
	/// </summary>
	Actors byte = 252

	/// <summary>
	///   The properties.
	/// </summary>
	Properties byte = 251

	/// <summary>
	///   The broadcast.
	/// </summary>
	Broadcast byte = 250

	/// <summary>
	///   The actor properties.
	/// </summary>
	ActorProperties byte = 249

	/// <summary>
	///   The game properties.
	/// </summary>
	GameProperties byte = 248

	/// <summary>
	///   Event parameter to indicate whether events are cached for new actors.
	/// </summary>
	Cache byte = 247

	/// <summary>
	///   Event parameter containing a <see cref="Photon.Hive.Operations.ReceiverGroup"/> value.
	/// </summary>
	ReceiverGroup byte = 246

	/// <summary>
	///   The data.
	/// </summary>
	Data byte = 245

	/// <summary>
	///   The paramter code for the <see cref="RaiseEventRequest">raise event</see> operations event code.
	/// </summary>
	Code byte = 244

	/// <summary>
	///   the flush event code for raise event.
	/// </summary>
	Flush byte = 243

	/// <summary>
	/// Event parameter to indicate whether cached events are deleted automaticly for actors leaving a room.
	/// </summary>
	DeleteCacheOnLeave byte = 241

	/// <summary>
	/// The group this event should be sent to. No error is happening if the group is empty or not existing.
	/// </summary>
	Group byte = 240

	/// <summary>
	/// Groups to leave. Null won't remove any groups. byte[0] will remove ALL groups. Otherwise only the groups listed will be removed.
	/// </summary>
	GroupsForRemove byte = 239
	/// <summary>
	/// Should or not JoinGame response and JoinEvent contain user ids
	/// </summary>
	PublishUserId byte = 239

	/// <summary>
	/// Groups to enter. Null won't add groups. byte[0] will add ALL groups. Otherwise only the groups listed will be added.
	/// </summary>
	GroupsForAdd byte = 238

	AddUsers byte = 238

	/// <summary>
	/// A parameter indicating if common room events (Join Leave) will be suppressed.
	/// </summary>
	SuppressRoomEvents byte = 237

	/// <summary>
	/// A parameter indicating how long a room instance should be kept alive in the
	/// room cache after all players left the room.
	/// </summary>
	EmptyRoomLiveTime byte = 236

	/// <summary>
	/// A parameter indicating how long a player is allowd to rejoin.
	/// </summary>
	PlayerTTL byte = 235

	/// <summary>
	/// A parameter indicating that content of room event should be forwarded to some server
	/// </summary>
	HttpForward byte = 234
	WebFlags    byte = 234

	/// <summary>
	/// Allows the player to re-join the game.
	/// </summary>
	IsInactive byte = 233

	/// <summary>
	/// Activates UserId checks on joining - allowing a users to be only once in the room.
	/// Default is deactivated for backwards compatibility but we recomend to use in future
	/// Note: Should be active for saved games.
	/// </summary>
	CheckUserOnJoin byte = 232

	/// <summary>
	/// Expected values for actor and game properties
	/// </summary>
	ExpectedValues byte = 231

	// load balancing project specific parameters
	Address     byte = 230
	PeerCount   byte = 229
	ForceRejoin byte = 229

	GameCount                  byte = 228
	MasterPeerCount            byte = 227
	UserId                     byte = 225
	ApplicationId              byte = 224
	Position                   byte = 223
	GameList                   byte = 222
	Token                      byte = 221
	AppVersion                 byte = 220
	NodeId                     byte = 219
	Info                       byte = 218
	ClientAuthenticationType   byte = 217
	ClientAuthenticationParams byte = 216
	CreateIfNotExists          byte = 215
	JoinMode                   byte = 215
	ClientAuthenticationData   byte = 214
	LobbyName                  byte = 213
	LobbyType                  byte = 212
	LobbyStats                 byte = 211
	Region                     byte = 210
	UriPath                    byte = 209

	RpcCallParams     byte = 208
	RpcCallRetCode    byte = 207
	RpcCallRetMessage byte = 206

	CacheSliceIndex byte = 205

	Plugins byte = 204

	MasterClientId byte = 203

	Nickname byte = 202

	PluginName    byte = 201
	PluginVersion byte = 200
	Flags         byte = 199
)
