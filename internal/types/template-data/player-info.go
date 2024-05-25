package templatedata

import "rcp/elite/internal/types"

type PlayerInfo struct {
	Player   *types.Player
	Score    Score
	TargetId string
}
