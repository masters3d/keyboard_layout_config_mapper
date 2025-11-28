package mappers

import (
	"strconv"
	"strings"
	"masters3d.com/keyboard_layout_config_mapper/internal/models"
)

// ZMKAdv360Mapper handles mapping between Advantage360 layout and IR
type ZMKAdv360Mapper struct {
	positionMapper *PositionMapper
}

// NewZMKAdv360Mapper creates a new mapper for Advantage360 keyboards
func NewZMKAdv360Mapper() *ZMKAdv360Mapper {
	mapper := &ZMKAdv360Mapper{
		positionMapper: &PositionMapper{
			KeyboardType:   models.KeyboardZMKAdv360,
			ToIRMapping:    createAdv360ToIRMapping(),
			FromIRMapping:  createAdv360FromIRMapping(),
			ZoneMapping:    createAdv360ZoneMapping(),
			KeyCodeMapping: createZMKKeyCodeMapping(),
		},
	}
	return mapper
}

// ToIR converts Advantage360 layout to intermediate representation
func (m *ZMKAdv360Mapper) ToIR(layout *models.KeyboardLayout) (*models.IRLayout, error) {
	irLayout := &models.IRLayout{
		Name:      layout.Name + "_IR",
		Source:    layout.Type,
		Layers:    make([]models.IRLayer, len(layout.Layers)),
		Behaviors: layout.Behaviors,
		Combos:    make([]models.IRCombo, len(layout.Combos)),
		Metadata:  layout.Metadata,
	}
	
	// Convert layers
	for i, layer := range layout.Layers {
		irLayer := models.IRLayer{
			Index:    layer.Index,
			Name:     layer.Name,
			Bindings: make([]models.IRKeyBinding, 0, len(layer.Bindings)),
			Metadata: layer.Metadata,
		}
		
		// Convert bindings
		for _, binding := range layer.Bindings {
			irPos, err := m.positionMapper.TranslateToIR(binding.Position)
			if err != nil {
				continue // Skip unmappable positions
			}
			
			irBinding := models.IRKeyBinding{
				Position: irPos,
				Value:    m.positionMapper.TranslateKeyCode(binding.Value),
				Layer:    binding.Layer,
				Type:     binding.Type,
				Metadata: binding.Metadata,
			}
			
			irLayer.Bindings = append(irLayer.Bindings, irBinding)
		}
		
		irLayout.Layers[i] = irLayer
	}
	
	// Convert combos
	for i, combo := range layout.Combos {
		irCombo := models.IRCombo{
			Name:    combo.Name,
			Keys:    make([]models.IRPosition, 0, len(combo.Keys)),
			Binding: m.positionMapper.TranslateKeyCode(combo.Binding),
			Layers:  combo.Layers,
			Timeout: combo.Timeout,
		}
		
		// Convert combo key positions
		for _, pos := range combo.Keys {
			if irPos, err := m.positionMapper.TranslateToIR(pos); err == nil {
				irCombo.Keys = append(irCombo.Keys, irPos)
			}
		}
		
		irLayout.Combos[i] = irCombo
	}
	
	return irLayout, nil
}

// FromIR converts intermediate representation to Advantage360 layout
func (m *ZMKAdv360Mapper) FromIR(irLayout *models.IRLayout) (*models.KeyboardLayout, error) {
	layout := &models.KeyboardLayout{
		Type:      models.KeyboardZMKAdv360,
		Name:      irLayout.Name,
		Layers:    make([]models.Layer, len(irLayout.Layers)),
		Behaviors: irLayout.Behaviors,
		Combos:    make([]models.Combo, 0, len(irLayout.Combos)),
		Metadata:  irLayout.Metadata,
	}
	
	// Convert layers
	for i, irLayer := range irLayout.Layers {
		layer := models.Layer{
			Index:    irLayer.Index,
			Name:     irLayer.Name,
			Bindings: make([]models.KeyBinding, 0, len(irLayer.Bindings)),
			Metadata: irLayer.Metadata,
		}
		
		// Convert bindings
		for _, irBinding := range irLayer.Bindings {
			pos, err := m.positionMapper.TranslateFromIR(irBinding.Position)
			if err != nil {
				continue // Skip unmappable positions
			}
			
			binding := models.KeyBinding{
				Position: pos,
				Value:    irBinding.Value, // Key codes should already be in ZMK format
				Layer:    irBinding.Layer,
				Type:     irBinding.Type,
				Metadata: irBinding.Metadata,
			}
			
			layer.Bindings = append(layer.Bindings, binding)
		}
		
		layout.Layers[i] = layer
	}
	
	// Convert combos
	for _, irCombo := range irLayout.Combos {
		combo := models.Combo{
			Name:    irCombo.Name,
			Keys:    make([]models.Position, 0, len(irCombo.Keys)),
			Binding: irCombo.Binding,
			Layers:  irCombo.Layers,
			Timeout: irCombo.Timeout,
		}
		
		// Convert combo key positions
		for _, irPos := range irCombo.Keys {
			if pos, err := m.positionMapper.TranslateFromIR(irPos); err == nil {
				combo.Keys = append(combo.Keys, pos)
			}
		}
		
		if len(combo.Keys) > 0 { // Only add combos with valid keys
			layout.Combos = append(layout.Combos, combo)
		}
	}
	
	return layout, nil
}

