# KLCM Agents Memory - Keyboard Integration Guide

## ZMK Advanced Mod (zmk_adv_mod) Integration

### Project Overview
**Objective**: Add support for the Pillz Mod keyboards (3 variants) using Nice!Nano with ZMK firmware, following the same build automation pattern as Adv360 and Glove80.

**Key Requirements**:
- Support Nice!Nano controller (nRF52840 wireless)
- Use dcpedit/pillzmod as the base hardware design
- Create zmk_adv_mod configuration following existing KLCM patterns
- Set up GitHub Actions for automated firmware builds
- Support 3 different keyboards using the same mod
- Follow exact same pattern as other ZMK keyboards in the system

### Hardware Analysis - Pillz Mod

**Source Repository**: https://github.com/dcpedit/pillzmod
**ZMK Shield**: `pillzmod_pro` (Pro-Micro footprint, supports Nice!Nano)
**ZMK Source Fork**: https://github.com/dcpedit/zmk/tree/pillzmod/app/boards/shields/pillzmod_pro

#### Key Features:
- **Pro PCB**: Supports Pro-Micro footprint (Elite-C, Elite-Pi, Nice!Nano)
- **Nice!Nano Support**: Full ZMK wireless support with 3.7V battery
- **Status LEDs**: 4 LEDs (1 fewer than other Kinesis mods, but has power LED option)
- **Existing Firmware**: `pro_zmk.uf2` and `pro_zmk_studio.uf2` available
- **Mill-Max 310 Sockets**: Required for Nice!Nano (thin pins for middle 3 pins)

#### Hardware Components for Nice!Nano Build:
- Nice!Nano development board
- 3.7V battery (2000mAh recommended)
- Power button (optional: 10mm or 19mm)
- Mill-Max 310 sockets (thin pins required)
- Status LEDs (4x 1.8mm LEDs + resistors)
- USB-C panel mount connector

### ZMK Shield Structure Analysis

**Files in `pillzmod_pro` shield**:
```
app/boards/shields/pillzmod_pro/
├── CMakeLists.txt              # Build configuration
├── Kconfig.defconfig           # Default configuration
├── Kconfig.shield             # Shield definition
├── leds.c                     # LED control code
├── pillzmod_pro-layouts.dtsi  # Physical layout definitions
├── pillzmod_pro.conf          # Shield configuration
├── pillzmod_pro.keymap        # Default keymap
├── pillzmod_pro.overlay       # Device tree overlay
└── pillzmod_pro.zmk.yml       # ZMK metadata
```

**Key Shield Configuration**:
- **Shield Name**: `SHIELD_PILLZMOD_PRO`
- **Compatible**: Pro-Micro footprint boards
- **Matrix**: Standard Kinesis Advantage matrix
- **Features**: Backlight, Bluetooth, wireless support

### KLCM Integration Plan

#### Phase 1: Core Configuration Setup
1. **Add zmk_adv_mod to KeyboardType enum** in `internal/models/types.go`
2. **Create configuration directory**: `configs/zmk_adv_mod/`
3. **Add parser support** in `internal/parsers/parser.go`
4. **Update all CLI commands** (pull, sync, download) to support zmk_adv_mod
5. **Create base keymap file**: `configs/zmk_adv_mod/adv_mod.keymap`

#### Phase 2: ZMK Config Repository Creation
1. **Create new repository**: `zmk-adv-mod-config` (following existing patterns)
2. **Set up repository structure**:
   ```
   zmk-adv-mod-config/
   ├── .github/
   │   └── workflows/
   │       └── build.yml          # GitHub Actions for firmware builds
   ├── config/
   │   ├── adv_mod.conf          # Board configuration
   │   └── adv_mod.keymap        # Main keymap file
   ├── build.yaml                # ZMK build configuration
   └── README.md                 # Documentation
   ```

#### Phase 3: GitHub Actions Build Setup
1. **Configure ZMK build matrix** for Nice!Nano + pillzmod_pro shield
2. **Set up artifact generation** (`.uf2` files)
3. **Enable build triggers** (push, PR, manual dispatch)
4. **Configure build cache** for faster builds

