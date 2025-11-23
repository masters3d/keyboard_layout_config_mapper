package mappers

import (
	"strconv"
	"strings"
	"masters3d.com/keyboard_layout_config_mapper/internal/models"
)

// ZMKAdvModMapper handles mapping between Advanced Mod (Pillz Mod) layout and IR
type ZMKAdvModMapper struct {
	positionMapper *PositionMapper
}

// NewZMKAdvModMapper creates a new mapper for Advanced Mod keyboards
func NewZMKAdvModMapper() *ZMKAdvModMapper {
	mapper := &ZMKAdvModMapper{
		positionMapper: &PositionMapper{
			KeyboardType:   models.KeyboardZMKAdvMod,
			ToIRMapping:    createAdvModToIRMapping(),
			FromIRMapping:  createAdvModFromIRMapping(),
			ZoneMapping:    createAdv360ZoneMapping(), // Reuse same zone mapping
			KeyCodeMapping: createZMKKeyCodeMapping(),  // Reuse same key code mapping
		},
	}
	return mapper
}

// ToIR converts Advanced Mod layout to intermediate representation
func (m *ZMKAdvModMapper) ToIR(layout *models.KeyboardLayout) (*models.IRLayout, error) {
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

// FromIR converts intermediate representation to Advanced Mod layout
func (m *ZMKAdvModMapper) FromIR(irLayout *models.IRLayout) (*models.KeyboardLayout, error) {
	layout := &models.KeyboardLayout{
		Type:      models.KeyboardZMKAdvMod,
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
func (m *ZMKAdvModMapper) GetKeyboardType() models.KeyboardType {
	return models.KeyboardZMKAdvMod
}

// GetPositionMapping returns the position mapping
func (m *ZMKAdvModMapper) GetPositionMapping() map[string]models.IRPosition {
	return m.positionMapper.ToIRMapping
}

// GetReverseMapping returns the reverse mapping
func (m *ZMKAdvModMapper) GetReverseMapping() map[string]models.Position {
	return m.positionMapper.FromIRMapping
}

// createAdvModToIRMapping creates the position mapping from Advanced Mod to IR coordinates
// The Advanced Mod is based on the Kinesis Advantage, so it should be very similar to Adv360
func createAdvModToIRMapping() map[string]models.IRPosition {
	// Start with the Adv360 mapping as a base since it's also a Kinesis Advantage-based layout
	mapping := createAdv360ToIRMapping()
	
	// Advanced Mod specific adjustments
	// The Pillz Mod maintains the same basic Kinesis Advantage layout but might have different matrix positions
	
	// Override any specific differences if needed
	// For now, we'll assume it's essentially the same as the Adv360 layout
	
	return mapping
}

// createAdvModFromIRMapping creates the reverse mapping from IR to Advanced Mod coordinates
func createAdvModFromIRMapping() map[string]models.Position {
	toIR := createAdvModToIRMapping()
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