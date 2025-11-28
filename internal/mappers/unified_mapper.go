package mappers

import (
	"fmt"
	"masters3d.com/keyboard_layout_config_mapper/internal/layouts"
	"masters3d.com/keyboard_layout_config_mapper/internal/models"
)

// UnifiedMapper provides translation between keyboard layouts using KeyID-based matching.
// Instead of mapping by grid coordinates, it maps by semantic key identity (KeyID).
// This allows accurate translation even between keyboards with different physical layouts.
type UnifiedMapper struct {
	sourceLayout *layouts.PhysicalLayout
	targetLayout *layouts.PhysicalLayout
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

// TranslateBindings translates bindings from source to target layout using KeyID matching
func (m *UnifiedMapper) TranslateBindings(sourceBindings []models.KeyBinding) []models.KeyBinding {
	// Build KeyID maps for both layouts
	sourceKeyIDMap := BuildKeyIDToIndexMap(m.sourceLayout)
	targetKeyIDMap := BuildKeyIDToIndexMap(m.targetLayout)
	
	// Build reverse map: source index -> KeyID
	sourceIndexToKeyID := make(map[int]string)
	for keyID, idx := range sourceKeyIDMap {
		sourceIndexToKeyID[idx] = keyID
	}
	
	// Initialize target bindings with &trans
	targetBindings := make([]models.KeyBinding, m.targetLayout.TotalKeys)
	for i := range targetBindings {
		irPos, _ := m.targetLayout.GetIRPosition(i)
		targetBindings[i] = models.KeyBinding{
			Position: models.Position{
				Row:   irPos.Row,
				Col:   irPos.Col,
				Side:  irPos.Hand,
				Zone:  irPos.Zone,
				KeyID: irPos.KeyID,
			},
			Value: "&trans",
			Type:  models.BindingBasic,
		}
	}
	
	// Map source bindings to target by KeyID
	for i, binding := range sourceBindings {
		// Find the KeyID for this source position
		keyID := binding.Position.KeyID
		if keyID == "" {
			// Try to look up by index
			if id, ok := sourceIndexToKeyID[i]; ok {
				keyID = id
			}
		}
		
		if keyID == "" {
			continue // Can't map without KeyID
		}
		
		// Find the target index for this KeyID
		if targetIdx, ok := targetKeyIDMap[keyID]; ok {
			targetBindings[targetIdx] = models.KeyBinding{
				Position: targetBindings[targetIdx].Position,
				Value:    binding.Value,
				Layer:    binding.Layer,
				Type:     binding.Type,
				Metadata: binding.Metadata,
			}
		}
		// If KeyID doesn't exist in target, the binding is dropped (stays &trans)
	}
	
	return targetBindings
}

// TranslateLayout translates an entire layout from source to target keyboard type
func (m *UnifiedMapper) TranslateLayout(source *models.KeyboardLayout) (*models.KeyboardLayout, error) {
	target := &models.KeyboardLayout{
		Type:      m.targetLayout.KeyboardType,
		Name:      source.Name,
		FilePath:  "",
		Layers:    make([]models.Layer, len(source.Layers)),
		Behaviors: source.Behaviors, // Copy behaviors as-is
		Combos:    []models.Combo{}, // Combos need special handling
		Metadata:  source.Metadata,
	}
	
	for i, srcLayer := range source.Layers {
		targetBindings := m.TranslateBindings(srcLayer.Bindings)
		target.Layers[i] = models.Layer{
			Index:    srcLayer.Index,
			Name:     srcLayer.Name,
			Bindings: targetBindings,
			Metadata: srcLayer.Metadata,
		}
	}
	
	return target, nil
}

// GetMappingSummary returns a summary of how keys map between keyboards
func (m *UnifiedMapper) GetMappingSummary() map[string]string {
	sourceKeyIDMap := BuildKeyIDToIndexMap(m.sourceLayout)
	targetKeyIDMap := BuildKeyIDToIndexMap(m.targetLayout)
	
	result := make(map[string]string)
	
	// Find common keys
	for keyID := range sourceKeyIDMap {
		if _, ok := targetKeyIDMap[keyID]; ok {
			result[keyID] = "mapped"
		} else {
			result[keyID] = "source_only"
		}
	}
	
	// Find target-only keys
	for keyID := range targetKeyIDMap {
		if _, ok := sourceKeyIDMap[keyID]; !ok {
			result[keyID] = "target_only"
		}
	}
	
	return result
}
