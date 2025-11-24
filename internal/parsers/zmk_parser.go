package parsers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"masters3d.com/keyboard_layout_config_mapper/internal/models"
)

// ZMKParser handles parsing of ZMK keymap files
type ZMKParser struct {
	keyboardType models.KeyboardType
}

// NewZMKParser creates a new ZMK parser
func NewZMKParser(keyboardType models.KeyboardType) *ZMKParser {
	return &ZMKParser{
		keyboardType: keyboardType,
	}
}

// GetKeyboardType returns the keyboard type this parser handles
func (p *ZMKParser) GetKeyboardType() models.KeyboardType {
	return p.keyboardType
}

// Parse parses a ZMK keymap file and returns a structured representation
func (p *ZMKParser) Parse(filePath string) (*models.KeyboardLayout, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	layout := &models.KeyboardLayout{
		Type:     p.keyboardType,
		FilePath: filePath,
		Layers:   []models.Layer{},
		Behaviors: []models.Behavior{},
		Combos:   []models.Combo{},
		Metadata: make(map[string]interface{}),
	}

	scanner := bufio.NewScanner(file)
	var currentSection string
	var layerIndex int = -1
	var inBindings bool
	var bindingsBuffer []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Detect sections
		if strings.Contains(line, "behaviors {") {
			currentSection = "behaviors"
			continue
		}
		
		if strings.Contains(line, "keymap {") {
			currentSection = "keymap"
			continue
		}

		if strings.Contains(line, "combos {") {
			currentSection = "combos"
			continue
		}
		
		// Exit sections when we hit a closing brace at the beginning
		if line == "}" || line == "};" {
			if currentSection == "behaviors" {
				currentSection = ""
			}
			continue
		}

		// Parse based on current section
		switch currentSection {
		case "keymap":
			// Match layer declarations: layer0_default, layer1_qwerty, etc.
			if (strings.Contains(line, "layer") && strings.Contains(line, "{")) {
				// Starting a new layer (but not a behavior definition or bindings)
				if !strings.Contains(line, ":") && !strings.Contains(line, "behavior") && !strings.Contains(line, "bindings") {
					if inBindings && len(bindingsBuffer) > 0 {
						// Process previous layer's bindings
						if layerIndex >= 0 {
							err := p.processLayerBindings(layout, layerIndex, bindingsBuffer)
							if err != nil {
								return nil, fmt.Errorf("error processing layer %d: %v", layerIndex, err)
							}
						}
						bindingsBuffer = []string{}
					}
					
					layerIndex++
					inBindings = false
					
					// Extract layer name
					layerName := p.extractLayerName(line)
					layer := models.Layer{
						Index:    layerIndex,
						Name:     layerName,
						Bindings: []models.KeyBinding{},
						Metadata: make(map[string]interface{}),
					}
					layout.Layers = append(layout.Layers, layer)
				}
			}

			// Only process bindings if we're in a layer (layerIndex >= 0)
			if layerIndex >= 0 && strings.Contains(line, "bindings") && strings.Contains(line, "<") {
				inBindings = true
				// Extract bindings from this line if they exist
				bindingsLine := p.extractBindingsFromLine(line)
				if bindingsLine != "" {
					bindingsBuffer = append(bindingsBuffer, bindingsLine)
				}
			} else if inBindings && layerIndex >= 0 {
				// Collect bindings lines
				if strings.Contains(line, ">;") {
					// End of bindings
					bindingsLine := p.extractBindingsFromLine(line)
					if bindingsLine != "" {
						bindingsBuffer = append(bindingsBuffer, bindingsLine)
					}
					// Process the collected bindings
					if layerIndex >= 0 {
						err := p.processLayerBindings(layout, layerIndex, bindingsBuffer)
						if err != nil {
							return nil, fmt.Errorf("error processing layer %d: %v", layerIndex, err)
						}
					}
					bindingsBuffer = []string{}
					inBindings = false
				} else {
					bindingsLine := p.extractBindingsFromLine(line)
					if bindingsLine != "" {
						bindingsBuffer = append(bindingsBuffer, bindingsLine)
					}
				}
			}

		case "behaviors":
			// Parse custom behaviors
			behavior := p.parseBehavior(line)
			if behavior != nil {
				layout.Behaviors = append(layout.Behaviors, *behavior)
			}

		case "combos":
			// Parse combos
			combo := p.parseCombo(line)
			if combo != nil {
				layout.Combos = append(layout.Combos, *combo)
			}
		}
	}

	// Process any remaining bindings
	if inBindings && len(bindingsBuffer) > 0 && layerIndex >= 0 {
		err := p.processLayerBindings(layout, layerIndex, bindingsBuffer)
		if err != nil {
			return nil, fmt.Errorf("error processing final layer %d: %v", layerIndex, err)
		}
	}

	return layout, nil
}

