package keyboard_config

var Adv2TopLayerLeft = KeycodeLayerHalf{
	KC_F1, KC_F2, KC_F3, KC_F4, KC_F5, KC_F6, KC_F7, // function row
	KC_EQUAL, KC_1, KC_2, KC_3, KC_4, KC_5, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
	KC_ESCAPE /*start of function row*/, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
}
var Adv2TopLayerRight = KeycodeLayerHalf{
	KC_F8, KC_F9, KC_F10, KC_F11, KC_F12, KC_PRINT_SCREEN, KC_SCROLL_LOCK, // function row
	KC_TRANSPARENT, KC_6, KC_7, KC_8, KC_9, KC_0, KC_MINUS,
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
	KC_PAUSE /*end of function row*/, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
}

type keycode_kinesis struct {
	description string
	tokenname   string
}

// KINESIS SmartSet
// https://cs.github.com/KinesisCorporation/Freestyle-Edge-Pro-SmartSet-App

var KC_VK_HYPER = keycode_kinesis{"VK_HYPER", "HYPER"} // Alt + Shift + Ctrl + Win HYPER

var KC_VK_MEH = keycode_kinesis{"VK_MEH", "MEH"} //Alt + Shift + Ctrl MEH
var KC_VK_LPEDAL = keycode_kinesis{"VK_LPEDAL", "lp-tab"}
var KC_VK_MPEDAL = keycode_kinesis{"VK_MPEDAL", "mp-kpshf"}
var KC_VK_RPEDAL = keycode_kinesis{"VK_RPEDAL", "rp-kpent"}
var KC_VK_KPSHIFT = keycode_kinesis{"VK_KPSHIFT", "KPSHIFT"}    // Temp shift. Both layers have to be mapped with the same key
var KC_VK_KPTOGGLE = keycode_kinesis{"VK_KPTOGGLE", "KPTOGGLE"} // This switches to the layer. you have to press again to switch. Both layers have to be mapped with the same key

func (self keycode_kinesis) String() string {
	return self.description
}

func KeyPadKinesis(input KeyCodeRepresentable) keycode_kinesis {
	value, isOk := kinesis_confirmed[input]
	if isOk {
		return keycode_kinesis{"KP_" + input.String(), "kp-" + value}
	}
	return keycode_kinesis{"UNKNOWN", "UNKNOWN"}
}

