package keyboard_config

// return []KeyCodeRepresentable{
// 	ST_MACRO_Screenshot, LSFT(KC_1), KC_GRAVE, KC_MINUS, KC_EQUAL, KC_SLASH, KC_TRANSPARENT, KC_TRANSPARENT, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, KC_LEFT_PAREN, KC_RIGHT_PAREN, KC_TRANSPARENT,
// 	KC_TRANSPARENT, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT, KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT,
// 	KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_H, KC_J, KC_K, KC_L, KC_DOT, KC_ENTER,
// 	KC_TRANSPARENT, KC_Z, KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT, KC_TRANSPARENT, KC_N, KC_M, KC_SEMICOLON, KC_QUOTE, KC_COMMA, KC_TRANSPARENT,
// 	KC_TRANSPARENT, TO(0), TO(1), TO(2), TO(3), KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, KC_TRANSPARENT,
// 	MO(2), KC_RIGHT_GUI, KC_ESCAPE, MO(1),
// 	KC_LEFT_CTRL, KC_DELETE,
// 	KC_LEFT_SHIFT, KC_LEFT_GUI, KC_LEFT_ALT, KC_TAB, KC_BACKSPACE, KC_SPACE,
// }

func convertLayerToErgodoxPrexy(input KeycodeLayerFull) []KeyCodeRepresentable {

	var collectedArray = []KeyCodeRepresentable{}

	indexesToIgnore := IntSlice{
		6 + (columnFullCount * 3), 7 + (columnFullCount * 3), // g h row. surprisingly the ErgodoxPrexy does not map anything on the middle here but instead it does to the next row
		5 + (columnFullCount * 5), 6 + (columnFullCount * 5), 7 + (columnFullCount * 5), 8 + (columnFullCount * 5), // arrow row. We need to skip two values here
	}
	firstRowCount := 1 // we will be starting at this row index. firstfunction layer is skipped since Ergodox does not have a function layer
	lastRowCount := 1  //we will skipping these many rows from the end
	for indexRow := firstRowCount; indexRow < rowCount-lastRowCount; indexRow++ {

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
		5 + (columnFullCount * 5), 5 + (columnFullCount * 6), // left side first two keys on thumb
		8 + (columnFullCount * 6), 8 + (columnFullCount * 5), // right side first two keys on thumb; The order here is swapped
		4 + (columnFullCount * 6),                                                       // left side second row single key
		9 + (columnFullCount * 6),                                                       // right side second row single key
		1 + (columnFullCount * 6), 2 + (columnFullCount * 6), 3 + (columnFullCount * 6), //left side bottom row single key
		10 + (columnFullCount * 6), 11 + (columnFullCount * 6), 12 + (columnFullCount * 6), //right side bottom row single key
	}

	for _, indexAsValue := range indexesForThumbCluster {
		currentValue := input[indexAsValue]
		collectedArray = append(collectedArray, currentValue)
	}

	return collectedArray
}
