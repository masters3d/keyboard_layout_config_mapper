package models

import "time"

// KeyboardType represents the different keyboard firmware systems
type KeyboardType string

const (
	KeyboardZMKAdv360   KeyboardType = "adv360"
	KeyboardZMKGlove80  KeyboardType = "glove80" 
	KeyboardZMKAdvMod   KeyboardType = "adv_mod"
	KeyboardQMKErgoDox  KeyboardType = "qmk_ergodox"
	KeyboardKinesis2    KeyboardType = "kinesis2"
)

// Position represents a physical key position on a keyboard
type Position struct {
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Side   string `json:"side"` // "left" or "right"
	Zone   string `json:"zone"` // "thumb", "main", "function"
	KeyID  string `json:"key_id"`
}

// KeyBinding represents a single key binding/mapping
type KeyBinding struct {
	Position Position    `json:"position"`
	Value    string      `json:"value"`     // The actual key code or binding
	Layer    int         `json:"layer"`     // Which layer this binding belongs to
	Type     BindingType `json:"type"`      // Type of binding
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// BindingType represents different types of key bindings
type BindingType string

const (
	BindingBasic     BindingType = "basic"      // Simple key press
	BindingModTap    BindingType = "mod_tap"    // Modifier when held, key when tapped  
	BindingLayerTap  BindingType = "layer_tap"  // Layer when held, key when tapped
	BindingMacro     BindingType = "macro"      // Macro or sequence
	BindingCombo     BindingType = "combo"      // Key combination
	BindingBehavior  BindingType = "behavior"   // Custom behavior
)

// Layer represents a complete keyboard layer
type Layer struct {
	Index    int          `json:"index"`
	Name     string       `json:"name"`
	Bindings []KeyBinding `json:"bindings"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// KeyboardLayout represents a complete keyboard layout
type KeyboardLayout struct {
	Type         KeyboardType `json:"type"`
	Name         string       `json:"name"`
	FilePath     string       `json:"file_path"`
	Layers       []Layer      `json:"layers"`
	Behaviors    []Behavior   `json:"behaviors"`
	Combos       []Combo      `json:"combos"`
	LastModified time.Time    `json:"last_modified"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// Behavior represents custom behaviors (ZMK specific)
type Behavior struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// Combo represents key combinations
type Combo struct {
	Name      string     `json:"name"`
	Keys      []Position `json:"keys"`
	Binding   string     `json:"binding"`
	Layers    []int      `json:"layers,omitempty"`
	Timeout   int        `json:"timeout,omitempty"`
}

// ChangeSet represents a set of changes between layouts
type ChangeSet struct {
	Source      KeyboardType `json:"source"`
	Target      KeyboardType `json:"target"`
	Changes     []Change     `json:"changes"`
	Conflicts   []Conflict   `json:"conflicts"`
	GeneratedAt time.Time    `json:"generated_at"`
}

// Change represents a single change to be applied
type Change struct {
	Type        ChangeType  `json:"type"`
	Position    Position    `json:"position"`
	Layer       int         `json:"layer"`
	OldValue    string      `json:"old_value,omitempty"`
	NewValue    string      `json:"new_value"`
	Description string      `json:"description"`
	Confidence  float64     `json:"confidence"` // 0.0 to 1.0
}

// ChangeType represents the type of change
type ChangeType string

const (
	ChangeAdd    ChangeType = "add"
	ChangeUpdate ChangeType = "update"
	ChangeDelete ChangeType = "delete"
	ChangeMove   ChangeType = "move"
)

// Conflict represents a mapping conflict that needs resolution
type Conflict struct {
	Position    Position `json:"position"`
	Layer       int      `json:"layer"`
	SourceValue string   `json:"source_value"`
	Reason      string   `json:"reason"`
	Suggestions []string `json:"suggestions"`
}