package cacheop

const (
	/// <summary>
	///   Don't cache the event. (default)
	/// </summary>
	DoNotCache byte = 0

	/// <summary>
	///   Merge cached event with data.
	/// </summary>
	MergeCache byte = 1

	/// <summary>
	///   Replace cached event with data.
	/// </summary>
	ReplaceCache byte = 2

	/// <summary>
	///   Remove cached event.
	/// </summary>
	RemoveCache byte = 3

	/// <summary>
	/// Add to the room cache.
	/// </summary>
	AddToRoomCache byte = 4

	AddToRoomCacheGlobal byte = 5

	RemoveFromRoomCache byte = 6

	RemoveFromCacheForActorsLeft byte = 7

	/// <summary>
	///   Increase the index of the sliced cache. (default)
	/// </summary>
	SliceIncreaseIndex byte = 10

	/// <summary>
	///   Set the index of the sliced cache. (default)
	/// </summary>
	SliceSetIndex byte = 11

	/// <summary>
	///   Purge cache slice with index.
	/// </summary>
	SlicePurgeIndex byte = 12

	/// <summary>
	///   Purge cache slice up to index.
	/// </summary>
	SlicePurgeUpToIndex byte = 13
)
