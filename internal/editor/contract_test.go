package editor_test

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

const configDir = "../../configs/editor"

type action struct {
	ID     string `json:"id"`
	Bank   string `json:"bank"`
	Slot   int    `json:"slot"`
	Key    string `json:"key"`
	Effect string `json:"effect"`
}

type contract struct {
	Version   int `json:"version"`
	Transport map[string]struct {
		Modifiers []string `json:"modifiers"`
	} `json:"transport"`
	Actions []action `json:"actions"`
}

func loadJSON(t *testing.T, name string, value any) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(configDir, name))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, value); err != nil {
		t.Fatal(err)
	}
}

func TestActionContract(t *testing.T) {
	var c contract
	loadJSON(t, "actions.json", &c)
	if c.Version != 1 || len(c.Actions) != 36 {
		t.Fatalf("unexpected contract version or action count: %d, %d", c.Version, len(c.Actions))
	}
	for bank, modifiers := range map[string][]string{
		"nav": {"ctrl", "alt"}, "select": {"ctrl", "alt", "shift"}, "edit": {"ctrl", "shift"},
	} {
		if !reflect.DeepEqual(c.Transport[bank].Modifiers, modifiers) {
			t.Fatalf("%s transport no longer matches the firmware encoding", bank)
		}
	}
	ids, signals := map[string]bool{}, map[string]bool{}
	gestures := map[string]map[int]string{}
	for _, a := range c.Actions {
		if a.ID == "" || a.Effect == "" || a.Key == "" || ids[a.ID] {
			t.Fatalf("invalid or duplicate action: %+v", a)
		}
		ids[a.ID] = true
		bank, ok := c.Transport[a.Bank]
		if !ok || len(bank.Modifiers) == 0 || a.Slot < 1 || a.Slot > 12 {
			t.Fatalf("invalid transport: %+v", a)
		}
		signal := fmt.Sprintf("%s-f%d", strings.Join(bank.Modifiers, "-"), a.Slot)
		if signals[signal] {
			t.Fatalf("duplicate signal %s", signal)
		}
		signals[signal] = true
		if gestures[a.Bank] == nil {
			gestures[a.Bank] = map[int]string{}
		}
		gestures[a.Bank][a.Slot] = a.Key
	}
	for slot := 1; slot <= 12; slot++ {
		if gestures["nav"][slot] == "" || gestures["nav"][slot] != gestures["select"][slot] {
			t.Fatalf("navigation/selection gesture differs at slot %d", slot)
		}
	}
}

// Empty headers let cpp check layer guards and macros, not Zephyr/ZMK compilation.
func preprocess(t *testing.T, enabled bool) string {
	t.Helper()
	cpp, err := exec.LookPath("cpp")
	if err != nil {
		t.Skip("cpp unavailable: preprocessing checks not run")
	}
	include := t.TempDir()
	for _, name := range []string{"behaviors.dtsi", "dt-bindings/zmk/keys.h", "dt-bindings/zmk/bt.h", "dt-bindings/zmk/outputs.h"} {
		path := filepath.Join(include, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0644); err != nil {
			t.Fatal(err)
		}
	}
	args := []string{"-P", "-x", "assembler-with-cpp", "-nostdinc", "-I", include}
	if enabled {
		args = append(args, "-DKLCM_EDITOR_ENABLE")
	}
	args = append(args, "../../configs/zmk_adv_mod/pillzmod_pro.keymap")
	out, err := exec.Command(cpp, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("cpp: %v\n%s", err, out)
	}
	return string(out)
}

var layersRE = regexp.MustCompile(`(?s)(\w+_layer)\s*\{.*?bindings\s*=\s*<(.*?)>;`)
var whitespaceRE = regexp.MustCompile(`\s+`)

func TestPillzModDefaultUnchanged(t *testing.T) {
	// Original layers after expanding layer constants, ignoring whitespace.
	expected := map[string]string{
		"default_layer": "327c4ff750845f5de54d226d0e1ec51a2f925872d4ffab23fbea86eb798ebd88",
		"keypad_layer":  "85ba6891e5ba5a72a86d052bc8ec8e8f1aaae7cbd0ec9731325f8ad8390a36a1",
		"cmd_layer":     "3538d9f8aa606c0dc0af5450f31c7c59e798be09dc02119659f2551a2f0ff120",
		"system_layer":  "e4c470616fe143a694c38e4e634c07776f8817def22523da8f5ad554ea7d8b77",
	}
	layers := layersRE.FindAllStringSubmatch(preprocess(t, false), -1)
	if len(layers) != len(expected) {
		t.Fatalf("disabled pilot changed layer count: %d", len(layers))
	}
	for _, layer := range layers {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(whitespaceRE.ReplaceAllString(layer[2], ""))))
		if hash != expected[layer[1]] {
			t.Errorf("original %s changed: %s", layer[1], hash)
		}
	}
}

