package room

import (
	"testing"

	"es2.uff/war-server/internal/domain/territory"
	"github.com/google/uuid"
)

func TestInstantiateGameTerritories(t *testing.T) {
	// Setup test players
	player1 := uuid.New()
	player2 := uuid.New()
	player3 := uuid.New()
	roomPlayers := []uuid.UUID{player1, player2, player3}

	// Call the function
	territories := InstantiateGameTerritories(roomPlayers)

	// Test 1: Verify all territories are created
	expectedTerritoryCount := len(territory.AllTerritories)
	if len(territories) != expectedTerritoryCount {
		t.Errorf("Expected %d territories, got %d", expectedTerritoryCount, len(territories))
	}

	// Test 2: Verify all territories have owners
	for _, terr := range territories {
		if terr.OwnerID == uuid.Nil {
			t.Errorf("Territory %d has no owner assigned", terr.TerritoryID)
		}
	}

	// Test 4: Verify all territories are assigned to valid players
	validPlayers := make(map[uuid.UUID]bool)
	for _, p := range roomPlayers {
		validPlayers[p] = true
	}

	for _, terr := range territories {
		if !validPlayers[terr.OwnerID] {
			t.Errorf("Territory %d assigned to invalid player %s", terr.TerritoryID, terr.OwnerID)
		}
	}

	// Test 5: Verify territories are distributed among players
	playerTerritoryCount := make(map[uuid.UUID]int)
	for _, terr := range territories {
		playerTerritoryCount[terr.OwnerID]++
	}

	// Check that all players have at least one territory
	for _, player := range roomPlayers {
		count := playerTerritoryCount[player]
		if count == 0 {
			t.Errorf("Player %s has no territories assigned", player)
		}
	}

	// Test 7: Verify each territory has a valid region
	for _, terr := range territories {
		expectedRegion := territory.TerritoryRegionMap[territory.TerritoryID(terr.TerritoryID)]
		if terr.RegionID != int(expectedRegion) {
			t.Errorf("Territory %d has region %d, expected %d", terr.TerritoryID, terr.RegionID, expectedRegion)
		}
	}
}

func TestInstantiateGameTerritoriesWithDifferentPlayerCounts(t *testing.T) {
	testCases := []struct {
		name        string
		playerCount int
	}{
		{"2 players", 2},
		{"3 players", 3},
		{"4 players", 4},
		{"5 players", 5},
		{"6 players", 6},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate players
			players := make([]uuid.UUID, tc.playerCount)
			for i := 0; i < tc.playerCount; i++ {
				players[i] = uuid.New()
			}

			// Call the function
			territories := InstantiateGameTerritories(players)

			// Verify all territories are assigned
			totalTerritories := len(territory.AllTerritories)
			if len(territories) != totalTerritories {
				t.Errorf("Expected %d territories, got %d", totalTerritories, len(territories))
			}

			// Verify all territories have owners
			unassignedCount := 0
			for _, terr := range territories {
				if terr.OwnerID == uuid.Nil {
					unassignedCount++
				}
			}

			if unassignedCount > 0 {
				t.Errorf("Found %d unassigned territories", unassignedCount)
			}

			// Verify balanced distribution
			playerTerritoryCount := make(map[uuid.UUID]int)
			for _, terr := range territories {
				playerTerritoryCount[terr.OwnerID]++
			}

			expectedPerPlayer := totalTerritories / tc.playerCount
			for playerID, count := range playerTerritoryCount {
				if count < expectedPerPlayer || count > expectedPerPlayer+1 {
					t.Errorf("Player %s has %d territories, expected %d or %d", playerID, count, expectedPerPlayer, expectedPerPlayer+1)
				}
			}
		})
	}
}
