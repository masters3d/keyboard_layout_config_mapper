# 🎹 KLCM (Keyboard Layout Configuration Mapper)

A powerful CLI tool for managing keyboard configuration files across different firmware systems (ZMK, QMK, Kinesis) with a robust intermediate representation system for seamless layout translation between keyboards.

## ✨ Features

- **🔄 Pull configurations** from remote repositories (like `git pull`)
- **🔍 Compare local vs remote** configurations with git-style diffs
- **🔄 Sync changes** between different keyboards (cross-vendor)
- **🌐 Universal layout translation** via intermediate representation (IR)
- **✅ Validate configurations** for syntax errors and compatibility
- **🚀 GitHub PR automation** for contributing changes back
- **🎯 Interactive workflow** guide for beginners
- **🛡️ Safe preview mode** for all operations

## 🚀 Quick Start

### **Option 1: Interactive Workflow (Recommended for Beginners)**
```bash
# Start the interactive guide that walks you through everything
klcm workflow
```

### **Option 2: Manual Commands (For Advanced Users)**
```bash
# 1. Update to latest configurations
klcm pull --preview              # Preview changes
klcm pull                        # Apply updates

# 2. Make your changes (edit .keymap files)

# 3. Validate your changes
klcm validate

# 4. Sync between keyboards (optional)
klcm sync adv360 glove80 --preview
klcm sync adv360 glove80

# 5. Create pull requests
klcm pr create --dry-run         # Preview PR creation
klcm pr create --apply           # Create actual PRs
```

## 📋 Complete Workflow Process

For detailed step-by-step instructions, see **[WORKFLOW.md](WORKFLOW.md)** - this contains the complete guide for making, validating, and contributing keyboard configuration changes.

## 🛠️ Commands

### **Core Commands**
| Command | Description | Example |
|---------|-------------|---------|
| `pull` | Update local files from remote (like git pull) | `klcm pull --preview` |
| `sync` | Copy changes between keyboards | `klcm sync adv360 glove80` |
| `translate` | Translate layouts between keyboards via IR | `klcm translate --from adv360 --to glove80` |
| `validate` | Check configurations for errors | `klcm validate adv360` |
| `compare-remote` | Compare local vs remote files | `klcm compare-remote` |

### **GitHub Integration**  
| Command | Description | Example |
|---------|-------------|---------|
| `pr create` | Create pull requests for your changes | `klcm pr create --dry-run` |
| `pr status` | Check status of your PRs | `klcm pr status` |
| `pr workflow` | Interactive PR workflow guide | `klcm pr workflow` |

### **Workflow Helpers**
| Command | Description | Example |
|---------|-------------|---------|
| `workflow` | Interactive guide for complete process | `klcm workflow` |
| `download` | Download specific configurations | `klcm download adv360` |

## 💡 Common Use Cases

### **Making Your First Change**
```bash
# Interactive guide walks you through everything
klcm workflow
```

### **Updating to Latest Remote Configurations**
```bash
# See what would change
klcm pull --preview

# Apply the updates
klcm pull
```

### **Syncing Same Layout to Both Keyboards**
```bash
# Edit configs/zmk_adv360/adv360.keymap
klcm validate adv360

# Copy changes to Glove80
klcm sync adv360 glove80 --preview
klcm sync adv360 glove80

# Validate both
klcm validate
```

### **Translating Layouts Between Keyboards**
```bash
# Translate from Advantage360 to Glove80 using IR system
klcm translate --from adv360 --to glove80

# Save to specific file
klcm translate --from glove80 --to adv_mod --output my_layout.keymap

# View the intermediate representation
klcm translate --from adv360 --show-ir
```

### **Contributing Changes Back**
```bash
# See what PRs would be created
klcm pr create --dry-run

# Create the actual PRs
klcm pr create --apply

# Check PR status later
klcm pr status
```

## 🔧 Installation

```bash
# Clone this repository
git clone <repository-url>
cd keyboard_layout_config_mapper

# Build the tool
go build -o klcm cmd/klcm/main.go

# Run commands
./klcm --help
```

## 📁 Supported Keyboards

### **ZMK Keyboards (Primary Focus)**
- **Advantage360 Pro** (`configs/zmk_adv360/`)
- **Glove80** (`configs/zmk_glove80/`)
- **Kinesis Advantage with Pillz Mod** (`configs/zmk_adv_mod/`)

### **Future Support**
- QMK keyboards
- Kinesis keyboards  
- Additional ZMK keyboards

## 🌐 Intermediate Representation (IR) System

KLCM includes a powerful intermediate representation system that enables seamless layout translation between different keyboard types.

### **Key Features**
- **Universal 10x10 grid per hand**: Standardized coordinate system for all keyboards
- **Zone-based mapping**: Logical grouping of keys (main, thumb, function, etc.)
- **Automatic translation**: Convert layouts between any supported keyboards
- **Key code normalization**: Handles firmware-specific syntax differences

### **Translation Examples**
```bash
# Convert Advantage360 layout to Glove80
klcm translate --from adv360 --to glove80

# Port Glove80 layout to Advanced Mod
klcm translate --from glove80 --to adv_mod

# Analyze layout structure
klcm translate --from adv360 --show-ir
```

### **Benefits**
- **Easy keyboard migration**: Keep your layout when switching keyboards
- **Cross-vendor compatibility**: Works between different manufacturers
- **Layout experimentation**: Try your layout on different keyboards
- **Consistent mapping**: Similar keys always map to the same positions

For detailed information, see **[docs/INTERMEDIATE_REPRESENTATION.md](docs/INTERMEDIATE_REPRESENTATION.md)**.

## 🎯 Key Benefits

### **🛡️ Safety First**
- Preview all changes before applying
- Interactive confirmations for destructive operations
- Git-style diffs show exactly what will change
- Validation catches errors before you commit

### **🚀 Productivity**
- Sync layouts between keyboards with one command
- Automated PR creation saves manual GitHub work
- Interactive workflow guides beginners
- Batch operations for multiple keyboards

### **🔧 Professional Workflow**
- Git-style commands (`pull`, `diff`, etc.)
- Proper branch management for PRs
- Descriptive commit messages
- Upstream contribution automation

## 🆘 Getting Help

### **Command Help**
```bash
klcm --help                    # General help
klcm pull --help              # Command-specific help
klcm pr create --help         # PR creation help
```

### **Verbose Output**
Add `--verbose` to any command for detailed information:
```bash
klcm pull --verbose --preview
klcm validate --verbose
```

### **Workflow Guide**
- **📖 Complete guide**: [WORKFLOW.md](WORKFLOW.md)
- **🎯 Interactive guide**: `klcm workflow`
- **🚀 PR guide**: `klcm pr workflow`

## 🤝 Contributing

1. Make your keyboard configuration changes
2. Run `klcm validate` to ensure they're correct
3. Use `klcm pr create` to submit pull requests
4. Monitor your PRs with `klcm pr status`

For contributing to this tool itself, see the standard GitHub workflow.

## 📝 License

[Add your license here]

---

**🎉 Happy keyboard customizing!** This tool makes it easy to manage configurations across multiple keyboards while safely contributing back to the community.