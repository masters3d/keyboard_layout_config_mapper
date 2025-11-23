package models

import "fmt"

// IRPosition represents a position in the intermediate representation grid
// Using a universal coordinate system that can map to any keyboard layout
type IRPosition struct {
	Hand   string `json:"hand"`   // "left" or "right"
	Row    int    `json:"row"`    // 0-9 (10 rows per hand)
	Col    int    `json:"col"`    // 0-9 (10 columns per hand)  
	Zone   string `json:"zone"`   // "thumb", "main", "function", "numpad"
	KeyID  string `json:"key_id"` // Unique identifier within IR space
}

// IRKeyBinding represents a key binding in the intermediate representation
type IRKeyBinding struct {
	Position IRPosition `json:"position"`
	Value    string     `json:"value"`    // Normalized key code
	Layer    int        `json:"layer"`    // Layer index
	Type     BindingType `json:"type"`    // Type of binding
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// IRLayer represents a layer in the intermediate representation
type IRLayer struct {
	Index    int            `json:"index"`
	Name     string         `json:"name"`
	Bindings []IRKeyBinding `json:"bindings"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// IRLayout represents a complete layout in intermediate representation
// This serves as the universal format that all keyboards can translate to/from
type IRLayout struct {
	Name         string    `json:"name"`
	Layers       []IRLayer `json:"layers"`
	Source       KeyboardType `json:"source"`       // Original keyboard type
	Behaviors    []Behavior   `json:"behaviors"`    // Custom behaviors
	Combos       []IRCombo    `json:"combos"`       // Key combinations
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// IRCombo represents key combinations in IR space
type IRCombo struct {
	Name      string       `json:"name"`
	Keys      []IRPosition `json:"keys"`
	Binding   string       `json:"binding"`
	Layers    []int        `json:"layers,omitempty"`
	Timeout   int          `json:"timeout,omitempty"`
}

// IRCoordinate creates a standardized coordinate for the IR position
func (pos IRPosition) IRCoordinate() string {
	return fmt.Sprintf("%s_%d_%d", pos.Hand, pos.Row, pos.Col)
}

// IRZone represents different zones within the keyboard layout
type IRZone string

const (
	IRZoneMain     IRZone = "main"     // Main typing area (QWERTY region)
	IRZoneThumb    IRZone = "thumb"    // Thumb cluster keys
	IRZoneFunction IRZone = "function" // Function row (F1-F12, etc.)
	IRZoneNumpad   IRZone = "numpad"   // Number pad area
	IRZoneNav      IRZone = "nav"      // Navigation keys (arrows, home, end)
	IRZoneModifier IRZone = "modifier" // Dedicated modifier key area
)

// IRGrid represents the standardized 10x10 grid for each hand
// This provides a consistent mapping target for all keyboards
type IRGrid struct {
	LeftHand  [10][10]*IRKeyBinding  `json:"left_hand"`
	RightHand [10][10]*IRKeyBinding  `json:"right_hand"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// GetPosition returns the key binding at the specified IR position
func (grid *IRGrid) GetPosition(pos IRPosition) *IRKeyBinding {
	var hand *[10][10]*IRKeyBinding
	switch pos.Hand {
	case "left":
		hand = &grid.LeftHand
	case "right":
		hand = &grid.RightHand
	default:
		return nil
	}
	
	if pos.Row < 0 || pos.Row >= 10 || pos.Col < 0 || pos.Col >= 10 {
		return nil
	}
	
	return hand[pos.Row][pos.Col]
}

// SetPosition sets a key binding at the specified IR position
func (grid *IRGrid) SetPosition(pos IRPosition, binding *IRKeyBinding) error {
	var hand *[10][10]*IRKeyBinding
	switch pos.Hand {
	case "left":
		hand = &grid.LeftHand
	case "right":
		hand = &grid.RightHand
	default:
		return fmt.Errorf("invalid hand: %s", pos.Hand)
	}
	
	if pos.Row < 0 || pos.Row >= 10 || pos.Col < 0 || pos.Col >= 10 {
		return fmt.Errorf("position out of bounds: row=%d, col=%d", pos.Row, pos.Col)
	}
	
	hand[pos.Row][pos.Col] = binding
	return nil
}

// GetAllPositions returns all non-nil key bindings in the grid
func (grid *IRGrid) GetAllPositions() []IRKeyBinding {
	var bindings []IRKeyBinding
	
	// Iterate through left hand
	for row := 0; row < 10; row++ {
		for col := 0; col < 10; col++ {
			if binding := grid.LeftHand[row][col]; binding != nil {
				bindings = append(bindings, *binding)
			}
		}
	}
	
	// Iterate through right hand
	for row := 0; row < 10; row++ {
		for col := 0; col < 10; col++ {
			if binding := grid.RightHand[row][col]; binding != nil {
				bindings = append(bindings, *binding)
			}
		}
	}
	
	return bindings
}