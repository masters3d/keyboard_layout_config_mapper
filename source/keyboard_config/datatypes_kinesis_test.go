package keyboard_config

import (
	"testing"
)

func Test_Creating_KeyPad_Descriptions(t *testing.T) {

	expectedValue := "kp-escape"
	actual := KeyPadKinesis(KC_ESCAPE).tokenname
	if expectedValue != actual {
		var message = "actual:" + actual + ", expected:" + expectedValue
		t.Error(message)
	}

}
func Test_Creating_KeyPad_Descriptions_FullArray(t *testing.T) {

	var fullLayerTarget = mergeHalfs(mainLayerLeft, mainLayerRight)
	var fullLayerTargetAdv2Default = mergeHalfs(Adv2TopLayerLeft, Adv2TopLayerRight)

	for index, keycode_each_target := range fullLayerTarget {
		if keycode_each_target.String() == "KC_TRANSPARENT" {
			continue // we don't have mapping for transparent
		}

		var keycode_each_source_default = fullLayerTargetAdv2Default[index]

		if keycode_each_target != keycode_each_source_default {
			var message = ":`" + keycode_each_target.String() + "` needs mapping `" + keycode_each_source_default.String() + "`"
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
