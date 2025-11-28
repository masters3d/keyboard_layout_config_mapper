package layouts

import "masters3d.com/keyboard_layout_config_mapper/internal/models"

// GetGlove80Layout returns the physical layout definition for Glove80
// Total: 80 keys per layer
//
// Physical row structure (bindings order in keymap file):
// Row 0: 5 left + 5 right = 10 keys (function row)
// Row 1: 6 left + 6 right = 12 keys (number/symbol row)
// Row 2: 6 left + 6 right = 12 keys (QWERTY row)
// Row 3: 6 left + 6 right = 12 keys (home row)
// Row 4: 6 left + 6 thumb + 6 right = 18 keys (bottom + inner thumb)
// Row 5: 6 left + 6 thumb + 6 right = 16 keys (modifier + outer thumb)
func GetGlove80Layout() *PhysicalLayout {
	positions := make([]models.IRPosition, 80)
	idx := 0

	// ============ Row 0: Function row (10 keys) ============
	// Left side: 5 keys (F1-ish area)
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 0, Col: col, Zone: "function", KeyID: funcLeftIDs[col]}
		idx++
	}
	// Right side: 5 keys (F6-ish area)
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 0, Col: col, Zone: "function", KeyID: funcRightIDs[col]}
		idx++
	}

	// ============ Row 1: Number/Symbol row (12 keys) ============
	// Left side: 6 keys
	for col := 0; col < 6; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 1, Col: col, Zone: "main", KeyID: g80NumberLeftIDs[col]}
		idx++
	}
	// Right side: 6 keys
	for col := 0; col < 6; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 1, Col: col, Zone: "main", KeyID: g80NumberRightIDs[col]}
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
	// Right side: 6 keys (H, J, K, L, dot, ENTER)
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 3, Col: col, Zone: "main", KeyID: homeRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 3, Col: 5, Zone: "main", KeyID: "ENTER"}
	idx++

	// ============ Row 4: Bottom + inner thumb (18 keys) ============
	// Left side: 6 keys (ctrl-bspc, Z, X, C, V, B)
	positions[idx] = models.IRPosition{Hand: "left", Row: 4, Col: 0, Zone: "main", KeyID: "CTRL_BSPC"}
	idx++
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 4, Col: col + 1, Zone: "main", KeyID: bottomLeftIDs[col]}
		idx++
	}
	// Inner thumb cluster: 6 keys (3 left + 3 right inner)
	for col := 0; col < 3; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 6, Col: col, Zone: "thumb", KeyID: g80ThumbInnerLeftIDs[col]}
		idx++
	}
	for col := 0; col < 3; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 6, Col: col, Zone: "thumb", KeyID: g80ThumbInnerRightIDs[col]}
		idx++
	}
	// Right side: 6 keys (N, M, comma, dot, morph_comma, TAB)
	for col := 0; col < 5; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 4, Col: col, Zone: "main", KeyID: bottomRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 4, Col: 5, Zone: "main", KeyID: "TAB"}
	idx++

	// ============ Row 5: Modifier + outer thumb (16 keys) ============
	// Left modifiers: 4 keys + magic key = 5
	positions[idx] = models.IRPosition{Hand: "left", Row: 5, Col: 0, Zone: "modifier", KeyID: "MAGIC_L"}
	idx++
	for col := 0; col < 4; col++ {
		positions[idx] = models.IRPosition{Hand: "left", Row: 5, Col: col + 1, Zone: "modifier", KeyID: g80ModLeftIDs[col]}
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
	// Right modifiers: 4 keys + magic key = 5
	for col := 0; col < 4; col++ {
		positions[idx] = models.IRPosition{Hand: "right", Row: 5, Col: col, Zone: "modifier", KeyID: g80ModRightIDs[col]}
		idx++
	}
	positions[idx] = models.IRPosition{Hand: "right", Row: 5, Col: 4, Zone: "modifier", KeyID: "MAGIC_R"}
	idx++

	return &PhysicalLayout{
		Name:         "Glove80",
		KeyboardType: models.KeyboardZMKGlove80,
		TotalKeys:    80,
		PositionMap:  positions,
	}
}

// Glove80-specific key IDs
var (
	funcLeftIDs  = []string{"F1", "F2", "F3", "F4", "F5"}
	funcRightIDs = []string{"F6", "F7", "F8", "F9", "F10"}

	g80NumberLeftIDs  = []string{"CAD_L", "QUOTE_S", "QUOTE_D", "MINUS", "EQUAL", "SLASH"}
	g80NumberRightIDs = []string{"EXCL", "LBKT", "RBKT", "LPAREN", "RPAREN", "MINUS_R"}

	g80ThumbInnerLeftIDs  = []string{"KEYPAD_L", "WIN", "CTRL_L"}
	g80ThumbInnerRightIDs = []string{"TAB_L", "ESC", "KEYPAD_R"}

	g80ModLeftIDs  = []string{"WIN_MOD", "ALT_L", "LAYER5", "LAYER6"}
	g80ModRightIDs = []string{"LEFT", "DOWN", "UP", "RIGHT"}
)
