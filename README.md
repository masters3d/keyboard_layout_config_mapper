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

## Experimental modeless editor controls

**Learn the keyboard gestures once; install an adapter for each editor.** This is
an opt-in vertical slice, not complete Vim parity. No firmware is flashed, no
editor configuration is installed automatically, and the default keyboard
behavior is unchanged.

The [versioned action contract](configs/editor/actions.json) defines 36 signals,
insertion-position semantics, selection, undo, failure behavior, and deferred
capabilities. ZMK only sends signals: it never tracks editor mode or emits a
document-editing keystroke script. The editor owns text, history, and prompts.

### Backend decision and capability audit

The programmable reference adapter is **Neovim**. **Zed native** is an explicitly
limited second adapter, not a hidden switch into Zed's Vim mode. Consult its
[per-action capability matrix](configs/editor/zed/capabilities.json) before use:
`native` means a native action exists, `different` identifies differing semantics,
and `unsupported` means the signal is intentionally blocked.

The Zed audit inspected upstream source at
[`d62802d456ac93a5abce5fecd5e98d2b8124485d`](https://github.com/zed-industries/zed/tree/d62802d456ac93a5abce5fecd5e98d2b8124485d),
not a running stable release. Check the installed release's action picker and
keymap diagnostics; source availability is not runtime verification.

| Capability family | Zed finding | Delivery status |
|---|---|---|
| Character, word, line, document navigation and selection | Native actions; boundary, word-end, and display-line semantics can differ | Supported subset with explicit differences |
| Paragraph movement | Native start/end of nonblank-line paragraph actions | Both adapters use their editor's paragraph rules |
| Inside/around delimiters | Native `SelectInsideDelimiters` / `SelectAroundDelimiters`; not typed objects | Typed-parentheses change is Neovim-only in this pilot |
| Operators and arbitrary counts | Vim-state actions, not interchangeable native commands | Deferred; fixed three-word deletion in Neovim |
| Repeat change | `vim::Repeat` depends on modal machinery | Neovim repeats adapter deletions, not inserted replacement text |
| Registers and macros | Vim-state register/record/replay actions | Deferred; pilot deletions preserve registers |
| Search and navigation history | Native actions | Both adapters, with editor-specific prompt handling |
| Definition navigation | Native language-service action | Both adapters; requires a working language service |
| General text objects, syntax navigation, view controls, workspace and other code actions | Mixed native/Vim capabilities requiring individual contracts | Deferred, not implicitly claimed |

Audit sources:

- [Native default bindings](https://github.com/zed-industries/zed/blob/d62802d456ac93a5abce5fecd5e98d2b8124485d/assets/keymaps/default-linux.json)
- [Action definitions](https://github.com/zed-industries/zed/blob/d62802d456ac93a5abce5fecd5e98d2b8124485d/crates/editor/src/actions.rs)
- [Vim bindings and contexts](https://github.com/zed-industries/zed/blob/d62802d456ac93a5abce5fecd5e98d2b8124485d/assets/keymaps/vim.json)
- [Keybinding documentation](https://github.com/zed-industries/zed/blob/d62802d456ac93a5abce5fecd5e98d2b8124485d/docs/src/key-bindings.md)
- [Public extension API](https://github.com/zed-industries/zed/blob/d62802d456ac93a5abce5fecd5e98d2b8124485d/crates/extension_api/src/extension_api.rs)

The inspected extension API does not expose general buffer-edit callbacks,
keyboard interception, or arbitrary editor-action registration. A conventional
Zed extension cannot be assumed to fill the gaps. Native forward-word-end is
**not** Vim next-word-start, and redo is **not** repeat-change. Neither is silently
substituted here.

### Physical language and signal allocation

All three banks use F1–F12 with explicit modifiers, leaving the existing plain
F13–F24 Keypad assignments alone. Signals still need host-level conflict testing.
**These are candidate signals, not universally safe shortcuts:** Linux desktops
may reserve Ctrl+Alt+function keys for virtual-console switching. Resolve such
OS bindings or choose another verified transport before enabling the pilot.

| Bank | Signal | Intended interaction |
|---|---|---|
| NAV | Ctrl+Alt+F1–F12 | Move insertion position |
| SELECT | Ctrl+Alt+Shift+F1–F12 | Extend anchored selection using the same target positions |
| EDIT | Ctrl+Shift+F1–F12 | Complete action or visible search prompt |

| Slot | NAV / SELECT physical key | Target | EDIT physical key | Action |
|---|---|---|---|---|
| F1 | H | Left | U | Undo |
| F2 | J | Down | R | Redo |
| F3 | K | Up | W | Delete to next word start |
| F4 | L | Right | E | Delete to third next word start |
| F5 | B | Previous word start | I | Change inside parentheses |
| F6 | W | Next word start | period | Repeat adapter edit |
| F7 | A | Absolute line start | slash | Open search |
| F8 | E | Line end | N | Next search match |
| F9 | G | Document start | P | Previous search match |
| F10 | T | Document end | D | Definition |
| F11 | U | Previous paragraph | B | Navigation history back |
| F12 | O | Next paragraph | upper-left Escape | Cancel selection / resume insertion |

The function row provides all twelve slots directly in each pilot layer.
Ordinary letters always type outside held layers, except when entering text in
an explicit search prompt. Selection persists visibly after layer release;
typing replaces it. Search, selection, and future counts are deliberate
interactions, not invisible operator-pending states.

### Install and verify adapters before enabling firmware

1. Back up the editor configuration. Keep ordinary Escape and normal editor
   navigation available.
2. For Neovim, add the following to your Lua configuration, adjusting the absolute
   checkout location when installing elsewhere. Use `require`, not a bare
   `dofile`: the adapter's queued continuation resolves the cached `klcm` module.
   `insert_on_enter` opts ordinary editable buffers into Insert mode on entry;
   leave it false to manage initial insertion yourself. Escape still deliberately
   exits Insert mode. When loading interactively, enter Insert mode once with `i`.
   The adapter's `teardown()` restores its mappings; do not install a blanket
   autocmd forcing Insert mode whenever it is left.

   ```lua
   package.path = package.path .. ";/home/runner/work/keyboard_layout_config_mapper/keyboard_layout_config_mapper/configs/editor/neovim/?.lua"
   local klcm = require("klcm")
   klcm.setup({ insert_on_enter = true })
   ```

   The Lua adapter is exercised on Neovim 0.9.5. Earlier versions are unverified.
   Advanced users can provide `setup({ mappings = { ... } })`: that table replaces,
   rather than merges with, default action-ID-to-shortcut mappings. Keep it in
   sync with the firmware transport. `dispatch(action_id)` exposes the same
   operations for editor-side customization.
3. For Zed, turn Vim and Helix modes off. Merge the objects from the
   [adapter keymap](configs/editor/zed/keymap.json) into the user keymap array;
   **do not overwrite existing bindings**. Resolve conflicting shortcuts and
   check the capability matrix. Unsupported signals are inert, not approximated.
   The native H/J/K/L bindings are explicitly **different**: horizontal movement
   crosses lines and vertical movement follows display rows. Cancel may dismiss
   a popup before collapsing selection. These are useful native controls, not
   passing reference-parity implementations. The search bar retains native
   Enter/Escape and focus behavior; semantic signals resume in the document.
4. Verify each modified function key arrives distinctly, first directly and
   then through any terminal, tmux, SSH, OS remapper, or remote desktop in use.
   A terminal without suitable extended-key support may not distinguish these
   banks. Do not flash until transport works; remap the banks consistently on
   both sides if the host intercepts them.
5. Test movement, selection replacement, undo, search confirm/cancel, and
   definition/back in disposable buffers. Never test destructive bindings for
   the first time in unsaved work.

### Enable the Pillz Mod pilot

Add `#define KLCM_EDITOR_ENABLE` before the layer definitions in the keymap used
by the firmware build. Do not assume a similarly named Kconfig option exists.
Without that preprocessor definition, the four original layers compile to their
original bindings.

The experimental access combos exist **only in the enabled pilot**:

- Press left Space + left Keypad together (within 50 ms) to hold NAV.
- Press right Space + right Keypad together (within 50 ms) to hold EDIT.
- Access combos start on the default layer only. Release **both** combo keys to
  leave that layer; after activation, either one can keep it held.
- While NAV/SELECT is held, right Space provides direct momentary EDIT access.
- While NAV is active, hold either physical Shift position for SELECT.
- EDIT has higher priority than SELECT, which has higher priority than NAV.
- The right-thumb Escape remains a literal emergency Escape in all three layers.
- Unassigned pilot positions are blocked. The raw default layer is unchanged as
  well as its compiled bindings, so KLCM's non-preprocessing parser still sees it.
  Existing layers remain unchanged, but
  pilot layers mask their keys while active. All three pedal positions inherit.

Activate NAV **before** its selection key. Holding Shift first sends a real host
modifier, not a selection-layer request. Release layer keys and modifiers before
switching banks; extra held modifiers can change a signal. Combo timing is only a
prototype and needs physical typing/repetition tests. Missing the combo window
falls back to the existing Space/Keypad behavior. The enabled pilot may briefly
delay participating keys while waiting for a possible combo.

The access mechanism uses standard [ZMK combos](https://github.com/zmkfirmware/zmk/blob/main/docs/docs/keymaps/combos.md)
with zero-based positions and `slow-release`; no custom firmware behavior is needed.

Remove the definition and restore known-good firmware to roll back. The enabled
matrix has 89 bindings per layer, including all three pedals; do not use historical
86-key assumptions for this keymap.

### Validation and rollout gates

The Go tests use the existing Go toolchain: contract uniqueness, matching
navigation/selection gestures, Zed capability/binding consistency, and C
preprocessing of both pilot states. Preprocessing uses empty ZMK headers to check
guards and layout structure; it is **not** a firmware compilation. Original
disabled-layer fingerprints protect default behavior.

Run the existing Go tooling and the adapter's dependency-free headless assertions:

```sh
cd /home/runner/work/keyboard_layout_config_mapper/keyboard_layout_config_mapper
go test ./...
go vet ./...
go build -o /tmp/klcm-editor ./cmd/klcm
/tmp/klcm-editor validate --all
nvim --headless -u NONE -i NONE -n -c "luafile /home/runner/work/keyboard_layout_config_mapper/keyboard_layout_config_mapper/configs/editor/neovim/klcm_test.lua"
```

Neovim character operations use UTF-8 codepoints, not grapheme clusters or visual
columns. Parenthesis operations are literal balanced-delimiter operations, not a
language parser. Combining marks, strings/comments containing parentheses, and
multicursor editing do not have full Vim/native-editor parity.
The Neovim search prompt accepts Vim regular expressions, but not native search
offsets or chained searches. Submitted queries run through protected editor APIs
so invalid or missing matches cannot flush queued typing into Normal-mode edits.
Existing buffer-local shortcuts take precedence over adapter mappings; resolve
such conflicts explicitly instead of silently replacing plugin bindings.

KLCM's `validate --all` is a syntax check. Its current `--compile` option merely
reports that local compilation is unimplemented; a real firmware build is still
required in the target firmware repository.

Before calling this portable, run the same scenarios in both editors:

| Scenario | Required observation |
|---|---|
| Empty line, end of file, Unicode text | Valid insertion boundaries; no text corruption |
| Repeated word/line navigation | Target and boundary behavior agree, or a documented gap blocks parity |
| Select forward/backward, then type | Visible selection and exactly one replacement |
| Delete one/three words; undo/redo | Expected text and separate history units |
| Change inside nested/missing parentheses | Delimiters preserved; missing target leaves text untouched |
| Repeat at a second location | Repeat the operation, not an old absolute text range |
| Search confirm/cancel and next/previous | No shortcut text inserted; deliberate prompt lifecycle |
| Definition and return, with/without LSP | Safe failure or correct history navigation |
| Combo timing, reversed release order, USB/Bluetooth | No stuck layer, key, or modifier |
| Existing CMD/Keypad/System and pedals | Existing functions remain available outside the pilot |

The complete demanding slice is a Neovim target; Zed's unsupported actions are
explicit portability failures, not passing tests. No live Zed, terminal transport,
firmware-build, or physical-keyboard results are implied by static checks.

Adv360/Glove80 rollout is intentionally gated on the Pillz Mod hardware and
transport checks. Their layer indices and matrices differ. KLCM's current sync
command copies default layers; it must not be used to propagate these new layers
or replace entire hardware layouts. Subsequent work should add capabilities from
the contract's deferred list individually, with action contracts and adapter tests,
before assigning additional physical gestures.

---

**🎉 Happy keyboard customizing!**