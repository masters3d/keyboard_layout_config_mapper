package keyboard_config

import (
	"testing"
)

var mainLayerLeft = KeycodeLayerHalf{
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, // function keys
	ST_MACRO_Screenshot, LSFT(KC_1), KC_GRAVE, KC_MINUS, KC_EQUAL, KC_SLASH, KC_TRANSPARENT, // number keys
	KC_TRANSPARENT, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT,
	KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_Z, KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT,
	KC_TRANSPARENT, TO(0), TO(1), TO(2), TO(3), MO(2), KC_TRANSPARENT,
	KC_TRANSPARENT, KC_LEFT_SHIFT, KC_LEFT_GUI, KC_LEFT_ALT, KC_LEFT_CTRL, KC_RIGHT_GUI, KC_TRANSPARENT,
}
var mainLayerRight = KeycodeLayerHalf{
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, // function keys
	KC_TRANSPARENT, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, KC_LEFT_PAREN, KC_RIGHT_PAREN, KC_TRANSPARENT, // number keys
	KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_H, KC_J, KC_K, KC_L, KC_DOT, KC_ENTER,
	KC_TRANSPARENT, KC_N, KC_M, KC_SEMICOLON, KC_QUOTE, KC_COMMA, KC_TRANSPARENT,
	KC_TRANSPARENT, MO(1), KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_ESCAPE, KC_DELETE, KC_TAB, KC_BACKSPACE, KC_SPACE, KC_TRANSPARENT,
}

var keyboardFullValidationSet = KeycodeLayerFull{
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
	ST_MACRO_Screenshot, LSFT(KC_1), KC_GRAVE, KC_MINUS, KC_EQUAL, KC_SLASH, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, LSFT(KC_9), LSFT(KC_0), KC_TRANSPARENT,
	KC_TRANSPARENT, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT,
	KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_H, KC_J, KC_K, KC_L, KC_DOT, KC_ENTER,
	KC_TRANSPARENT, KC_Z, KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_N, KC_M, KC_SEMICOLON, KC_QUOTE, KC_COMMA, KC_TRANSPARENT,
	KC_TRANSPARENT, TO(0), TO(1), TO(2), TO(3), MO(2), KC_TRANSPARENT,
	KC_TRANSPARENT, MO(1), KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_LEFT_SHIFT, KC_LEFT_GUI, KC_LEFT_ALT, KC_LEFT_CTRL, KC_RIGHT_GUI, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_ESCAPE, KC_DELETE, KC_TAB, KC_BACKSPACE, KC_SPACE, KC_TRANSPARENT,
}

func Test_Merging_Half_Layers(t *testing.T) {

	//fmt.Println("#####_________________________________________________#####")
	var actualFullLayer = mergeHalfs(mainLayerLeft, mainLayerRight)

	// for index, value := range actualFullLayer {
	// 	if (index)%rowCount == 0 && index != 0 {
	// 		fmt.Println("")
	// 	}
	// 	fmt.Print(value, ", ")
	// }
	// fmt.Println("")
	// fmt.Println("#####        end of Test_Merging_Half_Layers       #####")

	for index, expectedValue := range keyboardFullValidationSet {
		actual := actualFullLayer[index]
		if expectedValue != actual {
			var message = "actual:" + actual.String() + ", expected:" + expectedValue.String()
			t.Error(message)
		}
	}

}
