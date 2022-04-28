package keyboard_config

import (
	"testing"
)

func Test_Creating_KeyPad_Descriptions(t *testing.T) {

	expectedValue := "kp-escape"
	_, actual := KinesisKeypadLayerMapping(KC_ESCAPE)
	if expectedValue != actual.tokenname {
		var message = "actual:" + actual.tokenname + ", expected:" + expectedValue
		t.Error(message)
	}
}
func Test_Creating_KeyPad_Descriptions_FullArray(t *testing.T) {

	var fullLayerTargetAdv2Default = mergeHalfs(Adv2TopLayerLeft, Adv2TopLayerRight)

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
