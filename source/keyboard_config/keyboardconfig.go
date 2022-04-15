package keyboard_config

import "fmt"

const rowCount = 7
const columnHalfCount = 7
const columnFullCount = columnHalfCount * 2
const halfKeyboardKeyCount = columnHalfCount * rowCount
const fullKeyboardKeyCount = halfKeyboardKeyCount * 2

type KeycodeLayerHalf = [halfKeyboardKeyCount]KeyCodeRepresentable
type KeycodeLayerFull = [fullKeyboardKeyCount]KeyCodeRepresentable

var mainLayerLeft = KeycodeLayerHalf{
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, // function keys
	ST_MACRO_Screenshot, LSFT(KC_1), KC_GRAVE, KC_MINUS, KC_EQUAL, KC_SLASH, KC_TRANSPARENT, // number keys
	KC_TRANSPARENT, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT,
	KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_Z, KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT,
	KC_TRANSPARENT, TO(0), TO(1), TO(2), TO(3), MO(1), KC_TRANSPARENT,
	KC_TRANSPARENT, KC_LEFT_SHIFT, KC_LEFT_GUI, KC_LEFT_ALT, KC_LEFT_CTRL, KC_RIGHT_GUI, KC_TRANSPARENT,
}
var mainLayerRight = KeycodeLayerHalf{
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, // function keys
	KC_TRANSPARENT, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, KC_LEFT_PAREN, KC_RIGHT_PAREN, KC_TRANSPARENT, // number keys
	KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_H, KC_J, KC_K, KC_L, KC_DOT, KC_ENTER,
	KC_TRANSPARENT, KC_N, KC_M, KC_SEMICOLON, KC_QUOTE, KC_COMMA, KC_TRANSPARENT,
	KC_TRANSPARENT, MO(2), KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_ESCAPE, KC_DELETE, KC_TAB, KC_BACKSPACE, KC_SPACE, KC_TRANSPARENT,
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

type IntSlice []int

func (self IntSlice) Contains(i int) bool {
	for _, v := range self {
		if v == i {
			return true
		}
	}
	return false
}

func convertLayerToErgodoxPrexy(input KeycodeLayerFull) []KeyCodeRepresentable {

	var collectedArray = []KeyCodeRepresentable{}

	indexesToIgnore := IntSlice{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, //function layer since Ergodox does not have a function layer
		6 + (14 * 3), 7 + (14 * 3), // g h row. surprisingly the ErgodoxPrexy does not map anything on the middle here but intead it does to the next row
		5 + (14 * 5), 6 + (14 * 5), 7 + (14 * 5), 8 + (14 * 5), // arrow row. We need to skip two values here
	}
	lastRowCount := 1 //we will skipping these many rows from the end
	for indexRow := 0; indexRow < rowCount-lastRowCount; indexRow++ {

		for indexColumn := 0; indexColumn < columnFullCount; indexColumn++ {
			currentIndex := indexColumn + (indexRow * columnFullCount)
			if indexesToIgnore.Contains(currentIndex) {
				continue
			}
			currentValue := input[currentIndex]
			collectedArray = append(collectedArray, currentValue)
		}
	}
	indexesForThumbCluster := IntSlice{
		5 + (14 * 5), 5 + (14 * 6), // left side first two keys on thumb
		8 + (14 * 6), 8 + (14 * 5), // right side first two keys on thumb; The order here is swapped
		4 + (14 * 6),                             // left side second row single key
		9 + (14 * 6),                             // right side second row single key
		1 + (14 * 6), 2 + (14 * 6), 3 + (14 * 6), //left side bottom row single key
		10 + (14 * 6), 11 + (14 * 6), 12 + (14 * 6), //right side bottom row single key
	}

	for _, indexAsValue := range indexesForThumbCluster {
		currentValue := input[indexAsValue]
		collectedArray = append(collectedArray, currentValue)
	}

	return collectedArray
}

func mergeHalfs(left KeycodeLayerHalf, right KeycodeLayerHalf) KeycodeLayerFull {

	// To be able to use append this needs to be a non fixed slice
	var collectedArray = []KeyCodeRepresentable{}

	for indexRow := 0; indexRow < rowCount; indexRow++ {
		// left hand
		for indexColumn := 0; indexColumn < columnHalfCount; indexColumn++ {
			currentHalfIndex := indexColumn + (indexRow * columnHalfCount)
			currentHalfValue := left[currentHalfIndex]
			collectedArray = append(collectedArray, currentHalfValue)
		}
		// right hand
		for indexColumn := 0; indexColumn < columnHalfCount; indexColumn++ {
			currentHalfIndex := indexColumn + (indexRow * columnHalfCount)
			currentHalfValue := right[currentHalfIndex]
			collectedArray = append(collectedArray, currentHalfValue)
		}
	}
	var final_output = KeycodeLayerFull{}

	copy(final_output[:], collectedArray[:fullKeyboardKeyCount])

	return final_output
}
