package main

import (
	"io/ioutil"
	"log"

	"masters3d.com/keyboard_layout_config_mapper/source/keyboard_config"
)

func main() {

	qmk_ergodox_path := `..\qmk_ergodox\keymap.c`
	qmk_ergodox_source_bytes, _ := ioutil.ReadFile(qmk_ergodox_path)

	// Layer Zero
	layer0 := keyboard_config.MergeHalfLayers(keyboard_config.Layer0Left, keyboard_config.Layer0Right)

	qmk_ergodox_target_string := keyboard_config.Ergodox_replace_layer(string(qmk_ergodox_source_bytes), 0, layer0)

	// Layer One
	layer1 := keyboard_config.MergeHalfLayers(keyboard_config.Layer1Left, keyboard_config.Layer1Right)

	qmk_ergodox_target_string = keyboard_config.Ergodox_replace_layer(qmk_ergodox_target_string, 1, layer1)

	// Writting out Results
	err := ioutil.WriteFile(qmk_ergodox_path, []byte(qmk_ergodox_target_string), 0777)

	if err != nil {
		log.Fatalf("%v", err)
	}

	//kinesi2

	if false {

		path_file_2 := `..\kinesis2\querty_2.txt`
		keyboard_config.KinesisAdv2CreatedBlankFile(path_file_2)
	}
}
