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

func Test_Creating_KeyPad_DescriptionsArray(t *testing.T) {

	for index, keycode_each := range Adv2TopLayerLeft {

		println(keycode_each, index)
		keycode_each_asString := keycode_each.String()
		println(keycode_each_asString)

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
