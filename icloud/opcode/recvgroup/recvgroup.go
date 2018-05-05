package recvgroup

const (
	/// <summary>
	///   Send to all actors but the sender.
	/// </summary>
	Others byte = 0

	/// <summary>
	///   Send to all actors including the sender.
	/// </summary>
	All = 1

	/// <summary>
	///   Send to the peer with the lowest actor number.
	/// </summary>
	MasterClient = 2
)
