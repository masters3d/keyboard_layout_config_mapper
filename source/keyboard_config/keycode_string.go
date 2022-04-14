// Code generated by "stringer -type=keycode"; DO NOT EDIT.

package keyboard_config

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[KC_TRANSPARENT-0]
	_ = x[KC_ROLL_OVER-1]
	_ = x[KC_POST_FAIL-2]
	_ = x[KC_UNDEFINED-3]
	_ = x[KC_A-4]
	_ = x[KC_B-5]
	_ = x[KC_C-6]
	_ = x[KC_D-7]
	_ = x[KC_E-8]
	_ = x[KC_F-9]
	_ = x[KC_G-10]
	_ = x[KC_H-11]
	_ = x[KC_I-12]
	_ = x[KC_J-13]
	_ = x[KC_K-14]
	_ = x[KC_L-15]
	_ = x[KC_M-16]
	_ = x[KC_N-17]
	_ = x[KC_O-18]
	_ = x[KC_P-19]
	_ = x[KC_Q-20]
	_ = x[KC_R-21]
	_ = x[KC_S-22]
	_ = x[KC_T-23]
	_ = x[KC_U-24]
	_ = x[KC_V-25]
	_ = x[KC_W-26]
	_ = x[KC_X-27]
	_ = x[KC_Y-28]
	_ = x[KC_Z-29]
	_ = x[KC_1-30]
	_ = x[KC_2-31]
	_ = x[KC_3-32]
	_ = x[KC_4-33]
	_ = x[KC_5-34]
	_ = x[KC_6-35]
	_ = x[KC_7-36]
	_ = x[KC_8-37]
	_ = x[KC_9-38]
	_ = x[KC_0-39]
	_ = x[KC_ENTER-40]
	_ = x[KC_ESCAPE-41]
	_ = x[KC_BACKSPACE-42]
	_ = x[KC_TAB-43]
	_ = x[KC_SPACE-44]
	_ = x[KC_MINUS-45]
	_ = x[KC_EQUAL-46]
	_ = x[KC_LEFT_BRACKET-47]
	_ = x[KC_RIGHT_BRACKET-48]
	_ = x[KC_BACKSLASH-49]
	_ = x[KC_NONUS_HASH-50]
	_ = x[KC_SEMICOLON-51]
	_ = x[KC_QUOTE-52]
	_ = x[KC_GRAVE-53]
	_ = x[KC_COMMA-54]
	_ = x[KC_DOT-55]
	_ = x[KC_SLASH-56]
	_ = x[KC_CAPS_LOCK-57]
	_ = x[KC_F1-58]
	_ = x[KC_F2-59]
	_ = x[KC_F3-60]
	_ = x[KC_F4-61]
	_ = x[KC_F5-62]
	_ = x[KC_F6-63]
	_ = x[KC_F7-64]
	_ = x[KC_F8-65]
	_ = x[KC_F9-66]
	_ = x[KC_F10-67]
	_ = x[KC_F11-68]
	_ = x[KC_F12-69]
	_ = x[KC_PRINT_SCREEN-70]
	_ = x[KC_SCROLL_LOCK-71]
	_ = x[KC_PAUSE-72]
	_ = x[KC_INSERT-73]
	_ = x[KC_HOME-74]
	_ = x[KC_PAGE_UP-75]
	_ = x[KC_DELETE-76]
	_ = x[KC_END-77]
	_ = x[KC_PAGE_DOWN-78]
	_ = x[KC_RIGHT-79]
	_ = x[KC_LEFT-80]
	_ = x[KC_DOWN-81]
	_ = x[KC_UP-82]
	_ = x[KC_NUM_LOCK-83]
	_ = x[KC_KP_SLASH-84]
	_ = x[KC_KP_ASTERISK-85]
	_ = x[KC_KP_MINUS-86]
	_ = x[KC_KP_PLUS-87]
	_ = x[KC_KP_ENTER-88]
	_ = x[KC_KP_1-89]
	_ = x[KC_KP_2-90]
	_ = x[KC_KP_3-91]
	_ = x[KC_KP_4-92]
	_ = x[KC_KP_5-93]
	_ = x[KC_KP_6-94]
	_ = x[KC_KP_7-95]
	_ = x[KC_KP_8-96]
	_ = x[KC_KP_9-97]
	_ = x[KC_KP_0-98]
	_ = x[KC_KP_DOT-99]
	_ = x[KC_NONUS_BACKSLASH-100]
	_ = x[KC_APPLICATION-101]
	_ = x[KC_KB_POWER-102]
	_ = x[KC_KP_EQUAL-103]
	_ = x[KC_F13-104]
	_ = x[KC_F14-105]
	_ = x[KC_F15-106]
	_ = x[KC_F16-107]
	_ = x[KC_F17-108]
	_ = x[KC_F18-109]
	_ = x[KC_F19-110]
	_ = x[KC_F20-111]
	_ = x[KC_F21-112]
	_ = x[KC_F22-113]
	_ = x[KC_F23-114]
	_ = x[KC_F24-115]
	_ = x[KC_EXECUTE-116]
	_ = x[KC_HELP-117]
	_ = x[KC_MENU-118]
	_ = x[KC_SELECT-119]
	_ = x[KC_STOP-120]
	_ = x[KC_AGAIN-121]
	_ = x[KC_UNDO-122]
	_ = x[KC_CUT-123]
	_ = x[KC_COPY-124]
	_ = x[KC_PASTE-125]
	_ = x[KC_FIND-126]
	_ = x[KC_KB_MUTE-127]
	_ = x[KC_KB_VOLUME_UP-128]
	_ = x[KC_KB_VOLUME_DOWN-129]
	_ = x[KC_LOCKING_CAPS_LOCK-130]
	_ = x[KC_LOCKING_NUM_LOCK-131]
	_ = x[KC_LOCKING_SCROLL_LOCK-132]
	_ = x[KC_KP_COMMA-133]
	_ = x[KC_KP_EQUAL_AS400-134]
	_ = x[KC_INTERNATIONAL_1-135]
	_ = x[KC_INTERNATIONAL_2-136]
	_ = x[KC_INTERNATIONAL_3-137]
	_ = x[KC_INTERNATIONAL_4-138]
	_ = x[KC_INTERNATIONAL_5-139]
	_ = x[KC_INTERNATIONAL_6-140]
	_ = x[KC_INTERNATIONAL_7-141]
	_ = x[KC_INTERNATIONAL_8-142]
	_ = x[KC_INTERNATIONAL_9-143]
	_ = x[KC_LANGUAGE_1-144]
	_ = x[KC_LANGUAGE_2-145]
	_ = x[KC_LANGUAGE_3-146]
	_ = x[KC_LANGUAGE_4-147]
	_ = x[KC_LANGUAGE_5-148]
	_ = x[KC_LANGUAGE_6-149]
	_ = x[KC_LANGUAGE_7-150]
	_ = x[KC_LANGUAGE_8-151]
	_ = x[KC_LANGUAGE_9-152]
	_ = x[KC_ALTERNATE_ERASE-153]
	_ = x[KC_SYSTEM_REQUEST-154]
	_ = x[KC_CANCEL-155]
	_ = x[KC_CLEAR-156]
	_ = x[KC_PRIOR-157]
	_ = x[KC_RETURN-158]
	_ = x[KC_SEPARATOR-159]
	_ = x[KC_OUT-160]
	_ = x[KC_OPER-161]
	_ = x[KC_CLEAR_AGAIN-162]
	_ = x[KC_CRSEL-163]
	_ = x[KC_EXSEL-164]
	_ = x[KC_LEFT_CTRL-224]
	_ = x[KC_LEFT_SHIFT-225]
	_ = x[KC_LEFT_ALT-226]
	_ = x[KC_LEFT_GUI-227]
	_ = x[KC_RIGHT_CTRL-228]
	_ = x[KC_RIGHT_SHIFT-229]
	_ = x[KC_RIGHT_ALT-230]
	_ = x[KC_RIGHT_GUI-231]
	_ = x[KC_MS_UP-237]
	_ = x[KC_MS_DOWN-238]
	_ = x[KC_MS_LEFT-239]
	_ = x[KC_MS_RIGHT-240]
	_ = x[KC_MS_BTN1-241]
	_ = x[KC_MS_BTN2-242]
	_ = x[KC_MS_BTN3-243]
	_ = x[KC_MS_BTN4-244]
	_ = x[KC_MS_BTN5-245]
	_ = x[KC_MS_BTN6-246]
	_ = x[KC_MS_BTN7-247]
	_ = x[KC_MS_WH_UP-248]
	_ = x[KC_MS_WH_DOWN-249]
	_ = x[KC_MS_WH_LEFT-250]
	_ = x[KC_MS_WH_RIGHT-251]
	_ = x[KC_MS_ACCEL0-252]
	_ = x[KC_MS_ACCEL1-253]
	_ = x[KC_MS_ACCEL2-254]
}

