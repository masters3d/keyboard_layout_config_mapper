package keyboard_config

const rowCount = 7
const columnHalfCount = 7
const columnFullCount = columnHalfCount * 2
const halfKeyboardKeyCount = columnHalfCount * rowCount
const fullKeyboardKeyCount = halfKeyboardKeyCount * 2

type KeycodeLayerHalf = [halfKeyboardKeyCount]KeyCodeRepresentable
type KeycodeLayerFull = [fullKeyboardKeyCount]KeyCodeRepresentable

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
		6 + (columnFullCount * 3), 7 + (columnFullCount * 3), // g h row. surprisingly the ErgodoxPrexy does not map anything on the middle here but intead it does to the next row
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
