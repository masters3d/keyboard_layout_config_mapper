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

**PHASE 2: 🚧 IN PROGRESS - ZMK Config Repository Setup**

#### Step 1: Create zmk-adv-mod-config Repository
```bash
# Create new repository on GitHub
gh repo create zmk-adv-mod-config --public --description "ZMK configuration for Kinesis Advantage keyboards using Pillz Mod with Nice!Nano"

# Clone and set up locally  
git clone https://github.com/masters3d/zmk-adv-mod-config.git
cd zmk-adv-mod-config

# Copy template files (available in /tmp/zmk-adv-mod-config-template)
cp -r /tmp/zmk-adv-mod-config-template/* .

# Initial commit
git add .
git commit -m "Initial ZMK config setup for Advanced Mod

- GitHub Actions build workflow
- Custom zmk-adv-mod shield definition  
- Nice!Nano v2 build configuration
- Base keymap with HOME key positioning
- Shield files: Kconfig, overlay, keymap, conf"

git push origin main
```

#### Step 2: Test Initial Build
1. **Push repository** - Triggers GitHub Actions build
2. **Check build results** - Should produce `.uf2` artifacts
3. **Download artifacts** - Test firmware files are generated
4. **Fix any build issues** - Adjust pin mappings or configurations

#### Step 3: Hardware Validation (CRITICAL)
⚠️ **PIN MAPPING VERIFICATION REQUIRED**

The template overlay uses reasonable pin assignments, but **MUST BE VERIFIED** against actual Pillz Mod Pro PCB:

1. **Check PCB schematic** or Pillz Mod documentation for Nice!Nano pin mapping
2. **Verify row/column assignments** match physical PCB traces  
3. **Confirm LED pin assignments** for status indicators
4. **Test matrix scanning** - ensure all keys register correctly

**Pin Assignment Notes:**
- Rows: Uses pro_micro pins 0-10, 14-16, 18 (15 total)
- Columns: Uses pro_micro pins 19-21 + A0-A3 (7 total)  
- LEDs: Commented out - requires PCB verification

#### Step 4: Keymap Integration Testing
1. **Update KLCM download.go** to point to new repository
2. **Test keymap sync** from KLCM to ZMK config repo  
3. **Verify PR automation** works with new repository
4. **Test build trigger** on keymap changes

**PHASE 3: 📋 PLANNED - Production Validation**
- [ ] Flash firmware to actual hardware
- [ ] Test all keys and matrix positions  
- [ ] Verify wireless connectivity
- [ ] Test battery management
- [ ] Validate status LEDs
- [ ] Document flashing and usage process

**PHASE 4: 📋 PLANNED - Documentation & Maintenance**
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

*Last Updated: 2025-01-06*  
*Version: v7_target*

---

## Kinesis Advantage 2 SmartSet Configuration Guide

### Overview
SmartSet is the native programming system for Kinesis Advantage 2 keyboards. Configuration files use a text-based syntax stored on the keyboard's internal drive. KLCM supports SmartSet through the `kinesis2` keyboard type.

**Configuration Location**: `configs/kinesis2/1_qwerty.txt`

### SmartSet Syntax Reference

#### Basic Key Remapping
```
[source_key]>[target_key]
```
Example: `[caps]>[bspace]` - Remap Caps Lock to Backspace

#### Keypad Layer Remapping
Prefix with `kp-` for keypad layer:
```
[kp-caps]>[bspace]
```

#### Macro Definitions
Use curly braces `{}` for macros:
```
{trigger_key}>{speed9}{key1}{key2}{key3}
```
Example: `{=}>{speed5}{-lshift}{-lwin}{s}{+lshift}{+lwin}` - Screenshot macro

#### Modifier Syntax
- `{-modifier}` = Press and hold modifier
- `{+modifier}` = Release modifier
- Modifiers: `lshift`, `rshift`, `lctrl`, `rctrl`, `lalt`, `ralt`, `lwin`, `rwin`