var kinesisAdv2LayerMapping_confirmed = map[KeyCodeRepresentable]KeyCodeRepresentable{
	KC_ESCAPE:       KeyPadKinesis(KC_ESCAPE),
	KC_F1:           KeyPadKinesis(KC_LEFT_GUI),
	KC_F2:           KeyPadKinesis(KC_RIGHT_ALT),
	KC_F3:           KeyPadKinesis(KC_MENU),
	KC_F4:           KeyPadKinesis(KC_MEDIA_PLAY_PAUSE),
	KC_F5:           KeyPadKinesis(KC_MEDIA_PREV_TRACK),
	KC_F6:           KeyPadKinesis(KC_MEDIA_NEXT_TRACK),
	KC_F7:           KeyPadKinesis(KC_CALCULATOR),
	KC_F8:           KeyPadKinesis(KC_VK_KPSHIFT),
	KC_F9:           KeyPadKinesis(KC_F9),
	KC_F10:          KeyPadKinesis(KC_F10),
	KC_F11:          KeyPadKinesis(KC_F11),
	KC_F12:          KeyPadKinesis(KC_F12),
	KC_PRINT_SCREEN: KC_AUDIO_MUTE,
	KC_SCROLL_LOCK:  KC_AUDIO_VOL_UP,
	KC_PAUSE:        KC_AUDIO_VOL_DOWN,
	// number row
	KC_EQUAL: KeyPadKinesis(KC_EQUAL),
	KC_1:     KeyPadKinesis(KC_1),
	KC_2:     KeyPadKinesis(KC_2),
	KC_3:     KeyPadKinesis(KC_3),
	KC_4:     KeyPadKinesis(KC_4),
	KC_5:     KeyPadKinesis(KC_5),
	KC_6:     KeyPadKinesis(KC_6),
	KC_7:     KC_NUM_LOCK,
	KC_8:     KC_KP_EQUAL,
	KC_9:     KC_KP_SLASH,
	KC_0:     KC_KP_ASTERISK,
	KC_MINUS: KeyPadKinesis(KC_MINUS),

	// first alpha row
	KC_TAB:       KeyPadKinesis(KC_TAB),
	KC_Q:         KeyPadKinesis(KC_Q),
	KC_W:         KeyPadKinesis(KC_E),
	KC_E:         KeyPadKinesis(KC_W),
	KC_R:         KeyPadKinesis(KC_R),
	KC_T:         KeyPadKinesis(KC_T),
	KC_Y:         KeyPadKinesis(KC_Y),
	KC_U:         KC_KP_7,
	KC_I:         KC_KP_8,
	KC_O:         KC_KP_9,
	KC_P:         KC_KP_MINUS,
	KC_BACKSLASH: KeyPadKinesis(KC_BACKSLASH),

	// second alpha row
	KC_CAPS_LOCK: KeyPadKinesis(KC_CAPS_LOCK),
	KC_A:         KeyPadKinesis(KC_A),
	KC_S:         KeyPadKinesis(KC_S),
	KC_D:         KeyPadKinesis(KC_D),
	KC_F:         KeyPadKinesis(KC_F),
	KC_G:         KeyPadKinesis(KC_G),
	KC_H:         KeyPadKinesis(KC_H),
	KC_J:         KC_KP_4,
	KC_K:         KC_KP_5,
	KC_L:         KC_KP_6,
	KC_SEMICOLON: KC_KP_PLUS,
	KC_QUOTE:     KeyPadKinesis(KC_BACKSLASH),

	// third alpha row
	KC_Z:           KeyPadKinesis(KC_Z),
	KC_X:           KeyPadKinesis(KC_X),
	KC_C:           KeyPadKinesis(KC_C),
	KC_V:           KeyPadKinesis(KC_V),
	KC_B:           KeyPadKinesis(KC_B),
	KC_N:           KeyPadKinesis(KC_N),
	KC_M:           KC_KP_1,
	KC_COMMA:       KC_KP_2,
	KC_DOT:         KC_KP_3,
	KC_SLASH:       KC_KP_ENTER,
	KC_RIGHT_SHIFT: KeyPadKinesis(KC_RIGHT_SHIFT),

	// last  row
	KC_GRAVE:  KeyPadKinesis(KC_GRAVE),
	KC_INSERT: KeyPadKinesis(KC_INSERT),
	KC_LEFT:   KeyPadKinesis(KC_LEFT),
	KC_RIGHT:  KeyPadKinesis(KC_RIGHT),

	KC_UP:            KeyPadKinesis(KC_UP),
	KC_DOWN:          KeyPadKinesis(KC_DOWN),
	KC_LEFT_BRACKET:  KC_KP_DOT,
	KC_RIGHT_BRACKET: KC_RETURN,

	// thumb clusters
	KC_LEFT_CTRL:  KeyPadKinesis(KC_LEFT_CTRL),
	KC_LEFT_ALT:   KeyPadKinesis(KC_LEFT_ALT),
	KC_RIGHT_GUI:  KeyPadKinesis(KC_RIGHT_GUI),
	KC_RIGHT_CTRL: KeyPadKinesis(KC_RIGHT_CTRL),
	KC_BACKSPACE:  KeyPadKinesis(KC_BACKSPACE),
	KC_DELETE:     KeyPadKinesis(KC_DELETE),
	KC_HOME:       KeyPadKinesis(KC_HOME),
	KC_PAGE_UP:    KeyPadKinesis(KC_PAGE_UP),
	KC_ENTER:      KeyPadKinesis(KC_ENTER),
	KC_SPACE:      KC_KP_0,
	KC_END:        KeyPadKinesis(KC_END),
	KC_PAGE_DOWN:  KeyPadKinesis(KC_PAGE_DOWN),

	// pedals
	KC_VK_LPEDAL: KeyPadKinesis(KC_VK_LPEDAL),
	KC_VK_MPEDAL: KeyPadKinesis(KC_VK_MPEDAL),
	KC_VK_RPEDAL: KeyPadKinesis(KC_VK_RPEDAL),
}