// GetKeyboardType returns the keyboard type this mapper supports
func (m *ZMKAdv360Mapper) GetKeyboardType() models.KeyboardType {
	return models.KeyboardZMKAdv360
}

// GetPositionMapping returns the position mapping
func (m *ZMKAdv360Mapper) GetPositionMapping() map[string]models.IRPosition {
	return m.positionMapper.ToIRMapping
}

// GetReverseMapping returns the reverse mapping
func (m *ZMKAdv360Mapper) GetReverseMapping() map[string]models.Position {
	return m.positionMapper.FromIRMapping
}

// createAdv360ToIRMapping creates the position mapping from Adv360 to IR coordinates
func createAdv360ToIRMapping() map[string]models.IRPosition {
	mapping := make(map[string]models.IRPosition)
	
	// Left hand main area (rows 0-3, standard QWERTY-like layout)
	// Function row (F1-F6)
	mapping["left_0_0"] = models.IRPosition{Hand: "left", Row: 0, Col: 0, Zone: "function", KeyID: "F1"}
	mapping["left_0_1"] = models.IRPosition{Hand: "left", Row: 0, Col: 1, Zone: "function", KeyID: "F2"}
	mapping["left_0_2"] = models.IRPosition{Hand: "left", Row: 0, Col: 2, Zone: "function", KeyID: "F3"}
	mapping["left_0_3"] = models.IRPosition{Hand: "left", Row: 0, Col: 3, Zone: "function", KeyID: "F4"}
	mapping["left_0_4"] = models.IRPosition{Hand: "left", Row: 0, Col: 4, Zone: "function", KeyID: "F5"}
	mapping["left_0_5"] = models.IRPosition{Hand: "left", Row: 0, Col: 5, Zone: "function", KeyID: "F6"}
	
	// Number row (1-6)
	mapping["left_1_0"] = models.IRPosition{Hand: "left", Row: 1, Col: 0, Zone: "main", KeyID: "n1"}
	mapping["left_1_1"] = models.IRPosition{Hand: "left", Row: 1, Col: 1, Zone: "main", KeyID: "n2"}
	mapping["left_1_2"] = models.IRPosition{Hand: "left", Row: 1, Col: 2, Zone: "main", KeyID: "n3"}
	mapping["left_1_3"] = models.IRPosition{Hand: "left", Row: 1, Col: 3, Zone: "main", KeyID: "n4"}
	mapping["left_1_4"] = models.IRPosition{Hand: "left", Row: 1, Col: 4, Zone: "main", KeyID: "n5"}
	mapping["left_1_5"] = models.IRPosition{Hand: "left", Row: 1, Col: 5, Zone: "main", KeyID: "n6"}
	
	// Top row (Q-T)
	mapping["left_2_0"] = models.IRPosition{Hand: "left", Row: 2, Col: 0, Zone: "main", KeyID: "q"}
	mapping["left_2_1"] = models.IRPosition{Hand: "left", Row: 2, Col: 1, Zone: "main", KeyID: "w"}
	mapping["left_2_2"] = models.IRPosition{Hand: "left", Row: 2, Col: 2, Zone: "main", KeyID: "e"}
	mapping["left_2_3"] = models.IRPosition{Hand: "left", Row: 2, Col: 3, Zone: "main", KeyID: "r"}
	mapping["left_2_4"] = models.IRPosition{Hand: "left", Row: 2, Col: 4, Zone: "main", KeyID: "t"}
	
	// Home row (A-G)
	mapping["left_3_0"] = models.IRPosition{Hand: "left", Row: 3, Col: 0, Zone: "main", KeyID: "a"}
	mapping["left_3_1"] = models.IRPosition{Hand: "left", Row: 3, Col: 1, Zone: "main", KeyID: "s"}
	mapping["left_3_2"] = models.IRPosition{Hand: "left", Row: 3, Col: 2, Zone: "main", KeyID: "d"}
	mapping["left_3_3"] = models.IRPosition{Hand: "left", Row: 3, Col: 3, Zone: "main", KeyID: "f"}
	mapping["left_3_4"] = models.IRPosition{Hand: "left", Row: 3, Col: 4, Zone: "main", KeyID: "g"}
	
	// Bottom row (Z-B)
	mapping["left_4_0"] = models.IRPosition{Hand: "left", Row: 4, Col: 0, Zone: "main", KeyID: "z"}
	mapping["left_4_1"] = models.IRPosition{Hand: "left", Row: 4, Col: 1, Zone: "main", KeyID: "x"}
	mapping["left_4_2"] = models.IRPosition{Hand: "left", Row: 4, Col: 2, Zone: "main", KeyID: "c"}
	mapping["left_4_3"] = models.IRPosition{Hand: "left", Row: 4, Col: 3, Zone: "main", KeyID: "v"}
	mapping["left_4_4"] = models.IRPosition{Hand: "left", Row: 4, Col: 4, Zone: "main", KeyID: "b"}
	
	// Left thumb cluster
	mapping["left_5_0"] = models.IRPosition{Hand: "left", Row: 5, Col: 0, Zone: "thumb", KeyID: "thumb1"}
	mapping["left_5_1"] = models.IRPosition{Hand: "left", Row: 5, Col: 1, Zone: "thumb", KeyID: "thumb2"}
	mapping["left_5_2"] = models.IRPosition{Hand: "left", Row: 5, Col: 2, Zone: "thumb", KeyID: "thumb3"}
	mapping["left_5_3"] = models.IRPosition{Hand: "left", Row: 5, Col: 3, Zone: "thumb", KeyID: "thumb4"}
	
	// Right hand main area - mirror of left
	// Function row (F7-F12)
	mapping["right_0_0"] = models.IRPosition{Hand: "right", Row: 0, Col: 0, Zone: "function", KeyID: "F7"}
	mapping["right_0_1"] = models.IRPosition{Hand: "right", Row: 0, Col: 1, Zone: "function", KeyID: "F8"}
	mapping["right_0_2"] = models.IRPosition{Hand: "right", Row: 0, Col: 2, Zone: "function", KeyID: "F9"}
	mapping["right_0_3"] = models.IRPosition{Hand: "right", Row: 0, Col: 3, Zone: "function", KeyID: "F10"}
	mapping["right_0_4"] = models.IRPosition{Hand: "right", Row: 0, Col: 4, Zone: "function", KeyID: "F11"}
	mapping["right_0_5"] = models.IRPosition{Hand: "right", Row: 0, Col: 5, Zone: "function", KeyID: "F12"}
	
	// Number row (7-0)
	mapping["right_1_0"] = models.IRPosition{Hand: "right", Row: 1, Col: 0, Zone: "main", KeyID: "n7"}
	mapping["right_1_1"] = models.IRPosition{Hand: "right", Row: 1, Col: 1, Zone: "main", KeyID: "n8"}
	mapping["right_1_2"] = models.IRPosition{Hand: "right", Row: 1, Col: 2, Zone: "main", KeyID: "n9"}
	mapping["right_1_3"] = models.IRPosition{Hand: "right", Row: 1, Col: 3, Zone: "main", KeyID: "n0"}
	mapping["right_1_4"] = models.IRPosition{Hand: "right", Row: 1, Col: 4, Zone: "main", KeyID: "minus"}
	mapping["right_1_5"] = models.IRPosition{Hand: "right", Row: 1, Col: 5, Zone: "main", KeyID: "equal"}
	
	// Top row (Y-P)
	mapping["right_2_0"] = models.IRPosition{Hand: "right", Row: 2, Col: 0, Zone: "main", KeyID: "y"}
	mapping["right_2_1"] = models.IRPosition{Hand: "right", Row: 2, Col: 1, Zone: "main", KeyID: "u"}
	mapping["right_2_2"] = models.IRPosition{Hand: "right", Row: 2, Col: 2, Zone: "main", KeyID: "i"}
	mapping["right_2_3"] = models.IRPosition{Hand: "right", Row: 2, Col: 3, Zone: "main", KeyID: "o"}
	mapping["right_2_4"] = models.IRPosition{Hand: "right", Row: 2, Col: 4, Zone: "main", KeyID: "p"}
	mapping["right_2_5"] = models.IRPosition{Hand: "right", Row: 2, Col: 5, Zone: "main", KeyID: "lbkt"}
	
	// Home row (H-;)
	mapping["right_3_0"] = models.IRPosition{Hand: "right", Row: 3, Col: 0, Zone: "main", KeyID: "h"}
	mapping["right_3_1"] = models.IRPosition{Hand: "right", Row: 3, Col: 1, Zone: "main", KeyID: "j"}
	mapping["right_3_2"] = models.IRPosition{Hand: "right", Row: 3, Col: 2, Zone: "main", KeyID: "k"}
	mapping["right_3_3"] = models.IRPosition{Hand: "right", Row: 3, Col: 3, Zone: "main", KeyID: "l"}
	mapping["right_3_4"] = models.IRPosition{Hand: "right", Row: 3, Col: 4, Zone: "main", KeyID: "semi"}
	mapping["right_3_5"] = models.IRPosition{Hand: "right", Row: 3, Col: 5, Zone: "main", KeyID: "sqt"}
	
	// Bottom row (N-/)
	mapping["right_4_0"] = models.IRPosition{Hand: "right", Row: 4, Col: 0, Zone: "main", KeyID: "n"}
	mapping["right_4_1"] = models.IRPosition{Hand: "right", Row: 4, Col: 1, Zone: "main", KeyID: "m"}
	mapping["right_4_2"] = models.IRPosition{Hand: "right", Row: 4, Col: 2, Zone: "main", KeyID: "comma"}
	mapping["right_4_3"] = models.IRPosition{Hand: "right", Row: 4, Col: 3, Zone: "main", KeyID: "dot"}
	mapping["right_4_4"] = models.IRPosition{Hand: "right", Row: 4, Col: 4, Zone: "main", KeyID: "fslh"}
	
	// Right thumb cluster
	mapping["right_5_0"] = models.IRPosition{Hand: "right", Row: 5, Col: 0, Zone: "thumb", KeyID: "thumb1"}
	mapping["right_5_1"] = models.IRPosition{Hand: "right", Row: 5, Col: 1, Zone: "thumb", KeyID: "thumb2"}
	mapping["right_5_2"] = models.IRPosition{Hand: "right", Row: 5, Col: 2, Zone: "thumb", KeyID: "thumb3"}
	mapping["right_5_3"] = models.IRPosition{Hand: "right", Row: 5, Col: 3, Zone: "thumb", KeyID: "thumb4"}
	
	return mapping
}

