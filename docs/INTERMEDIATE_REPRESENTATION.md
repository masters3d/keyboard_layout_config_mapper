# Intermediate Representation (IR) System

## Overview

The Keyboard Layout Config Mapper (KLCM) now includes a robust intermediate representation (IR) system that enables seamless translation between different keyboard layouts. This system provides a universal coordinate mapping that allows layouts to be easily ported between compatible keyboards.

## Key Features

### Universal Coordinate System
- **10x10 grid per hand**: Each hand has a standardized 10x10 coordinate system
- **Zone-based mapping**: Keys are categorized into logical zones (main, thumb, function, etc.)
- **Consistent positioning**: Similar keys across keyboards map to the same IR coordinates

### Supported Keyboards
Currently supports ZMK-based keyboards:
- **Kinesis Advantage 360** (`adv360`)
- **MoErgo Glove80** (`glove80`) 
- **Kinesis Advantage with Pillz Mod** (`adv_mod`)

### Translation Capabilities
- **Bidirectional mapping**: Convert any supported keyboard to/from IR
- **Layout preservation**: Maintains layer structure, combos, and behaviors
- **Automatic key code normalization**: Handles keyboard-specific key code differences

## Architecture

### Core Components

#### 1. IR Models (`internal/models/ir.go`)
- `IRPosition`: Universal coordinate system (hand, row, col, zone)
- `IRLayout`: Complete layout in intermediate representation
- `IRGrid`: 10x10 grid structure for each hand
- `IRKeyBinding`: Individual key mappings in IR space

#### 2. Mappers (`internal/mappers/`)
- `KeyboardMapper`: Interface for keyboard-specific mapping
- `PositionMapper`: Handles coordinate transformations
- `LayoutTranslator`: High-level translation between keyboards

#### 3. Keyboard-Specific Mappers
- `ZMKAdv360Mapper`: Advantage 360 ↔ IR mapping
- `ZMKGlove80Mapper`: Glove80 ↔ IR mapping  
- `ZMKAdvModMapper`: Advanced Mod ↔ IR mapping

## Coordinate System

### IR Grid Layout (per hand)
```
    0   1   2   3   4   5   6   7   8   9
0 [ F1] [F2] [F3] [F4] [F5] [F6] [ ] [ ] [ ] [ ]  (Function Row)
1 [ 1] [ 2] [ 3] [ 4] [ 5] [ 6] [ ] [ ] [ ] [ ]   (Number Row)  
2 [ Q] [ W] [ E] [ R] [ T] [ ] [ ] [ ] [ ] [ ]     (Top Row)
3 [ A] [ S] [ D] [ F] [ G] [ ] [ ] [ ] [ ] [ ]     (Home Row)
4 [ Z] [ X] [ C] [ V] [ B] [ ] [ ] [ ] [ ] [ ]     (Bottom Row)
5 [TH1][TH2][TH3][TH4][ ] [ ] [ ] [ ] [ ] [ ]     (Thumb Cluster)
6 [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]           (Reserved)
7 [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]           (Reserved)
8 [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]           (Reserved)
9 [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ] [ ]           (Reserved)
```

### Zone Classifications
- **`main`**: Primary typing area (QWERTY region)
- **`function`**: Function row (F1-F12)
- **`thumb`**: Thumb cluster keys
- **`modifier`**: Dedicated modifier key areas
- **`nav`**: Navigation keys (arrows, home, end)
- **`numpad`**: Number pad area

## Usage Examples

### Basic Translation
```bash
# Translate from Advantage360 to Glove80
klcm translate --from adv360 --to glove80

# Translate from Glove80 to Advanced Mod
klcm translate --from glove80 --to adv_mod --output my_adv_mod.keymap
```

### Intermediate Representation Analysis
```bash
# View the IR representation of a layout
klcm translate --from adv360 --show-ir

# Save IR to file for analysis
klcm translate --from glove80 --show-ir --output glove80_ir.json
```

## Adding New Keyboards

### Step 1: Create Position Mapping
Define the mapping between your keyboard's matrix positions and IR coordinates:

```go
func createMyKeyboardToIRMapping() map[string]models.IRPosition {
    mapping := make(map[string]models.IRPosition)
    
    // Map each physical key position to IR coordinate
    mapping["left_0_0"] = models.IRPosition{
        Hand: "left", Row: 2, Col: 0, Zone: "main", KeyID: "q"
    }
    // ... more mappings
    
    return mapping
}
```

### Step 2: Implement KeyboardMapper Interface
```go
type MyKeyboardMapper struct {
    positionMapper *PositionMapper
}

func (m *MyKeyboardMapper) ToIR(layout *models.KeyboardLayout) (*models.IRLayout, error) {
    // Convert keyboard-specific layout to IR
}

func (m *MyKeyboardMapper) FromIR(irLayout *models.IRLayout) (*models.KeyboardLayout, error) {
    // Convert IR back to keyboard-specific layout
}
```

### Step 3: Register with Translator
```go
translator := mappers.NewLayoutTranslator()
translator.RegisterMapper(NewMyKeyboardMapper())
```

## Key Code Normalization

The IR system normalizes key codes to standard identifiers:

### ZMK → IR Examples
- `&kp A` → `a`
- `&kp SPACE` → `space`
- `&kp ENTER` → `enter`
- `&mt LSHIFT A` → `mt(lshift, a)` (future enhancement)

### Benefits
- **Cross-firmware compatibility**: Same logical keys work across different firmware
- **Simplified mapping**: Focus on key function rather than syntax
- **Future extensibility**: Easy to add new keyboard types and firmware

## Technical Details

### Position Mapping Algorithm
1. **Parse source layout**: Extract all key bindings and positions
2. **Translate to IR**: Map each keyboard position to IR coordinate
3. **Normalize key codes**: Convert keyboard-specific codes to IR format
4. **Generate IR layout**: Create universal representation
5. **Translate from IR**: Map IR coordinates to target keyboard positions
6. **Generate target layout**: Create keyboard-specific output

### Mapping Validation
The system validates that mappings are:
- **Bidirectional**: Every keyboard→IR mapping has a corresponding IR→keyboard mapping
- **Consistent**: Round-trip translations preserve key assignments
- **Complete**: All essential keys are mapped

### Error Handling
- **Unmappable positions**: Keys that don't exist on target keyboard are skipped
- **Zone conflicts**: Graceful handling of zone mismatches
- **Key code fallbacks**: Unknown key codes pass through unchanged

## Future Enhancements

### Planned Features
- **QMK keyboard support**: Extend to QMK-based keyboards
- **Advanced behavior mapping**: Support for complex mod-tap and layer-tap behaviors
- **Layout optimization**: Suggest layout improvements based on usage patterns
- **Visual mapping**: Generate visual representations of key mappings
- **Conflict resolution**: Interactive handling of mapping conflicts

### Extension Points
- **Custom mappers**: Plugin system for user-defined keyboard mappers
- **Transformation rules**: Define custom key transformation logic
- **Layout templates**: Pre-defined layout templates for common configurations
- **Migration tools**: Automated migration from old layouts to new keyboards

## Validation

The system includes comprehensive validation:

```bash
# Test the translation system
go test ./internal/mappers/... -v

# Validate specific keyboard mappings
klcm translate --from adv360 --to glove80 --validate

# Check mapping consistency
klcm validate-mappings
```

---

This intermediate representation system makes KLCM a powerful tool for keyboard enthusiasts who want to maintain consistent layouts across multiple keyboards or easily migrate to new hardware while preserving their carefully crafted key bindings.