func TestPillzModPilot(t *testing.T) {
	var c contract
	loadJSON(t, "actions.json", &c)
	preprocessed := preprocess(t, true)
	layers := layersRE.FindAllStringSubmatch(preprocessed, -1)
	if len(layers) != 7 {
		t.Fatalf("expected seven layers, got %d", len(layers))
	}
	byName := map[string]string{}
	for _, layer := range layers {
		if count := strings.Count(layer[2], "&"); count != 89 {
			t.Errorf("%s has %d bindings, expected 89 including three pedals", layer[1], count)
		}
		byName[layer[1]] = layer[2]
	}
	base := byName["default_layer"]
	for _, access := range []string{
		"key-positions=<8074>;bindings=<&mo4>;layers=<0>;slow-release;",
		"key-positions=<8577>;bindings=<&mo6>;layers=<0>;slow-release;",
	} {
		if !strings.Contains(whitespaceRE.ReplaceAllString(preprocessed, ""), access) {
			t.Fatal("pilot must have same-hand Space/Keypad combo access")
		}
	}
	for _, oldLayer := range layersRE.FindAllStringSubmatch(preprocess(t, false), -1) {
		active := whitespaceRE.ReplaceAllString(byName[oldLayer[1]], "")
		if active != whitespaceRE.ReplaceAllString(oldLayer[2], "") {
			t.Errorf("enabled pilot changed original %s", oldLayer[1])
		}
	}
	bankLayers := map[string]string{"nav": "editor_nav_layer", "select": "editor_select_layer", "edit": "editor_edit_layer"}
	wrappers := map[string]string{"nav": "LC(LA(F%d))", "select": "LC(LA(LS(F%d)))", "edit": "LC(LS(F%d))"}
	baseBindings := strings.Split(whitespaceRE.ReplaceAllString(base, ""), "&")[1:]
	positions := map[string]int{}
	for pos, binding := range baseBindings {
		if strings.HasPrefix(binding, "kp") {
			key := strings.TrimPrefix(binding, "kp")
			if _, found := positions[key]; !found {
				positions[key] = pos
			}
		}
	}
	for _, a := range c.Actions {
		bindings := whitespaceRE.ReplaceAllString(byName[bankLayers[a.Bank]], "")
		signal := "&kp" + fmt.Sprintf(wrappers[a.Bank], a.Slot)
		// One function-row binding and one mnemonic binding for every action.
		if strings.Count(bindings, signal) != 2 {
			t.Errorf("%s must be reachable twice as %s", a.ID, signal)
		}
		pos, ok := positions[a.Key]
		if !ok || "&"+strings.Split(bindings, "&")[1:][pos] != signal {
			t.Errorf("%s is not at its documented physical %s key", a.ID, a.Key)
		}
	}
	for _, name := range bankLayers {
		bindings := whitespaceRE.ReplaceAllString(byName[name], "")
		if !strings.HasSuffix(bindings, "&trans&trans&trans") {
			t.Errorf("%s must preserve pedal inheritance", name)
		}
		if !strings.Contains(bindings, "&kpESC") {
			t.Errorf("%s must preserve emergency Escape", name)
		}
	}
}

func TestRawDefaultBindingsRemainParseable(t *testing.T) {
	data, err := os.ReadFile("../../configs/zmk_adv_mod/pillzmod_pro.keymap")
	if err != nil {
		t.Fatal(err)
	}
	layers := layersRE.FindAllStringSubmatch(string(data), -1)
	for _, layer := range layers {
		if layer[1] == "default_layer" {
			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(whitespaceRE.ReplaceAllString(layer[2], ""))))
			if hash != "a2b2b87ba14b1093dc6a3a7aaf508de3026e017e163e9cac2388baf8926148dd" {
				t.Fatal("raw default bindings changed; KLCM sync does not preprocess macros")
			}
			return
		}
	}
	t.Fatal("default layer missing")
}

func TestZedCapabilitiesMatchBindings(t *testing.T) {
	var c contract
	loadJSON(t, "actions.json", &c)
	var keymap []struct {
		Context  string                     `json:"context"`
		Bindings map[string]json.RawMessage `json:"bindings"`
	}
	loadJSON(t, "zed/keymap.json", &keymap)
	var audit struct {
		Actions map[string]struct {
			Status string `json:"status"`
			Reason string `json:"reason"`
		} `json:"actions"`
	}
	loadJSON(t, "zed/capabilities.json", &audit)
	if len(keymap) == 0 || len(audit.Actions) != len(c.Actions) {
		t.Fatal("missing keymap or incomplete capability matrix")
	}
	if len(keymap) != 1 || keymap[0].Context != "Editor && mode == full && !multibuffer && !vim_mode" {
		t.Fatal("adapter must exclude modal editing and inline prompts")
	}
	if len(keymap[0].Bindings) != len(c.Actions) {
		t.Fatal("every signal must be implemented or explicitly blocked")
	}
	for _, a := range c.Actions {
		key := fmt.Sprintf("%s-f%d", strings.Join(c.Transport[a.Bank].Modifiers, "-"), a.Slot)
		binding, ok := keymap[0].Bindings[key]
		capability, declared := audit.Actions[a.ID]
		if !ok || !declared || capability.Reason == "" {
			t.Fatalf("missing binding or audit entry for %s", a.ID)
		}
		switch capability.Status {
		case "unsupported":
			if string(binding) != "null" {
				t.Errorf("unsupported %s must be blocked, not approximated silently", a.ID)
			}
		case "native", "different":
			if string(binding) == "null" {
				t.Errorf("supported %s has no binding", a.ID)
			}
		default:
			t.Errorf("invalid status for %s", a.ID)
		}
		if strings.Contains(string(binding), "vim::") || strings.Contains(string(binding), "SendKeystrokes") {
			t.Errorf("native adapter must not simulate hidden modes: %s", a.ID)
		}
	}
}
