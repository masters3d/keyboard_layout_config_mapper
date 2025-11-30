# 🎹 KLCM (Keyboard Layout Configuration Mapper)

A CLI tool for managing ZMK keyboard configurations across multiple keyboards (Advantage360, Glove80, Pillz Mod) with GitHub PR automation.

## ✨ Features

- **🔄 Pull configurations** from remote repositories
- **🔍 Compare local vs remote** configurations with git-style diffs
- **🔄 Sync changes** between ZMK keyboards
- **✅ Validate configurations** for syntax errors
- **🚀 GitHub PR automation** for contributing changes back
- **🎯 Interactive workflow** guide

## 🚀 Quick Start

```bash
# Build the tool
go build -o klcm cmd/klcm/main.go

# Interactive workflow (recommended)
./klcm workflow

# Or manual commands:
./klcm pull --preview           # Preview changes from remote
./klcm pull                     # Apply updates
./klcm validate                 # Validate all configs
./klcm sync adv360 glove80      # Sync between keyboards
./klcm pr create --dry-run      # Preview PR creation
./klcm pr create --apply        # Create actual PRs
```

## 📁 Supported Keyboards

| Keyboard | Config Path | Description |
|----------|-------------|-------------|
| `adv360` | `configs/zmk_adv360/adv360.keymap` | Kinesis Advantage360 Pro |
| `glove80` | `configs/zmk_glove80/glove80.keymap` | MoErgo Glove80 |
| `adv_mod` | `configs/zmk_adv_mod/pillzmod_pro.keymap` | Kinesis Advantage with Pillz Mod (Nice!Nano) |

## 🛠️ Commands

| Command | Description |
|---------|-------------|
| `pull` | Update local files from remote repos |
| `sync` | Copy changes between keyboards |
| `validate` | Check configurations for syntax errors |
| `compare-remote` | Compare local vs remote files |
| `download` | Download configurations |
| `pr create` | Create GitHub PRs for changes |
| `pr status` | Check status of PRs |
| `workflow` | Interactive guide |

## 📂 Project Structure

```
configs/
├── zmk_adv360/          # Advantage360 ZMK config
├── zmk_glove80/         # Glove80 ZMK config
├── zmk_adv_mod/         # Pillz Mod ZMK config
└── archived/            # Archived non-ZMK configs (kinesis2, qmk_ergodox)
```

## 🔗 Related Repositories

- [Adv360-Pro-ZMK](https://github.com/masters3d/Adv360-Pro-ZMK) - Advantage360 firmware
- [glove80-zmk-config](https://github.com/masters3d/glove80-zmk-config) - Glove80 firmware
- [zmk-config-pillzmod-nicenano](https://github.com/masters3d/zmk-config-pillzmod-nicenano) - Pillz Mod firmware

---

**🎉 Happy keyboard customizing!**