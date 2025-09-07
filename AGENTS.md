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

1. **Immediate**: Create zmk_adv_mod configuration in KLCM
2. **Setup**: Create zmk-adv-mod-config repository
3. **Integration**: Test build pipeline end-to-end
4. **Validation**: Flash and test firmware on actual hardware
5. **Documentation**: Complete setup and usage documentation

---

*Last Updated: 2025-01-06*  
*Version: v7_target*