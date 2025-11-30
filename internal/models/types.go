package models

import "time"

// KeyboardType represents the different keyboard firmware systems
type KeyboardType string

const (
	KeyboardZMKAdv360   KeyboardType = "adv360"
	KeyboardZMKGlove80  KeyboardType = "glove80" 
	KeyboardZMKAdvMod   KeyboardType = "adv_mod"
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
	Value    string      `json:"value"`
	Layer    int         `json:"layer"`
	Type     BindingType `json:"type"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// BindingType represents different types of key bindings
type BindingType string

const (
	BindingBasic     BindingType = "basic"
	BindingModTap    BindingType = "mod_tap"
	BindingLayerTap  BindingType = "layer_tap"
	BindingMacro     BindingType = "macro"
	BindingCombo     BindingType = "combo"
	BindingBehavior  BindingType = "behavior"
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

// Behavior represents custom ZMK behaviors
type Behavior struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// Combo represents key combinations
type Combo struct {
	Name    string     `json:"name"`
	Keys    []Position `json:"keys"`
	Binding string     `json:"binding"`
	Layers  []int      `json:"layers,omitempty"`
	Timeout int        `json:"timeout,omitempty"`
}