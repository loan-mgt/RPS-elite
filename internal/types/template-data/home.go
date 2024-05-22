package templatedata

import "rcp/elite/internal/types"

type Home struct {
	Player   *types.Player
	Opponent *types.Player
	Messenger
}
