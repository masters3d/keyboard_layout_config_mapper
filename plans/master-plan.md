# 🎹 Keyboard Layout Configuration Mapper - Master Plan

> **Living Document** - Updated as project progresses  
> **Last Updated**: 2024-12-19  
> **Version**: 1.0  

## 📊 Project Overview

### **Mission Statement**
Create an intelligent system to synchronize keyboard layout configurations across 4 different keyboard types, with primary focus on ZMK-based keyboards and automated PR creation to upstream repositories.

### **Current Inventory**
- **Primary Keyboards** (ZMK-based, user's main keyboards):
  - ✅ Advantage360 (`configs/zmk_adv360/adv360.keymap`)
  - ✅ Glove80 (`configs/zmk_glove80/glove80.keymap`)
- **Secondary Keyboards**:
  - ✅ QMK ErgoDox (`configs/qmk_ergodx/keymap.c`)
  - ✅ Kinesis Advantage 2 (`configs/kinesis2/1_qwerty.txt`)

### **Upstream Repositories**
- **Advantage360**: `masters3d/Adv360-Pro-ZMK` (cheyo branch)
- **Glove80**: `masters3d/glove80-zmk-config` (cheyo branch)  
- **QMK ErgoDox**: `masters3d/qmk_firmware` (masters3d branch)
- **Kinesis2**: `masters3d/supportfiles` (master branch)

---

## 🎯 Strategic Priorities

### **Priority 1: ZMK Compatibility** (🔥 HIGH)
- Seamless sync between Advantage360 ↔ Glove80
- These are user's primary keyboards - must work flawlessly
- Focus on common ZMK features: layers, behaviors, combos, home row mods

### **Priority 2: Change Propagation** (🟡 MEDIUM)
- ZMK → QMK translation
- ZMK → Kinesis2 basic remapping
- Handle incompatible features gracefully

### **Priority 3: Automation & Integration** (🔥 HIGH)
- CLI tools for daily workflow
- Automated PR creation to upstream repos
- LLM integration for complex scenarios

---

## 🏗️ Architecture Design

### **Project Structure**
```
keyboard_layout_config_mapper/
├── plans/                    # Project documentation (this file)
├── configs/                  # Downloaded keyboard configurations
├── src/                      # Source code (to be created)
│   ├── parsers/             # Configuration file parsers
│   ├── mappers/             # Key mapping logic
│   ├── generators/          # Configuration generators
│   ├── diff/                # Change detection
│   ├── cli/                 # Command-line interface
│   ├── github/              # GitHub integration
│   └── llm/                 # LLM integration helpers
├── tests/                    # Test suites
├── docs/                     # Technical documentation
└── CLI Download Integration   # ✅ COMPLETED: Integrated into klcm CLI
```

### **Core Components**

#### **1. Configuration Parsers** (`src/parsers/`)
- **ZMK Parser**: Parse `.keymap` files (devicetree syntax)
- **QMK Parser**: Parse `.c` files (C macros and functions)
- **Kinesis2 Parser**: Parse `.txt` files (simple remapping format)
- **Output**: Unified internal representation (JSON/YAML)

#### **2. Layout Mapping Engine** (`src/mappers/`)
- **Physical Mappings**: Key position translations between layouts
- **Logical Mappings**: Function mappings (e.g., home row mods)
- **Layer Mappings**: Layer structure translations
- **Feature Compatibility**: Handle keyboard-specific features

#### **3. Change Detection System** (`src/diff/`)
- **Git Integration**: Track changes in config files
- **Semantic Analysis**: Understand what changed functionally
- **Conflict Detection**: Identify incompatible changes
- **Change Propagation**: Determine what needs updating

#### **4. Configuration Generators** (`src/generators/`)
- **ZMK Generator**: Create valid `.keymap` files
- **QMK Generator**: Create valid `.c` files
- **Kinesis2 Generator**: Create valid `.txt` files
- **Validation**: Ensure generated configs are syntactically correct

#### **5. CLI Interface** (`src/cli/`)
- **Sync Commands**: Sync between keyboards
- **Interactive Mode**: Handle complex mappings
- **Preview Mode**: Dry-run capabilities
- **Validation Commands**: Check config validity

#### **6. GitHub Integration** (`src/github/`)
- **PR Creation**: Automated pull request creation
- **Branch Management**: Handle feature branches
- **Status Tracking**: Monitor PR status
- **Conflict Resolution**: Handle merge conflicts

---

## 📋 Implementation Roadmap

### **Phase 1: Foundation & Analysis** 
- [x] ✅ Repository setup with config downloads
- [ ] 📝 Deep analysis of existing ZMK configurations
- [ ] 📝 Create physical key mapping matrices
- [ ] 📝 Document common patterns and differences
- [ ] 📝 Identify ZMK features used in current configs

### **Phase 2: Core Parsing Infrastructure**
- [ ] 🔨 Build ZMK parser for .keymap files
- [ ] 🔨 Build QMK parser for .c files
- [ ] 🔨 Build Kinesis2 parser for .txt files
- [ ] 🔨 Design unified internal representation format
- [ ] 🔨 Create parser tests and validation

### **Phase 3: ZMK-to-ZMK Mapping** (Priority 1)
- [ ] 🔨 Analyze Advantage360 vs Glove80 layouts
- [ ] 🔨 Create physical key mapping
- [ ] 🔨 Build ZMK-to-ZMK translator
- [ ] 🔨 Handle layer mapping between keyboards
- [ ] 🔨 Test bidirectional sync

### **Phase 4: Cross-Platform Mapping**
- [ ] 🔨 ZMK → QMK translation engine
- [ ] 🔨 ZMK → Kinesis2 translation engine
- [ ] 🔨 Handle feature incompatibilities
- [ ] 🔨 Create fallback strategies

### **Phase 5: CLI Development**
- [ ] 🔨 Basic CLI framework
- [ ] 🔨 Sync commands (`klcm sync`)
- [ ] 🔨 Interactive mode (`klcm sync --interactive`)
- [ ] 🔨 Preview mode (`klcm sync --preview`)
- [ ] 🔨 Validation commands (`klcm validate`)

### **Phase 6: GitHub Integration**
- [ ] 🔨 GitHub API client
- [ ] 🔨 PR creation automation
- [ ] 🔨 Branch management
- [ ] 🔨 Status monitoring

### **Phase 7: LLM Integration**
- [ ] 🔨 LLM prompt generation for complex mappings
- [ ] 🔨 Structured input/output for LLM
- [ ] 🔨 Integration with CLI for unsupported scenarios
- [ ] 🔨 Template system for common patterns

---

## 🔄 Proposed Workflows

### **Daily Workflow (Target State)**
```bash
# 1. Make changes to primary ZMK keyboard (Advantage360)
# 2. Sync changes to other keyboards
./klcm sync --from zmk_adv360 --to all --preview
./klcm sync --from zmk_adv360 --to all --apply

# 3. Create PRs to upstream repositories
./klcm pr create --all
./klcm pr status  # Check PR status
```

### **Interactive Workflow (Complex Changes)**
```bash
# For changes requiring manual decisions
./klcm sync --interactive
./klcm resolve-conflicts
./klcm validate --all
```

### **LLM-Assisted Workflow (Edge Cases)**
```bash
# When automatic mapping fails
./klcm sync --llm-assist
./klcm generate-prompt --change-summary
# Copy prompt to LLM, get response, import result
./klcm import-llm-result --file response.json
```

---

## 🎮 Keyboard-Specific Technical Details

### **ZMK Keyboards (Advantage360 & Glove80)**

#### **Common ZMK Features to Support:**
- **Layers**: Multiple keyboard layers
- **Behaviors**: `&kp`, `&mo`, `&lt`, `&mt` (mod-tap), `&sk` (sticky key)
- **Combos**: Key combinations for additional functions
- **Home Row Mods**: Modifier keys on home row
- **Tap-Dance**: Multiple tap behaviors
- **Unicode**: Special character input

#### **Key Mapping Challenges:**
- **Different Physical Layouts**: Advantage360 (split) vs Glove80 (different thumb clusters)
- **Key Count Differences**: Need to handle missing/extra keys
- **Thumb Cluster Variations**: Different thumb key arrangements

### **QMK (ErgoDox)**

#### **QMK → ZMK Translation Patterns:**
- `KC_*` → `&kp` behaviors
- `LT()` → `&lt` (layer-tap)
- `MT()` → `&mt` (mod-tap)
- Custom functions → ZMK macros or combos
- Layer handling differences

### **Kinesis Advantage 2**

#### **Limitations:**
- Basic remapping only (no layers, no complex behaviors)
- Limited to simple key-to-key mapping
- Text-based configuration format

#### **Translation Strategy:**
- Extract base layer mappings from ZMK
- Ignore complex behaviors
- Focus on basic alphanumeric and modifier remapping

---

## 🤖 LLM Integration Strategy

### **When to Use LLM Assistance:**
1. **Complex Feature Translation**: When automatic mapping fails
2. **New ZMK Features**: Features not yet supported by the tool
3. **Custom Layouts**: Highly customized configurations
4. **Conflict Resolution**: When multiple strategies are possible

### **LLM Prompt Templates:**
- Current keyboard configurations
- Specific change being made
- Target keyboard capabilities
- Desired outcome specification

### **Structured Input/Output:**
- JSON format for configuration changes
- Standardized change description format
- Validation rules for LLM outputs

---

## 📊 Success Metrics & Testing

### **Success Criteria:**
- ✅ **ZMK ↔ ZMK sync**: 95%+ accuracy for common changes
- ✅ **ZMK → QMK**: 80%+ accuracy for translatable features
- ✅ **ZMK → Kinesis2**: 70%+ accuracy for basic remapping
- ✅ **PR Automation**: 95%+ success rate for PR creation
- ✅ **Manual Intervention**: Only needed for complex/new features (<10% of cases)

### **Test Strategy:**
- **Unit Tests**: Parser and generator components
- **Integration Tests**: Full sync workflows
- **Regression Tests**: Ensure existing configs don't break
- **User Acceptance Tests**: Real-world usage scenarios

---

## 📝 Decision Log

### **Technical Decisions Made:**
1. **✅ 2024-12-19**: Organized configs into `configs/` directory
2. **✅ 2024-12-19**: Created reusable download script for all 4 keyboards
3. **✅ 2024-12-19**: Prioritized ZMK-to-ZMK compatibility as primary goal

### **Technical Decisions Pending:**
- [ ] **Programming Language Choice**: Go, Python, or Rust for CLI tool?
- [ ] **Internal Representation Format**: JSON, YAML, or custom format?
- [ ] **LLM Integration Method**: API calls, local model, or external prompts?
- [ ] **Testing Framework**: Language-specific or cross-platform?

---

## 🔧 Development Environment

### **Required Tools:**
- Git (for repository management)
- GitHub CLI (`gh`) for PR automation
- Language-specific tools (TBD based on language choice)
- ZMK compilation tools (for validation)
- QMK compilation tools (for validation)

### **Development Setup:**
```bash
# Current setup
git clone https://github.com/[your-repo]/keyboard_layout_config_mapper
cd keyboard_layout_config_mapper

# Build the CLI tool
go build -o klcm cmd/klcm/main.go

# Download configurations and start using
./klcm download
./klcm validate --all
./klcm diff adv360 glove80
klcm setup    # Initialize configuration
```

---

## 🚀 Next Immediate Steps

### **Current Focus** (Next 1-2 Sessions):
1. **📝 Deep Analysis Phase**: 
   - Analyze existing ZMK configurations in detail
   - Create physical key mapping matrices
   - Document current layer structures and behaviors used

2. **🔨 Choose Technical Stack**:
   - Decide on programming language for CLI tool
   - Set up development environment
   - Create basic project structure

3. **🔨 Build First Parser**:
   - Start with ZMK parser (highest priority)
   - Create tests for parsing existing configs
   - Design internal representation format

---

## 📚 References & Resources

### **ZMK Documentation:**
- [ZMK Firmware Documentation](https://zmk.dev/)
- [ZMK Keycodes](https://zmk.dev/docs/codes)
- [ZMK Behaviors](https://zmk.dev/docs/behaviors)

### **QMK Documentation:**
- [QMK Firmware Documentation](https://docs.qmk.fm/)
- [QMK Keycodes](https://docs.qmk.fm/#/keycodes)

### **Kinesis Documentation:**
- [Kinesis Advantage2 Manual](https://kinesis-ergo.com/support/advantage2/)

---

*This document is actively maintained and updated as the project progresses. All checkboxes (✅/❌/📝/🔨) indicate current status and should be updated regularly.*