package keyboard_config

var Adv2TopLayerLeft = KeycodeLayerHalf{
	KC_F1, KC_F2, KC_F3, KC_F4, KC_F5, KC_F6, KC_F7, // function row
	KC_EQUAL, KC_1, KC_2, KC_3, KC_4, KC_5, KC_TRANSPARENT,
	KC_TAB, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT,
	KC_CAPS_LOCK, KC_A, KC_S, KC_D, KC_F, KC_G, KC_TRANSPARENT,
	KC_LEFT_SHIFT, KC_Z, KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT,
	/*from fn row -->*/ KC_ESCAPE /*<-- from fn row*/, KC_GRAVE, KC_INSERT, KC_LEFT, KC_RIGHT, KC_LEFT_CTRL, KC_TRANSPARENT,
	KC_TRANSPARENT, KC_BACKSPACE, KC_DELETE, KC_END, KC_HOME, KC_LEFT_ALT, KC_TRANSPARENT,
}
var Adv2TopLayerRight = KeycodeLayerHalf{
	KC_F8, KC_F9, KC_F10, KC_F11, KC_F12, KC_PRINT_SCREEN, KC_SCROLL_LOCK, // function row
	KC_VK_LPEDAL, KC_6, KC_7, KC_8, KC_9, KC_0, KC_MINUS,
	KC_VK_MPEDAL, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_BACKSLASH,
	KC_VK_RPEDAL, KC_H, KC_J, KC_K, KC_L, KC_SEMICOLON, KC_QUOTE,
	KC_TRANSPARENT, KC_N, KC_M, KC_COMMA, KC_DOT, KC_SLASH, KC_RIGHT_SHIFT,
	KC_TRANSPARENT, KC_RIGHT_CTRL, KC_UP, KC_DOWN, KC_LEFT_BRACKET, KC_RIGHT_BRACKET /*from fn row -->*/, KC_PAUSE, /*<--- from fn row*/
	KC_TRANSPARENT, KC_RIGHT_GUI, KC_PAGE_UP, KC_PAGE_DOWN, KC_ENTER, KC_SPACE, KC_TRANSPARENT,
}