var kinesis_confirmed = map[KeyCodeRepresentable]string{
	KC_F1:               "F1",
	KC_F2:               "F2",
	KC_F3:               "F3",
	KC_F4:               "F4",
	KC_F5:               "F5",
	KC_F6:               "F6",
	KC_F7:               "F7",
	KC_F8:               "F8",
	KC_F9:               "F9",
	KC_F10:              "F10",
	KC_F11:              "F11",
	KC_F12:              "F12",
	KC_F13:              "F13",
	KC_F14:              "F14",
	KC_F15:              "F15",
	KC_F16:              "F16",
	KC_F17:              "F17",
	KC_F18:              "F18",
	KC_F19:              "F19",
	KC_F20:              "F20",
	KC_F21:              "F21",
	KC_F22:              "F22",
	KC_F23:              "F23",
	KC_F24:              "F24",
	KC_TRANSPARENT:      "null", // maybe we need a different sigil for null
	KC_ESCAPE:           "escape",
	KC_ENTER:            "enter",
	KC_CAPS_LOCK:        "caps",
	KC_TAB:              "tab",
	KC_LEFT_CTRL:        "lctrl",
	KC_LEFT_SHIFT:       "lshift",
	KC_LEFT_ALT:         "lalt",
	KC_LEFT_GUI:         "lwin",
	KC_RIGHT_CTRL:       "rctrl",
	KC_RIGHT_SHIFT:      "rshift",
	KC_RIGHT_ALT:        "ralt",
	KC_RIGHT_GUI:        "rwin",
	KC_MINUS:            "HYPHEN",
	KC_SPACE:            "SPACE",
	KC_BACKSPACE:        "BSPACE",
	KC_PRINT_SCREEN:     "prtscr",
	KC_LEFT_BRACKET:     "oBRACK",
	KC_RIGHT_BRACKET:    "cBRACk",
	KC_PAUSE:            "PAUSE",
	KC_HOME:             "HOME",
	KC_PAGE_UP:          "pup",
	KC_PAGE_DOWN:        "pdown",
	KC_SCROLL_LOCK:      "SCROLL",
	KC_END:              "END",
	KC_DELETE:           "DELETE",
	KC_RIGHT:            "RIGHT",
	KC_LEFT:             "LEFT",
	KC_DOWN:             "DOWN",
	KC_UP:               "UP",
	KC_INSERT:           "INSERT", // This is only on the keypad layer and it labeled as kp-insert
	KC_NONUS_BACKSLASH:  `intl-\`, // Kinesis calls this international which is different from internal 1,2,3, etc:   ConfigKeys.Add(TKey.Create(VK_OEM_102, 'intl-\', '', 'intl\', 'intl-\', 'intl-\', true, true)); //International <> key between Left Shift and Z
	KC_MENU:             "MENU",
	KC_MS_BTN1:          "LMOUSE",
	KC_MS_BTN2:          "RMOUSE",
	KC_MS_BTN3:          "MMOUSE",
	KC_AUDIO_VOL_UP:     "VOL+",
	KC_AUDIO_VOL_DOWN:   "VOL-",
	KC_AUDIO_MUTE:       "MUTE",
	KC_MEDIA_NEXT_TRACK: "NEXT",
	KC_MEDIA_PREV_TRACK: "PREV",
	KC_MEDIA_PLAY_PAUSE: "PLAY",
	KC_CALCULATOR:       "CALC",
	//
	// TypeableChars
	//
	KC_A:         "A",
	KC_B:         "B",
	KC_C:         "C",
	KC_D:         "D",
	KC_E:         "E",
	KC_F:         "F",
	KC_G:         "G",
	KC_H:         "H",
	KC_I:         "I",
	KC_J:         "J",
	KC_K:         "K",
	KC_L:         "L",
	KC_M:         "M",
	KC_N:         "N",
	KC_O:         "O",
	KC_P:         "P",
	KC_Q:         "Q",
	KC_R:         "R",
	KC_S:         "S",
	KC_T:         "T",
	KC_U:         "U",
	KC_V:         "V",
	KC_W:         "W",
	KC_X:         "X",
	KC_Y:         "Y",
	KC_Z:         "Z",
	KC_1:         "1",
	KC_2:         "2",
	KC_3:         "3",
	KC_4:         "4",
	KC_5:         "5",
	KC_6:         "6",
	KC_7:         "7",
	KC_8:         "8",
	KC_9:         "9",
	KC_0:         "0",
	KC_EQUAL:     "=",
	KC_GRAVE:     "`",
	KC_COMMA:     ",",
	KC_DOT:       ".",
	KC_SLASH:     "/",
	KC_BACKSLASH: `\`, // this needs to be escaped
	KC_QUOTE:     "'",
	KC_SEMICOLON: ";",
	//
	// custom characters only on kinesis
	//
	KC_VK_HYPER:    KC_VK_HYPER.tokenname, // Alt + Shift + Ctrl + Win HYPER
	KC_VK_MEH:      KC_VK_MEH.tokenname,   //Alt + Shift + Ctrl MEH
	KC_VK_LPEDAL:   KC_VK_LPEDAL.tokenname,
	KC_VK_MPEDAL:   KC_VK_MPEDAL.tokenname,
	KC_VK_RPEDAL:   KC_VK_RPEDAL.tokenname,
	KC_VK_KPSHIFT:  KC_VK_KPSHIFT.tokenname,  // Temp shift. Both layers have to be mapped with the same key
	KC_VK_KPTOGGLE: KC_VK_KPTOGGLE.tokenname, // This switches to the layer. you have to press again to switch. Both layers have to be mapped with the same key
	//
	// Keypad
	//
	KC_NUM_LOCK:    "numlk",
	KC_KP_SLASH:    "kpdiv",
	KC_KP_ASTERISK: "kpmult",
	KC_KP_MINUS:    "kpmin",
	KC_KP_PLUS:     "kpplus",
	KC_KP_ENTER:    "kpenter1",
	KC_RETURN:      "kpenter2", // not sure about this mapping
	KC_KP_1:        "kp1",
	KC_KP_2:        "kp2",
	KC_KP_3:        "kp3",
	KC_KP_4:        "kp4",
	KC_KP_5:        "kp5",
	KC_KP_6:        "kp6",
	KC_KP_7:        "kp7",
	KC_KP_8:        "kp8",
	KC_KP_9:        "kp9",
	KC_KP_0:        "kp0",
	KC_KP_DOT:      "kp.",
	KC_KP_EQUAL:    "kp=",
	KC_KP_COMMA:    "kp,",    // This is not part of the kinesis codebase
	KC_KB_MUTE:     "kpMUTE", // This is not part of the kinesis codebase
}

var kinesis_NotUsedOrConfirmed = map[KeyCodeRepresentable]string{
	KC_NONUS_HASH:          "",
	KC_APPLICATION:         "",
	KC_KB_POWER:            "",
	KC_EXECUTE:             "",
	KC_HELP:                "",
	KC_SELECT:              "",
	KC_STOP:                "",
	KC_AGAIN:               "",
	KC_UNDO:                "",
	KC_CUT:                 "",
	KC_COPY:                "",
	KC_PASTE:               "",
	KC_FIND:                "",
	KC_KB_VOLUME_UP:        "",
	KC_KB_VOLUME_DOWN:      "",
	KC_LOCKING_CAPS_LOCK:   "",
	KC_LOCKING_NUM_LOCK:    "",
	KC_LOCKING_SCROLL_LOCK: "",
	KC_KP_EQUAL_AS400:      "",
	KC_INTERNATIONAL_1:     "",
	KC_INTERNATIONAL_2:     "",
	KC_INTERNATIONAL_3:     "",
	KC_INTERNATIONAL_4:     "",
	KC_INTERNATIONAL_5:     "",
	KC_INTERNATIONAL_6:     "",
	KC_INTERNATIONAL_7:     "",
	KC_INTERNATIONAL_8:     "",
	KC_INTERNATIONAL_9:     "",
	KC_LANGUAGE_1:          "",
	KC_LANGUAGE_2:          "",
	KC_LANGUAGE_3:          "",
	KC_LANGUAGE_4:          "",
	KC_LANGUAGE_5:          "",
	KC_LANGUAGE_6:          "",
	KC_LANGUAGE_7:          "",
	KC_LANGUAGE_8:          "",
	KC_LANGUAGE_9:          "",
	KC_ALTERNATE_ERASE:     "",
	KC_SYSTEM_REQUEST:      "",
	KC_CANCEL:              "",
	KC_CLEAR:               "",
	KC_PRIOR:               "",
	KC_SEPARATOR:           "",
	KC_OUT:                 "",
	KC_OPER:                "",
	KC_CLEAR_AGAIN:         "",
	KC_CRSEL:               "",
	KC_EXSEL:               "",

	/* Generic Desktop Page (0x01) */
	KC_SYSTEM_POWER: "",
	KC_SYSTEM_SLEEP: "",
	KC_SYSTEM_WAKE:  "",

	/* Consumer Page (0x0C) */
	KC_MEDIA_STOP:   "",
	KC_MEDIA_SELECT: "",

	// These overlap with defined code like keypad 00
	KC_MEDIA_EJECT:        "", // 0xB0
	KC_MAIL:               "",
	KC_MY_COMPUTER:        "",
	KC_WWW_SEARCH:         "",
	KC_WWW_HOME:           "",
	KC_WWW_BACK:           "",
	KC_WWW_FORWARD:        "",
	KC_WWW_STOP:           "",
	KC_WWW_REFRESH:        "",
	KC_WWW_FAVORITES:      "",
	KC_MEDIA_FAST_FORWARD: "",
	KC_MEDIA_REWIND:       "",
	KC_BRIGHTNESS_UP:      "",
	KC_BRIGHTNESS_DOWN:    "",

	/* Mouse Buttons */
	KC_MS_UP:    "",
	KC_MS_DOWN:  "",
	KC_MS_LEFT:  "",
	KC_MS_RIGHT: "",
	KC_MS_BTN4:  "",
	KC_MS_BTN5:  "",
	KC_MS_BTN6:  "",
	KC_MS_BTN7:  "",
	/* Mouse Wheel */
	KC_MS_WH_UP:    "",
	KC_MS_WH_DOWN:  "",
	KC_MS_WH_LEFT:  "",
	KC_MS_WH_RIGHT: "",
	/* Acceleration */
	KC_MS_ACCEL0: "",
	KC_MS_ACCEL1: "",
	KC_MS_ACCEL2: "",
}