Example: `{-lctrl}{bspace}{+lctrl}` = Ctrl+Backspace

### ⚠️ Critical: Macro + Null Pattern

**IMPORTANT FOR AI AGENTS**: When defining a macro `{key}>{...}`, you MUST also set `[key]>[null]` so the base key outputs nothing and the macro is triggered instead.

**Without `[key]>[null]`**: The key outputs its default character before/instead of executing the macro.

**Correct Pattern**:
```
*# Comment explaining the macro
[key]>[null]
{key}>{speed9}{macro_sequence}
```

**Example - Screenshot on `=` key**:
```
*# This will take a screenshot (matches ZMK LS(LG(S)))
[=]>[null]
[kp-=]>[null]
{=}>{speed5}{-lshift}{-lwin}{s}{+lshift}{+lwin}
{kp-=}>{speed5}{-lshift}{-lwin}{s}{+lshift}{+lwin}
```

**Example - Bracket macros**:
```
*# macro_brackets: []←
[kp=]>[null]
{kp=}>{speed9}{obrack}{cbrack}{left}

*# macro_braces: {}←
[kpdiv]>[null]
{kpdiv}>{speed9}{-lshift}{obrack}{+lshift}{-lshift}{cbrack}{+lshift}{left}
```

### Speed Settings
- `{speed1}` to `{speed9}` - Macro playback speed (9 = fastest)
- Use `{speed9}` for simple key sequences
- Use `{speed5}` for complex multi-key operations

### Common Key Names

| SmartSet Name | Key |
|---------------|-----|
| `lshift`, `rshift` | Shift keys |
| `lctrl`, `rctrl` | Control keys |
| `lalt`, `ralt` | Alt/Option keys |
| `lwin`, `rwin` | Windows/Command keys |
| `bspace` | Backspace |
| `escape` | Escape |
| `obrack`, `cbrack` | [ and ] |
| `hyphen` | - (minus) |
| `pup`, `pdown` | Page Up, Page Down |
| `kpshft` | Keypad layer toggle |
| `intl-\` | International backslash |
| `null` | No output (for macro triggers) |

### Comments
Lines starting with `*#` are comments:
```
*# This is a comment
*# Aligned with ZMK - description of what this does
```

### Layer System
SmartSet has two layers:
1. **Default Layer**: Normal key mappings
2. **Keypad Layer**: Accessed via `kpshft`, prefixed with `kp-`

Both layers should be configured together:
```
[caps]>[bspace]
[kp-caps]>[bspace]
```

### Aligning SmartSet with ZMK

When maintaining both SmartSet (Kinesis Advantage 2) and ZMK (Pillz Mod) configurations:

1. **ZMK is the source of truth** for layout decisions
2. **SmartSet should mirror ZMK** where hardware allows
3. **Use F13-F24** for software-overridable keys
4. **Document differences** with comments

**Example - F13-F24 for software override**:
```
*# left keypad row 1
*# Aligned with ZMK - F13-F16 for software override
[kp-q]>[F13]
[kp-w]>[F14]
[kp-e]>[F15]
[kp-r]>[F16]
```

### Validation Checklist

When updating SmartSet configs:
- [ ] All macros have corresponding `[key]>[null]` entries
- [ ] Both default and keypad layer versions exist (`[key]` and `[kp-key]`)
- [ ] Comments explain non-obvious mappings
- [ ] Configuration aligns with ZMK source of truth
- [ ] Modifier macros use correct `{-mod}...{+mod}` syntax

---

## Session Context - 2025-01-28

### Key Learnings from This Session

#### 1. ZMK is the Source of Truth
- All keyboard configs (adv360, glove80, pillzmod_pro) should be kept in sync
- SmartSet (kinesis2) mirrors ZMK where hardware allows
- The user upgraded their Kinesis Advantage 2 to run ZMK via Pillz Mod, so SmartSet config is now backup/reference only

