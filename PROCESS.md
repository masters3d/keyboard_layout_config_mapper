# 🎹 KLCM Development & Usage Process Guide

> **For Future Development and Daily Usage**  
> **Last Updated**: 2024-12-25  
> **Version**: 1.0  

## 📋 Quick Start for Daily Usage

### **Essential Commands**
```bash
# Build the CLI tool (after any code changes)
go build -o klcm cmd/klcm/main.go

# Pull latest configurations from remote repositories
./klcm pull --preview                    # Preview changes first
./klcm pull                              # Apply updates

# Validate configurations
./klcm validate                          # Validate all keyboards
./klcm validate adv360                   # Validate specific keyboard

# Sync changes between keyboards  
./klcm sync adv360 kinesis2 --preview    # Preview cross-format sync
./klcm sync adv360 kinesis2              # Apply sync changes
./klcm sync adv360 glove80 --preview     # Preview ZMK-to-ZMK sync

# Compare with remote versions
./klcm compare-remote                    # Compare all with remote
./klcm compare-remote adv360             # Compare specific keyboard

# Download configurations (initial setup)
./klcm download --preview                # Preview what would be downloaded
./klcm download                          # Download all configurations

# Interactive workflow (guided experience)
./klcm workflow                          # Step-by-step guided process
```

---

## 🔄 Daily Workflow Process

### **Scenario A: Making Changes to Advantage360 and Syncing to Kinesis2**

```bash
# 1. Update to latest configurations
./klcm pull --preview                    # See what's changed remotely
./klcm pull                              # Apply updates

# 2. Edit your Advantage360 layout
# Edit configs/zmk_adv360/adv360.keymap manually

# 3. Validate your changes
./klcm validate adv360                   # Ensure syntax is correct

# 4. Preview sync to Kinesis2  
./klcm sync adv360 kinesis2 --preview    # See what would change

# 5. Apply sync to Kinesis2
./klcm sync adv360 kinesis2              # Apply the changes

# 6. Validate everything
./klcm validate                          # Ensure all configs are valid

# 7. Create pull requests (when ready)
./klcm pr create --dry-run               # Preview PR creation
./klcm pr create                         # Create actual PRs
```

### **Scenario B: Syncing Between ZMK Keyboards**

```bash
# 1. Make changes to one ZMK keyboard (e.g., adv360)
# 2. Sync to the other ZMK keyboard
./klcm sync adv360 glove80 --preview     # Preview changes
./klcm sync adv360 glove80               # Apply sync

# 3. Validate both keyboards
./klcm validate adv360 glove80           # Validate specific keyboards
```

---

## 🛠️ Development Process

### **Adding New Keyboard Support**

1. **Add to KeyboardType enum** in `internal/models/types.go`:
   ```go
   const (
       KeyboardNewType KeyboardType = "newtype"
   )
   ```

2. **Create parser** in `internal/parsers/`:
   - Implement the `Parser` interface
   - Add `Parse()` and `Validate()` methods
   - Handle the specific file format

3. **Add to CLI mappings** in `internal/cli/sync.go`:
   - Update `getKeyboardConfigPath()` function
   - Add sync support in sync functions

4. **Add to config paths** in `internal/parsers/parser.go`:
   - Update `GetConfigPath()` function

### **Improving Sync Logic**

1. **Enhance ZMK analysis** in `analyzeZMKChanges()`:
   - Add detection for new key patterns
   - Detect new behaviors and mod-morphs
   - Handle complex layer switching

2. **Improve mapping** in `mapZMKToKinesis2()`:
   - Add new keycode mappings
   - Handle new ZMK behaviors
   - Convert complex macros

3. **Test thoroughly**:
   ```bash
   ./klcm sync adv360 kinesis2 --preview   # Always preview first
   cp configs/kinesis2/1_qwerty.txt backup # Backup before testing
   ./klcm sync adv360 kinesis2             # Test actual sync
   ./klcm validate                         # Validate results
   diff backup configs/kinesis2/1_qwerty.txt # Check changes
   ```

### **Adding New Commands**

1. **Create command file** in `internal/cli/`:
   ```go
   // newcommand.go
   var newCmd = &cobra.Command{
       Use:   "new",
       Short: "Description",
       RunE:  runNew,
   }
   
   func init() {
       rootCmd.AddCommand(newCmd)
   }
   ```

2. **Follow existing patterns**:
   - Use consistent flag names (`--preview`, `--verbose`, etc.)
   - Provide helpful output with emojis and clear formatting
   - Add validation and error handling
   - Support dry-run/preview modes

---

## 🎯 Key Features & Capabilities

### **Current Sync Support Matrix**

