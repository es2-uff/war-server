package battle

import "math/rand/v2"

// RollDice simulates rolling n dice and returns sorted results (highest first)
func RollDice(n int) []int {
	dice := make([]int, n)

	for i := range n {
		dice[i] = rand.IntN(6) + 1
	}

	// Sort descending
	for i := range len(dice) {
		for j := i + 1; j < len(dice); j++ {
			if dice[j] > dice[i] {
				dice[i], dice[j] = dice[j], dice[i]
			}
		}
	}
	return dice
}

// CompareDice compares attacker and defender dice, returns (attackerLosses, defenderLosses)
// Rule: Compare highest with highest, second highest with second highest
// Ties go to defender
func CompareDice(attackerDice, defenderDice []int) (int, int) {
	attackerLosses := 0
	defenderLosses := 0

	comparisons := min(len(attackerDice), len(defenderDice))

	for i := range comparisons {
		if attackerDice[i] > defenderDice[i] {
			defenderLosses++
		} else {
			// Ties go to defender
			attackerLosses++
		}
	}

	return attackerLosses, defenderLosses
}
