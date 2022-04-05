package main

import (
	"fmt"

	"masters3d.com/keyboard_layout_config_mapper/source/keyboard_config"
)

func main() {
	fmt.Printf(fmt.Sprint(len(keyboard_config.FullKeyboard)))

	for index, element := range keyboard_config.FullKeyboard {
		fmt.Println(index, element)
	}
}
