package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)
import (
	"github.com/gonutz/w32/v2"
	"gopkg.in/yaml.v3"
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

// This mini config is returned if we can't load a valid file
// and cannot write the detailed example yaml config.example.yaml
// into the expected path at %HOME%
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

//go:embed config.example.yaml
var configExampleYaml []byte

// Expected config path at %HOME%/.config/RectangleWin/config.yaml
var DEFAULT_CONF_PATH_PREFIX = ".config/RectangleWin/"
var DEFAULT_CONF_NAME = "config.yaml"

func convertModifier(keyName string) (int32, error) {
	switch strings.ToLower(keyName) {
	case "ctrl":
		return MOD_CONTROL, nil
	case "alt":
		return MOD_ALT, nil
	case "shift":
		return MOD_SHIFT, nil
	case "win", "meta", "super":
		return MOD_WIN, nil
	default:
		return 0, fmt.Errorf("invalid keyname: %s", keyName)
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
	default:
		return 0, fmt.Errorf("Unknown key %s", key)
	}
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

func getValidConfigPathOrCreate() (string, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = os.Getenv("USERPROFILE")
	}
	if homeDir == "" {
		// Give up generating a valid path.
		// read or write the conf in current folder.
		return DEFAULT_CONF_NAME, errors.New("Failed to find user home directory")
	}
	configDir := filepath.Join(homeDir, filepath.FromSlash(DEFAULT_CONF_PATH_PREFIX))
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory under user's home folder: %s", err)
		// read or write the conf in current folder
		return DEFAULT_CONF_NAME, fmt.Errorf("Failed to create folders under user's home directory: %s", configDir)
	}
	configPath := filepath.Join(configDir, DEFAULT_CONF_NAME)
	return configPath, nil
}

func maybeDropExampleConfigFile(target string) {
	// Check if the file exists, if not, create it with some content
	if _, err := os.Stat(target); os.IsNotExist(err) {
		// Create the file and write the sample content
		err := ioutil.WriteFile(target, configExampleYaml, 0644)
		if err != nil {
			fmt.Println("Failed to create file created: %s %v", target, err)
		}
		fmt.Println("File created: %s", target)
	}
}

func fetchConfiguration() Configuration {
	// Create a Configuration file.
	myConfig := Configuration{}

	// Yaml feeder
	configFilePath, err := getValidConfigPathOrCreate()
	if err == nil {
		maybeDropExampleConfigFile(configFilePath)
	}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Failed to load config file at expected path %s\n", configFilePath)
		// use the last-ditch config
		return DEFAULT_CONF
	}

	if err := yaml.Unmarshal(data, &myConfig); err != nil {
		showMessageBox("Failed to parse config file at %s.\n")
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
	return myConfig
}
