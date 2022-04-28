package keyboard_config

import (
	"testing"
)

func Test_Ergodox_Array(t *testing.T) {

	var expectedValue = []KeyCodeRepresentable{
		ST_MACRO_Screenshot, LSFT(KC_1), KC_GRAVE, KC_MINUS, KC_EQUAL, KC_SLASH, KC_TRANSPARENT, KC_TRANSPARENT, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, KC_LEFT_PAREN, KC_RIGHT_PAREN, KC_TRANSPARENT,
		KC_TRANSPARENT, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT, KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT,
		KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_H, KC_J, KC_K, KC_L, KC_DOT, KC_ENTER,
		KC_TRANSPARENT, KC_Z, KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT, KC_TRANSPARENT, KC_N, KC_M, KC_SEMICOLON, KC_QUOTE, KC_COMMA, KC_TRANSPARENT,
		KC_TRANSPARENT, TO(0), TO(1), TO(2), TO(3), KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, KC_TRANSPARENT,
		MO(2), KC_RIGHT_GUI, KC_ESCAPE, MO(1),
		KC_LEFT_CTRL, KC_DELETE,
		KC_LEFT_SHIFT, KC_LEFT_GUI, KC_LEFT_ALT, KC_TAB, KC_BACKSPACE, KC_SPACE,
	}

	actual := convertLayerToErgodoxPrexy(keyboardFullValidationSet)
	for index, each_actual_keycode := range actual {

		var expected_keycode_from_ergodox = expectedValue[index]

		if each_actual_keycode != expected_keycode_from_ergodox {
			var message = ":`" + each_actual_keycode.String() + "` should match `" + expected_keycode_from_ergodox.String() + "`"
			t.Error(message)
		}
	}

}

func Test_layout_ergodox_pretty(t *testing.T) {

	expected_layout_ergodox_pretty := []KeyCodeRepresentable{
		ST_MACRO_Screenshot, LSFT(KC_1), KC_GRAVE, KC_MINUS, KC_EQUAL, KC_SLASH, KC_TRANSPARENT, KC_TRANSPARENT, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, KC_LEFT_PAREN, KC_RIGHT_PAREN, KC_TRANSPARENT,
		KC_TRANSPARENT, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT, KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT,
		KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_H, KC_J, KC_K, KC_L, KC_DOT, KC_ENTER,
		KC_TRANSPARENT, KC_Z, KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT, KC_TRANSPARENT, KC_N, KC_M, KC_SEMICOLON, KC_QUOTE, KC_COMMA, KC_TRANSPARENT,
		KC_TRANSPARENT, TO(0), TO(1), TO(2), TO(3), KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, KC_TRANSPARENT,
		MO(1), KC_RIGHT_GUI, KC_ESCAPE, MO(2),
		KC_LEFT_CTRL, KC_DELETE,
		KC_LEFT_SHIFT, KC_LEFT_GUI, KC_LEFT_ALT, KC_TAB, KC_BACKSPACE, KC_SPACE,
	}
	got := convertLayerToErgodoxPrexy(keyboardFullValidationSet)

	for index, expectedValue := range expected_layout_ergodox_pretty {
		actual := got[index]
		if expectedValue != actual {
			var message = "actual:" + actual.String() + ", expected:" + expectedValue.String()
			t.Error(message)
		}
	}
}