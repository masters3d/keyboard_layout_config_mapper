# Navigation Layer Quick Reference

## Overview

The Navigation Layer provides comprehensive navigation and editing motions for modeless editing in code editors like Zed, VSCode, and others. This layer eliminates the need for arrow keys on the base layer and provides efficient access to advanced navigation commands.

## Activation

**Hold any of the 4 bottom-right positions** (previously arrow keys) to activate the navigation layer momentarily. Release to return to the default layer.

## Quick Reference Card

### Right Hand - Navigation

| Key | Function | Shortcut | Priority |
|-----|----------|----------|----------|
| **Y** | Go Back | `Ctrl + -` | ⭐⭐⭐ |
| **U** | Word Left | `Option + ←` | ⭐⭐⭐ |
| **I** | Up Arrow | `↑` | ⭐⭐⭐ |
| **O** | Word Right | `Option + →` | ⭐⭐⭐ |
| **P** | Go Forward | `Ctrl + Shift + -` | ⭐⭐ |
| | | | |
| **H** | Line Start | `Cmd + ←` | ⭐⭐⭐ |
| **J** | Left Arrow | `←` | ⭐⭐⭐ |
| **K** | Down Arrow | `↓` | ⭐⭐⭐ |
| **L** | Right Arrow | `→` | ⭐⭐⭐ |
| **;** | Line End | `Cmd + →` | ⭐⭐⭐ |
| | | | |
| **N** | Delete Word Left | `Option + Backspace` | ⭐⭐⭐ |
| **M** | Backspace | `Backspace` | ⭐⭐⭐ |
| **,** | Page Down | `Page Down` | ⭐⭐ |
| **.** | Page Up | `Page Up` | ⭐⭐ |
| **/** | Delete | `Delete` | ⭐⭐⭐ |

### Left Hand - Editing

| Key | Function | Shortcut | Priority |
|-----|----------|----------|----------|
| **Q** | Redo | `Cmd + Shift + Z` | ⭐⭐⭐ |
| **W** | Undo | `Cmd + Z` | ⭐⭐⭐ |
| **E** | Delete Line | `Cmd + Shift + K` | ⭐⭐ |
| **R** | Duplicate Line | `Cmd + Shift + D` | ⭐⭐ |
| | | | |
| **A** | Select All | `Cmd + A` | ⭐⭐ |
| **S** | Cut | `Cmd + X` | ⭐⭐⭐ |
| **D** | Copy | `Cmd + C` | ⭐⭐⭐ |
| **F** | Paste | `Cmd + V` | ⭐⭐⭐ |
| **G** | Select Line | `Cmd + L` | ⭐⭐⭐ |
| | | | |
| **Z** | Move Line Up | `Option + ↑` | ⭐⭐ |
| **X** | Move Line Down | `Option + ↓` | ⭐⭐ |
| **C** | Go To Definition | `F12` | ⭐⭐⭐ |
| **V** | Find References | `Shift + F12` | ⭐⭐ |
| **B** | Command Palette | `Cmd + Shift + P` | ⭐⭐⭐ |

## Common Workflows

### 1. Navigate to End of Line and Delete Word
1. Hold nav layer key (bottom-right)
2. Tap **L** (→) to move right OR tap **;** (End) to jump to end
3. Release nav layer key

### 2. Select Entire Line and Copy
1. Hold nav layer key
2. Tap **G** (Select line)
3. Tap **D** (Copy)
4. Release nav layer key

### 3. Jump to Definition and Return
1. Hold nav layer key
2. Tap **C** (F12 - Go to definition)
3. Release nav layer key
4. *(Review code at definition)*
5. Hold nav layer key again
6. Tap **Y** (Go back)
7. Release nav layer key

## Selection Mode

To select while navigating, **hold Shift** in addition to the nav layer key:

- `Shift + →` = Select character right
- `Shift + Word→` = Select word right  
- `Shift + End` = Select to line end

All navigation motions become selection motions when Shift is held.

## Design Principles

1. **Home row priority** - Most frequent actions on home row
2. **Vim-inspired** - IJKL forms inverted-T arrow cluster
3. **Symmetric access** - Navigation on right, editing on left
4. **No mode switching** - Momentary activation only

## Benefits

✅ **No hand movement** - Navigate without leaving home position  
✅ **Fast text editing** - Clipboard operations on left home row  
✅ **IDE integration** - Code navigation built-in  
✅ **No arrow key dependency** - Complete navigation without arrow cluster  
✅ **Vim-friendly** - IJKL arrow cluster familiar to Vim users  
✅ **Ergonomic** - Reduces hand strain  

## Keyboard Compatibility

- ✅ Kinesis Advantage 360 (adv360) - Layer 9
- ✅ MoErgo Glove80 (glove80) - Layer 9
- ✅ Kinesis Advantage with Pillz Mod Pro (pillzmod_pro) - Layer 4

---

*For complete documentation, see [AGENTS.md](../AGENTS.md)*
