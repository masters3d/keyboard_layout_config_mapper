package keyboard_config

import (
	"strconv"
	"strings"
	"testing"
)

func Test_Kinesis_Insert_Data_Into_TokenAreas(t *testing.T) {

	rangeToIgnore := " 	\n"
	expectedDoc :=
		`
*#def KC_CAPS_LOCK# caps, kp-caps
[caps]>[escape]
[kp-caps]>[escape]
*#end KC_CAPS_LOCK#
*ignore me 
*#def KC_QUOTE# apos, kp-apos
[apos]>[enter]
[kp-apos]>[enter]
*#end KC_QUOTE#
	`
	inputDoc :=
		`
*#def KC_CAPS_LOCK# caps, kp-caps
*#end KC_CAPS_LOCK#
*ignore me 
*#def KC_QUOTE# apos, kp-apos
*#end KC_QUOTE#
`

	actualDoc := Kinesis_ParseAndFill_SpecialTokens(inputDoc, keyboardFullValidationSet, keyboardFullValidationSet)

	if strings.Trim(actualDoc, rangeToIgnore) != strings.Trim(expectedDoc, rangeToIgnore) {
		var message = "actual:`" + actualDoc + "`, expected:" + expectedDoc + "`"
		t.Error(message)
	}

}

func Test_Kinesis_Parse_SpecialToken(t *testing.T) {

	inputDoc :=
		`
		*#def KC_TAB# tab,kp-tab
HELLO FROM THE OTHER SIDE
		*#end KC_TAB#
		*I should be ignored
		*#def VK_MPEDAL# mp-kpshf, kp-mp-kpshf
		HELLO FROM THE OTHER SIDE
HELLO FROM THE OTHER SIDE
		HELLO FROM THE OTHER SIDE
		*#end VK_MPEDAL#
	`

	var actualMap = Kinesis_Parse_SpecialToken(inputDoc)

	// KC_TAB
	actualTab := actualMap[KC_TAB]
	expectedTabRange := Token_index_range{start: 27, end: 56}
	expectedTabMessage := "HELLO FROM THE OTHER SIDE"
	actualTabMessage := strings.Trim(inputDoc[actualTab.start:actualTab.end], " 	\n")
	if expectedTabMessage != actualTabMessage {
		var message = KC_VK_MPEDAL.String() + "actual:`" + actualTabMessage + "`, expected:" + expectedTabMessage
		t.Error(message)
	}
	if expectedTabRange.start != actualTab.start {
		var message = KC_TAB.String() + " Range Start actual:" + strconv.FormatInt(int64(actualTab.start), 10) + ", expected:" + strconv.FormatInt(int64(expectedTabRange.start), 10)
		t.Error(message)
	}
	if expectedTabRange.end != actualTab.end {
		var message = KC_TAB.String() + " Range end actual:" + strconv.FormatInt(int64(actualTab.end), 10) + ", expected:" + strconv.FormatInt(int64(expectedTabRange.end), 10)
		t.Error(message)
	}

	// VK_MPEDAL
	actualVK_MPEDAL := actualMap[KC_VK_MPEDAL]
	expectedVK_MPEDALRange := Token_index_range{start: 133, end: 218}
	expectedVK_MPEDALMessage := `HELLO FROM THE OTHER SIDE
HELLO FROM THE OTHER SIDE
		HELLO FROM THE OTHER SIDE`

	actualVK_MPEDALMessage := strings.Trim(inputDoc[actualVK_MPEDAL.start:actualVK_MPEDAL.end], " 	\n")
	if expectedVK_MPEDALMessage != actualVK_MPEDALMessage {
		var message = KC_VK_MPEDAL.String() + "actual:`" + actualVK_MPEDALMessage + "`, expected:" + expectedVK_MPEDALMessage
		t.Error(message)
	}
	if expectedVK_MPEDALRange.start != actualVK_MPEDAL.start {
		var message = KC_VK_MPEDAL.String() + " Range Start actual:" + strconv.FormatInt(int64(actualVK_MPEDAL.start), 10) + ", expected:" + strconv.FormatInt(int64(expectedVK_MPEDALRange.start), 10)
		t.Error(message)
	}
	if expectedVK_MPEDALRange.end != actualVK_MPEDAL.end {
		var message = KC_VK_MPEDAL.String() + " Range End actual:" + strconv.FormatInt(int64(actualVK_MPEDAL.end), 10) + ", expected:" + strconv.FormatInt(int64(expectedVK_MPEDALRange.end), 10)
		t.Error(message)
	}
}

func Test_Creating_KeyPad_Descriptions(t *testing.T) {

	expectedValue := "kp-escape"
	_, actual := KinesisKeypadLayerMapping(KC_ESCAPE)
	if expectedValue != actual.tokenname {
		var message = "actual:" + actual.tokenname + ", expected:" + expectedValue
		t.Error(message)
	}
}

func Test_Creating_KeyPad_Descriptions_FullArray(t *testing.T) {

	var fullLayerTargetAdv2Default = MergeHalfLayers(Adv2TopLayerLeft, Adv2TopLayerRight)

	for index, keycode_each_target := range fullLayerTargetAdv2Default {

		expected := Adv2KeypadValidation[index]

		_, value := KinesisKeypadLayerMapping(keycode_each_target)
		got := value.tokenname

		if expected != got {
			var message = ":`" + keycode_each_target.String() + "` should map to `" + expected + "` but instead we got `" + got + "`"
			t.Error(message)
		}

	}

}

func Test_Creating_KeyPad_DescriptionsArray(t *testing.T) {

	for _, keycode_each := range Adv2TopLayerLeft {

		keycode_each_asString := keycode_each.String()

		if keycode_each_asString == "KC_TRANSPARENT" {
			continue // we don't have mapping for transparent
		}
		_, isOk := kinesisAdv2ndLayerMapping[keycode_each]

		if !isOk {
			var message = "expected:" + keycode_each.String() + " to have a valid keypad layer mapping"
			t.Error(message)
		}

	}
}
