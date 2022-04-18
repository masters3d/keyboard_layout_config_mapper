package keyboard_config

import "strconv"

//go:generate stringer -type=macrocode
type macrocode uint

const (
	PLACEHOLDER macrocode = iota
	ST_MACRO_Screenshot
	ST_MACRO_Anglebrakets
	ST_MACRO_Parenthesis
	ST_MACRO_SquareBraces
	ST_MACRO_CurlyBraces
)

//go:generate stringer -type=keycode
type keycode uint

const (
	KC_NO          keycode = iota // 0x00
	KC_TRANSPARENT                // KC_TRANSPARENT and KC_ROLL_OVER share the same code. KC_ROLL_OVER not used but
	KC_POST_FAIL                  // not used
	KC_UNDEFINED                  // not used
	KC_A
	KC_B
	KC_C
	KC_D
	KC_E
	KC_F
	KC_G
	KC_H
	KC_I
	KC_J
	KC_K
	KC_L
	KC_M // 0x10
	KC_N
	KC_O
	KC_P
	KC_Q
	KC_R
	KC_S
	KC_T
	KC_U
	KC_V
	KC_W
	KC_X
	KC_Y
	KC_Z
	KC_1
	KC_2
	KC_3 // 0x20
	KC_4
	KC_5
	KC_6
	KC_7
	KC_8
	KC_9
	KC_0
	KC_ENTER
	KC_ESCAPE
	KC_BACKSPACE
	KC_TAB
	KC_SPACE
	KC_MINUS
	KC_EQUAL
	KC_LEFT_BRACKET
	KC_RIGHT_BRACKET // 0x30
	KC_BACKSLASH
	KC_NONUS_HASH
	KC_SEMICOLON
	KC_QUOTE
	KC_GRAVE
	KC_COMMA
	KC_DOT
	KC_SLASH
	KC_CAPS_LOCK
	KC_F1
	KC_F2
	KC_F3
	KC_F4
	KC_F5
	KC_F6
	KC_F7 // 0x40
	KC_F8
	KC_F9
	KC_F10
	KC_F11
	KC_F12
	KC_PRINT_SCREEN
	KC_SCROLL_LOCK
	KC_PAUSE
	KC_INSERT
	KC_HOME
	KC_PAGE_UP
	KC_DELETE
	KC_END
	KC_PAGE_DOWN
	KC_RIGHT
	KC_LEFT // 0x50
	KC_DOWN
	KC_UP
	KC_NUM_LOCK
	KC_KP_SLASH
	KC_KP_ASTERISK
	KC_KP_MINUS
	KC_KP_PLUS
	KC_KP_ENTER
	KC_KP_1
	KC_KP_2
	KC_KP_3
	KC_KP_4
	KC_KP_5
	KC_KP_6
	KC_KP_7
	KC_KP_8 // 0x60
	KC_KP_9
	KC_KP_0
	KC_KP_DOT
	KC_NONUS_BACKSLASH
	KC_APPLICATION
	KC_KB_POWER
	KC_KP_EQUAL
	KC_F13
	KC_F14
	KC_F15
	KC_F16
	KC_F17
	KC_F18
	KC_F19
	KC_F20
	KC_F21 // 0x70
	KC_F22
	KC_F23
	KC_F24
	KC_EXECUTE
	KC_HELP
	KC_MENU
	KC_SELECT
	KC_STOP
	KC_AGAIN
	KC_UNDO
	KC_CUT
	KC_COPY
	KC_PASTE
	KC_FIND
	KC_KB_MUTE
	KC_KB_VOLUME_UP // 0x80
	KC_KB_VOLUME_DOWN
	KC_LOCKING_CAPS_LOCK
	KC_LOCKING_NUM_LOCK
	KC_LOCKING_SCROLL_LOCK
	KC_KP_COMMA
	KC_KP_EQUAL_AS400
	KC_INTERNATIONAL_1
	KC_INTERNATIONAL_2
	KC_INTERNATIONAL_3
	KC_INTERNATIONAL_4
	KC_INTERNATIONAL_5
	KC_INTERNATIONAL_6
	KC_INTERNATIONAL_7
	KC_INTERNATIONAL_8
	KC_INTERNATIONAL_9
	KC_LANGUAGE_1 // 0x90
	KC_LANGUAGE_2
	KC_LANGUAGE_3
	KC_LANGUAGE_4
	KC_LANGUAGE_5
	KC_LANGUAGE_6
	KC_LANGUAGE_7
	KC_LANGUAGE_8
	KC_LANGUAGE_9
	KC_ALTERNATE_ERASE
	KC_SYSTEM_REQUEST
	KC_CANCEL
	KC_CLEAR
	KC_PRIOR
	KC_RETURN
	KC_SEPARATOR
	KC_OUT // 0xA0
	KC_OPER
	KC_CLEAR_AGAIN
	KC_CRSEL
	KC_EXSEL
)