#### Phase 4: Testing & Validation
1. **Verify firmware builds successfully**
2. **Test keymap changes propagation**
3. **Validate KLCM tool integration**
4. **Document flashing process**

### Implementation Steps

#### Step 1: Update KLCM Core Types
```go
// File: internal/models/types.go
const (
    KeyboardZMKAdv360   KeyboardType = "adv360"
    KeyboardZMKGlove80  KeyboardType = "glove80" 
    KeyboardZMKAdvMod   KeyboardType = "adv_mod"    // NEW
    KeyboardQMKErgoDox  KeyboardType = "qmk_ergodox"
    KeyboardKinesis2    KeyboardType = "kinesis2"
)
```

#### Step 2: Update Parser Configuration
```go
// File: internal/parsers/parser.go
func GetConfigPath(keyboardType models.KeyboardType) (string, error) {
    switch keyboardType {
    case models.KeyboardZMKAdv360:
        return filepath.Join(configsDir, "zmk_adv360", "adv360.keymap"), nil
    case models.KeyboardZMKGlove80:
        return filepath.Join(configsDir, "zmk_glove80", "glove80.keymap"), nil
    case models.KeyboardZMKAdvMod:                                            // NEW
        return filepath.Join(configsDir, "zmk_adv_mod", "adv_mod.keymap"), nil // NEW
    // ... existing cases
}

func NewParser(keyboardType models.KeyboardType) (Parser, error) {
    switch keyboardType {
    case models.KeyboardZMKAdv360:
        return NewZMKParser(keyboardType), nil
    case models.KeyboardZMKGlove80:
        return NewZMKParser(keyboardType), nil
    case models.KeyboardZMKAdvMod:        // NEW
        return NewZMKParser(keyboardType), nil // NEW
    // ... existing cases
}
```

#### Step 3: Update CLI Commands
Files to update:
- `internal/cli/pull.go` - Add zmk_adv_mod URL mapping
- `internal/cli/sync.go` - Add zmk_adv_mod sync configuration  
- `internal/cli/download.go` - Add zmk_adv_mod download configuration
- `internal/parsers/parser.go` - Add validation support

#### Step 4: Create Base Keymap Configuration
Create `configs/zmk_adv_mod/adv_mod.keymap` based on pillzmod_pro default keymap but adapted for KLCM structure.

#### Step 5: Repository Setup
1. Create `zmk-adv-mod-config` repository
2. Copy GitHub Actions workflow from existing ZMK configs (Adv360/Glove80)
3. Adapt `build.yaml` for pillzmod_pro shield + nice_nano board
4. Set up proper shield/board combinations

### ZMK Build Configuration

**Target Configuration**:
- **Board**: `nice_nano_v2` (or `nice_nano`)
- **Shield**: `pillzmod_pro`
- **ZMK Fork**: Use official ZMK with pillzmod_pro shield integrated, or dcpedit fork

**Build Matrix Example**:
```yaml
strategy:
  matrix:
    board: [nice_nano_v2]
    shield: [pillzmod_pro]
```

**Build Command**:
```bash
west build -p -d build/adv_mod -b nice_nano_v2 -- -DSHIELD=pillzmod_pro
```

### Key Considerations

#### 1. ZMK Shield Integration
**Question**: Does the official ZMK repository include pillzmod_pro shield, or do we need to use dcpedit's fork?

**Research Needed**:
- Check if pillzmod_pro exists in upstream ZMK
- If not, determine how to integrate dcpedit's shield
- Consider submitting upstream PR or using fork

#### 2. Keymap Compatibility  
**Challenge**: Ensure pillzmod_pro keymap structure matches KLCM expectations

**Solution Approach**:
- Analyze existing adv360.keymap and glove80.keymap structures
- Map pillzmod_pro matrix to KLCM key position expectations
- Create converter/adapter if needed

#### 3. Build Automation Reliability
**Requirement**: GitHub Actions must build successfully on every commit

