// Package layouts defines the physical layout mappings for each keyboard type.
// Each keyboard defines how its linear binding sequence maps to the universal IR grid.
package layouts

import "masters3d.com/keyboard_layout_config_mapper/internal/models"

// PhysicalLayout defines how a keyboard's bindings map to the universal IR grid.
// The key insight: ZMK keymap files list bindings in a linear sequence, but
// each keyboard has a different physical arrangement. This struct defines
// how to interpret that linear sequence for each keyboard type.
type PhysicalLayout struct {
	Name           string
	KeyboardType   models.KeyboardType
	TotalKeys      int                    // Total keys per layer
	PositionMap    []models.IRPosition    // Index i → IRPosition for binding i
}

// GetIRPosition returns the IR position for a given binding index
func (pl *PhysicalLayout) GetIRPosition(index int) (models.IRPosition, bool) {
	if index < 0 || index >= len(pl.PositionMap) {
		return models.IRPosition{}, false
	}
	return pl.PositionMap[index], true
}

// GetLayoutForKeyboard returns the physical layout definition for a keyboard type
func GetLayoutForKeyboard(keyboardType models.KeyboardType) (*PhysicalLayout, error) {
	switch keyboardType {
	case models.KeyboardZMKAdv360:
		return GetAdv360Layout(), nil
	case models.KeyboardZMKGlove80:
		return GetGlove80Layout(), nil
	case models.KeyboardZMKAdvMod:
		return GetAdvModLayout(), nil
	default:
		return nil, nil
	}
}

// Universal IR Grid Layout (conceptual):
// 
// The universal grid is 10 rows × 10 cols per hand (left/right).
// We use semantic positions that are consistent across keyboards:
//
// Row 0: Function keys (F1-F12, etc.)
// Row 1: Number row (1-0, symbols)
// Row 2: Top alpha row (Q-P)
// Row 3: Home row (A-;)
// Row 4: Bottom alpha row (Z-/)
// Row 5: Modifier row (Ctrl, Alt, etc.)
// Row 6-8: Thumb cluster
// Row 9: Extra/reserved
//
// Columns 0-5 are main keys, 6-9 are edge/extra keys
//
// Example mapping for QWERTY Q key:
//   Hand: "left", Row: 2, Col: 0, Zone: "main", KeyID: "Q"