const (
	/* Modifiers */
	KC_LEFT_CTRL keycode = iota + 0xE0 //0xE0
	KC_LEFT_SHIFT
	KC_LEFT_ALT
	KC_LEFT_GUI
	KC_RIGHT_CTRL
	KC_RIGHT_SHIFT
	KC_RIGHT_ALT
	KC_RIGHT_GUI
)

const (
	/* Mouse Buttons */
	KC_MS_UP keycode = iota + 0xED // 0xED
	KC_MS_DOWN
	KC_MS_LEFT
	KC_MS_RIGHT // 0xF0
	KC_MS_BTN1
	KC_MS_BTN2
	KC_MS_BTN3
	KC_MS_BTN4
	KC_MS_BTN5
	KC_MS_BTN6
	KC_MS_BTN7
	/* Mouse Wheel */
	KC_MS_WH_UP
	KC_MS_WH_DOWN
	KC_MS_WH_LEFT
	KC_MS_WH_RIGHT
	/* Acceleration */
	KC_MS_ACCEL0
	KC_MS_ACCEL1
	KC_MS_ACCEL2 // 0xFF
)

type KeyCodeRepresentable interface {
	String() string
}

// Using this value to separate the values in different grouping so combos are unique
type markercombovalue uint

// The max value of regular keycodes
const max_value_for_keycode markercombovalue = 0xFF // see KC_MS_ACCEL2

var LayerShifedKeys markercombovalue = max_value_for_keycode * 2 // used for shifted keys.
var LayerSwitchKeys markercombovalue = max_value_for_keycode * 3 // used for switching layer
var LayerToggleKeys markercombovalue = max_value_for_keycode * 4 // used for toggleing a layer

