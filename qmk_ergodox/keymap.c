#include QMK_KEYBOARD_H
#include "debug.h"
#include "action_layer.h"

#define KC_MAC_UNDO LGUI(KC_Z)
#define KC_MAC_CUT LGUI(KC_X)
#define KC_MAC_COPY LGUI(KC_C)
#define KC_MAC_PASTE LGUI(KC_V)
#define KC_PC_UNDO LCTL(KC_Z)
#define KC_PC_CUT LCTL(KC_X)
#define KC_PC_COPY LCTL(KC_C)
#define KC_PC_PASTE LCTL(KC_V)
#define ES_LESS_MAC KC_GRAVE
#define ES_GRTR_MAC LSFT(KC_GRAVE)
#define ES_BSLS_MAC ALGR(KC_6)
#define NO_PIPE_ALT KC_GRAVE
#define NO_BSLS_ALT KC_EQUAL
#define LSA_T(kc) MT(MOD_LSFT | MOD_LALT, kc)
#define BP_NDSH_MAC ALGR(KC_8)
#define SE_SECT_MAC ALGR(KC_6)
#define KC_ControlAltDelete LALT(LCTL(KC_DEL))

enum custom_keycodes {
  PLACEHOLDER = SAFE_RANGE, // can always be here
  ST_MACRO_Screenshot,
  ST_MACRO_Anglebrakets,
  ST_MACRO_Parenthesis,
  ST_MACRO_SquareBraces,
  ST_MACRO_CurlyBraces,
};


const uint16_t PROGMEM keymaps[][MATRIX_ROWS][MATRIX_COLS] = {
  [0] = LAYOUT_ergodox_pretty(
ST_MACRO_Screenshot, KC_QUOTE, LSFT(KC_QUOTE), KC_MINUS, KC_EQUAL, KC_SLASH, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_BACKSLASH, KC_LEFT_BRACKET, KC_RIGHT_BRACKET, LSFT(KC_9), LSFT(KC_0), KC_TRANSPARENT, 
KC_ControlAltDelete, KC_Q, KC_W, KC_E, KC_R, KC_T, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_Y, KC_U, KC_I, KC_O, KC_P, KC_TRANSPARENT, 
KC_ESCAPE, KC_A, KC_S, KC_D, KC_F, KC_G, KC_H, 
KC_J, KC_K, KC_L, KC_DOT, KC_ENTER, TO(0), KC_Z, 
KC_X, KC_C, KC_V, KC_B, KC_TRANSPARENT, KC_TRANSPARENT, KC_N, 
KC_M, KC_QUOTE, LSFT(KC_QUOTE), KC_COMMA, TO(2), KC_TRANSPARENT, RGUI(KC_Z), 
RGUI(KC_S), RGUI(KC_C), RGUI(KC_V), KC_LEFT, KC_DOWN, KC_UP, KC_RIGHT, 
KC_TRANSPARENT, MO(1), KC_LEFT_GUI, KC_ESCAPE, MO(1), KC_LEFT_CTRL, KC_TAB, 
KC_BACKSPACE, KC_LEFT_SHIFT, KC_LEFT_ALT, KC_RIGHT_GUI, KC_RIGHT_SHIFT, KC_SPACE
  ), //[0] = GENERATED
  [1] = LAYOUT_ergodox_pretty(
KC_TRANSPARENT, KC_F5, KC_F6, KC_F7, KC_F8, LSFT(KC_5), KC_TRANSPARENT, 
KC_TRANSPARENT, LSFT(KC_6), ST_MACRO_SquareBraces, ST_MACRO_CurlyBraces, ST_MACRO_Parenthesis, ST_MACRO_Anglebrakets, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_F9, KC_F10, KC_F11, KC_F12, LSFT(KC_4), KC_TRANSPARENT, 
KC_TRANSPARENT, LSFT(KC_7), KC_1, KC_2, KC_3, LSFT(KC_QUOTE), KC_TRANSPARENT, 
KC_TRANSPARENT, KC_F13, KC_F14, KC_F15, KC_F16, LSFT(KC_3), LSFT(KC_8), 
KC_4, KC_5, KC_6, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_F17, 
KC_F18, KC_F19, KC_F20, LSFT(KC_2), KC_TRANSPARENT, KC_TRANSPARENT, LSFT(KC_BACKSLASH), 
KC_7, KC_8, KC_9, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_F21, 
KC_F22, KC_F23, KC_F24, KC_TRANSPARENT, KC_0, KC_TRANSPARENT, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT
  ), //[1] = GENERATED
  [2] = LAYOUT_ergodox_pretty(
KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_V, KC_W, KC_D, KC_L, KC_QUOTE, KC_TRANSPARENT, 
KC_TRANSPARENT, LSFT(KC_QUOTE), KC_U, KC_O, KC_Y, KC_B, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_N, KC_S, KC_T, KC_R, KC_COMMA, KC_DOT, 
KC_A, KC_E, KC_I, KC_C, KC_TRANSPARENT, KC_TRANSPARENT, KC_P, 
KC_F, KC_M, KC_H, KC_Z, KC_TRANSPARENT, KC_TRANSPARENT, KC_Q, 
KC_K, KC_J, KC_X, KC_G, KC_TRANSPARENT, KC_TRANSPARENT, KC_Z, 
KC_Q, KC_V, KC_B, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, 
KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT
  ), //[2] = GENERATED
  [3] = LAYOUT_ergodox_pretty(
    KC_F12,         KC_F1,          KC_F2,          KC_F3,          KC_F4,          KC_F5,          KC_TRANSPARENT,                                 KC_TRANSPARENT, KC_F6,          KC_F7,          KC_F8,          KC_F9,          KC_F10,         KC_F11,
    KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,                                 KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
    KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,                                                                 KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
    KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,                                 KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
    KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,                                                                                                 KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
                                                                                                    KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT,
                                                                                                                    KC_TRANSPARENT, KC_TRANSPARENT,
                                                                                    KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT
  ), //[3] = GENERATED
};


