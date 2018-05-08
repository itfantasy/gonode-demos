package cacheop

const (
	/// <summary>
	///   Don't cache the event. (default)
	/// </summary>
	DoNotCache byte = 0

	/// <summary>
	///   Merge cached event with data.
	/// </summary>
	MergeCache = 1

	/// <summary>
	///   Replace cached event with data.
	/// </summary>
	ReplaceCache = 2

	/// <summary>
	///   Remove cached event.
	/// </summary>
	RemoveCache = 3

	/// <summary>
	/// Add to the room cache.
	/// </summary>
	AddToRoomCache = 4

	AddToRoomCacheGlobal = 5

	RemoveFromRoomCache = 6

	RemoveFromCacheForActorsLeft = 7

	/// <summary>
	///   Increase the index of the sliced cache. (default)
	/// </summary>
	SliceIncreaseIndex = 10

	/// <summary>
	///   Set the index of the sliced cache. (default)
	/// </summary>
	SliceSetIndex = 11

	/// <summary>
	///   Purge cache slice with index.
	/// </summary>
	SlicePurgeIndex = 12

	/// <summary>
	///   Purge cache slice up to index.
	/// </summary>
	SlicePurgeUpToIndex = 13
)
