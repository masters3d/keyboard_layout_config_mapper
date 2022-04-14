package main

import (
	"fmt"

	keycode "masters3d.com/keyboard_layout_config_mapper/source/keyboard_config"
)

func main() {
	fmt.Printf(fmt.Sprint(len(keycode.FullKeyboard)))

	for index, element := range keycode.FullKeyboard {
		fmt.Println(index, element)
	}
	fmt.Println(keycode.KC_GRAVE)
	fmt.Println(keycode.KC_MINUS)
}
