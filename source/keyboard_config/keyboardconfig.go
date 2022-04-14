package keyboard_config

const rowCount = 5
const columnCount = 8
const splitKeyboardCount = 2

var FullKeyboard = [columnCount * rowCount * splitKeyboardCount]string{
	LSFT(KC_1).String(), KC_GRAVE.String(),
	// ST_MACRO_0.String(), LSFT(KC_1).String(), KC_GRAVE.String(), KC_MINUS.String(), KC_EQUAL.String(), KC_SLASH.String(), KC_TRANSPARENT.String(), KC_TRANSPARENT.String(), KC_BSLASH.String(), KC_LBRACKET.String(), KC_RBRACKET.String(), KC_LPRN.String(), KC_RPRN.String(), KC_TRANSPARENT.String(),
	// KC_TRANSPARENT.String(), KC_Q.String(), KC_W.String(), KC_E.String(), KC_R.String(), KC_T.String(), KC_TRANSPARENT.String(), KC_TRANSPARENT.String(), KC_Y.String(), KC_U.String(), KC_I.String(), KC_O.String(), KC_P.String(), KC_TRANSPARENT.String(),
	// KC_ESCAPE.String(), KC_A.String(), KC_S.String(), KC_D.String(), KC_F.String(), KC_G.String(), KC_H.String(), KC_J.String(), KC_K.String(), KC_L.String(), KC_DOT.String(), KC_ENTER.String(),
	// KC_TRANSPARENT.String(), KC_Z.String(), KC_X.String(), KC_C.String(), KC_V.String(), KC_B.String(), KC_TRANSPARENT.String(), KC_TRANSPARENT.String(), KC_N.String(), KC_M.String(), KC_SCOLON.String(), KC_QUOTE.String(), KC_COMMA.String(), KC_TRANSPARENT.String(),
	// KC_TRANSPARENT.String(), TO(0).String(), TO(1).String(), TO(2).String(), TO(3).String(), KC_LEFT.String(), KC_DOWN.String(), KC_UP.String(), KC_RIGHT.String(), KC_TRANSPARENT.String(),
	// MO(1).String(), KC_RIGHT_GUI.String(), KC_ESCAPE.String(), MO(2).String(),
	// KC_LEFT_CTRL.String(), KC_DELETE.String(),
	// KC_LEFT_SHIFT.String(), KC_LEFT_GUI.String(), KC_LEFT_ALT.String(), KC_TAB.String(), KC_BACKSPACE.String(), KC_SPACE.String(),
}
