package common

const (
	PlayerMoveChange int = 0
	NPCPMoveChange   int = 1
)

// GameState is a model that gets passed between server and clients.
type GameState struct {
	ChangeType int
	X          int
	Y          int
	Dir        int
}

// PlayerState is input that comes from a client.
// this is what they're wanting to do.
type PlayerState struct {

	// PlayerID simple (naive) way to identify players :)
	PlayerID   int
	ChangeType int
	X          int
	Y          int
	Dir        int
}
