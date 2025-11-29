package mappers

import (
	"fmt"
	"strings"
	"masters3d.com/keyboard_layout_config_mapper/internal/layouts"
	"masters3d.com/keyboard_layout_config_mapper/internal/models"
)

// UnifiedMapper provides translation between keyboard layouts using KeyID-based matching.
// It uses the "mapping layer" in each keyboard's keymap to determine the logical identity
// of each physical key position, enabling accurate translation between different keyboards.
type UnifiedMapper struct {
	sourceLayout *layouts.PhysicalLayout
	targetLayout *layouts.PhysicalLayout
	sourceMapping *layouts.KeyMapping
	targetMapping *layouts.KeyMapping
}

// NewUnifiedMapper creates a mapper for translating between two keyboard types
func NewUnifiedMapper(sourceType, targetType models.KeyboardType) (*UnifiedMapper, error) {
	sourceLayout, err := layouts.GetLayoutForKeyboard(sourceType)
	if err != nil || sourceLayout == nil {
		return nil, fmt.Errorf("no physical layout defined for source keyboard: %s", sourceType)
	}
	
	targetLayout, err := layouts.GetLayoutForKeyboard(targetType)
	if err != nil || targetLayout == nil {
		return nil, fmt.Errorf("no physical layout defined for target keyboard: %s", targetType)
	}
	
	return &UnifiedMapper{
		sourceLayout: sourceLayout,
		targetLayout: targetLayout,
	}, nil
}

// NewUnifiedMapperFromMappingLayers creates a mapper using mapping layer bindings from keymaps
func NewUnifiedMapperFromMappingLayers(sourceBindings, targetBindings []string) *UnifiedMapper {
	return &UnifiedMapper{
		sourceMapping: layouts.NewKeyMapping(sourceBindings),
		targetMapping: layouts.NewKeyMapping(targetBindings),
	}
}

// SetMappingLayers sets the mapping layer bindings for translation
func (m *UnifiedMapper) SetMappingLayers(sourceBindings, targetBindings []string) {
	m.sourceMapping = layouts.NewKeyMapping(sourceBindings)
	m.targetMapping = layouts.NewKeyMapping(targetBindings)
}

// TranslateBindings translates bindings from source to target layout using KeyID matching
func (m *UnifiedMapper) TranslateBindings(sourceBindings []models.KeyBinding) []models.KeyBinding {
	// Determine which mapping to use
	var sourceIndexToKeyID map[int]string
	var targetKeyIDToIndex map[string]int
	var targetTotalKeys int
	
	if m.sourceMapping != nil && m.targetMapping != nil {
		// Use mapping layer-based translation
		sourceIndexToKeyID = m.sourceMapping.IndexToKeyID
		targetKeyIDToIndex = m.targetMapping.KeyIDToIndex
		targetTotalKeys = len(m.targetMapping.IndexToKeyID)
		// Estimate total keys from max index
		for idx := range m.targetMapping.IndexToKeyID {
			if idx >= targetTotalKeys {
				targetTotalKeys = idx + 1
			}
		}
	} else if m.sourceLayout != nil && m.targetLayout != nil {
		// Fall back to physical layout-based translation
		// Build index→keyID map for source
		sourceIndexToKeyID = make(map[int]string)
		for i := 0; i < m.sourceLayout.TotalKeys; i++ {
			if irPos, ok := m.sourceLayout.GetIRPosition(i); ok {
				sourceIndexToKeyID[i] = irPos.KeyID
			}
		}
		// Build keyID→index map for target
		targetKeyIDToIndex = make(map[string]int)
		for i := 0; i < m.targetLayout.TotalKeys; i++ {
			if irPos, ok := m.targetLayout.GetIRPosition(i); ok && irPos.KeyID != "" {
				targetKeyIDToIndex[irPos.KeyID] = i
			}
		}
		targetTotalKeys = m.targetLayout.TotalKeys
	} else {
		// No mapping available
		return sourceBindings
	}
	
	// Initialize target bindings with &trans, setting position based on index
	targetBindings := make([]models.KeyBinding, targetTotalKeys)
	for i := range targetBindings {
		// Calculate position from index - this needs to match the generator's expectations
		pos := indexToPosition(i, targetTotalKeys)
		targetBindings[i] = models.KeyBinding{
			Position: pos,
			Value: "&trans",
			Type:  models.BindingBasic,
		}
	}
	
	// Map source bindings to target by KeyID
	for i, binding := range sourceBindings {
		// Find the KeyID for this source position
		keyID, ok := sourceIndexToKeyID[i]
		if !ok || keyID == "" {
			continue
		}
		
		// Find the target index for this KeyID
		if targetIdx, ok := targetKeyIDToIndex[keyID]; ok && targetIdx < len(targetBindings) {
			// Preserve the position from initialization
			pos := targetBindings[targetIdx].Position
			targetBindings[targetIdx] = models.KeyBinding{
				Position: pos,
				Value:    binding.Value,
				Layer:    binding.Layer,
				Type:     binding.Type,
				Metadata: binding.Metadata,
			}
		}
	}
	
	return targetBindings
}