#### 2. Repositories Structure
- **zmk-config-pillzmod-nicenano**: Actual ZMK firmware config (in `/Volumes/ExternalCheyo/source/`)
- **keyboard_layout_config_mapper**: KLCM tool with configs in `configs/` subdirectory
- The `pillzmod_pro.keymap` exists in BOTH locations - keep them in sync

#### 3. SmartSet Macro Pattern (CRITICAL)
When defining macros in SmartSet, you MUST use this pattern:
```
[key]>[null]           # Suppress default key output
{key}>{macro_sequence} # Define the macro
```
Without `[key]>[null]`, the key outputs its default character before/instead of the macro.

#### 4. F13-F24 Strategy
- Left keypad rows 1-3 use F13-F24 for software override (Karabiner, BetterTouchTool)
- This allows custom per-app shortcuts without changing firmware
- SmartSet and ZMK are aligned on this

#### 5. Thumb Cluster Unified Mapping
All keyboards use identical logical positions:
```
L1=SPACE  L2=LSHIFT  L3=LALT  L4=LCTRL  L5=LWIN  L6=mo_KEYPAD
R1=SPACE  R2=RSHIFT  R3=mo_CMD  R4=LCTRL  R5=ESC  R6=mo_KEYPAD
```
See `configs/THUMB_CLUSTER_MAPPING.md` for full documentation.

#### 6. CMD Layer has Vim-style Navigation
Right nav cluster in cmd_layer: `RC(LEFT) RC(DOWN) RC(UP) RC(RIGHT)`
This enables word-by-word navigation on Mac (Cmd+arrows).

#### 7. Morph Behaviors
Custom shifted outputs (same across all ZMK configs):
- `.` → `:` (morph_dot)
- `,` → `;` (morph_comma)
- `(` → `<` (morph_parens_left)
- `)` → `>` (morph_parens_right)
- `\` → `!` (morph_exclamation)
- `'` → `` ` `` (morph_quote_single)
- `"` → `~` (morph_quote_double)

#### 8. Hardware Context
- User has Kinesis Advantage 2 upgraded with Pillz Mod Pro + Nice!Nano v2
- Running ZMK firmware (not SmartSet anymore)
- SmartSet config kept for reference/backup in case needed

#### 9. Key Files to Keep in Sync
1. `zmk-config-pillzmod-nicenano/config/pillzmod_pro.keymap` (actual firmware)
2. `keyboard_layout_config_mapper/configs/zmk_adv_mod/pillzmod_pro.keymap` (KLCM copy)
3. `keyboard_layout_config_mapper/configs/zmk_adv360/adv360.keymap`
4. `keyboard_layout_config_mapper/configs/zmk_glove80/glove80.keymap`
5. `keyboard_layout_config_mapper/configs/kinesis2/1_qwerty.txt` (SmartSet backup)

#### 10. Pedal Configuration (Pillz Mod Pro)
The Pillz Mod Pro PCB supports 3 foot pedals via matrix positions `RC(0, 6)`, `RC(1, 6)`, and `RC(2, 6)`.

**Matrix Position in Keymap**:
The pedal row appears at the end of each layer's bindings, after the bottom row of the thumb cluster:
```
// ... thumb cluster row
&kp SPACE   &kp LEFT_SHIFT      &kp LEFT_ALT    &mo LAYER_CMD  &kp RIGHT_SHIFT &kp SPACE
&kp X  &kp ESC  &kp ESC   // <-- Pedal row: Pedal 1, Pedal 2, Pedal 3
```

**Default Configuration**:
- **Pedal 1**: `X` key (for software-defined actions)
- **Pedal 2**: `ESC` key
- **Pedal 3**: `ESC` key

**Layer Behavior**:
- In non-default layers (keypad, cmd, system), pedals use `&trans` to inherit from the default layer.

**Historical Note**:
The cheyo branch initially omitted the pedal keys, which was an error. The main branch and dcpedit's reference implementation include all 3 pedals. This was corrected to ensure keymap matrix completeness.