# KLCM - Keyboard Layout Configuration Mapper

A unified CLI tool to synchronize keyboard configurations across different firmware systems (ZMK, QMK, Kinesis). Priority focus on ZMK keyboards (Advantage360 & Glove80) with automatic sync, change detection, and GitHub PR automation.

## Quick Start

### Installation
```bash
# Build the CLI tool
go build -o klcm cmd/klcm/main.go

# Or build and install to PATH
go install ./cmd/klcm
```

### Installation

**Prerequisites:** Go 1.21+ required

**Build the CLI tool:**
```bash
# Clone the repository
git clone https://github.com/yourusername/keyboard_layout_config_mapper.git
cd keyboard_layout_config_mapper

# Build the CLI tool
go build -o klcm cmd/klcm/main.go

# Verify installation
./klcm --help
```

### Usage
```bash
# Download all keyboard configurations
./klcm download

# Compare local changes with remote before downloading  
./klcm compare-remote

# Download with preview (shows changes and asks confirmation)
./klcm download --preview

# Validate configurations
./klcm validate --all

# Compare keyboards  
./klcm diff adv360 glove80

# Preview sync changes with detailed diffs
./klcm sync --preview

# Compare local configs with remote versions  
./klcm compare-remote

# Preview download changes before applying
./klcm download --preview

# Apply sync changes
./klcm sync --from adv360 --to glove80

# Get help for any command
./klcm --help
./klcm sync --help
```

## Supported Keyboards

- **ZMK Advantage360**: `configs/zmk_adv360/adv360.keymap`
- **ZMK Glove80**: `configs/zmk_glove80/glove80.keymap`  
- **QMK ErgoDox**: `configs/qmk_ergodox/keymap.c`
- **Kinesis Advantage 2**: `configs/kinesis2/1_qwerty.txt`

## Features

- 🔄 **Automatic Sync**: Bidirectional configuration synchronization
- 🔍 **Smart Validation**: Syntax checking and compatibility analysis
- 📊 **Layout Comparison**: Detailed diff reports between keyboards
- 📥 **GitHub Integration**: Download configs and create PRs automatically
- 🛡️ **Safe Operations**: Preview mode and change detection
- 🔍 **Remote Comparison**: Compare local files with remote versions before downloading
- 🎯 **Priority Support**: Optimized for ZMK keyboards (Advantage360 ↔ Glove80)
- ⚡ **Git-Style Diffs**: Professional unified diff output with context lines and change highlighting

### Enhanced Diff Functionality

All comparison operations now include git-style unified diffs showing:

- **Line-by-line changes** with proper context (3 lines before/after)
- **Professional formatting** with `+` (additions) and `-` (removals)  
- **Chunk headers** showing exact line ranges affected
- **Summary statistics** displaying impact (e.g., "+5 -2 ~3 lines")
- **Truncated output** for large changes to keep output manageable
- **Color-coded display** for improved readability

```bash
# Example diff output
@@ -182,4 +182,4 @@
         };
 
     };
-};// Local change
+};

📈 Summary: ~1 lines
```

## Layout

We are going to base this layout on the ergodox but it also needs to overlap with glove80 and the kinesis360 as i have these on preorder (as of 2022) :). 


ROWS:

For this we are going to add a sixth row when usually we would only need five rows since the glove80 has a function row. In addition I also have a kinesis advantage 2 which includes all the function keys. To support thumb clusters we are going to add another row at the bottom for a total of seven rows. 


COLUMNS:
Initially there was going to be 8x2=16 columns since we need both sides. Six for the keys themselves, one column for the additional function keys that are supported by ergodox and kinesses 360, and one column to accomodate the thumb cluster overflow.
I have seen decided that seven would be enough for our usecase which the advantage of this giving us a square grid. 

Total amount keys 7x7x2= 98

## Keycodes

These are all the keycodes that are compatible with the HID spec
https://github.com/qmk/qmk_firmware/blob/master/quantum/keycode.h


Additional alias which builds on top of keycode. These is there unshifted keys are defined
https://github.com/qmk/qmk_firmware/blob/master/quantum/quantum_keycodes.h

See the HID defined ones used hex. Fox example define KC_A has a HID of hex 0x04

https://github.com/qmk/libqmk/blob/master/include/qmk/keycodes/basic.h
https://github.com/qmk/qmk_firmware/blob/master/docs/keycodes.md
https://github.com/qmk/qmk_firmware/blob/master/docs/keycodes_basic.md


These are from ZMK
https://github.com/zmkfirmware/zmk/blob/main/app/include/dt-bindings/zmk/keys.h
https://github.com/zmkfirmware/zmk/blob/main/app/include/dt-bindings/zmk/hid_usage.h

These are not part of the HID spec but are used by qmk to depict different features.
KC_NO                  == 0x0000
KC_TRANSPARENT         == 0x0001



Some other cool keyboards:

http://xahlee.info/kbd/fancy_keyboards.html
https://www.reddit.com/r/MechanicalKeyboards/comments/91wwse/gmk_9009_ortho_just_arrived_today/?st=JK219TUF&sh=297dd81e
https://www.allthingsergo.com/dyi-ergonomic-keyboard/
https://trulyergonomic.com/ergonomic-keyboards/best-truly-ergonomic-mechanical-keyboard/
https://x-bows.com/collections/keyboards/products/x-bows-lite-ergonomic-mechanical-keyboard // The light is not using open source software to program
https://x-bows.com/products/x-bows-nature-ergonomic-mechanical-keyboard // this uses qmc
https://www.reddit.com/r/MechanicalKeyboards/comments/dk9b34/tractyl_split_keyboard_with_trackball/?utm_source=ifttt
https://drop.com/buy/x-bows-knight-plus-ergonomic-mechanical-keyboard
https://www.maltron.com/store/p20/Maltron_L90_dual_hand_fully_ergonomic_%283D%29_keyboard_-_US_English.html
https://geekhack.org/index.php?topic=46015.0 // Maniform

Very cool ready made keyboards
https://bastardkb.com/

Almost Cool:
https://www.zergotech.com/products/zergotech-freedom-mechanical-ergonomic-keyboard
https://www.pcbway.com/project/shareproject/ErgoDox_Creation___Infinity_ErgoDox_Mod.html

Some old keyboards:
http://xahlee.info/kbd/nec_ergofit_keyboard.html
http://xahlee.info/kbd/Truly_Ergonomic_keyboard.html
http://xahlee.info/kbd/ergonomic_keyboard_history_index.html
