package main

import (
	"errors"
	"fmt"
	"strings"
)
import (
	"github.com/davecgh/go-spew/spew"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/gonutz/w32/v2"
)

type KeyBinding struct {
	// A repeated value of key modifiers.
	// Valid values include:
	//   SHIFT, ALT, CTRL, WIN (SUPER/META).
	Modifier []string `yaml: "modifier"`
	// When this is set, this overrides Modifier
	ModifierCode []int32
	// Calculated bitwise OR result of modifiers
	CombinedMod int32
	// Valid values are:
	//   A - Z, 0 - 9, UP_ARROW, =, -
	// Anything not covered here could be set directly via KeyCode
	Key string `yaml: "key"`
	// Automatically converted from Key.
	// When this is set, it overrides Key.
	KeyCode int32 `yaml: "key_code"`
	// The feature in RectangleWin to bind to.
	// Valid values:
	//   moveToTop
	//   moveToBottom
	//   moveToLeft
	//   moveToRight
	//   moveToTopLeft
	//   moveToTopRight
	//   moveToBottomLeft
	//   moveToBottomRight
	//   makeSmaller
	//   makeLarger
	//   makeFullHeight
	//
	BindFeature string `yaml: "bindfeature"`
}

type Configuration struct {
	Keybindings []KeyBinding `yaml: "key_binding"`
}

var DEFAULT_CONF = Configuration{
	Keybindings: []KeyBinding{
		{
			Modifier:    []string{"Ctrl", "Alt"},
			Key:         "UP_ARROW",
			KeyCode:     0x26,
			BindFeature: "moveToTop",
		},
	},
}

var DEFAULT_CONF_NAME = "config.yaml"

func convertModifier(keyName string) (int32, error) {
	switch strings.ToLower(keyName) {
	case "ctrl":
		return MOD_CONTROL, nil
	case "alt":
		return MOD_ALT, nil
	case "shift":
		return MOD_SHIFT, nil
	case "win":
	case "meta":
	case "super":
		return MOD_WIN, nil
	default:
		return 0, errors.New("invalid keyname")
	}
	return 0, errors.New("unreachable")
}

func convertKeyCode(key string) (int32, error) {
	k := strings.ToLower(key)
	if len(k) == 1 {
		if k[0] >= 'a' && k[0] <= 'z' {
			return int32(k[0]) - 32, nil
		}
		if k[0] >= '0' && k[0] <= '9' {
			return int32(k[0]), nil
		}
	}
	switch k {
	case "up_arrow":
		return w32.VK_UP, nil
	case "down_arrow":
		return w32.VK_DOWN, nil
	case "left_arrow":
		return w32.VK_LEFT, nil
	case "right_arrow":
		return w32.VK_RIGHT, nil
	case "-":
		return 189, nil
	case "=":
		return 187, nil
	}
	for id, v := range keyNames {
		lv := strings.ToLower(v)
		if lv == k || lv == (k+" key") {
			return int32(id), nil
		}
	}
	return 0, errors.New("Unknown key")
}

func bitwiseOr(nums []int32) int32 {
	if len(nums) == 0 {
		return 0
	}
	result := nums[0]
	for _, n := range nums[1:] {
		result |= n // bitwise OR
	}
	return result
}

func fetchConfiguration() Configuration {
	spew.Dump(DEFAULT_CONF)
	// Create a Configuration file.
	myConfig := Configuration{}

	// Yaml feeder
	yamlFeeder := feeder.Yaml{Path: DEFAULT_CONF_NAME}
	c := config.New()
	c.AddFeeder(yamlFeeder)
	c.AddStruct(&myConfig)

	err := c.Feed()
	if err != nil {
		fmt.Printf("warn: invalid config files found: %s %v\n", DEFAULT_CONF_NAME, err)
		return DEFAULT_CONF
	}

	for i := range myConfig.Keybindings {
		if len(myConfig.Keybindings[i].ModifierCode) == 0 {
			for _, mod := range myConfig.Keybindings[i].Modifier {
				if modCode, err := convertModifier(mod); err == nil {
					myConfig.Keybindings[i].ModifierCode = append(myConfig.Keybindings[i].ModifierCode, modCode)
				} else {
					fmt.Printf("warn: invalid key name %s", mod)
					continue
				}
			}
		}
		myConfig.Keybindings[i].CombinedMod = bitwiseOr(myConfig.Keybindings[i].ModifierCode)
		if myConfig.Keybindings[i].KeyCode == 0 {
			if key, err := convertKeyCode(myConfig.Keybindings[i].Key); err == nil {
				myConfig.Keybindings[i].KeyCode = key
			} else {
				fmt.Printf("warn: invalid key string %s", myConfig.Keybindings[i].Key)
				continue
			}
		}
	}
	spew.Dump(myConfig)
	return myConfig
}
