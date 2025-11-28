# Thumb Cluster Mapping Documentation

Last Updated: 2025-01-28

## Physical Layout

All three keyboards (adv360, adv_mod/pillzmod_pro, glove80) share the same 6-key thumb cluster arrangement per hand.

```
        Left Thumb Cluster              Right Thumb Cluster
        (Left hand)                     (Right hand)
    ┌───┬───┬───┐                      ┌───┬───┬───┐
    │ 6 │ 5 │ 4 │                      │ 4 │ 5 │ 6 │
    ├───┼───┼───┤                      ├───┼───┼───┤
    │ 1 │ 2 │ 3 │                      │ 3 │ 2 │ 1 │
    └───┴───┴───┘                      └───┴───┴───┘
   (thumb)    (far)                  (far)    (thumb)
```

**Numbering Convention:**
- Start from the key closest to your thumb (position 1)
- Bottom row: 1, 2, 3 (thumb to far)
- Top row: 4, 5, 6 (directly above 3, 2, 1 respectively)
- So 4 is above 3 (farthest), 5 is above 2 (middle), 6 is above 1 (closest to thumb)

## Logical Mapping (Unified Layout)

| Position | Key Binding | Notes |
|----------|-------------|-------|
| **L1** | `SPACE` | Primary space (left thumb, closest) |
| **L2** | `LEFT_SHIFT` | Left shift modifier |
| **L3** | `LEFT_ALT` | Alt/Option modifier (farthest from thumb) |
| **L4** | `LEFT_CONTROL` | Control modifier (above L3) |
| **L5** | `LEFT_WIN` | Win/GUI/Command modifier (above L2) |
| **L6** | `mo KEYPAD` | Momentary keypad layer (above L1, closest to thumb) |
| **R1** | `SPACE` | Primary space (right thumb, closest) |
| **R2** | `RIGHT_SHIFT` | Right shift modifier |
| **R3** | `mo LAYER_CMD` | Momentary command layer (farthest from thumb) |
| **R4** | `LEFT_CONTROL` | Control modifier (above R3) |
| **R5** | `ESCAPE` | Escape key (above R2) |
| **R6** | `mo KEYPAD` | Momentary keypad layer (above R1, closest to thumb) |

## macOS Karabiner Configuration

On macOS, use Karabiner-Elements to remap:
- **Right Control → Left Command**: Makes `RC()` bindings in ZMK act as Cmd on macOS
- **Screenshot shortcut**: Override to match Windows `Ctrl+Alt+Del` style screenshot command

This allows the same ZMK keymap to work consistently across Windows and macOS:
- On Windows: `RC()` = Right Control (native)
- On macOS: `RC()` = Left Command (via Karabiner remap)

## Parser Position Reference

For ZMK keymap parsing, thumb cluster keys appear in specific positions within the binding rows.
Each keyboard has a different matrix layout, but the logical positions map as follows:

### adv360 (7-column layout with extra inner columns)
```
Row 3 (homerow level): ... [to 0] [L6:KEYPAD] [L5:WIN] ... [R5:ESC] [R6:KEYPAD] [to 0] ...
Row 4 (below homerow): ...        [L4:CTRL]            ... [R4:CTRL]                  ...
Row 5 (bottom):        ...        [L1:SPC] [L2:LSHFT] [L3:LALT] | [R3:CMD] [R2:RSHFT] [R1:SPC] ...
```

### adv_mod/pillzmod_pro (Kinesis Advantage matrix)
```
Row 6 (after arrows): [L6:KEYPAD] [L5:WIN]            ... [R5:ESC] [R6:KEYPAD]
Row 7 (single key):              [L4:CTRL]            ... [R4:CTRL]
Row 8 (bottom):       [L1:SPC] [L2:LSHFT] [L3:LALT]   ... [R3:CMD] [R2:RSHFT] [R1:SPC]
```

### glove80 (6-column main + 6 thumb inline)
```
Row 4 (Z row): ... [L6:KEYPAD] [L5:WIN] [L4:CTRL] [R4:CTRL] [R5:ESC] [R6:KEYPAD] ...
Row 5 (bottom): ... [L1:SPC] [L2:LSHFT] [L3:LALT] | [R3:CMD] [R2:RSHFT] [R1:SPC] ...
```

## Historical Differences (Now Unified)

Prior to unification, R4 had different bindings:
- adv360/glove80: `TAB`
- adv_mod: `LEFT_CONTROL`

All keyboards now use `LEFT_CONTROL` on R4 for consistency.
