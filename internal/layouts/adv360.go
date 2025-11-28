package layouts

import "masters3d.com/keyboard_layout_config_mapper/internal/models"

// GetAdv360Layout returns the physical layout definition for Advantage 360
// Total: 76 keys per layer
//
// Physical row structure (bindings order in keymap file):
// Row 0: 7 left + 7 right = 14 keys (number/symbol row)
// Row 1: 7 left + 7 right = 14 keys (QWERTY row with layer keys)
// Row 2: 7 left + 4 inner + 7 right = 18 keys (home row + inner modifiers)
// Row 3: 6 left + 2 center + 6 right = 14 keys (bottom alpha row)
// Row 4: 5 left + 3 outer + 3 outer + 5 right = 16 keys (modifier + outer thumb)
func GetAdv360Layout() *PhysicalLayout {
	positions := make([]models.IRPosition, 76)
	idx := 0

	// ============ Row 0: Number/Symbol row (14 keys) ============
	// Left side: 7 keys (positions 0-6)
	for col := 0; col < 7; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 1, Col: col, Zone: "main", KeyID: leftNumberRowIDs[col]}
		idx++
	}
	// Right side: 7 keys (positions 7-13)
	for col := 0; col < 7; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 1, Col: col, Zone: "main", KeyID: rightNumberRowIDs[col]}
		idx++
	}

	// ============ Row 1: QWERTY row (14 keys) ============
	// Left side: 7 keys (HOME, Q, W, E, R, T, layer)
	positions[idx] = models.IRPosition{Hand: "left", Row: 2, Col: 0, Zone: "main", KeyID: "HOME_L"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 2, Col: col + 1, Zone: "main", KeyID: qwertyLeftIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "left", Row: 2, Col: 6, Zone: "main", KeyID: "LAYER_L1"}
	idx++
	// Right side: 7 keys (layer, Y, U, I, O, P, DEL)
	positions[idx] = models.IRPosition{Hand: "right", Row: 2, Col: 0, Zone: "main", KeyID: "LAYER_R1"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 2, Col: col + 1, Zone: "main", KeyID: qwertyRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 2, Col: 6, Zone: "main", KeyID: "DEL"}
	idx++

	// ============ Row 2: Home row + inner modifiers (18 keys) ============
	// Left main: 7 keys (BSPC, A, S, D, F, G, layer)
	positions[idx] = models.IRPosition{Hand: "left", Row: 3, Col: 0, Zone: "main", KeyID: "BSPC"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 3, Col: col + 1, Zone: "main", KeyID: homeLeftIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "left", Row: 3, Col: 6, Zone: "main", KeyID: "LAYER_L2"}
	idx++
	// Inner modifiers: 4 keys (KEYPAD, WIN, CTRL, TAB on right side inner)
	for col := 0; col < 4; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 6, Col: col, Zone: "thumb", KeyID: thumbInnerLeftIDs[col]}
		idx++
	}
	// Right main: 7 keys (layer, H, J, K, L, dot, ENTER)
	positions[idx] = models.IRPosition{Hand: "right", Row: 3, Col: 0, Zone: "main", KeyID: "LAYER_R2"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 3, Col: col + 1, Zone: "main", KeyID: homeRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 3, Col: 6, Zone: "main", KeyID: "ENTER"}
	idx++

	// ============ Row 3: Bottom alpha row (14 keys) ============
	// Left side: 6 keys (ctrl-bspc, Z, X, C, V, B)
	positions[idx] = models.IRPosition{Hand: "left", Row: 4, Col: 0, Zone: "main", KeyID: "CTRL_BSPC"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 4, Col: col + 1, Zone: "main", KeyID: bottomLeftIDs[col]}
		idx++
	}
	// Center: 2 keys (inner modifiers)
	positions[idx] = models.IRPosition{Hand: "left", Row: 5, Col: 5, Zone: "modifier", KeyID: "LCTRL"}
	idx++
	positions[idx] = models.IRPosition{Hand: "right", Row: 5, Col: 5, Zone: "modifier", KeyID: "TAB_R"}
	idx++
	// Right side: 6 keys (N, M, comma, dot, morph_comma, TAB)
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 4, Col: col + 1, Zone: "main", KeyID: bottomRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 4, Col: 6, Zone: "main", KeyID: "TAB"}
	idx++

	// ============ Row 4: Modifier + outer thumb (16 keys) ============
	// Left modifiers: 5 keys
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 5, Col: col, Zone: "modifier", KeyID: modLeftIDs[col]}
		idx++
	}
	// Left outer thumb: 3 keys
	for col := 0; col < 3; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 7, Col: col, Zone: "thumb", KeyID: thumbOuterLeftIDs[col]}
		idx++
	}
	// Right outer thumb: 3 keys
	for col := 0; col < 3; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 7, Col: col, Zone: "thumb", KeyID: thumbOuterRightIDs[col]}
		idx++
	}
	// Right modifiers: 5 keys
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 5, Col: col, Zone: "modifier", KeyID: modRightIDs[col]}
		idx++
	}

	return &PhysicalLayout{
		Name:         "Advantage 360",
		KeyboardType: models.KeyboardZMKAdv360,
		TotalKeys:    76,
		PositionMap:  positions,
	}
}

// Key ID constants for consistent naming across keyboards
var (
	leftNumberRowIDs  = []string{"SPECIAL", "QUOTE_S", "QUOTE_D", "MINUS", "EQUAL", "SLASH", "CAD_L"}
	rightNumberRowIDs = []string{"CAD_R", "EXCL", "LBKT", "RBKT", "LPAREN", "RPAREN", "MINUS_R"}

	qwertyLeftIDs  = []string{"Q", "W", "E", "R", "T"}
	qwertyRightIDs = []string{"Y", "U", "I", "O", "P"}

	homeLeftIDs  = []string{"A", "S", "D", "F", "G"}
	homeRightIDs = []string{"H", "J", "K", "L", "DOT_MORPH"}

	bottomLeftIDs  = []string{"Z", "X", "C", "V", "B"}
	bottomRightIDs = []string{"N", "M", "COMMA", "DOT", "COMMA_MORPH"}

	thumbInnerLeftIDs  = []string{"KEYPAD_L", "WIN", "CTRL_L", "TAB_L"}
	thumbInnerRightIDs = []string{"ESC", "KEYPAD_R", "CMD", "SHIFT_R"}

	thumbOuterLeftIDs  = []string{"SPACE_L", "SHIFT_L", "ALT_L"}
	thumbOuterRightIDs = []string{"CMD_R", "SHIFT_R2", "SPACE_R"}

	modLeftIDs  = []string{"MOD_L", "HOME_MOD", "PGDN", "LAYER5", "LAYER6"}
	modRightIDs = []string{"LEFT", "DOWN", "UP", "RIGHT", "MOD_R"}
)