func (i markercombovalue) String() string {
	switch {
	case i == 0xFF*2:
		return "LSFT"
	case i == 0xFF*3:
		return "TO"
	case i == 0xFF*4:
		return "MO"
	default:
		return "makercombovalue(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

type keycombo struct {
	value keycode
	combo markercombovalue
}

func (self keycombo) String() string {
	return self.combo.String() + "(" + self.value.String() + ")"
}

type layercombo struct {
	value int
	combo markercombovalue
}

func (self layercombo) String() string {
	return self.combo.String() + "(" + strconv.FormatInt(int64(self.value), 10) + ")"
}

func LSFT(i keycode) keycombo {
	return keycombo{i, LayerShifedKeys}
}

func TO(i int) layercombo {
	return layercombo{i, LayerSwitchKeys}
}

func MO(i int) layercombo {
	return layercombo{i, LayerToggleKeys}
}

// US ANSI shifted keycode aliases

var KC_TILDE = LSFT(KC_GRAVE)                     // ~
var KC_EXCLAIM = LSFT(KC_1)                       // !
var KC_AT = LSFT(KC_2)                            // @
var KC_HASH = LSFT(KC_3)                          // #
var KC_DOLLAR = LSFT(KC_4)                        // $
var KC_PERCENT = LSFT(KC_5)                       // %
var KC_CIRCUMFLEX = LSFT(KC_6)                    // ^
var KC_AMPERSAND = LSFT(KC_7)                     // &
var KC_ASTERISK = LSFT(KC_8)                      // *
var KC_LEFT_PAREN = LSFT(KC_9)                    // (
var KC_RIGHT_PAREN = LSFT(KC_0)                   // )
var KC_UNDERSCORE = LSFT(KC_MINUS)                // _
var KC_PLUS = LSFT(KC_EQUAL)                      // +
var KC_LEFT_CURLY_BRACE = LSFT(KC_LEFT_BRACKET)   // {
var KC_RIGHT_CURLY_BRACE = LSFT(KC_RIGHT_BRACKET) // }
var KC_LEFT_ANGLE_BRACKET = LSFT(KC_COMMA)        // <
var KC_RIGHT_ANGLE_BRACKET = LSFT(KC_DOT)         // >
var KC_COLON = LSFT(KC_SEMICOLON)                 // :
var KC_PIPE = LSFT(KC_BACKSLASH)                  // |
var KC_QUESTION = LSFT(KC_SLASH)                  // ?
var KC_DOUBLE_QUOTE = LSFT(KC_QUOTE)              // "

// KINESIS SmartSet
// https://cs.github.com/KinesisCorporation/Freestyle-Edge-Pro-SmartSet-App

var kinesisAdv2_TypeableChars = map[KeyCodeRepresentable]string{
	KC_A:         "A",
	KC_B:         "B",
	KC_C:         "C",
	KC_D:         "D",
	KC_E:         "E",
	KC_F:         "F",
	KC_G:         "G",
	KC_H:         "H",
	KC_I:         "I",
	KC_J:         "J",
	KC_K:         "K",
	KC_L:         "L",
	KC_M:         "M",
	KC_N:         "N",
	KC_O:         "O",
	KC_P:         "P",
	KC_Q:         "Q",
	KC_R:         "R",
	KC_S:         "S",
	KC_T:         "T",
	KC_U:         "U",
	KC_V:         "V",
	KC_W:         "W",
	KC_X:         "X",
	KC_Y:         "Y",
	KC_Z:         "Z",
	KC_1:         "1",
	KC_2:         "2",
	KC_3:         "3",
	KC_4:         "4",
	KC_5:         "5",
	KC_6:         "6",
	KC_7:         "7",
	KC_8:         "8",
	KC_9:         "9",
	KC_0:         "0",
	KC_EQUAL:     "=",
	KC_GRAVE:     "`",
	KC_COMMA:     ",",
	KC_DOT:       ".",
	KC_SLASH:     "/",
	KC_BACKSLASH: `\`, // this needs to be escaped
	KC_QUOTE:     "'",
	KC_SEMICOLON: ";",
}

var kinesisAdv2_Keypad = map[KeyCodeRepresentable]string{
	KC_NUM_LOCK:    "numlk",
	KC_KP_SLASH:    "kpdiv",
	KC_KP_ASTERISK: "kpmult",
	KC_KP_MINUS:    "kpmin",
	KC_KP_PLUS:     "kpplus",
	KC_KP_ENTER:    "kpenter1",
	KC_RETURN:      "kpenter2", // not sure about this mapping
	KC_KP_1:        "kp1",
	KC_KP_2:        "kp2",
	KC_KP_3:        "kp3",
	KC_KP_4:        "kp4",
	KC_KP_5:        "kp5",
	KC_KP_6:        "kp6",
	KC_KP_7:        "kp7",
	KC_KP_8:        "kp8",
	KC_KP_9:        "kp9",
	KC_KP_0:        "kp0",
	KC_KP_DOT:      "kp.",
	KC_KP_EQUAL:    "kp=",
	KC_KP_COMMA:    "kp,", // This is not part of the kinesis codebase
}

var kinesisAdv2_Confirmed = map[KeyCodeRepresentable]string{
	KC_F1:            "F1",
	KC_F2:            "F2",
	KC_F3:            "F3",
	KC_F4:            "F4",
	KC_F5:            "F5",
	KC_F6:            "F6",
	KC_F7:            "F7",
	KC_F8:            "F8",
	KC_F9:            "F9",
	KC_F10:           "F10",
	KC_F11:           "F11",
	KC_F12:           "F12",
	KC_F13:           "F13",
	KC_F14:           "F14",
	KC_F15:           "F15",
	KC_F16:           "F16",
	KC_F17:           "F17",
	KC_F18:           "F18",
	KC_F19:           "F19",
	KC_F20:           "F20",
	KC_F21:           "F21",
	KC_F22:           "F22",
	KC_F23:           "F23",
	KC_F24:           "F24",
	KC_TRANSPARENT:   "null", // maybe we need a different sigil for null
	KC_ESCAPE:        "escape",
	KC_ENTER:         "enter",
	KC_CAPS_LOCK:     "caps",
	KC_TAB:           "tab",
	KC_LEFT_CTRL:     "lctrl",
	KC_LEFT_SHIFT:    "lshift",
	KC_LEFT_ALT:      "lalt",
	KC_LEFT_GUI:      "lwin",
	KC_RIGHT_CTRL:    "rctrl",
	KC_RIGHT_SHIFT:   "rshift",
	KC_RIGHT_ALT:     "ralt",
	KC_RIGHT_GUI:     "rwin",
	KC_MINUS:         "HYPHEN",
	KC_SPACE:         "SPACE",
	KC_BACKSPACE:     "BSPACE",
	KC_PRINT_SCREEN:  "prtscr",
	KC_LEFT_BRACKET:  "oBRACK",
	KC_RIGHT_BRACKET: "cBRACk",
	KC_PAUSE:         "PAUSE",
	KC_HOME:          "HOME",
	KC_PAGE_UP:       "pup",
	KC_PAGE_DOWN:     "pdown",
	KC_SCROLL_LOCK:   "SCROLL",
	KC_END:           "END",
	KC_DELETE:        "DELETE",
	KC_RIGHT:         "RIGHT",
	KC_LEFT:          "LEFT",
	KC_DOWN:          "DOWN",
	KC_UP:            "UP",
	KC_INSERT:        "INSERT", // This is only on the keypad layer and it labeled as kp-insert
}

// Kinesis Mapping
var kinesisAdv2_NotUsedOrConfirmed = map[KeyCodeRepresentable]string{

	KC_NONUS_HASH:          "NONUS_HASH",
	KC_NONUS_BACKSLASH:     "NONUS_BACKSLASH",
	KC_APPLICATION:         "APPLICATION",
	KC_KB_POWER:            "KB_POWER",
	KC_EXECUTE:             "EXECUTE",
	KC_HELP:                "HELP",
	KC_MENU:                "MENU",
	KC_SELECT:              "SELECT",
	KC_STOP:                "STOP",
	KC_AGAIN:               "AGAIN",
	KC_UNDO:                "UNDO",
	KC_CUT:                 "CUT",
	KC_COPY:                "COPY",
	KC_PASTE:               "PASTE",
	KC_FIND:                "FIND",
	KC_KB_MUTE:             "KB_MUTE",
	KC_KB_VOLUME_UP:        "KB_VOLUME_UP",
	KC_KB_VOLUME_DOWN:      "KB_VOLUME_DOWN",
	KC_LOCKING_CAPS_LOCK:   "LOCKING_CAPS_LOCK",
	KC_LOCKING_NUM_LOCK:    "LOCKING_NUM_LOCK",
	KC_LOCKING_SCROLL_LOCK: "LOCKING_SCROLL_LOCK",
	KC_KP_EQUAL_AS400:      "KP_EQUAL_AS400",
	KC_INTERNATIONAL_1:     "INTERNATIONAL_1", //    ConfigKeys.Add(TKey.Create(VK_OEM_102, 'intl-\', '', 'intl\', 'intl-\', 'intl-\', true, true)); //International <> key between Left Shift and Z

	KC_INTERNATIONAL_2: "INTERNATIONAL_2",
	KC_INTERNATIONAL_3: "INTERNATIONAL_3",
	KC_INTERNATIONAL_4: "INTERNATIONAL_4",
	KC_INTERNATIONAL_5: "INTERNATIONAL_5",
	KC_INTERNATIONAL_6: "INTERNATIONAL_6",
	KC_INTERNATIONAL_7: "INTERNATIONAL_7",
	KC_INTERNATIONAL_8: "INTERNATIONAL_8",
	KC_INTERNATIONAL_9: "INTERNATIONAL_9",
	KC_LANGUAGE_1:      "LANGUAGE_1",
	KC_LANGUAGE_2:      "LANGUAGE_2",
	KC_LANGUAGE_3:      "LANGUAGE_3",
	KC_LANGUAGE_4:      "LANGUAGE_4",
	KC_LANGUAGE_5:      "LANGUAGE_5",
	KC_LANGUAGE_6:      "LANGUAGE_6",
	KC_LANGUAGE_7:      "LANGUAGE_7",
	KC_LANGUAGE_8:      "LANGUAGE_8",
	KC_LANGUAGE_9:      "LANGUAGE_9",
	KC_ALTERNATE_ERASE: "ALTERNATE_ERASE",
	KC_SYSTEM_REQUEST:  "SYSTEM_REQUEST",
	KC_CANCEL:          "CANCEL",
	KC_CLEAR:           "CLEAR",
	KC_PRIOR:           "PRIOR",
	KC_SEPARATOR:       "SEPARATOR",
	KC_OUT:             "OUT",
	KC_OPER:            "OPER",
	KC_CLEAR_AGAIN:     "CLEAR_AGAIN",
	KC_CRSEL:           "CRSEL",
	KC_EXSEL:           "EXSEL",

	/* Mouse Buttons */
	KC_MS_UP:    "MS_UP",
	KC_MS_DOWN:  "MS_DOWN",
	KC_MS_LEFT:  "MS_LEFT",
	KC_MS_RIGHT: "MS_RIGHT",
	KC_MS_BTN1:  "MS_BTN1",
	KC_MS_BTN2:  "MS_BTN2",
	KC_MS_BTN3:  "MS_BTN3",
	KC_MS_BTN4:  "MS_BTN4",
	KC_MS_BTN5:  "MS_BTN5",
	KC_MS_BTN6:  "MS_BTN6",
	KC_MS_BTN7:  "MS_BTN7",
	/* Mouse Wheel */
	KC_MS_WH_UP:    "MS_WH_UP",
	KC_MS_WH_DOWN:  "MS_WH_DOWN",
	KC_MS_WH_LEFT:  "MS_WH_LEFT",
	KC_MS_WH_RIGHT: "MS_WH_RIGHT",
	/* Acceleration */
	KC_MS_ACCEL0: "MS_ACCEL0",
	KC_MS_ACCEL1: "MS_ACCEL1",
	KC_MS_ACCEL2: "MS_ACCEL2",
}
