package protocol

type GameProvider interface {
	// Returns the unique identifier for this game type (e.g. blackjack, poker, etc.)
	GetID() string
	// Returns a displayable name for the game
	GetName() string
}

type GameInstance interface {
	// Returns the ID of the individual game instance (e.g. slot-001, slot-002, etc.)
	GetID() string
	// Returns the ID of the GameProvider this game belongs to
	GetProviderID() string
}
