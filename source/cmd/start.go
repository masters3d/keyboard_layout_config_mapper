package main

import (
	"io/ioutil"
	"log"

	"masters3d.com/keyboard_layout_config_mapper/source/keyboard_config"
)

func main() {

	qmk_ergodox_path := `..\qmk_ergodox\keymap.c`
	qmk_ergodox_source_bytes, _ := ioutil.ReadFile(qmk_ergodox_path)

	layer0 := keyboard_config.MergeHalfLayers(keyboard_config.Layer0Left, keyboard_config.Layer0Right)

	qmk_ergodox_target_string := keyboard_config.Ergodox_replace_layer(string(qmk_ergodox_source_bytes), 0, layer0)

	err := ioutil.WriteFile(qmk_ergodox_path, []byte(qmk_ergodox_target_string), 0777)

	if err != nil {
		log.Fatalf("%v", err)
	}

	//kinesi2

	// kinesis2_path := `..\kinesis2\querty_1.txt`
	// kinesis2_source_bytes, kinesis2_path_reading_error := ioutil.ReadFile(kinesis2_path)
	// if kinesis2_path_reading_error != nil {
	// 	log.Fatalf("%v", kinesis2_path_reading_error)
	// }
	kinesis_layer0 := keyboard_config.MergeHalfLayers(keyboard_config.Adv2TopLayerLeft, keyboard_config.Adv2TopLayerRight)

	if false {
		kinesis2_target_string2 := ""

		for index, keyLayer0 := range kinesis_layer0 {

			if keyLayer0 == keyboard_config.KC_TRANSPARENT {
				continue
			}
			keyLayer1 := keyboard_config.Adv2KeypadValidation[index]

			var string_key = "*#  " + keyLayer0.String()

			isOkay, value := keyboard_config.KinesisMainLayerMapping(keyLayer0)

			if !isOkay {
				panic("error " + string_key)
			}
			kinesis2_target_string2 += string_key + ":" + value.GetTokenname() + "|" + keyLayer1 + "\n"

		}
		err := ioutil.WriteFile(`..\kinesis2\querty_2.txt`, []byte(kinesis2_target_string2), 0777)

		if err != nil {
			log.Fatalf("%v", err)
		}

	}

}
