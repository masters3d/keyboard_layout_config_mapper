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