// createAdv360FromIRMapping creates the reverse mapping from IR to Adv360 coordinates
func createAdv360FromIRMapping() map[string]models.Position {
	toIR := createAdv360ToIRMapping()
	fromIR := make(map[string]models.Position)
	
	for kbPos, irPos := range toIR {
		// Parse keyboard position string: "left_0_0" -> side="left", row=0, col=0
		parts := strings.Split(kbPos, "_")
		if len(parts) == 3 {
			side := parts[0]
			row, err1 := strconv.Atoi(parts[1])
			col, err2 := strconv.Atoi(parts[2])
			
			if err1 == nil && err2 == nil {
				fromIR[irPos.IRCoordinate()] = models.Position{
					Row:   row,
					Col:   col,
					Side:  side,
					Zone:  string(irPos.Zone),
					KeyID: irPos.KeyID,
				}
			}
		}
	}
	
	return fromIR
}

// createAdv360ZoneMapping creates zone mapping from Adv360 to IR zones
func createAdv360ZoneMapping() map[string]models.IRZone {
	return map[string]models.IRZone{
		"function": models.IRZoneFunction,
		"main":     models.IRZoneMain,
		"thumb":    models.IRZoneThumb,
		"modifier": models.IRZoneModifier,
	}
}

// createZMKKeyCodeMapping creates key code mapping for ZMK keyboards
// IR is format-preserving - no normalization needed
func createZMKKeyCodeMapping() map[string]string {
	// Return empty map - preserve original ZMK syntax in IR
	// This allows compound behaviors like &kp LS(LG(S)) to pass through unchanged
	return map[string]string{}
}