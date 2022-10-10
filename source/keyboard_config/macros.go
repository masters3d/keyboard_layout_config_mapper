package keyboard_config

import "strconv"

// ST_MACRO_Screenshot

// macro_quotes: macro_quotes {
//     compatible = "zmk,behavior-macro";
//     label = "macro_quotes";
//     #binding-cells = <0>;
//     bindings = <&kp SQT>, <&kp SQT>, <&kp LEFT>;
//   };

// type macro_keycodes struct {
// 	label   string
// 	binding []KeyCodeRepresentable
// }

// func (self macro_keycodes) String() string {
// 	return self.label
// }

// var ST_MACRO_Screenshot = macro_keycodes{label: "ST_MACRO_Screenshot"}

type macrocode uint

const (
	PLACEHOLDER macrocode = iota
	ST_MACRO_Screenshot
	ST_MACRO_Anglebrakets
	ST_MACRO_Parenthesis
	ST_MACRO_SquareBraces
	ST_MACRO_CurlyBraces
	// these are defines
	KC_ControlAltDelete
)

func (i macrocode) String() string {
	switch {
	case i == ST_MACRO_Anglebrakets:
		return "ST_MACRO_Anglebrakets"
	case i == ST_MACRO_Screenshot:
		return "ST_MACRO_Screenshot"
	case i == ST_MACRO_Parenthesis:
		return "ST_MACRO_Parenthesis"
	case i == ST_MACRO_SquareBraces:
		return "ST_MACRO_SquareBraces"
	case i == ST_MACRO_CurlyBraces:
		return "ST_MACRO_CurlyBraces"
	case i == ST_MACRO_CurlyBraces:
		return "ST_MACRO_CurlyBraces"
	case i == KC_ControlAltDelete:
		return "KC_ControlAltDelete"
	default:
		return "MACRO_UNKNOWN(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
