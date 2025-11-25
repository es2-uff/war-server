package territory

import (
	"github.com/google/uuid"
)

type Territory struct {
	TerritoryID int
	RegionID    int
	OwnerID     uuid.UUID
	OwnerColor  string

	ArmyQuantity int
}
