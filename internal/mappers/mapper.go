package mappers

import (
	"fmt"
	"masters3d.com/keyboard_layout_config_mapper/internal/models"
)

// KeyboardMapper defines the interface for converting between keyboard-specific
// layouts and the intermediate representation (IR)
type KeyboardMapper interface {
	// ToIR converts a keyboard-specific layout to intermediate representation
	ToIR(layout *models.KeyboardLayout) (*models.IRLayout, error)
	
	// FromIR converts an intermediate representation to keyboard-specific layout
	FromIR(irLayout *models.IRLayout) (*models.KeyboardLayout, error)
	
	// GetKeyboardType returns the keyboard type this mapper supports
	GetKeyboardType() models.KeyboardType
	
	// GetPositionMapping returns the mapping from keyboard positions to IR positions
	GetPositionMapping() map[string]models.IRPosition
	
	// GetReverseMapping returns the mapping from IR positions to keyboard positions
	GetReverseMapping() map[string]models.Position
}

// PositionMapper handles the coordinate transformations between keyboard-specific
// positions and the universal IR coordinate system
type PositionMapper struct {
	KeyboardType     models.KeyboardType
	ToIRMapping      map[string]models.IRPosition // keyboard pos -> IR pos
	FromIRMapping    map[string]models.Position   // IR pos -> keyboard pos
	ZoneMapping      map[string]models.IRZone     // keyboard zone -> IR zone
	KeyCodeMapping   map[string]string            // keyboard key code -> IR key code
}

// TranslateToIR converts a keyboard-specific position to IR position
func (pm *PositionMapper) TranslateToIR(pos models.Position) (models.IRPosition, error) {
	key := fmt.Sprintf("%s_%d_%d", pos.Side, pos.Row, pos.Col)
	irPos, exists := pm.ToIRMapping[key]
	if !exists {
		return models.IRPosition{}, fmt.Errorf("no IR mapping found for position %s", key)
	}
	return irPos, nil
}

// TranslateFromIR converts an IR position to keyboard-specific position
func (pm *PositionMapper) TranslateFromIR(irPos models.IRPosition) (models.Position, error) {
	key := irPos.IRCoordinate()
	pos, exists := pm.FromIRMapping[key]
	if !exists {
		return models.Position{}, fmt.Errorf("no keyboard mapping found for IR position %s", key)
	}
	return pos, nil
}

// TranslateKeyCode converts keyboard-specific key codes to normalized IR key codes
func (pm *PositionMapper) TranslateKeyCode(keyCode string) string {
	if irCode, exists := pm.KeyCodeMapping[keyCode]; exists {
		return irCode
	}
	// Return original if no mapping found (pass-through for standard codes)
	return keyCode
}

// LayoutTranslator provides high-level translation between different keyboard layouts
// via the intermediate representation
type LayoutTranslator struct {
	mappers map[models.KeyboardType]KeyboardMapper
}

// NewLayoutTranslator creates a new layout translator with registered mappers
func NewLayoutTranslator() *LayoutTranslator {
	return &LayoutTranslator{
		mappers: make(map[models.KeyboardType]KeyboardMapper),
	}
}

// RegisterMapper registers a keyboard mapper for a specific keyboard type
func (lt *LayoutTranslator) RegisterMapper(mapper KeyboardMapper) {
	lt.mappers[mapper.GetKeyboardType()] = mapper
}

// Translate converts a layout from one keyboard type to another via IR
func (lt *LayoutTranslator) Translate(source *models.KeyboardLayout, targetType models.KeyboardType) (*models.KeyboardLayout, error) {
	// Get source mapper
	sourceMapper, exists := lt.mappers[source.Type]
	if !exists {
		return nil, fmt.Errorf("no mapper found for source keyboard type: %s", source.Type)
	}
	
	// Get target mapper
	targetMapper, exists := lt.mappers[targetType]
	if !exists {
		return nil, fmt.Errorf("no mapper found for target keyboard type: %s", targetType)
	}
	
	// Convert source to IR
	irLayout, err := sourceMapper.ToIR(source)
	if err != nil {
		return nil, fmt.Errorf("failed to convert source to IR: %w", err)
	}
	
	// Convert IR to target
	targetLayout, err := targetMapper.FromIR(irLayout)
	if err != nil {
		return nil, fmt.Errorf("failed to convert IR to target: %w", err)
	}
	
	return targetLayout, nil
}

// GetSupportedTypes returns all supported keyboard types
func (lt *LayoutTranslator) GetSupportedTypes() []models.KeyboardType {
	types := make([]models.KeyboardType, 0, len(lt.mappers))
	for keyboardType := range lt.mappers {
		types = append(types, keyboardType)
	}
	return types
}

// ValidateMapping checks if a position mapping is valid and complete
func ValidateMapping(toIRMapping map[string]models.IRPosition, fromIRMapping map[string]models.Position) error {
	// Check that mappings are consistent (bidirectional)
	for kbPos, irPos := range toIRMapping {
		irKey := irPos.IRCoordinate()
		if mappedKbPos, exists := fromIRMapping[irKey]; exists {
			kbKey := fmt.Sprintf("%s_%d_%d", mappedKbPos.Side, mappedKbPos.Row, mappedKbPos.Col)
			if kbKey != kbPos {
				return fmt.Errorf("inconsistent mapping: KB pos %s maps to IR %s, but IR %s maps back to KB %s", 
					kbPos, irKey, irKey, kbKey)
			}
		} else {
			return fmt.Errorf("missing reverse mapping for IR position %s", irKey)
		}
	}
	
	return nil
}