const (
	_keycode_name_0 = "KC_TRANSPARENTKC_ROLL_OVERKC_POST_FAILKC_UNDEFINEDKC_AKC_BKC_CKC_DKC_EKC_FKC_GKC_HKC_IKC_JKC_KKC_LKC_MKC_NKC_OKC_PKC_QKC_RKC_SKC_TKC_UKC_VKC_WKC_XKC_YKC_ZKC_1KC_2KC_3KC_4KC_5KC_6KC_7KC_8KC_9KC_0KC_ENTERKC_ESCAPEKC_BACKSPACEKC_TABKC_SPACEKC_MINUSKC_EQUALKC_LEFT_BRACKETKC_RIGHT_BRACKETKC_BACKSLASHKC_NONUS_HASHKC_SEMICOLONKC_QUOTEKC_GRAVEKC_COMMAKC_DOTKC_SLASHKC_CAPS_LOCKKC_F1KC_F2KC_F3KC_F4KC_F5KC_F6KC_F7KC_F8KC_F9KC_F10KC_F11KC_F12KC_PRINT_SCREENKC_SCROLL_LOCKKC_PAUSEKC_INSERTKC_HOMEKC_PAGE_UPKC_DELETEKC_ENDKC_PAGE_DOWNKC_RIGHTKC_LEFTKC_DOWNKC_UPKC_NUM_LOCKKC_KP_SLASHKC_KP_ASTERISKKC_KP_MINUSKC_KP_PLUSKC_KP_ENTERKC_KP_1KC_KP_2KC_KP_3KC_KP_4KC_KP_5KC_KP_6KC_KP_7KC_KP_8KC_KP_9KC_KP_0KC_KP_DOTKC_NONUS_BACKSLASHKC_APPLICATIONKC_KB_POWERKC_KP_EQUALKC_F13KC_F14KC_F15KC_F16KC_F17KC_F18KC_F19KC_F20KC_F21KC_F22KC_F23KC_F24KC_EXECUTEKC_HELPKC_MENUKC_SELECTKC_STOPKC_AGAINKC_UNDOKC_CUTKC_COPYKC_PASTEKC_FINDKC_KB_MUTEKC_KB_VOLUME_UPKC_KB_VOLUME_DOWNKC_LOCKING_CAPS_LOCKKC_LOCKING_NUM_LOCKKC_LOCKING_SCROLL_LOCKKC_KP_COMMAKC_KP_EQUAL_AS400KC_INTERNATIONAL_1KC_INTERNATIONAL_2KC_INTERNATIONAL_3KC_INTERNATIONAL_4KC_INTERNATIONAL_5KC_INTERNATIONAL_6KC_INTERNATIONAL_7KC_INTERNATIONAL_8KC_INTERNATIONAL_9KC_LANGUAGE_1KC_LANGUAGE_2KC_LANGUAGE_3KC_LANGUAGE_4KC_LANGUAGE_5KC_LANGUAGE_6KC_LANGUAGE_7KC_LANGUAGE_8KC_LANGUAGE_9KC_ALTERNATE_ERASEKC_SYSTEM_REQUESTKC_CANCELKC_CLEARKC_PRIORKC_RETURNKC_SEPARATORKC_OUTKC_OPERKC_CLEAR_AGAINKC_CRSELKC_EXSEL"
	_keycode_name_1 = "KC_LEFT_CTRLKC_LEFT_SHIFTKC_LEFT_ALTKC_LEFT_GUIKC_RIGHT_CTRLKC_RIGHT_SHIFTKC_RIGHT_ALTKC_RIGHT_GUI"
	_keycode_name_2 = "KC_MS_UPKC_MS_DOWNKC_MS_LEFTKC_MS_RIGHTKC_MS_BTN1KC_MS_BTN2KC_MS_BTN3KC_MS_BTN4KC_MS_BTN5KC_MS_BTN6KC_MS_BTN7KC_MS_WH_UPKC_MS_WH_DOWNKC_MS_WH_LEFTKC_MS_WH_RIGHTKC_MS_ACCEL0KC_MS_ACCEL1KC_MS_ACCEL2"
)