// TranslateLayout translates an entire layout from source to target keyboard type
// It filters out padding layers and only translates essential layers (default, keypad, cmd, mapping).
// System/bootloader layers are NOT translated - they must be preserved from the target keyboard
// because they contain hardware-specific configurations.
func (m *UnifiedMapper) TranslateLayout(source *models.KeyboardLayout) (*models.KeyboardLayout, error) {
	var targetType models.KeyboardType
	if m.targetLayout != nil {
		targetType = m.targetLayout.KeyboardType
	}
	
	// Filter layers: skip padding layers, translate essential layers
	var translatedLayers []models.Layer
	newIndex := 0
	
	for _, srcLayer := range source.Layers {
		layerName := strings.ToLower(srcLayer.Name)
		
		// Skip padding layers
		if strings.Contains(layerName, "padding") {
			continue
		}
		
		// Skip system/bootloader/mod layers - they should be preserved from target, not translated
		// These typically contain hardware-specific bootloader and Bluetooth configurations
		if strings.Contains(layerName, "system") || strings.Contains(layerName, "boot") {
			continue
		}
		// "mod" layers are often system layers in adv360 (layer8_mod is the bootloader layer)
		if strings.HasSuffix(layerName, "_mod") || strings.HasSuffix(layerName, "mod") {
			continue
		}
		
		// Translate the layer
		targetBindings := m.TranslateBindings(srcLayer.Bindings)
		
		// Rename layers to semantic names
		newName := m.getSemanticLayerName(srcLayer.Name, newIndex)
		
		translatedLayers = append(translatedLayers, models.Layer{
			Index:    newIndex,
			Name:     newName,
			Bindings: targetBindings,
			Metadata: srcLayer.Metadata,
		})
		newIndex++
	}
	
	target := &models.KeyboardLayout{
		Type:      targetType,
		Name:      source.Name,
		FilePath:  "",
		Layers:    translatedLayers,
		Behaviors: source.Behaviors,
		Combos:    []models.Combo{},
		Metadata:  source.Metadata,
	}
	
	return target, nil
}

// getSemanticLayerName returns a semantic name for the layer
func (m *UnifiedMapper) getSemanticLayerName(originalName string, index int) string {
	nameLower := strings.ToLower(originalName)
	
	switch {
	case index == 0 || strings.Contains(nameLower, "default"):
		return "default"
	case strings.Contains(nameLower, "qwerty"):
		return "qwerty"
	case strings.Contains(nameLower, "keypad") || strings.Contains(nameLower, "numpad"):
		return "keypad"
	case strings.Contains(nameLower, "cmd") || strings.Contains(nameLower, "fn"):
		return "cmd"
	case strings.Contains(nameLower, "mapping"):
		return "mapping"
	default:
		return originalName
	}
}

// BuildKeyIDToIndexMap creates a map from KeyID to physical index for a layout
func BuildKeyIDToIndexMap(layout *layouts.PhysicalLayout) map[string]int {
	result := make(map[string]int)
	for i, pos := range layout.PositionMap {
		if pos.KeyID != "" {
			result[pos.KeyID] = i
		}
	}
	return result
}

