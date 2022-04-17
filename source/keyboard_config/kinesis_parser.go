package keyboard_config

var Adv2TopLayerLeft = KeycodeLayerHalf{
	KC_F1, KC_F2, KC_F3, KC_F4, KC_F5, KC_F6, KC_F7, // function keys
	ST_MACRO_Screenshot, LSFT(KC_1), KC_GRAVE, KC_MINUS, KC_EQUAL, KC_SLASH, KC_ESCAPE, // This escape here is because the full 16 to key are too long for the width so we are going to staff the first and last on the inner side of the secord row
	KC_TRANSPARENT, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT,
	KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_Z, KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT,
	KC_TRANSPARENT, TO(0), TO(1), TO(2), TO(3), MO(1), KC_TRANSPARENT,
	KC_TRANSPARENT, KC_LEFT_SHIFT, KC_LEFT_GUI, KC_LEFT_ALT, KC_LEFT_CTRL, KC_RIGHT_GUI, KC_TRANSPARENT,
}
var Adv2TopLayerRight = KeycodeLayerHalf{
	KC_F8, KC_F9, KC_F10, KC_F11, KC_F12, KC_PRINT_SCREEN, KC_SCROLL_LOCK, // function keys
	KC_PAUSE, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, KC_LEFT_PAREN, KC_RIGHT_PAREN, KC_TRANSPARENT, // number keys
	KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_H, KC_J, KC_K, KC_L, KC_DOT, KC_ENTER,
	KC_TRANSPARENT, KC_N, KC_M, KC_SEMICOLON, KC_QUOTE, KC_COMMA, KC_TRANSPARENT,
	KC_TRANSPARENT, MO(2), KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_ESCAPE, KC_DELETE, KC_TAB, KC_BACKSPACE, KC_SPACE, KC_TRANSPARENT,
}
