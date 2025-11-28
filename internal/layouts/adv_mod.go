package layouts

import "masters3d.com/keyboard_layout_config_mapper/internal/models"

// GetAdvModLayout returns the physical layout definition for Kinesis Advantage with Pillz Mod
// Total: 86 keys per layer
//
// Physical row structure (bindings order in keymap file):
// Row 0: 18 keys (function row - unique to adv_mod!)
// Row 1: 6 left + 6 right = 12 keys (number/symbol row)
// Row 2: 6 left + 6 right = 12 keys (QWERTY row)
// Row 3: 6 left + 6 right = 12 keys (home row)
// Row 4: 6 left + 6 right = 12 keys (bottom alpha row)
// Row 5: 4 left + 4 right = 8 keys (modifier row)
// Row 6: 2 left + 2 right = 4 keys (thumb top)
// Row 7: 1 left + 1 right = 2 keys (thumb middle)
// Row 8: 3 left + 3 right = 6 keys (thumb bottom)
func GetAdvModLayout() *PhysicalLayout {
	positions := make([]models.IRPosition, 86)
	idx := 0

	// ============ Row 0: Function row (18 keys) - UNIQUE TO ADV_MOD ============
	// This is the Kinesis Advantage's function key row that spans the entire top
	// HOME, F1-F12, PSCRN, SLCK, PAUSE, LAYER_TOGGLE, LAYER_SYSTEM
	advModFuncIDs := []string{
		"HOME_FUNC", "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8",
		"F9", "F10", "F11", "F12", "PSCRN", "SLCK", "PAUSE", "LAYER_TOG", "LAYER_SYS",
	}
	// First 9 on left, next 9 on right
	for col := 0; col < 9; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 0, Col: col, Zone: "function", KeyID: advModFuncIDs[col]}
		idx++
	}
	for col := 0; col < 9; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 0, Col: col, Zone: "function", KeyID: advModFuncIDs[col+9]}
		idx++
	}

	// ============ Row 1: Number/Symbol row (12 keys) ============
	// Left side: 6 keys (quote_single, quote_double, minus, equal, slash, ctrl-alt-del)
	advModNumLeftIDs := []string{"QUOTE_S", "QUOTE_D", "MINUS", "EQUAL", "SLASH", "CAD_L"}
	for col := 0; col < 6; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 1, Col: col, Zone: "main", KeyID: advModNumLeftIDs[col]}
		idx++
	}
	// Right side: 6 keys (excl, lbkt, rbkt, lparen, rparen, minus)
	advModNumRightIDs := []string{"EXCL", "LBKT", "RBKT", "LPAREN", "RPAREN", "MINUS_R"}
	for col := 0; col < 6; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 1, Col: col, Zone: "main", KeyID: advModNumRightIDs[col]}
		idx++
	}

	// ============ Row 2: QWERTY row (12 keys) ============
	// Left side: 6 keys (HOME, Q, W, E, R, T)
	positions[idx] = models.IRPosition{Hand: "left", Row: 2, Col: 0, Zone: "main", KeyID: "HOME_L"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 2, Col: col + 1, Zone: "main", KeyID: qwertyLeftIDs[col]}
		idx++
	}
	// Right side: 6 keys (Y, U, I, O, P, DEL)
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 2, Col: col, Zone: "main", KeyID: qwertyRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 2, Col: 5, Zone: "main", KeyID: "DEL"}
	idx++

	// ============ Row 3: Home row (12 keys) ============
	// Left side: 6 keys (BSPC, A, S, D, F, G)
	positions[idx] = models.IRPosition{Hand: "left", Row: 3, Col: 0, Zone: "main", KeyID: "BSPC"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 3, Col: col + 1, Zone: "main", KeyID: homeLeftIDs[col]}
		idx++
	}
	// Right side: 6 keys (H, J, K, L, dot_morph, ENTER)
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 3, Col: col, Zone: "main", KeyID: homeRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 3, Col: 5, Zone: "main", KeyID: "ENTER"}
	idx++

	// ============ Row 4: Bottom alpha row (12 keys) ============
	// Left side: 6 keys (ctrl-bspc, Z, X, C, V, B)
	positions[idx] = models.IRPosition{Hand: "left", Row: 4, Col: 0, Zone: "main", KeyID: "CTRL_BSPC"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 4, Col: col + 1, Zone: "main", KeyID: bottomLeftIDs[col]}
		idx++
	}
	// Right side: 6 keys (N, M, comma, dot, comma_morph, TAB)
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 4, Col: col, Zone: "main", KeyID: bottomRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 4, Col: 5, Zone: "main", KeyID: "TAB"}
	idx++

	// ============ Row 5: Modifier row (8 keys) ============
	// Left side: 4 keys (WIN, PGDN, ALT, KEYPAD)
	advModModLeftIDs := []string{"WIN_MOD", "PGDN", "ALT_L", "KEYPAD_L"}
	for col := 0; col < 4; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 5, Col: col, Zone: "modifier", KeyID: advModModLeftIDs[col]}
		idx++
	}
	// Right side: 4 keys (LEFT, DOWN, UP, RIGHT)
	advModModRightIDs := []string{"LEFT", "DOWN", "UP", "RIGHT"}
	for col := 0; col < 4; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 5, Col: col, Zone: "modifier", KeyID: advModModRightIDs[col]}
		idx++
	}

	// ============ Row 6: Thumb top (4 keys) ============
	// Left: 2 keys (CTRL, TAB)
	positions[idx] = models.IRPosition{Hand: "left", Row: 6, Col: 0, Zone: "thumb", KeyID: "CTRL_L"}
	idx++
	positions[idx] = models.IRPosition{Hand: "left", Row: 6, Col: 1, Zone: "thumb", KeyID: "TAB_L"}
	idx++
	// Right: 2 keys (CMD, SHIFT)
	positions[idx] = models.IRPosition{Hand: "right", Row: 6, Col: 0, Zone: "thumb", KeyID: "CMD"}
	idx++
	positions[idx] = models.IRPosition{Hand: "right", Row: 6, Col: 1, Zone: "thumb", KeyID: "SHIFT_R"}
	idx++

	// ============ Row 7: Thumb middle (2 keys) ============
	positions[idx] = models.IRPosition{Hand: "left", Row: 7, Col: 0, Zone: "thumb", KeyID: "HOME_THUMB"}
	idx++
	positions[idx] = models.IRPosition{Hand: "right", Row: 7, Col: 0, Zone: "thumb", KeyID: "PGUP"}
	idx++

	// ============ Row 8: Thumb bottom (6 keys) ============
	// Left: 3 keys (SPACE, SHIFT, ALT)
	positions[idx] = models.IRPosition{Hand: "left", Row: 8, Col: 0, Zone: "thumb", KeyID: "SPACE_L"}
	idx++
	positions[idx] = models.IRPosition{Hand: "left", Row: 8, Col: 1, Zone: "thumb", KeyID: "SHIFT_L"}
	idx++
	positions[idx] = models.IRPosition{Hand: "left", Row: 8, Col: 2, Zone: "thumb", KeyID: "ALT_THUMB"}
	idx++
	// Right: 3 keys (CMD, SHIFT, SPACE)
	positions[idx] = models.IRPosition{Hand: "right", Row: 8, Col: 0, Zone: "thumb", KeyID: "CMD_R"}
	idx++
	positions[idx] = models.IRPosition{Hand: "right", Row: 8, Col: 1, Zone: "thumb", KeyID: "SHIFT_R2"}
	idx++
	positions[idx] = models.IRPosition{Hand: "right", Row: 8, Col: 2, Zone: "thumb", KeyID: "SPACE_R"}
	idx++

	return &PhysicalLayout{
		Name:         "Kinesis Advantage (Pillz Mod)",
		KeyboardType: models.KeyboardZMKAdvMod,
		TotalKeys:    86,
		PositionMap:  positions,
	}
}
