# ZMK Configuration Deep Analysis

## Configuration Overview

### **Advantage360 (zmk_adv360/adv360.keymap)**
- **Format**: ZMK keymap (.keymap)
- **Layers**: 9 layers (0-8)
- **Physical Layout**: Split ergonomic with thumb clusters
- **Key Features**:
  - Complex behavior definitions (mo_key, hm homerow mods)
  - Custom mod-morph behaviors (morph_dot, morph_comma, etc.)
  - Macros for bracket pairs
  - Layer structure: default, qwerty, padding layers, keypad, fn, mod

### **Glove80 (zmk_glove80/glove80.keymap)**  
- **Format**: ZMK keymap (.keymap)
- **Layers**: 9 layers (0-8)
- **Physical Layout**: Curved ergonomic split with extensive thumb clusters
- **Key Features**:
  - Similar behavior patterns to Advantage360
  - Same custom mod-morph behaviors
  - Additional macro definitions (quotes, brackets)
  - Bluetooth management macros
  - RGB and backlight controls

## Shared ZMK Patterns

### **Common Behaviors**:
1. **mod-morph behaviors**: Same definitions across both keyboards
   - `morph_dot`: Period → Colon with shift
   - `morph_comma`: Comma → Semicolon with shift
   - `morph_parens_left/right`: Parentheses → Angle brackets with shift
   - `morph_exclamation`: Backslash → Exclamation with shift
   - `morph_quote_single/double`: Quotes → Grave/Tilde with shift

2. **Hold-tap behaviors**:
   - `mo_key`: Layer hold, key tap
   - `hr_mod/hm`: Homerow modifier implementation
   - Different timing configurations

3. **Macros**:
   - Bracket pair macros (type opening+closing, cursor between)
   - Bluetooth selection macros

### **Key Differences**:

| Feature | Advantage360 | Glove80 |
|---------|-------------|---------|
| Physical Keys | ~76 keys | ~80 keys |
| Thumb Cluster | 6 keys per side | 6 keys per side |
| Combo Support | Basic | More extensive |
| Magic Layer | Basic RGB/BT | Advanced RGB/BT controls |
| Default Layout | Custom layout | Similar but adapted to Glove80 |

## Layer Structure Mapping

### **Layer 0 (Default)**:
- Both use similar QWERTY-based layouts
- Home row modifications
- Custom symbol positioning
- Thumb cluster for modifiers/layers

### **Layer 1 (QWERTY)**:
- Standard QWERTY layout
- Fallback for compatibility

### **Layers 2-5 (Padding/Symbols)**:
- Symbol and function access
- Number pad simulation
- Special character input

### **Layer 6 (Keypad)**:
- Number pad with symbols
- Function key access
- Macro shortcuts

### **Layer 7 (Function)**:
- Function keys with modifiers
- System controls
- Application shortcuts

### **Layer 8 (Magic/System)**:
- Bluetooth management
- RGB/backlight controls
- Bootloader access
- System functions

## Synchronization Compatibility

### **High Compatibility (90-95%)**:
- Basic key mappings (letters, numbers, common symbols)
- Layer structure and logic
- Modifier key assignments
- Most custom behaviors

### **Medium Compatibility (70-80%)**:
- Physical layout differences require key position mapping
- Thumb cluster assignments need adaptation
- Some combos may not translate directly

### **Low Compatibility (Manual)**:
- Hardware-specific features (RGB patterns, specific GPIO)
- Bootloader/reset sequences
- Hardware-dependent macros

## Recommended Sync Strategy

1. **Primary Reference**: Use Advantage360 as primary (user's main keyboard)
2. **Auto-sync**: Basic layout, behaviors, macros
3. **Manual review**: Physical position mappings, hardware features
4. **Validation**: Test compilation and functionality

## Technical Implementation Notes

### **Parser Requirements**:
- ZMK DTS syntax parsing
- Behavior extraction and mapping
- Layer structure analysis
- Keycode translation

### **Key Position Mapping**:
- Create physical layout matrices for both keyboards
- Map logical positions (home row, number row, etc.)
- Handle thumb cluster differences

### **Behavior Translation**:
- Most behaviors can be copied directly
- Timing parameters may need adjustment
- Hardware-specific behaviors require adaptation