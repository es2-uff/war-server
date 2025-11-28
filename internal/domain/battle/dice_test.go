package battle

import "testing"

func TestRollDice(t *testing.T) {
	tests := []struct {
		name     string
		numDice  int
		wantSize int
	}{
		{"Roll 1 die", 1, 1},
		{"Roll 2 dice", 2, 2},
		{"Roll 3 dice", 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dice := RollDice(tt.numDice)

			if len(dice) != tt.wantSize {
				t.Errorf("RollDice(%d) returned %d dice, want %d", tt.numDice, len(dice), tt.wantSize)
			}

			// Check all dice values are between 1 and 6
			for i, val := range dice {
				if val < 1 || val > 6 {
					t.Errorf("dice[%d] = %d, want value between 1 and 6", i, val)
				}
			}

			// Check dice are sorted in descending order
			for i := 0; i < len(dice)-1; i++ {
				if dice[i] < dice[i+1] {
					t.Errorf("dice not sorted descending: dice[%d]=%d < dice[%d]=%d", i, dice[i], i+1, dice[i+1])
				}
			}
		})
	}
}

func TestCompareDice_AttackerWins(t *testing.T) {
	attackerDice := []int{6, 5, 4}
	defenderDice := []int{3, 2, 1}

	attackerLosses, defenderLosses := CompareDice(attackerDice, defenderDice)

	if attackerLosses != 0 {
		t.Errorf("CompareDice() attackerLosses = %d, want 0", attackerLosses)
	}
	if defenderLosses != 3 {
		t.Errorf("CompareDice() defenderLosses = %d, want 3", defenderLosses)
	}
}

func TestCompareDice_DefenderWins(t *testing.T) {
	attackerDice := []int{3, 2, 1}
	defenderDice := []int{6, 5, 4}

	attackerLosses, defenderLosses := CompareDice(attackerDice, defenderDice)

	if attackerLosses != 3 {
		t.Errorf("CompareDice() attackerLosses = %d, want 3", attackerLosses)
	}
	if defenderLosses != 0 {
		t.Errorf("CompareDice() defenderLosses = %d, want 0", defenderLosses)
	}
}

func TestCompareDice_TiesGoToDefender(t *testing.T) {
	attackerDice := []int{5, 5, 5}
	defenderDice := []int{5, 5, 5}

	attackerLosses, defenderLosses := CompareDice(attackerDice, defenderDice)

	if attackerLosses != 3 {
		t.Errorf("CompareDice() with ties: attackerLosses = %d, want 3 (ties go to defender)", attackerLosses)
	}
	if defenderLosses != 0 {
		t.Errorf("CompareDice() with ties: defenderLosses = %d, want 0", defenderLosses)
	}
}

func TestCompareDice_Mixed(t *testing.T) {
	attackerDice := []int{6, 3}
	defenderDice := []int{5, 4}

	attackerLosses, defenderLosses := CompareDice(attackerDice, defenderDice)

	// 6 > 5: defender loses
	// 3 < 4: attacker loses
	if attackerLosses != 1 {
		t.Errorf("CompareDice() attackerLosses = %d, want 1", attackerLosses)
	}
	if defenderLosses != 1 {
		t.Errorf("CompareDice() defenderLosses = %d, want 1", defenderLosses)
	}
}

func TestCompareDice_UnequalDice(t *testing.T) {
	attackerDice := []int{6, 5, 4}
	defenderDice := []int{5}

	attackerLosses, defenderLosses := CompareDice(attackerDice, defenderDice)

	// Only compare first dice: 6 > 5
	if attackerLosses != 0 {
		t.Errorf("CompareDice() attackerLosses = %d, want 0", attackerLosses)
	}
	if defenderLosses != 1 {
		t.Errorf("CompareDice() defenderLosses = %d, want 1", defenderLosses)
	}
}
