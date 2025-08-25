package parsers

import (
	"fmt"
	"os"
	"path/filepath"

	"masters3d.com/keyboard_layout_config_mapper/internal/models"
)

// Parser interface defines methods that all keyboard parsers must implement
type Parser interface {
	Parse(filePath string) (*models.KeyboardLayout, error)
	Validate(filePath string) error
	GetKeyboardType() models.KeyboardType
}

// NewParser creates a parser based on the keyboard type
func NewParser(keyboardType models.KeyboardType) (Parser, error) {
	switch keyboardType {
	case models.KeyboardZMKAdv360:
		return NewZMKParser(keyboardType), nil
	case models.KeyboardZMKGlove80:
		return NewZMKParser(keyboardType), nil
	case models.KeyboardQMKErgoDox:
		return NewQMKParser(), nil
	case models.KeyboardKinesis2:
		return NewKinesis2Parser(), nil
	default:
		return nil, fmt.Errorf("unsupported keyboard type: %s", keyboardType)
	}
}

// GetConfigPath returns the configuration file path for a keyboard type
func GetConfigPath(keyboardType models.KeyboardType) (string, error) {
	configsDir := "configs"
	
	switch keyboardType {
	case models.KeyboardZMKAdv360:
		return filepath.Join(configsDir, "zmk_adv360", "adv360.keymap"), nil
	case models.KeyboardZMKGlove80:
		return filepath.Join(configsDir, "zmk_glove80", "glove80.keymap"), nil
	case models.KeyboardQMKErgoDox:
		return filepath.Join(configsDir, "qmk_ergodox", "keymap.c"), nil
	case models.KeyboardKinesis2:
		return filepath.Join(configsDir, "kinesis2", "1_qwerty.txt"), nil
	default:
		return "", fmt.Errorf("unsupported keyboard type: %s", keyboardType)
	}
}

// Validator provides validation functionality for keyboard configurations
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateAll validates all keyboard configurations
func (v *Validator) ValidateAll(compileCheck bool) error {
	keyboards := []models.KeyboardType{
		models.KeyboardZMKAdv360,
		models.KeyboardZMKGlove80,
		models.KeyboardQMKErgoDox,
		models.KeyboardKinesis2,
	}

	var errors []string
	
	for _, keyboard := range keyboards {
		err := v.ValidateKeyboard(string(keyboard), compileCheck)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", keyboard, err))
		} else {
			fmt.Printf("✅ %s configuration is valid\n", keyboard)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors:\n%v", errors)
	}

	fmt.Println("🎉 All keyboard configurations are valid!")
	return nil
}

// ValidateKeyboard validates a specific keyboard configuration
func (v *Validator) ValidateKeyboard(keyboard string, compileCheck bool) error {
	keyboardType := models.KeyboardType(keyboard)
	
	configPath, err := GetConfigPath(keyboardType)
	if err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Create parser and validate
	parser, err := NewParser(keyboardType)
	if err != nil {
		return err
	}

	err = parser.Validate(configPath)
	if err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	if compileCheck {
		return v.validateCompilation(keyboardType, configPath)
	}

	return nil
}

// validateCompilation attempts to validate that the configuration can be compiled
func (v *Validator) validateCompilation(keyboardType models.KeyboardType, configPath string) error {
	switch keyboardType {
	case models.KeyboardZMKAdv360, models.KeyboardZMKGlove80:
		// For ZMK, we could potentially set up a docker container to test compilation
		// For now, just do syntax validation
		fmt.Printf("⚠️  Compilation check not yet implemented for %s\n", keyboardType)
		return nil
		
	case models.KeyboardQMKErgoDox:
		// For QMK, we could potentially validate with qmk compile
		fmt.Printf("⚠️  Compilation check not yet implemented for %s\n", keyboardType)
		return nil
		
	case models.KeyboardKinesis2:
		// Kinesis2 is just text mapping, no compilation needed
		return nil
		
	default:
		return fmt.Errorf("compilation check not supported for %s", keyboardType)
	}
}