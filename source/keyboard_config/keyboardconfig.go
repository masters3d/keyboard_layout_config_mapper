package keyboard_config

import "fmt"

const rowCount = 6
const columnCount = 8
const halfKeyboardKeyCount = columnCount * rowCount
const fullKeyboardKeyCount = halfKeyboardKeyCount * 2

type KeycodeLayerHalf = [halfKeyboardKeyCount]KeyCodeRepresentable
type KeycodeLayerFull = [fullKeyboardKeyCount]KeyCodeRepresentable

var mainLayerLeft = KeycodeLayerHalf{
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, // function keys
	ST_MACRO_Screenshot, LSFT(KC_1), KC_GRAVE, KC_MINUS, KC_EQUAL, KC_SLASH, KC_TRANSPARENT, KC_TRANSPARENT, // number keys
	KC_TRANSPARENT, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT, KC_TRANSPARENT,
	KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_TRANSPARENT, KC_RIGHT_GUI,
	KC_TRANSPARENT, KC_Z, KC_X, KC_C, KC_V, KC_B, MO(1), KC_LEFT_CTRL,
	KC_TRANSPARENT, TO(0), TO(1), TO(2), TO(3), KC_LEFT_SHIFT, KC_LEFT_GUI, KC_LEFT_ALT,
}
var mainLayerRight = KeycodeLayerHalf{
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, // function keys
	KC_TRANSPARENT, KC_TRANSPARENT, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, KC_LEFT_PAREN, KC_RIGHT_PAREN, KC_TRANSPARENT, // number keys
	KC_TRANSPARENT, KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT,
	KC_ESCAPE, KC_TRANSPARENT, KC_H, KC_J, KC_K, KC_L, KC_DOT, KC_ENTER,
	KC_DELETE, MO(2), KC_N, KC_M, KC_SEMICOLON, KC_QUOTE, KC_COMMA, KC_TRANSPARENT,
	KC_TAB, KC_BACKSPACE, KC_SPACE, KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, KC_TRANSPARENT,
}

func MergeTest() {
	fmt.Println("#####        start        #####")
	for index, value := range mergeHalfs(mainLayerLeft, mainLayerRight) {
		if (index)%rowCount == 0 && index != 0 {
			fmt.Println("")
		}
		fmt.Print(value, ", ")
	}
	fmt.Println("")
	fmt.Println("#####        end        #####")
}

func convertLayerToErgodoxPrexy(input KeycodeLayerFull) []KeyCodeRepresentable {
	return input[:]
}

func mergeHalfs(left KeycodeLayerHalf, right KeycodeLayerHalf) KeycodeLayerFull {

	// To be able to use append this needs to be a non fixed slice
	var collectedArray = []KeyCodeRepresentable{}

	for indexRow := 0; indexRow < rowCount; indexRow++ {
		// left hand
		for indexColumn := 0; indexColumn < columnCount; indexColumn++ {
			currentHalfIndex := indexColumn + (indexRow * columnCount)
			currentHalfValue := left[currentHalfIndex]
			collectedArray = append(collectedArray, currentHalfValue)
		}
		// right hand
		for indexColumn := 0; indexColumn < columnCount; indexColumn++ {
			currentHalfIndex := indexColumn + (indexRow * columnCount)
			currentHalfValue := right[currentHalfIndex]
			collectedArray = append(collectedArray, currentHalfValue)
		}
	}
	var final_output = KeycodeLayerFull{}

	copy(final_output[:], collectedArray[:fullKeyboardKeyCount])

	return final_output
}