// GetMappingSummary returns a summary of how keys map between keyboards
func (m *UnifiedMapper) GetMappingSummary() map[string]string {
	result := make(map[string]string)
	
	var sourceKeys, targetKeys map[string]int
	
	if m.sourceMapping != nil && m.targetMapping != nil {
		sourceKeys = m.sourceMapping.KeyIDToIndex
		targetKeys = m.targetMapping.KeyIDToIndex
	} else if m.sourceLayout != nil && m.targetLayout != nil {
		sourceKeys = BuildKeyIDToIndexMap(m.sourceLayout)
		targetKeys = BuildKeyIDToIndexMap(m.targetLayout)
	} else {
		return result
	}
	
	// Find common keys
	for keyID := range sourceKeys {
		if _, ok := targetKeys[keyID]; ok {
			result[keyID] = "mapped"
		} else {
			result[keyID] = "source_only"
		}
	}
	
	// Find target-only keys
	for keyID := range targetKeys {
		if _, ok := sourceKeys[keyID]; !ok {
			result[keyID] = "target_only"
		}
	}
	
	return result
}

// indexToPosition converts a linear index to a Position struct for adv_mod layout
// This matches the row structure expected by the generator:
// Row 0: 18 keys (function row) - indices 0-17
// Row 1: 12 keys (number row)   - indices 18-29
// Row 2: 12 keys (QWERTY)       - indices 30-41
// Row 3: 12 keys (home row)     - indices 42-53
// Row 4: 12 keys (bottom row)   - indices 54-65
// Row 5: 8 keys (modifier)      - indices 66-73
// Row 6: 4 keys (thumb top)     - indices 74-77
// Row 7: 2 keys (thumb mid)     - indices 78-79
// Row 8: 6 keys (thumb bottom)  - indices 80-85
func indexToPosition(index, totalKeys int) models.Position {
	var row, col int
	var side string
	
	switch {
	case index < 18:
		// Function row (9 left + 9 right)
		row = 0
		if index < 9 {
			side = "left"
			col = index
		} else {
			side = "right"
			col = index - 9
		}
	case index < 30:
		// Number row (6 left + 6 right)
		row = 1
		idx := index - 18
		if idx < 6 {
			side = "left"
			col = idx
		} else {
			side = "right"
			col = idx - 6
		}
	case index < 42:
		// QWERTY row (6 left + 6 right)
		row = 2
		idx := index - 30
		if idx < 6 {
			side = "left"
			col = idx
		} else {
			side = "right"
			col = idx - 6
		}
	case index < 54:
		// Home row (6 left + 6 right)
		row = 3
		idx := index - 42
		if idx < 6 {
			side = "left"
			col = idx
		} else {
			side = "right"
			col = idx - 6
		}
	case index < 66:
		// Bottom row (6 left + 6 right)
		row = 4
		idx := index - 54
		if idx < 6 {
			side = "left"
			col = idx
		} else {
			side = "right"
			col = idx - 6
		}
	case index < 74:
		// Modifier row (4 left + 4 right)
		row = 5
		idx := index - 66
		if idx < 4 {
			side = "left"
			col = idx
		} else {
			side = "right"
			col = idx - 4
		}
	case index < 78:
		// Thumb top (2 left + 2 right)
		row = 6
		idx := index - 74
		if idx < 2 {
			side = "left"
			col = idx
		} else {
			side = "right"
			col = idx - 2
		}
	case index < 80:
		// Thumb middle (1 left + 1 right)
		row = 7
		idx := index - 78
		if idx < 1 {
			side = "left"
			col = 0
		} else {
			side = "right"
			col = 0
		}
	default:
		// Thumb bottom (3 left + 3 right)
		row = 8
		idx := index - 80
		if idx < 3 {
			side = "left"
			col = idx
		} else {
			side = "right"
			col = idx - 3
		}
	}
	
	return models.Position{
		Row:  row,
		Col:  col,
		Side: side,
		Zone: "main",
	}
}
