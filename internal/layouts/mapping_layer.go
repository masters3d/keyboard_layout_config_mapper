// Package layouts defines mapping layer conventions for keyboard translation.
//
// MAPPING LAYER CONCEPT:
// Each keyboard keymap includes a special "mapping layer" (layer2_mapping or mapping_layer)
// that defines the logical identity of each physical key position. The binding VALUE in this
// layer serves as the KeyID, while its POSITION defines where that key lives on this keyboard.
//
// This allows us to:
// 1. Parse the mapping layer to understand the keyboard's physical layout
// 2. Match keys between keyboards by their logical KeyID
// 3. Translate layouts accurately regardless of different physical structures
//
// KEYID EXTRACTION:
// The KeyID is extracted from the binding value:
// - "&kp Q" → "Q"
// - "&kp KP_N1" → "KP_N1" (for thumb cluster identification)
// - "&none" → "" (no mapping)
// - "&trans" → "" (no mapping)

package layouts

import (
	"regexp"
	"strings"
)

// MappingLayerNames are the conventional names for mapping layers
var MappingLayerNames = []string{"layer2_mapping", "mapping_layer", "layer_mapping"}

// KeyIDFromBinding extracts the KeyID from a ZMK binding string
// Examples:
//   "&kp Q" -> "Q"
//   "&kp KP_N1" -> "KP_N1"
//   "&kp LC(A)" -> "LC_A" (normalized)
//   "&mo 1" -> "MO_1"
//   "&none" -> ""
//   "&trans" -> ""
func KeyIDFromBinding(binding string) string {
	binding = strings.TrimSpace(binding)
	
	// Skip non-mappable bindings
	if binding == "&none" || binding == "&trans" || binding == "" {
		return ""
	}
	
	// Handle &kp bindings
	if strings.HasPrefix(binding, "&kp ") {
		keyCode := strings.TrimPrefix(binding, "&kp ")
		// Normalize modifiers like LC(A) to LC_A
		keyCode = strings.ReplaceAll(keyCode, "(", "_")
		keyCode = strings.ReplaceAll(keyCode, ")", "")
		return strings.TrimSpace(keyCode)
	}
	
	// Handle &mo bindings (layer momentary)
	if strings.HasPrefix(binding, "&mo ") {
		layer := strings.TrimPrefix(binding, "&mo ")
		return "MO_" + strings.TrimSpace(layer)
	}
	
	// Handle &to bindings (layer toggle)
	if strings.HasPrefix(binding, "&to ") {
		layer := strings.TrimPrefix(binding, "&to ")
		return "TO_" + strings.TrimSpace(layer)
	}
	
	// Handle &tog bindings (layer toggle)
	if strings.HasPrefix(binding, "&tog ") {
		layer := strings.TrimPrefix(binding, "&tog ")
		return "TOG_" + strings.TrimSpace(layer)
	}
	
	// Handle custom behaviors (like &morph_dot)
	if strings.HasPrefix(binding, "&") {
		// Extract behavior name
		re := regexp.MustCompile(`^&(\w+)`)
		if matches := re.FindStringSubmatch(binding); len(matches) > 1 {
			return "BEH_" + strings.ToUpper(matches[1])
		}
	}
	
	return ""
}

// ParseMappingLayerBindings parses the bindings from a mapping layer
// Returns a map from index to KeyID
func ParseMappingLayerBindings(bindings []string) map[int]string {
	result := make(map[int]string)
	
	for i, binding := range bindings {
		keyID := KeyIDFromBinding(binding)
		if keyID != "" {
			result[i] = keyID
		}
	}
	
	return result
}

// BuildIndexToKeyIDMap creates a bidirectional mapping for a keyboard's mapping layer
type KeyMapping struct {
	IndexToKeyID map[int]string    // Physical position → logical key ID
	KeyIDToIndex map[string]int    // Logical key ID → physical position
}

