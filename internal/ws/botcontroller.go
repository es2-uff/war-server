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

		borderTerritories := g.getBotBorderTerritories(botID)
		var territory *Territory
		if len(borderTerritories) > 0 {
			territory = borderTerritories[rand.Intn(len(borderTerritories))]
		} else {
			territory = ownedTerritories[rand.Intn(len(ownedTerritories))]
		}

		if err := g.GameState.Deploy(botID, territory.ID); err != nil {
			log.Printf("Bot deploy error: %v", err)
			break
		}

		g.sendBotAction("troop_assign", botID, map[string]any{
			"territory_id": territory.ID,
		})
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

		maxArmies := attack.from.Armies - 1
		if maxArmies < 1 {
			continue
		}
		attackingArmies := rand.Intn(min(maxArmies, 3)) + 1

		victory, err := g.GameState.Attack(botID, attack.from.ID, attack.to.ID, attackingArmies)
		if err != nil {
			log.Printf("Bot attack error: %v", err)
			continue
		}

		g.sendBotAction("attack", botID, map[string]any{
			"from":             attack.from.ID,
			"to":               attack.to.ID,
			"attacking_armies": attackingArmies,
		})

		time.Sleep(500 * time.Millisecond)

		// If conquered, update attack options
		if victory {
			attackOptions = g.getBotAttackOptions(botID)
		}
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
			to := g.getTerritoryByID(adjID)
			if to != nil && to.Owner == botID {
				movingArmies := (from.Armies - 1) / 2
				if movingArmies < 1 {
					continue
				}

				if err := g.GameState.Move(botID, from.ID, to.ID, movingArmies); err != nil {
					log.Printf("Bot move error: %v", err)
					continue
				}

				g.sendBotAction("troop_move", botID, map[string]any{
					"from":          from.ID,
					"to":            to.ID,
					"moving_armies": movingArmies,
				})

				time.Sleep(500 * time.Millisecond)
				return
			}
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

	// Send to broadcast channel (same as human messages)
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

func (g *Game) getBotBorderTerritories(botID string) []*Territory {
	g.GameState.RLock()
	defer g.GameState.RUnlock()

	borders := make([]*Territory, 0)
	for _, t := range g.GameState.Territories {
		if t.Owner != botID {
			continue
		}

		// Check if any adjacent territory is owned by someone else
		for _, adjID := range t.Adjacent {
			adjTerritory := g.getTerritoryByIDLocked(adjID)
			if adjTerritory != nil && adjTerritory.Owner != botID {
				borders = append(borders, t)
				break
			}
		}
	}
	return borders
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

		// Check adjacent territories for enemies
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

func (g *Game) getTerritoryByID(territoryID string) *Territory {
	g.GameState.RLock()
	defer g.GameState.RUnlock()
	return g.getTerritoryByIDLocked(territoryID)
}

func (g *Game) getTerritoryByIDLocked(territoryID string) *Territory {
	for _, t := range g.GameState.Territories {
		if t.ID == territoryID {
			return t
		}
	}
	return nil
}