var (
	_keycode_index_0 = [...]uint16{0, 14, 26, 38, 50, 54, 58, 62, 66, 70, 74, 78, 82, 86, 90, 94, 98, 102, 106, 110, 114, 118, 122, 126, 130, 134, 138, 142, 146, 150, 154, 158, 162, 166, 170, 174, 178, 182, 186, 190, 194, 202, 211, 223, 229, 237, 245, 253, 268, 284, 296, 309, 321, 329, 337, 345, 351, 359, 371, 376, 381, 386, 391, 396, 401, 406, 411, 416, 422, 428, 434, 449, 463, 471, 480, 487, 497, 506, 512, 524, 532, 539, 546, 551, 562, 573, 587, 598, 608, 619, 626, 633, 640, 647, 654, 661, 668, 675, 682, 689, 698, 716, 730, 741, 752, 758, 764, 770, 776, 782, 788, 794, 800, 806, 812, 818, 824, 834, 841, 848, 857, 864, 872, 879, 885, 892, 900, 907, 917, 932, 949, 969, 988, 1010, 1021, 1038, 1056, 1074, 1092, 1110, 1128, 1146, 1164, 1182, 1200, 1213, 1226, 1239, 1252, 1265, 1278, 1291, 1304, 1317, 1335, 1352, 1361, 1369, 1377, 1386, 1398, 1404, 1411, 1425, 1433, 1441}
	_keycode_index_1 = [...]uint8{0, 12, 25, 36, 47, 60, 74, 86, 98}
	_keycode_index_2 = [...]uint8{0, 8, 18, 28, 39, 49, 59, 69, 79, 89, 99, 109, 120, 133, 146, 160, 172, 184, 196}
)

func (i keycode) String() string {
	switch {
	case i <= 164:
		return _keycode_name_0[_keycode_index_0[i]:_keycode_index_0[i+1]]
	case 224 <= i && i <= 231:
		i -= 224
		return _keycode_name_1[_keycode_index_1[i]:_keycode_index_1[i+1]]
	case 237 <= i && i <= 254:
		i -= 237
		return _keycode_name_2[_keycode_index_2[i]:_keycode_index_2[i+1]]
	default:
		return "keycode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
