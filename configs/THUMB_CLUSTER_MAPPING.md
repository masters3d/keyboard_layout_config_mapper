# Thumb Cluster Mapping Documentation

Last Updated: 2025-11-28

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

## Logical Mapping (Unified Layout)

| Position | Key Binding | Notes |
|----------|-------------|-------|
| **L1** | `SPACE` | Primary space (left thumb) |
| **L2** | `LEFT_SHIFT` | Left shift modifier |
| **L3** | `LEFT_ALT` | Alt/Option modifier |
| **L4** | `LEFT_CONTROL` | Control modifier |
| **L5** | `LEFT_WIN` | Win/GUI/Command modifier |
| **L6** | `mo KEYPAD` | Momentary keypad layer |
| **R1** | `SPACE` | Primary space (right thumb) |
| **R2** | `RIGHT_SHIFT` | Right shift modifier |
| **R3** | `mo LAYER_CMD` | Momentary command layer |
| **R4** | `LEFT_CONTROL` | Control modifier |
| **R5** | `ESCAPE` | Escape key |
| **R6** | `mo KEYPAD` | Momentary keypad layer |

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
- Thumb keys span across rows 3-5 in the matrix
- Row 3 (homerow level): L6, L5 ... R5, R6
- Row 4 (below homerow): L4 ... R4
- Row 5 (bottom): L1, L2, L3 ... R3, R2, R1

### adv_mod/pillzmod_pro (Kinesis Advantage matrix)
- Thumb keys in dedicated thumb cluster rows
- Row after arrow keys: L6, L5 ... R5, R6
- Single key row: L4 ... R4
- Bottom row: L1, L2, L3 ... R3, R2, R1

### glove80 (6-column main + 6 thumb inline)
- Row 4 (Z row level): includes L6, L5, L4, R4, R5, R6 inline
- Row 5 (bottom): L1, L2, L3 ... R3, R2, R1

## Historical Differences (Now Unified)

Prior to unification, R4 had different bindings:
- adv360/glove80: `TAB`
- adv_mod: `LEFT_CONTROL`

All keyboards now use `LEFT_CONTROL` on R4 for consistency.