// Validate performs syntax validation on a ZMK keymap file
func (p *ZMKParser) Validate(filePath string) error {
	// Basic file existence and readability check
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %v", err)
	}
	defer file.Close()

	// Basic ZMK syntax validation
	scanner := bufio.NewScanner(file)
	lineNum := 0
	var braceCount int

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Count braces for basic syntax checking
		braceCount += strings.Count(line, "{")
		braceCount -= strings.Count(line, "}")

		// Basic ZMK syntax checks
		if strings.Contains(line, "bindings") && !strings.Contains(line, "<") && !strings.Contains(line, "=") {
			return fmt.Errorf("line %d: bindings should be assigned with = and wrapped in < >", lineNum)
		}
	}

	if braceCount != 0 {
		return fmt.Errorf("unmatched braces in file")
	}

	return nil
}

// Helper methods for parsing

func (p *ZMKParser) extractLayerName(line string) string {
	// Extract layer name from lines like "layer0_default {" or "default_layer {"
	// First try layer0_name format
	re := regexp.MustCompile(`layer\d+_(\w+)\s*{`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	
	// Try name_layer format
	re = regexp.MustCompile(`(\w+)_layer\s*{`)
	matches = re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	
	// Extract whole name before {
	parts := strings.Split(strings.TrimSpace(line), "{")
	if len(parts) > 0 {
		name := strings.TrimSpace(parts[0])
		// Remove "layer" prefix if exists
		name = strings.TrimPrefix(name, "layer")
		// Remove leading digits and underscore
		re = regexp.MustCompile(`^\d+_`)
		name = re.ReplaceAllString(name, "")
		if name != "" {
			return name
		}
	}
	
	return "unnamed"
}

func (p *ZMKParser) extractBindingsFromLine(line string) string {
	// Extract the part between < and > or just collect the line content
	start := strings.Index(line, "<")
	end := strings.LastIndex(line, ">")
	
	if start >= 0 && end > start {
		return strings.TrimSpace(line[start+1 : end])
	}
	
	// If no brackets, return the whole line (might be a continuation)
	trimmed := strings.TrimSpace(line)
	if trimmed != "bindings" && trimmed != "=" && trimmed != ";" && trimmed != ">;" {
		return trimmed
	}
	
	return ""
}

func (p *ZMKParser) processLayerBindings(layout *models.KeyboardLayout, layerIndex int, bindingsLines []string) error {
	if layerIndex < 0 || layerIndex >= len(layout.Layers) {
		return fmt.Errorf("invalid layer index: %d", layerIndex)
	}

	// Join all bindings lines
	allBindings := strings.Join(bindingsLines, " ")
	allBindings = strings.ReplaceAll(allBindings, ",", " ")
	
	// Remove common non-binding words
	allBindings = strings.ReplaceAll(allBindings, "bindings", "")
	allBindings = strings.ReplaceAll(allBindings, "=", "")
	
	// Parse bindings as semantic units (handles &kp LS(LG(S)) correctly)
	validBindings := p.parseZMKBindings(allBindings)

	for i, binding := range validBindings {
		keyBinding := models.KeyBinding{
			Position: p.getPositionForIndex(i),
			Value:    binding,
			Layer:    layerIndex,
			Type:     p.determineBindingType(binding),
			Metadata: make(map[string]interface{}),
		}
		
		layout.Layers[layerIndex].Bindings = append(layout.Layers[layerIndex].Bindings, keyBinding)
	}

	return nil
}

// parseZMKBindings parses ZMK bindings as semantic units
// Handles compound behaviors like &kp LS(LG(S)) correctly by treating
// &behavior + arguments as a single unit instead of splitting on whitespace
func (p *ZMKParser) parseZMKBindings(text string) []string {
	var bindings []string
	text = strings.TrimSpace(text)
	
	i := 0
	for i < len(text) {
		// Skip whitespace
		for i < len(text) && (text[i] == ' ' || text[i] == '\t' || text[i] == '\n') {
			i++
		}
		if i >= len(text) {
			break
		}
		
		// Check if this is a binding (starts with &)
		if text[i] == '&' {
			start := i
			i++ // Skip the &
			
			// Read behavior name (until space or &)
			for i < len(text) && text[i] != ' ' && text[i] != '\t' && text[i] != '&' {
				i++
			}
			
			// Skip whitespace between behavior and args
			for i < len(text) && (text[i] == ' ' || text[i] == '\t') {
				// Peek ahead - if next char is &, we're done with this binding
				if i+1 < len(text) && text[i+1] == '&' {
					break
				}
				i++
			}
			
			// Read arguments (until next & or end)
			// This captures things like LS(LG(S)) or LAYER_KEYPAD as part of the binding
			for i < len(text) {
				// Stop if we hit the start of next binding
				if text[i] == '&' {
					break
				}
				// Check for space followed by & (end of this binding's args)
				if text[i] == ' ' || text[i] == '\t' {
					// Look ahead for &
					j := i + 1
					for j < len(text) && (text[j] == ' ' || text[j] == '\t') {
						j++
					}
					if j < len(text) && text[j] == '&' {
						// Next binding starts soon, trim trailing space
						break
					}
				}
				i++
			}
			
			// Build the complete binding with arguments
			binding := text[start:i]
			binding = strings.TrimSpace(binding)
			
			// Filter out syntax elements
			if binding != "" && binding != "<" && binding != ">" && binding != ";" && 
			   binding != "=>" && !strings.HasPrefix(binding, "bindings") {
				bindings = append(bindings, binding)
			}
		} else {
			// Skip non-binding characters
			i++
		}
	}
	
	return bindings
}

func (p *ZMKParser) getPositionForIndex(index int) models.Position {
	// This is a simplified position mapping - should be enhanced with actual keyboard layout
	row := index / 12 // Assuming roughly 12 keys per row
	col := index % 12
	side := "left"
	if col >= 6 {
		side = "right"
		col -= 6
	}

	zone := "main"
	if row >= 5 { // Thumb cluster typically at bottom
		zone = "thumb"
	}

	return models.Position{
		Row:   row,
		Col:   col,
		Side:  side,
		Zone:  zone,
		KeyID: fmt.Sprintf("%s_%d_%d", side, row, col),
	}
}

func (p *ZMKParser) determineBindingType(binding string) models.BindingType {
	if strings.Contains(binding, "&mt") {
		return models.BindingModTap
	}
	if strings.Contains(binding, "&lt") {
		return models.BindingLayerTap
	}
	if strings.Contains(binding, "&mo") || strings.Contains(binding, "&to") {
		return models.BindingLayerTap
	}
	if strings.Contains(binding, "&kp") {
		return models.BindingBasic
	}
	if strings.HasPrefix(binding, "&") {
		return models.BindingBehavior
	}
	
	return models.BindingBasic
}

func (p *ZMKParser) parseBehavior(line string) *models.Behavior {
	// Basic behavior parsing - can be enhanced
	if strings.Contains(line, ":") {
		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			return &models.Behavior{
				Name:       strings.TrimSpace(parts[0]),
				Type:       "custom",
				Properties: make(map[string]interface{}),
			}
		}
	}
	return nil
}

func (p *ZMKParser) parseCombo(line string) *models.Combo {
	// Basic combo parsing - can be enhanced  
	if strings.Contains(line, "combo") {
		return &models.Combo{
			Name:    "parsed_combo",
			Keys:    []models.Position{},
			Binding: line,
		}
	}
	return nil
}