// NewKeyMapping creates a new key mapping from binding values
func NewKeyMapping(bindings []string) *KeyMapping {
	km := &KeyMapping{
		IndexToKeyID: make(map[int]string),
		KeyIDToIndex: make(map[string]int),
	}
	
	for i, binding := range bindings {
		keyID := KeyIDFromBinding(binding)
		if keyID != "" {
			km.IndexToKeyID[i] = keyID
			km.KeyIDToIndex[keyID] = i
		}
	}
	
	return km
}

// StandardKeyIDs defines the canonical set of KeyIDs used across all keyboards
var StandardKeyIDs = map[string]string{
	// Alpha keys (direct from ZMK key codes)
	"Q": "Q", "W": "W", "E": "E", "R": "R", "T": "T",
	"Y": "Y", "U": "U", "I": "I", "O": "O", "P": "P",
	"A": "A", "S": "S", "D": "D", "F": "F", "G": "G",
	"H": "H", "J": "J", "K": "K", "L": "L",
	"Z": "Z", "X": "X", "C": "C", "V": "V", "B": "B",
	"N": "N", "M": "M",

	// Number row
	"N1": "N1", "N2": "N2", "N3": "N3", "N4": "N4", "N5": "N5",
	"N6": "N6", "N7": "N7", "N8": "N8", "N9": "N9", "N0": "N0",

	// Function keys
	"F1": "F1", "F2": "F2", "F3": "F3", "F4": "F4", "F5": "F5", "F6": "F6",
	"F7": "F7", "F8": "F8", "F9": "F9", "F10": "F10", "F11": "F11", "F12": "F12",
	"F13": "F13", "F14": "F14", "F15": "F15", "F16": "F16", "F17": "F17", "F18": "F18",

	// Symbols
	"MINUS": "MINUS", "EQUAL": "EQUAL", "LBKT": "LBKT", "RBKT": "RBKT",
	"BSLH": "BSLH", "SEMI": "SEMI", "SQT": "SQT", "GRAVE": "GRAVE",
	"COMMA": "COMMA", "DOT": "DOT", "SLASH": "SLASH",

	// Modifiers
	"LSHFT": "LSHFT", "RSHFT": "RSHFT", "LCTRL": "LCTRL", "RCTRL": "RCTRL",
	"LALT": "LALT", "RALT": "RALT", "LGUI": "LGUI", "RGUI": "RGUI",

	// Special keys
	"SPACE": "SPACE", "BSPC": "BSPC", "DEL": "DEL", "RET": "RET",
	"TAB": "TAB", "ESC": "ESC", "CAPS": "CAPS",
	"HOME": "HOME", "END": "END", "PG_UP": "PG_UP", "PG_DN": "PG_DN",

	// Arrow keys
	"LEFT": "LEFT", "RIGHT": "RIGHT", "UP": "UP", "DOWN": "DOWN",

	// Thumb cluster (using keypad as unique identifiers)
	// These provide consistent naming for thumb keys across keyboards
	"KP_N1": "THUMB_L1", "KP_N2": "THUMB_L2", "KP_N3": "THUMB_L3",
	"KP_N4": "THUMB_L4", "KP_N5": "THUMB_L5", "KP_N6": "THUMB_L6",
	"KP_N7": "THUMB_R1", "KP_N8": "THUMB_R2", "KP_N9": "THUMB_R3",
	"KP_MULTIPLY": "THUMB_R4", "KP_N0": "THUMB_R5", "KP_PLUS": "THUMB_R6",
}

// ThumbKeyOrder defines a consistent ordering for thumb keys across keyboards
var ThumbKeyOrder = []string{
	// Left thumb: inner to outer, top to bottom
	"THUMB_L1", "THUMB_L2", "THUMB_L3",
	"THUMB_L4", "THUMB_L5", "THUMB_L6",
	// Right thumb: inner to outer, top to bottom  
	"THUMB_R1", "THUMB_R2", "THUMB_R3",
	"THUMB_R4", "THUMB_R5", "THUMB_R6",
}
