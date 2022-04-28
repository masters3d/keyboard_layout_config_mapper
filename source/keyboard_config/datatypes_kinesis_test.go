package keyboard_config

import (
	"fmt"
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

	var fullLayerTargetAdv2Default = mergeHalfs(Adv2TopLayerLeft, Adv2TopLayerRight)
	//var fullLayerTargetAdv2_Keypad = mergeHalfs(Adv2KeypadLayerLeft, Adv2KeypadLayerRight)

	for index, keycode_each_target := range fullLayerTargetAdv2Default {
		if (index)%rowCount == 0 && index != 0 {
			fmt.Println("")
		}
		value := KeyPadKinesis(keycode_each_target).tokenname
		if keycode_each_target.String() == "KC_TRANSPARENT" {
			value = `_`
		}

		toprint := "`" + value + "`"
		fmt.Print(toprint, ", ")

		//var expected = kinesisAdv2ndLayer(keycode_each_target)
		//var got = fullLayerTargetAdv2_Keypad[index]

		// if expected != got {
		// 	var message = ":`" + keycode_each_target.String() + "` should map to `" + expected.String() + "` but instead we got `" + got.String() + "`"
		// 	t.Error(message)
		// }

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
