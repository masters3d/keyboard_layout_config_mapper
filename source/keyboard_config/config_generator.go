package keyboard_config

import (
	"log"
	"strconv"
	"strings"
)

func ConvertLayerToErgodoxPrexyAsString(input KeycodeLayerFull) string {

	var resultArray = convertLayerToErgodoxPrexy(input)

	var resultString = ""
	for index, value := range resultArray {
		if (index)%rowCount == 0 && index != 0 {
			resultString += "\n"
		}
		separator := ", "
		// if this is the last one we will not add a separator as
		if index == len(resultArray)-1 {
			// the qmk macro doesn't like trailing commas as they are considred as new item
			separator = ""
		}
		resultString += (value.String() + separator)
	}
	return resultString

}

func find_index_of_previous_line(text string, indexEnd int) int {
	const line = "\n"
	for index := indexEnd; index > 0 && index < len(text); index += -1 {
		value := string(text[index])
		if value == line {
			return index
		}
	}
	return -1
}

func find_index_of_next_line(text string, indexStart int) int {
	const line = "\n"
	for index := indexStart; index < len(text) && index > 0; index += 1 {
		value := string(text[index])
		if value == line {
			return index
		}
	}
	return -1
}

func Ergodox_replace_layer(template string, layer int, input KeycodeLayerFull) string {
	return ergodox_replace_layer_specific(template, layer, input, "[{layer}] = LAYOUT_ergodox_pretty(", "[{layer}] = GENERATED")
}

// Generate String
func ergodox_replace_layer_specific(template string, layer int, input KeycodeLayerFull, startPatternTemplate string, endPatternTemplate string) string {

	startPattern := strings.Replace(startPatternTemplate, "{layer}", strconv.FormatInt(int64(layer), 10), 1) //"[" + strconv.FormatInt(int64(layer), 10) + "] = LAYOUT_ergodox_pretty("
	endPattern := strings.Replace(endPatternTemplate, "{layer}", strconv.FormatInt(int64(layer), 10), 1)     //"[" + strconv.FormatInt(int64(layer), 10) + "] = GENERATED"

	log.Print("startPattern = `" + startPattern + "`")
	log.Print("endPattern = `" + endPattern + "`")

	startIndex := strings.Index(template, startPattern)
	endIndex := strings.Index(template, endPattern)

	log.Print("startIndex = " + strconv.FormatInt(int64(startIndex), 10))
	log.Print("endIndex = " + strconv.FormatInt(int64(endIndex), 10))

	startIndexNewLine := find_index_of_next_line(template, startIndex)
	endIndexNewLine := find_index_of_previous_line(template, endIndex)

	log.Print("startIndexNewLine = " + strconv.FormatInt(int64(startIndexNewLine), 10))
	log.Print("endIndexNewLine = " + strconv.FormatInt(int64(endIndexNewLine), 10))

	if startIndexNewLine == -1 || endIndexNewLine == -1 {
		panic("startIndexNewLine == " + strconv.FormatInt(int64(startIndexNewLine), 10) + " endIndexNewLine == " + strconv.FormatInt(int64(endIndexNewLine), 10))
	}

	layerAsString := ConvertLayerToErgodoxPrexyAsString(input)

	// the end of the slide range is not inclusive
	inclusiveStartIndexNewLine := startIndexNewLine + 1
	valueToReturn := template[:inclusiveStartIndexNewLine] + layerAsString + template[endIndexNewLine:]
	return strings.Trim(valueToReturn, " \n")

}

// Generate Array
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