bool process_record_user(uint16_t keycode, keyrecord_t *record) {
  switch (keycode) {
    case ST_MACRO_Screenshot:
    if (record->event.pressed) {
      SEND_STRING(SS_LSFT(SS_LGUI(SS_TAP(X_S))));

    }
    break;
    case ST_MACRO_Anglebrakets:
    if (record->event.pressed) {
      SEND_STRING(SS_LSFT(SS_TAP(X_COMMA)) SS_DELAY(100) SS_LSFT(SS_TAP(X_DOT)) SS_DELAY(100) SS_TAP(X_LEFT));

    }
    break;
    case ST_MACRO_Parenthesis:
    if (record->event.pressed) {
      SEND_STRING(SS_LSFT(SS_TAP(X_9)) SS_DELAY(100) SS_LSFT(SS_TAP(X_0)) SS_DELAY(100) SS_TAP(X_LEFT));

    }
    break;
    case ST_MACRO_SquareBraces:
    if (record->event.pressed) {
      SEND_STRING(SS_TAP(X_LBRACKET) SS_DELAY(100) SS_TAP(X_RBRACKET) SS_DELAY(100) SS_TAP(X_LEFT));

    }
    break;
    case ST_MACRO_CurlyBraces:
    if (record->event.pressed) {
      SEND_STRING(SS_LSFT(SS_TAP(X_LBRACKET)) SS_DELAY(100) SS_LSFT(SS_TAP(X_RBRACKET)) SS_DELAY(100) SS_TAP(X_LEFT));

    }
    break;
  }
  return true;
};

const key_override_t dot_override = ko_make_basic(MOD_MASK_SHIFT, KC_DOT, LSFT(KC_SCOLON));
const key_override_t comma_override = ko_make_basic(MOD_MASK_SHIFT, KC_COMMA, KC_SCOLON);
const key_override_t parens_left_override = ko_make_basic(MOD_MASK_SHIFT, LSFT(KC_9), LSFT(KC_COMMA));
const key_override_t parens_right_override = ko_make_basic(MOD_MASK_SHIFT, KC_RPRN, LSFT(KC_DOT));
const key_override_t exclamation_override = ko_make_basic(MOD_MASK_SHIFT, KC_BACKSLASH, LSFT(KC_1));
// ' KC_QUOTE and ` KC_GRAVE
const key_override_t quote_single_override = ko_make_basic(MOD_MASK_SHIFT, KC_QUOTE, KC_GRAVE);
// " LSFT(KC_QUOTE) and ~ LSFT(KC_GRAVE)
const key_override_t quote_double_override = ko_make_basic(MOD_MASK_SHIFT, LSFT(KC_QUOTE), LSFT(KC_GRAVE));


// This globally defines all key overrides to be used
const key_override_t **key_overrides = (const key_override_t *[]){
    &dot_override,
    &comma_override,
    &parens_left_override,
    &parens_right_override,
    &exclamation_override,
    &quote_single_override,
    &quote_double_override,
    NULL // Null terminate the array of overrides!
};

uint32_t layer_state_set_user(uint32_t state) {
  uint8_t layer = biton32(state);
  ergodox_board_led_off();
  ergodox_right_led_1_off();
  ergodox_right_led_2_off();
  ergodox_right_led_3_off();
  switch (layer) {
    case 1:
      ergodox_right_led_1_on();
      break;
    case 2:
      ergodox_right_led_2_on();
      break;
    case 3:
      ergodox_right_led_3_on();
      break;
    case 4:
      ergodox_right_led_1_on();
      ergodox_right_led_2_on();
      break;
    case 5:
      ergodox_right_led_1_on();
      ergodox_right_led_3_on();
      break;
    case 6:
      ergodox_right_led_2_on();
      ergodox_right_led_3_on();
      break;
    case 7:
      ergodox_right_led_1_on();
      ergodox_right_led_2_on();
      ergodox_right_led_3_on();
      break;
    default:
      break;
  }
  return state;
};