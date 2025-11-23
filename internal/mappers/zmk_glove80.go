package mappers

import (
	"strconv"
	"strings"
	"masters3d.com/keyboard_layout_config_mapper/internal/models"
)

// ZMKGlove80Mapper handles mapping between Glove80 layout and IR
// Since Glove80 and Adv360 are basically the same layout, this reuses most of the Adv360 mapping
type ZMKGlove80Mapper struct {
	positionMapper *PositionMapper
}

// NewZMKGlove80Mapper creates a new mapper for Glove80 keyboards
func NewZMKGlove80Mapper() *ZMKGlove80Mapper {
	mapper := &ZMKGlove80Mapper{
		positionMapper: &PositionMapper{
			KeyboardType:   models.KeyboardZMKGlove80,
			ToIRMapping:    createGlove80ToIRMapping(),
			FromIRMapping:  createGlove80FromIRMapping(),
			ZoneMapping:    createAdv360ZoneMapping(), // Reuse same zone mapping
			KeyCodeMapping: createZMKKeyCodeMapping(),  // Reuse same key code mapping
		},
	}
	return mapper
}

// ToIR converts Glove80 layout to intermediate representation
func (m *ZMKGlove80Mapper) ToIR(layout *models.KeyboardLayout) (*models.IRLayout, error) {
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

// FromIR converts intermediate representation to Glove80 layout
func (m *ZMKGlove80Mapper) FromIR(irLayout *models.IRLayout) (*models.KeyboardLayout, error) {
	layout := &models.KeyboardLayout{
		Type:      models.KeyboardZMKGlove80,
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
func (m *ZMKGlove80Mapper) GetKeyboardType() models.KeyboardType {
	return models.KeyboardZMKGlove80
}

// GetPositionMapping returns the position mapping
func (m *ZMKGlove80Mapper) GetPositionMapping() map[string]models.IRPosition {
	return m.positionMapper.ToIRMapping
}

// GetReverseMapping returns the reverse mapping
func (m *ZMKGlove80Mapper) GetReverseMapping() map[string]models.Position {
	return m.positionMapper.FromIRMapping
}

// createGlove80ToIRMapping creates the position mapping from Glove80 to IR coordinates
// The Glove80 has a very similar layout to the Adv360, so we reuse most of the mapping
func createGlove80ToIRMapping() map[string]models.IRPosition {
	// Start with the Adv360 mapping as a base since they're essentially the same layout
	mapping := createAdv360ToIRMapping()
	
	// Glove80-specific adjustments (if any)
	// The Glove80 might have slightly different matrix positions but similar physical layout
	
	// Add any Glove80-specific keys or adjustments here
	// For now, we'll use the same mapping as Adv360
	
	return mapping
}

// createGlove80FromIRMapping creates the reverse mapping from IR to Glove80 coordinates
func createGlove80FromIRMapping() map[string]models.Position {
	toIR := createGlove80ToIRMapping()
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