**Validation Steps**:
- Test build with default keymap
- Test build with custom keymap changes
- Verify artifact generation (`.uf2` files)
- Test firmware flashing process

#### 4. Multi-Keyboard Support
**Requirement**: Support 3 different keyboards with same mod

**Options**:
1. **Single repository with variants**: Use build matrix with different configurations
2. **Separate repositories**: Create individual repos for each keyboard variant
3. **Configuration-driven**: Use single shield with configuration parameters

**Recommended**: Single repository with build variants (most maintainable)

### Success Criteria

#### Technical Validation
- [ ] KLCM tool recognizes `adv_mod` keyboard type
- [ ] `klcm generate adv_mod` creates valid keymap files
- [ ] `klcm validate adv_mod --compile` passes successfully
- [ ] GitHub Actions builds firmware without errors
- [ ] Generated `.uf2` files are valid and flashable

#### Functional Validation
- [ ] Firmware flashes successfully to Nice!Nano
- [ ] All keys function as expected
- [ ] Wireless connectivity works properly
- [ ] Battery management functions correctly
- [ ] Status LEDs work as expected

#### Integration Validation
- [ ] PR creation tool works with zmk-adv-mod-config repository
- [ ] Keymap changes propagate from KLCM to ZMK config automatically
- [ ] Build artifacts are generated and downloadable
- [ ] Documentation is complete and accurate

### Troubleshooting Guide

#### Common Issues

**1. Shield Not Found Error**
```
Error: Shield 'pillzmod_pro' not found
```
**Solution**: Verify shield is available in ZMK tree or fork

**2. Build Matrix Issues**
```
Error: No matching board/shield combination
```
**Solution**: Check board/shield compatibility in ZMK documentation

**3. Keymap Syntax Errors** 
```
Error: Invalid keymap binding
```
**Solution**: Validate keymap syntax against ZMK keymap documentation

**4. Nice!Nano Pin Mapping Issues**
```
Error: Pin assignment conflict
```
**Solution**: Verify pillzmod_pro overlay matches Nice!Nano pinout

### References & Resources