var AdvKeypad = []string{
	`kp-F1`, `kp-F2`, `kp-F3`, `kp-F4`, `kp-F5`, `kp-F6`, `kp-F7`,
	`kp-F8`, `kp-F9`, `kp-F10`, `kp-F11`, `kp-F12`, `kp-prtscr`, `kp-SCROLL`,
	`kp-=`, `kp-1`, `kp-2`, `kp-3`, `kp-4`, `kp-5`, `_`,
	`kp-lp-tab`, `kp-6`, `kp-7`, `kp-8`, `kp-9`, `kp-0`, `kp-HYPHEN`,
	`kp-tab`, `kp-Q`, `kp-W`, `kp-E`, `kp-R`, `kp-T`, `_`,
	`kp-mp-kpshf`, `kp-Y`, `kp-U`, `kp-I`, `kp-O`, `kp-P`, `kp-\`,
	`kp-caps`, `kp-A`, `kp-S`, `kp-D`, `kp-F`, `kp-G`, `_`,
	`kp-rp-kpent`, `kp-H`, `kp-J`, `kp-K`, `kp-L`, `kp-;`, `kp-'`,
	`kp-lshift`, `kp-Z`, `kp-X`, `kp-C`, `kp-V`, `kp-B`, `_`,
	`_`, `kp-N`, `kp-M`, `kp-,`, `kp-.`, `kp-/`, `kp-rshift`,
	`kp-escape`, "kp-`", `kp-INSERT`, `kp-LEFT`, `kp-RIGHT`, `kp-lctrl`, `_`,
	`_`, `kp-rctrl`, `kp-UP`, `kp-DOWN`, `kp-oBRACK`, `kp-cBRACk`, `kp-PAUSE`,
	`_`, `kp-BSPACE`, `kp-DELETE`, `kp-END`, `kp-HOME`, `kp-lalt`, `_`,
	`_`, `kp-rwin`, `kp-pup`, `kp-pdown`, `kp-enter`, `kp-SPACE`, `_`,
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

// TODO: We are calling two in two nested functions
// kinesis_confirmed calls this function again. I hate recursion

const unknown_sigil = "UNKNOWN"

func _keyPadKinesisHelper(input KeyCodeRepresentable) keycode_kinesis {
	value, isOk := kinesis_confirmed[input]
	if isOk {
		return keycode_kinesis{"KP_" + input.String(), "kp-" + value}
	}
	// Handle Keypad
	return keycode_kinesis{unknown_sigil, unknown_sigil}
}

func KinesisKeypayLayerMapping(input KeyCodeRepresentable) (bool, keycode_kinesis) {
	value, isOk := kinesisAdv2ndLayerMapping[input]
	if isOk && value.tokenname != unknown_sigil {
		return true, value
	}
	// Handle Keypad
	return false, value
}

var kinesisAdv2ndLayerMapping = map[KeyCodeRepresentable]keycode_kinesis{
	KC_ESCAPE:       _keyPadKinesisHelper(KC_ESCAPE),
	KC_F1:           _keyPadKinesisHelper(KC_LEFT_GUI),
	KC_F2:           _keyPadKinesisHelper(KC_RIGHT_ALT),
	KC_F3:           _keyPadKinesisHelper(KC_MENU),
	KC_F4:           _keyPadKinesisHelper(KC_MEDIA_PLAY_PAUSE),
	KC_F5:           _keyPadKinesisHelper(KC_MEDIA_PREV_TRACK),
	KC_F6:           _keyPadKinesisHelper(KC_MEDIA_NEXT_TRACK),
	KC_F7:           _keyPadKinesisHelper(KC_CALCULATOR),
	KC_F8:           _keyPadKinesisHelper(KC_VK_KPSHIFT),
	KC_F9:           _keyPadKinesisHelper(KC_F9),
	KC_F10:          _keyPadKinesisHelper(KC_F10),
	KC_F11:          _keyPadKinesisHelper(KC_F11),
	KC_F12:          _keyPadKinesisHelper(KC_F12),
	KC_PRINT_SCREEN: {description: KC_AUDIO_MUTE.String(), tokenname: "mute"},
	KC_SCROLL_LOCK:  {description: KC_AUDIO_VOL_UP.String(), tokenname: "vol+"},
	KC_PAUSE:        {description: KC_AUDIO_VOL_DOWN.String(), tokenname: "vol-"},
	// number row
	KC_EQUAL: _keyPadKinesisHelper(KC_EQUAL),
	KC_1:     _keyPadKinesisHelper(KC_1),
	KC_2:     _keyPadKinesisHelper(KC_2),
	KC_3:     _keyPadKinesisHelper(KC_3),
	KC_4:     _keyPadKinesisHelper(KC_4),
	KC_5:     _keyPadKinesisHelper(KC_5),
	KC_6:     _keyPadKinesisHelper(KC_6),
	KC_7:     {description: KC_NUM_LOCK.String(), tokenname: "numlk"},
	KC_8:     {description: KC_KP_EQUAL.String(), tokenname: "kp="},
	KC_9:     {description: KC_KP_SLASH.String(), tokenname: "kpdiv"},
	KC_0:     {description: KC_KP_ASTERISK.String(), tokenname: "kpmult"},
	KC_MINUS: _keyPadKinesisHelper(KC_MINUS),

	// first alpha row
	KC_TAB:       _keyPadKinesisHelper(KC_TAB),
	KC_Q:         _keyPadKinesisHelper(KC_Q),
	KC_W:         _keyPadKinesisHelper(KC_E),
	KC_E:         _keyPadKinesisHelper(KC_W),
	KC_R:         _keyPadKinesisHelper(KC_R),
	KC_T:         _keyPadKinesisHelper(KC_T),
	KC_Y:         _keyPadKinesisHelper(KC_Y),
	KC_U:         {description: KC_KP_7.String(), tokenname: "kp7"},
	KC_I:         {description: KC_KP_8.String(), tokenname: "kp8"},
	KC_O:         {description: KC_KP_9.String(), tokenname: "kp9"},
	KC_P:         {description: KC_KP_MINUS.String(), tokenname: "kpmin"},
	KC_BACKSLASH: _keyPadKinesisHelper(KC_BACKSLASH),

	// second alpha row
	KC_CAPS_LOCK: _keyPadKinesisHelper(KC_CAPS_LOCK),
	KC_A:         _keyPadKinesisHelper(KC_A),
	KC_S:         _keyPadKinesisHelper(KC_S),
	KC_D:         _keyPadKinesisHelper(KC_D),
	KC_F:         _keyPadKinesisHelper(KC_F),
	KC_G:         _keyPadKinesisHelper(KC_G),
	KC_H:         _keyPadKinesisHelper(KC_H),
	KC_J:         {description: KC_KP_4.String(), tokenname: "kp4"},
	KC_K:         {description: KC_KP_5.String(), tokenname: "kp5"},
	KC_L:         {description: KC_KP_6.String(), tokenname: "kp6"},
	KC_SEMICOLON: {description: KC_KP_PLUS.String(), tokenname: "kpplus"},
	KC_QUOTE:     _keyPadKinesisHelper(KC_BACKSLASH),

	// third alpha row
	KC_LEFT_SHIFT:  _keyPadKinesisHelper(KC_LEFT_SHIFT),
	KC_Z:           _keyPadKinesisHelper(KC_Z),
	KC_X:           _keyPadKinesisHelper(KC_X),
	KC_C:           _keyPadKinesisHelper(KC_C),
	KC_V:           _keyPadKinesisHelper(KC_V),
	KC_B:           _keyPadKinesisHelper(KC_B),
	KC_N:           _keyPadKinesisHelper(KC_N),
	KC_M:           {description: KC_KP_1.String(), tokenname: "kp1"},
	KC_COMMA:       {description: KC_KP_2.String(), tokenname: "kp2"},
	KC_DOT:         {description: KC_KP_3.String(), tokenname: "kp3"},
	KC_SLASH:       {description: KC_KP_ENTER.String(), tokenname: "kpenter-1"},
	KC_RIGHT_SHIFT: _keyPadKinesisHelper(KC_RIGHT_SHIFT),

	// last  row
	KC_GRAVE:  _keyPadKinesisHelper(KC_GRAVE),
	KC_INSERT: _keyPadKinesisHelper(KC_INSERT),
	KC_LEFT:   _keyPadKinesisHelper(KC_LEFT),
	KC_RIGHT:  _keyPadKinesisHelper(KC_RIGHT),

	KC_UP:            _keyPadKinesisHelper(KC_UP),
	KC_DOWN:          _keyPadKinesisHelper(KC_DOWN),
	KC_LEFT_BRACKET:  {description: KC_KP_DOT.String(), tokenname: "kp."},
	KC_RIGHT_BRACKET: {description: KC_RETURN.String(), tokenname: "kpenter-2"},

	// thumb clusters
	KC_LEFT_CTRL:  _keyPadKinesisHelper(KC_LEFT_CTRL),
	KC_LEFT_ALT:   _keyPadKinesisHelper(KC_LEFT_ALT),
	KC_RIGHT_GUI:  _keyPadKinesisHelper(KC_RIGHT_GUI),
	KC_RIGHT_CTRL: _keyPadKinesisHelper(KC_RIGHT_CTRL),
	KC_BACKSPACE:  _keyPadKinesisHelper(KC_BACKSPACE),
	KC_DELETE:     _keyPadKinesisHelper(KC_DELETE),
	KC_HOME:       _keyPadKinesisHelper(KC_HOME),
	KC_PAGE_UP:    _keyPadKinesisHelper(KC_PAGE_UP),
	KC_ENTER:      _keyPadKinesisHelper(KC_ENTER),
	KC_SPACE:      {description: KC_KP_0.String(), tokenname: "kp0"},
	KC_END:        _keyPadKinesisHelper(KC_END),
	KC_PAGE_DOWN:  _keyPadKinesisHelper(KC_PAGE_DOWN),

	// pedals
	KC_VK_LPEDAL: _keyPadKinesisHelper(KC_VK_LPEDAL),
	KC_VK_MPEDAL: _keyPadKinesisHelper(KC_VK_MPEDAL),
	KC_VK_RPEDAL: _keyPadKinesisHelper(KC_VK_RPEDAL),
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
