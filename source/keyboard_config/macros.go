package keyboard_config

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