| Source | Target | Status | Description |
|--------|--------|--------|-------------|
| adv360 | kinesis2 | ✅ Full | ZMK→Kinesis2 with mod-morph conversion |
| adv360 | glove80 | 🟡 Basic | ZMK→ZMK (needs enhancement) |
| glove80 | adv360 | 🟡 Basic | ZMK→ZMK (needs enhancement) |
| adv360 | qmk_ergodx | ❌ None | Not implemented |
| kinesis2 | adv360 | ❌ None | Not implemented |

### **Supported ZMK Features**

- ✅ **Basic keycodes**: `&kp KEY` → Kinesis2 `[key]>[target]`
- ✅ **Mod-morphs**: `morph_dot` → Kinesis2 `{lshift}{.}>{macro}`
- ✅ **Layer access**: `&mo LAYER` → Function key mappings
- ✅ **Basic behaviors**: Custom behaviors converted to macros
- 🟡 **Combos**: Detected but not yet converted
- ❌ **Hold-taps**: Not yet implemented
- ❌ **Complex macros**: Partial support

### **Kinesis2 Features Generated**

- ✅ **Basic remapping**: `[source]>[target]`
- ✅ **Shift macros**: `{lshift}{key}>{speed9}{macro}`
- ✅ **Speed control**: `{speed9}` for consistent timing
- ✅ **Modifier handling**: `{-lshift}`, `{+lshift}` for modifier control
- 🟡 **Complex macros**: Basic support, needs enhancement
- ❌ **Pedal mappings**: Not yet implemented

---

## 🔍 Validation & Testing

### **Before Making Changes**
```bash
./klcm validate                          # Baseline validation
./klcm pull --preview                    # Check for remote updates
```

### **After Making Changes**
```bash
./klcm validate                          # Ensure syntax correctness
./klcm sync source target --preview      # Preview sync effects
```

### **Before Committing**
```bash
git status                               # See what files changed
git diff                                 # Review exact changes
./klcm validate                          # Final validation
```

---

## 📊 Understanding Sync Output

### **Sync Preview Interpretation**

When you run `./klcm sync adv360 kinesis2 --preview`, you'll see:

```
📊 Found X mappable changes:

1. Map Keypad Tab to Backspace (matches ZMK top-left change)
   Current: [kp-tab]>[escape]           # Current Kinesis2 mapping
   Proposed: [kp-tab]>[bspace]          # Proposed change

2. Add macro for Shift+Period → Colon (matches ZMK mod-morph)
   Current:                             # Empty = new line will be added
   Proposed: {lshift}{.}>{speed9}{-lshift}{;}{+lshift}  # Kinesis2 macro
```

**Understanding the mappings**:
- **Basic remaps**: `[key]>[target]` - direct key remapping
- **Macros**: `{modifier}{key}>{macro}` - complex key sequences
- **Speed control**: `{speed9}` - consistent macro timing
- **Modifier control**: `{-lshift}`, `{+lshift}` - manage shift state

---

## 🚨 Important Considerations

### **Backup Strategy**
```bash
# Always backup before applying changes
cp configs/kinesis2/1_qwerty.txt configs/kinesis2/1_qwerty.txt.backup
./klcm sync adv360 kinesis2
# If issues occur:
cp configs/kinesis2/1_qwerty.txt.backup configs/kinesis2/1_qwerty.txt
```

### **Limitations to Keep in Mind**

1. **ZMK→Kinesis2 sync limitations**:
   - Complex hold-taps may not convert perfectly
   - Some ZMK behaviors have no Kinesis2 equivalent
   - Layer switching is mapped to function keys (limited)

2. **Kinesis2 format constraints**:
   - Limited macro complexity compared to ZMK
   - No direct equivalent for some ZMK features
   - Speed control needed for reliable macro execution

3. **Cross-platform considerations**:
   - Different physical layouts between keyboards
   - Not all key positions have direct equivalents
   - Some features are firmware-specific

### **When to Use Manual Editing vs Sync**

**Use sync when**:
- ✅ Copying basic layout changes
- ✅ Converting simple mod-morphs
- ✅ Applying widespread key remappings

**Use manual editing when**:
- ❌ Creating keyboard-specific optimizations
- ❌ Implementing complex behaviors
- ❌ Fine-tuning for specific use cases

---

## 🎉 Success Indicators

### **Signs of Successful Sync**
```bash
# These should all pass:
./klcm validate                          # ✅ All configurations valid
./klcm sync source target --preview      # 📊 Reasonable number of changes
diff backup current_file                 # 🔍 Changes look correct
```

### **Red Flags**
- Validation failures after sync
- Hundreds of changes for simple layout modifications
- Syntax errors in generated Kinesis2 macros
- Missing or duplicated mappings

---

This process document should serve as your go-to reference for both daily usage and future development of the KLCM tool. Keep it updated as you add new features and discover better workflows! 🚀