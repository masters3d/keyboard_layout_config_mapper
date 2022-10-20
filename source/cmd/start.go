package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"masters3d.com/keyboard_layout_config_mapper/source/keyboard_config"
)

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(filepath.Dir(filepath.Dir(ex)))

	// shared
	layer0 := keyboard_config.MergeHalfLayers(keyboard_config.Layer0Left, keyboard_config.Layer0Right)
	layer1 := keyboard_config.MergeHalfLayers(keyboard_config.Layer1Left, keyboard_config.Layer1Right)
	layer2 := keyboard_config.MergeHalfLayers(keyboard_config.Layer2Left, keyboard_config.Layer2Right)

	// Ergodox
	qmk_ergodox_path := filepath.Join(exPath, "qmk_ergodox", "keymap.c")
	qmk_ergodox_source_bytes, error_qmk_ergodox_path := ioutil.ReadFile(qmk_ergodox_path)
	if error_qmk_ergodox_path != nil {
		log.Fatalf("%v", error_qmk_ergodox_path)
	}

	// Layer Zero
	qmk_ergodox_target_string := keyboard_config.Ergodox_replace_layer(string(qmk_ergodox_source_bytes), 0, layer0)

	// Layer One
	qmk_ergodox_target_string = keyboard_config.Ergodox_replace_layer(qmk_ergodox_target_string, 1, layer1)

	// Layer Two
	qmk_ergodox_target_string = keyboard_config.Ergodox_replace_layer(qmk_ergodox_target_string, 2, layer2)

	// Writting out Results
	errErgodox := ioutil.WriteFile(qmk_ergodox_path, []byte(qmk_ergodox_target_string), 0777)

	if errErgodox != nil {
		log.Fatalf("%v", errErgodox)
	}

	//kinesi2

	path_file_2 := filepath.Join(exPath, "kinesis2", "querty_2.txt")

	kinesi2_source_bytes, error_kinesi2_path := ioutil.ReadFile(path_file_2)
	if error_kinesi2_path != nil {
		log.Fatalf("%v", error_kinesi2_path)
	}

	// use the below to reset the file
	// keyboard_config.KinesisAdv2CreatedBlankFile(path_file_2)

	kinesi2_target := keyboard_config.Kinesis_ParseAndFill_SpecialTokens(string(kinesi2_source_bytes), layer0, layer1)

	// Writting out Results
	errKinesis2 := ioutil.WriteFile(path_file_2, []byte(kinesi2_target), 0777)

	if errKinesis2 != nil {
		log.Fatalf("%v", errKinesis2)
	}

}
