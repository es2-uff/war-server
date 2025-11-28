package ws

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

func (g *Game) executeBotTurn(botID string) {
	time.Sleep(1 * time.Second)

	g.botDeployPhase(botID)
	time.Sleep(500 * time.Millisecond)

	g.botAttackPhase(botID)
	time.Sleep(500 * time.Millisecond)

	g.botMovePhase(botID)
	time.Sleep(500 * time.Millisecond)

	g.botFinishTurn(botID)
}

func (g *Game) botDeployPhase(botID string) {
	for {
		g.GameState.RLock()
		bot := g.GameState.Players[botID]
		if bot == nil || bot.Armies <= 0 {
			g.GameState.RUnlock()
			break
		}
		g.GameState.RUnlock()

		ownedTerritories := g.getBotOwnedTerritories(botID)
		if len(ownedTerritories) == 0 {
			break
		}

		territory := ownedTerritories[rand.Intn(len(ownedTerritories))]

		g.sendBotAction("troop_assign", botID, map[string]any{
			"territory_id": territory.ID,
		})

		time.Sleep(100 * time.Millisecond)
	}
}

func (g *Game) botAttackPhase(botID string) {
	attackOptions := g.getBotAttackOptions(botID)
	if len(attackOptions) == 0 {
		return
	}

	numAttacks := rand.Intn(3) + 1
	for i := 0; i < numAttacks && i < len(attackOptions); i++ {
		attackIdx := rand.Intn(len(attackOptions))
		attack := attackOptions[attackIdx]

		g.GameState.RLock()
		fromTerritory := g.getTerritoryByIDLocked(attack.from.ID)
		toTerritory := g.getTerritoryByIDLocked(attack.to.ID)

		if fromTerritory == nil || fromTerritory.Owner != botID {
			g.GameState.RUnlock()
			continue
		}

		if toTerritory == nil || toTerritory.Owner == botID {
			g.GameState.RUnlock()
			continue
		}

		maxArmies := fromTerritory.Armies - 1
		g.GameState.RUnlock()

		if maxArmies < 1 {
			continue
		}
		attackingArmies := rand.Intn(min(maxArmies, 3)) + 1

		g.sendBotAction("attack", botID, map[string]any{
			"from":             attack.from.ID,
			"to":               attack.to.ID,
			"attacking_armies": attackingArmies,
		})

		time.Sleep(500 * time.Millisecond)

		attackOptions = g.getBotAttackOptions(botID)
	}
}

func (g *Game) botMovePhase(botID string) {
	ownedTerritories := g.getBotOwnedTerritories(botID)
	if len(ownedTerritories) < 2 {
		return
	}

	for _, from := range ownedTerritories {
		if from.Armies <= 1 {
			continue
		}

		for _, adjID := range from.Adjacent {
			g.GameState.RLock()
			fromFresh := g.getTerritoryByIDLocked(from.ID)
			toFresh := g.getTerritoryByIDLocked(adjID)

			if fromFresh == nil || fromFresh.Owner != botID || fromFresh.Armies <= 1 {
				g.GameState.RUnlock()
				continue
			}
			if toFresh == nil || toFresh.Owner != botID {
				g.GameState.RUnlock()
				continue
			}

			movingArmies := (fromFresh.Armies - 1) / 2
			g.GameState.RUnlock()

			if movingArmies < 1 {
				continue
			}

			g.sendBotAction("troop_move", botID, map[string]any{
				"from":          from.ID,
				"to":            adjID,
				"moving_armies": movingArmies,
			})

			time.Sleep(500 * time.Millisecond)
			return
		}
	}
}

func (g *Game) botFinishTurn(botID string) {
	g.sendBotAction("finish_turn", botID, nil)
}

func (g *Game) sendBotAction(actionType string, botID string, params map[string]any) {
	msg := map[string]any{
		"type":      actionType,
		"player_id": botID,
	}

	for k, v := range params {
		msg[k] = v
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling bot action: %v", err)
		return
	}

	g.broadcast <- jsonMsg
}

func (g *Game) getBotOwnedTerritories(botID string) []*Territory {
	g.GameState.RLock()
	defer g.GameState.RUnlock()

	owned := make([]*Territory, 0)
	for _, t := range g.GameState.Territories {
		if t.Owner == botID {
			owned = append(owned, t)
		}
	}
	return owned
}

type attackOption struct {
	from *Territory
	to   *Territory
}

func (g *Game) getBotAttackOptions(botID string) []attackOption {
	g.GameState.RLock()
	defer g.GameState.RUnlock()

	options := make([]attackOption, 0)
	for _, t := range g.GameState.Territories {
		if t.Owner != botID || t.Armies <= 1 {
			continue
		}

		for _, adjID := range t.Adjacent {
			adjTerritory := g.getTerritoryByIDLocked(adjID)
			if adjTerritory != nil && adjTerritory.Owner != botID {
				options = append(options, attackOption{
					from: t,
					to:   adjTerritory,
				})
			}
		}
	}
	return options
}

func (g *Game) getTerritoryByIDLocked(territoryID string) *Territory {
	for _, t := range g.GameState.Territories {
		if t.ID == territoryID {
			return t
		}
	}
	return nil
}
