package protocol

type GameProvider interface {
	// Returns the unique identifier for this game type (e.g. blackjack, poker, etc.)
	GetID() string
	// Returns a displayable name for the game
	GetName() string
	// Get all running game instances
	GetInstances() []GameInstance
}

type GameInstance interface {
	// Returns the ID of the individual game instance (e.g. slot-001, slot-002, etc.)
	GetID() string
	// Returns the ID of the GameProvider this game belongs to
	GetProviderID() string

	//// Backend Interaction ////

	SetAdapter(adapter CasinoAdapter)
	GetAdapter() CasinoAdapter

	//// User Management ////

	UserJoin(userID string)
	UserLeave(userID string)
	GetUsers() []GameUserAssociation

	//// Client Packets ////

	HandleClientJoin(client GameClient)
	HandleClientLeave(clientID string)
	HandlePacket(packet GamePacket)

	//// Game Loop ////

	Tick()
}

type GameUserAssociation struct {
	UserID string `json:"user_id"`
	GameID string `json:"game_id"` // ID of the game instance
}
