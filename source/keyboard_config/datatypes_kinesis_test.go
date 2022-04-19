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