#### Documentation
- [ZMK Documentation](https://zmk.dev/)
- [Nice!Nano Documentation](https://nicekeyboards.com/docs/)
- [Pillz Mod Documentation](https://github.com/dcpedit/pillzmod/blob/main/README.MD)

#### Repositories
- **Hardware Design**: https://github.com/dcpedit/pillzmod
- **ZMK Fork**: https://github.com/dcpedit/zmk/tree/pillzmod
- **Example Config**: https://github.com/keepitsimplejim/zmk-config-pillzmod-kinesis-adv

#### Hardware Suppliers
- **Nice!Nano**: https://nicekeyboards.com/nice-nano/
- **Batteries**: 2000mAh 3.7V LiPo recommended
- **Mill-Max Sockets**: 310 series (thin pins for Nice!Nano)

### Next Actions

**PHASE 1: ✅ COMPLETED - KLCM Core Integration**
- [x] Create zmk_adv_mod configuration in KLCM ✅
- [x] Update all parser and CLI components ✅ 
- [x] Create base keymap configuration ✅
- [x] Test KLCM tool recognition and validation ✅

**PHASE 1.5: ✅ COMPLETED - Configuration Sync**
- [x] Sync QWERTY layer from adv360 to adv_mod ✅
- [x] Sync keypad layer from adv360 to adv_mod ✅
- [x] Sync CMD/FN layer from adv360 to adv_mod ✅
- [x] Add macros (brackets, braces, parens, angle brackets) ✅
- [x] Add mod-morph behaviors (dot→colon, comma→semicolon, etc.) ✅
- [x] Preserve hardware-specific top row (F1-F12) ✅
- [x] Preserve bootloader layer (system_layer) ✅

#### Configuration Sync Summary

**Synced from adv360 to adv_mod:**

1. **Default Layer Changes:**
   - Number row replaced with morphed special characters
   - QWERTY layout with mod-morph keys (. → :, , → ;)
   - Parentheses morph to angle brackets when shifted
   - Quote keys with grave/tilde morphs
   - Control key positions synced (HOME, BSPC, LC(BSPC), DEL, ENTER, TAB)

2. **Keypad Layer (LAYER_KEYPAD = 1):**
   - F13-F24 function keys on left side
   - Number pad 1-9, 0 on right side
   - Symbols: %, $, #, @ (left); ^, &, *, |, dot, comma (right)
   - Macro shortcuts for brackets/braces in top right

3. **CMD Layer (LAYER_CMD = 2) - NEW:**
   - Control key combinations for all QWERTY keys
   - Uses RC() (Right Control) instead of RG() (Right GUI/Command)
   - Matches glove80 implementation for cross-keyboard consistency

4. **System Layer (LAYER_SYSTEM = 3):**
   - **PRESERVED** - bootloader access is hardware-specific
   - Bluetooth configuration unchanged
   - Critical for firmware updates and system functions

5. **Behaviors Added:**
   - `morph_dot` - Period → Colon when shifted
   - `morph_comma` - Comma → Semicolon when shifted
   - `morph_parens_left/right` - Parentheses → Angle brackets when shifted
   - `morph_exclamation` - Backslash → Exclamation when shifted
   - `morph_quote_single/double` - Quote morphing with grave/tilde

6. **Macros Added:**
   - `macro_brackets` - Types [] with cursor in middle
   - `macro_braces` - Types {} with cursor in middle
   - `macro_parens` - Types () with cursor in middle
   - `macro_angle_brackets` - Types <> with cursor in middle

**Hardware-Specific Preservation:**
- Top function key row (18 keys: HOME, F1-F12, PSCRN, SLCK, PAUSE, Layer toggle, System) kept as-is
- System layer bootloader access unchanged (critical for Nice!Nano firmware updates)

**PHASE 2: ✅ COMPLETED - Branch References Updated**
- [x] Update v7_target → main in pull.go ✅
- [x] Update v7_target → main in download.go ✅

**PHASE 3: ✅ COMPLETED - Intermediate Representation System**
- [x] Create IR models and universal coordinate system ✅
- [x] Implement KeyboardMapper interface ✅
- [x] Create ZMK mappers for adv360, glove80, and adv_mod ✅
- [x] Add translate CLI command ✅
- [x] Create comprehensive documentation ✅

**PHASE 4: ✅ COMPLETED - ZMK Config Repository Setup**
- [x] Repository created: `masters3d/zmk-config-pillzmod-nicenano` ✅
- [x] GitHub Actions build workflow configured ✅
- [x] Uses official dcpedit `pillzmod_pro` shield ✅
- [x] Firmware builds successfully (.uf2 artifacts generated) ✅
- [x] KLCM download.go points to correct repository/branch ✅

**Repository Details:**
- **URL**: https://github.com/masters3d/zmk-config-pillzmod-nicenano
- **Branch**: `cheyo` (production branch)
- **Shield**: `pillzmod_pro` (official dcpedit shield)
- **Board**: `nice_nano_v2`

**PHASE 5: 🚧 IN PROGRESS - Keymap Sync to Remote**

The local KLCM source of truth (`configs/zmk_adv_mod/adv_mod.keymap`) has been enhanced with:
- [x] Macros (brackets, braces, parens, angle_brackets)
- [x] Mod-morph behaviors (dot→colon, comma→semicolon, etc.)
- [x] CMD layer with RC() control key combinations
- [x] Synced default layer from adv360

**Remaining Task:**
- [ ] **Push updated keymap to remote repository** - The local KLCM keymap needs to be synced to `masters3d/zmk-config-pillzmod-nicenano` cheyo branch

To sync, use the PR command or manual push:
```bash
# Option 1: Use KLCM PR workflow
klcm pr adv_mod

# Option 2: Manual sync
gh repo clone masters3d/zmk-config-pillzmod-nicenano /tmp/zmk-adv-mod
cd /tmp/zmk-adv-mod
git checkout cheyo
cp /path/to/keyboard_layout_config_mapper/configs/zmk_adv_mod/adv_mod.keymap config/adv_mod.keymap
git add config/adv_mod.keymap
git commit -m "Sync keymap with KLCM source of truth - Add macros and mod-morphs"
git push origin cheyo
```

**PHASE 6: 📋 PLANNED - Production Validation**
- [ ] Flash firmware to actual hardware
- [ ] Test all keys and matrix positions  
- [ ] Verify wireless connectivity
- [ ] Test battery management
- [ ] Validate status LEDs
- [ ] Document flashing and usage process

**PHASE 7: 📋 PLANNED - Documentation & Maintenance**
- [ ] Complete hardware assembly documentation
- [ ] Create troubleshooting guides
- [ ] Set up automated testing workflow
- [ ] Create release management process

### Critical Dependencies

#### Hardware Requirements
- **Pillz Mod Pro PCB** with Nice!Nano socket installation
- **Nice!Nano v2** controller board
- **Mill-Max 310 sockets** (thin pins required)
- **3.7V LiPo battery** (2000mAh recommended)
- **Status LEDs and resistors** (per Pillz Mod BOM)

#### Software Dependencies
- **ZMK main branch** (no fork required with custom shield approach)
- **GitHub Actions** for automated builds
- **KLCM tool** for keymap management

#### Unknown/To Be Verified
1. **Exact pin mapping** for Nice!Nano on Pillz Mod Pro PCB
2. **LED pin assignments** for status indicators
3. **Matrix scan reliability** with Nice!Nano timing
4. **Power management** configuration specifics

---

*Last Updated: 2025-11-28*  
*Version: v8_mapping_layer*

## Mapping Layer Translation System

### Current State (v8_mapping_layer)

The translation system now uses **mapping layers** embedded in each keyboard's keymap file. Each key's binding VALUE in the mapping layer defines its logical identity (KeyID), enabling accurate translation between keyboards.

### Mapping Layer Structure

Each keyboard has a `layer2_mapping` or `mapping_layer` that defines key identities:

```zmk
layer2_mapping {
  bindings = <
    &kp GRAVE  &kp N1    &kp N2    ...  // Number row
    &kp TAB    &kp Q     &kp W     ...  // QWERTY row
    &kp ESC    &kp A     &kp S     ...  // Home row
    &kp LSHFT  &kp Z     &kp X     ...  // Bottom row
    &kp KP_N3  &kp KP_N4 &kp KP_N5 ...  // Thumb keys (use KP_N* for unique IDs)
  >;
};
```

### Thumb Key Identifiers

Thumb keys use keypad numbers as unique identifiers (since they don't conflict with main keyboard):

| KeyID | Function | Notes |
|-------|----------|-------|
| KP_N1 | LAYER_KEYPAD (left) | Inner thumb |
| KP_N2 | WIN/GUI (left) | |
| KP_N3 | SPACE (left) | Outer thumb |
| KP_N4 | SHIFT (left) | |
| KP_N5 | ALT (left) | |
| KP_N7 | ESC (right) | Inner thumb |
| KP_N8 | LAYER_KEYPAD (right) | |
| KP_N9 | LAYER_CMD (right) | Outer thumb |
| KP_N0 | SHIFT (right) | |
| KP_PLUS | SPACE (right) | |
| LCTRL | CTRL key | Center |
| KP_DIVIDE | TAB key | Center |

### Translation Coverage

| Source → Target | Keys Matched | Coverage |
|----------------|--------------|----------|
| adv360 → adv_mod | 70/76 | 92% |
| adv360 → glove80 | 70/80 | 88% |
| glove80 → adv_mod | 75/86 | 87% |

## Glove80 Generator Requirements

### Current Limitation

The glove80 output uses `writeGenericBindings` which doesn't preserve the proper row structure. To generate valid glove80 keymaps, we need a dedicated `writeGlove80Bindings` function.

### Required Implementation

**File**: `internal/generators/zmk_generator.go`

Add a new function `writeGlove80Bindings` that outputs the proper glove80 format:

```go
// writeGlove80Bindings writes bindings in Glove80 format
func (g *ZMKGenerator) writeGlove80Bindings(sb *strings.Builder, layer *models.Layer) {
    // Glove80 layout: 80 keys total
    // Row 0: 5 left + 5 right = 10 (function row)
    // Row 1: 6 left + 6 right = 12 (number row)
    // Row 2: 6 left + 6 right = 12 (QWERTY)
    // Row 3: 6 left + 6 right = 12 (home row)
    // Row 4: 6 left + 6 inner + 6 right = 18 (bottom + inner thumb)
    // Row 5: 5 left + 3 outer + 3 outer + 5 right = 16 (modifier + outer thumb)
    
    bindingMap := make(map[string]string)
    for _, binding := range layer.Bindings {
        key := fmt.Sprintf("%s_%d_%d", binding.Position.Side, binding.Position.Row, binding.Position.Col)
        bindingMap[key] = binding.Value
    }
    
    // Row 0: Function row (5+5)
    writeRow(sb, bindingMap, 0, 5, 5, "")
    
    // Rows 1-3: Main alpha rows (6+6 each)
    for row := 1; row <= 3; row++ {
        writeRow(sb, bindingMap, row, 6, 6, "")
    }
    
    // Row 4: Bottom + inner thumb (6 left + 6 inner + 6 right)
    // Format: 6 left main, 6 inner thumb, 6 right main
    writeRowWithThumb(sb, bindingMap, 4, 6, 6, 6)
    
    // Row 5: Modifier + outer thumb
    // Format: 5 left, 3 left thumb, 3 right thumb, 5 right
    writeModifierRow(sb, bindingMap, 5)
}
```

### Position Calculation for Glove80

Add `indexToPositionGlove80` function:

```go
func indexToPositionGlove80(index, totalKeys int) models.Position {
    // Glove80: 80 keys
    // Row 0: 10 keys (5+5 function)
    // Row 1-3: 12 keys each (6+6)
    // Row 4: 18 keys (6+6+6 with inner thumb)
    // Row 5: 16 keys (5+3+3+5 with outer thumb)
    
    var row, col int
    var side string
    
    switch {
    case index < 10:
        // Function row
        row = 0
        if index < 5 {
            side = "left"
            col = index
        } else {
            side = "right"
            col = index - 5
        }
    case index < 22:
        // Number row
        row = 1
        idx := index - 10
        // ... similar logic
    // ... remaining rows
    }
    
    return models.Position{Row: row, Col: col, Side: side}
}
```

### Generator Selection

Update `writeBindings` to select the correct generator:

```go
func (g *ZMKGenerator) writeBindings(sb *strings.Builder, layer *models.Layer) {
    switch g.keyboardType {
    case models.KeyboardZMKAdvMod:
        g.writeAdvModBindings(sb, layer)
    case models.KeyboardZMKGlove80:
        g.writeGlove80Bindings(sb, layer)  // NEW
    case models.KeyboardZMKAdv360:
        g.writeAdv360Bindings(sb, layer)   // NEW (optional)
    default:
        g.writeGenericBindings(sb, layer)
    }
}
```

### Testing Checklist

- [ ] `klcm translate --from adv360 --to glove80` produces valid glove80 format
- [ ] Row structure matches original glove80.keymap
- [ ] All 80 positions are correctly filled
- [ ] Thumb cluster keys in proper positions
- [ ] Inner thumb (row 4) formatted correctly with 18 keys
- [ ] Outer thumb (row 5) formatted correctly with 16 keys

## Sync Workflow Documentation

For future agents working on keyboard configuration syncing:

### Automated Sync Process (Current)

The sync from adv360 to adv_mod is now fully automated using the IR (Intermediate Representation) system:

1. **IR Translation** - `klcm translate --from adv360 --to adv_mod`
2. **Universal Mapping** - Automatic key position translation via 10x10 grid
3. **Behavior Transfer** - Mod-morphs and macros automatically included
4. **Hardware Preservation** - Template system preserves keyboard-specific features
5. **Validation** - Built-in syntax checking ensures correctness

### Key Learnings

- **Physical layout matters**: adv_mod has 18-key function row that adv360 doesn't have
- **Layer numbering**: Must update LAYER_* defines when adding new layers
- **Bootloader layers**: Never sync bootloader/system layers - they're hardware-specific
- **Modifier choices**: Use RC() (Control) not RG() (GUI) for CMD layer to match glove80
- **Morphs vs Macros**: Morphs change behavior with modifiers, macros type sequences

### Automation Features

The IR translation system provides:
- Multi-layer sync support ✅
- Behavior/macro detection and copying ✅
- Layer mapping configuration ✅
- Physical layout awareness (preserve hardware-specific keys) ✅
- Zone-based intelligent key mapping ✅

---

## Translation Tool Issues (2025-11-28)

### Current Problems with `klcm translate`

When running `klcm translate --from adv360 --to adv_mod`, the generated output has several critical issues:

#### Issue 1: Wrong Layer Structure
**Problem**: Generated output has 9 layers with wrong numbering (LAYER_KEYPAD=6, LAYER_FN=7, LAYER_MOD=8)
**Expected**: 4 functional layers (LAYER_KEYPAD=1, LAYER_CMD=2, LAYER_SYSTEM=3) + mapping_layer
**Root Cause**: The translation copies all source layers verbatim without filtering or renumbering
**Fix Location**: `internal/mappers/unified_mapper.go` - `TranslateLayout()` needs layer filtering

#### Issue 2: Macros in Wrong Section
**Problem**: Macros are embedded in the behaviors section instead of a separate macros section
**Expected**: ZMK format requires `macros { }` block before `behaviors { }`
**Fix Location**: `internal/generators/zmk_generator.go` - `writeBehaviors()` should not include macros

#### Issue 3: Missing Morph Behaviors
**Problem**: `morph_parens_left`, `morph_parens_right`, `macro_parens`, `macro_angle_brackets` are missing
**Expected**: All morph behaviors from source should be detected and included
**Fix Location**: `internal/generators/zmk_generator.go` - `writeMorphBehaviors()` needs to check for all morph types

#### Issue 4: Wrong Layer Names
**Problem**: Generated uses `layer0_default`, `layer7_fn` instead of `default_layer`, `cmd_layer`
**Expected**: Use semantic names that match the source keymap naming conventions
**Fix Location**: `internal/generators/zmk_generator.go` - `writeLayer()` should use proper names

#### Issue 5: Uses RG() Instead of RC()
**Problem**: The CMD layer uses `RG(Q)` (Right GUI/Command) instead of `RC(Q)` (Right Control)
**Expected**: adv_mod should use RC() for the CMD layer to match glove80 behavior
**Fix Location**: Either in translation or as a configurable option per keyboard

#### Issue 6: Source Metadata Shows `%!s(<nil>)`
**Problem**: The header shows `Source: %!s(<nil>)` instead of proper source info
**Fix Location**: `internal/generators/zmk_generator.go` - `writeHeader()` needs nil check

#### Issue 7: Padding Layers Included
**Problem**: Empty padding layers (3, 4, 5) are included in output
**Expected**: Padding layers should be filtered out during translation
**Fix Location**: `internal/mappers/unified_mapper.go` - filter layers by name/content

### Fix Priority Order

1. ✅ **Layer filtering** - Remove padding layers, only include functional layers
2. ✅ **Layer renumbering** - Update LAYER_* defines to match target expectations
3. ✅ **Macros section** - Move macros to separate section before behaviors
4. ✅ **Missing behaviors** - Ensure all morph/macro types are detected
5. ✅ **Layer names** - Use semantic names (`default_layer`, `keypad_layer`, etc.)
6. ✅ **Header metadata** - Fix nil source reference with proper keyboard-specific headers
7. ⏳ **RG→RC conversion** - Add per-keyboard modifier mapping (optional, documented)

### Fixes Implemented (2025-11-28)

**Generator improvements:**
- Header now shows keyboard-specific info (repo URL, branch for adv_mod)
- Macros section is separate from behaviors section (ZMK format requirement)
- All morph behaviors include comment hints (e.g., `// &dot_override,`)
- Semantic layer names: `default_layer`, `keypad_layer`, `cmd_layer`, `system_layer`
- Descriptive layer comments
- Padding layers are filtered out during generation
- Fixed layer defines for adv_mod: LAYER_KEYPAD=1, LAYER_CMD=2, LAYER_SYSTEM=3

**Remaining items:**
- RG()→RC() conversion for CMD layer (optional, adv360 uses RG, glove80 uses RC)
- Per-keyboard modifier preference configuration

### Critical Translation Issues Discovered (2025-11-28)

#### Physical Layout Mismatch

| Row | adv360 | adv_mod | Impact |
|-----|--------|---------|--------|
| Number row | 14 keys (7+7) | 12 keys (6+6) | Macros in positions 11-14 don't map |
| Function row | 12 keys (6+6) | 18 keys | adv_mod has 6 extra keys |
| System layer | Adv360-specific | Nice!Nano-specific | Should NOT translate |

#### Missing Macros in Generated Output

**Problem**: `macro_parens` and `macro_angle_brackets` exist in adv360 keypad layer but don't appear in generated adv_mod output.

**Root Cause**: The adv360 number row has 14 keys per row (7 left + 7 right), but adv_mod only has 12 keys per row (6 left + 6 right). The macros are in positions that correspond to keys N9, N0, MINUS, EQUAL on adv360's mapping layer, but adv_mod's mapping layer only goes up to N0.

**Source (adv360 keypad layer)**:
```
&kp CARET  &macro_brackets  &macro_braces  &macro_parens  &macro_angle_brackets
```

**Generated (adv_mod)**: Only `&macro_brackets` and `&macro_braces` appear; the others map to `&trans`.

#### System Layer Should Be Preserved

**Problem**: Generated output overwrites adv_mod's system_layer with adv360's system layer, which has different:
- Bootloader positions
- Bluetooth configuration
- RGB/backlight settings

**Required Behavior**: System layer (and other hardware-specific layers) should be **preserved from target**, not translated from source.

#### RC vs RG Modifier Difference

**Remote (correct for adv_mod)**:
```
&kp RC(Q)  &kp RC(W)  &kp RC(E)  // Right Control
```

**Generated (from adv360)**:
```
&kp RG(Q)  &kp RG(W)  &kp RG(E)  // Right GUI/Command
```

### Required Fixes (Priority Order)

1. **Layer Preservation** - Add logic to preserve target's hardware-specific layers:
   - `system_layer` - bootloader, BT, RGB (hardware-specific)
   - `mapping_layer` - already being preserved
   - Only translate: `default_layer`, `keypad_layer`, `cmd_layer`

2. **Column Count Handling** - Handle physical layout differences:
   - adv360: 7 columns per side on main rows
   - adv_mod: 6 columns per side on main rows
   - Keys in "extra" columns need alternative placement or explicit drop

3. **Modifier Mapping** - Add per-keyboard modifier preferences:
   - adv360 uses `RG()` (Right GUI/Command)
   - adv_mod/glove80 should use `RC()` (Right Control)
   - Make this configurable per target keyboard

4. **Macro Detection Fix** - Ensure all macros from source are detected even if they map to `&trans` positions (so they can be manually placed)

### Files to Modify

1. `internal/mappers/unified_mapper.go`
   - Add layer filtering (remove padding)
   - Add layer renumbering based on target layout

2. `internal/generators/zmk_generator.go`
   - Separate macros section from behaviors
   - Fix layer naming convention
   - Fix header metadata nil check
   - Add all morph behavior types

3. `internal/cli/translate.go`
   - Pass target layer configuration to mapper
   - Support layer mapping (source layer X → target layer Y)