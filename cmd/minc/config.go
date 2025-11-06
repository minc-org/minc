package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/minc-org/minc/pkg/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration for minc",
}

// defaultConfig holds all the default configuration values
var defaultConfig = map[string]interface{}{
	"provider":              "podman",
	"log-level":             "info",
	"microshift-version":    constants.UShiftVersion,
	"https-port":            "9443",
	"http-port":             "9080",
	"microshift-config":     "",
	"disable-overlay-cache": false,
}

// getDefaults returns the default configuration values
func getDefaults() map[string]interface{} {
	return defaultConfig
}

// writeConfigWithoutDefaults writes only non-default values to the config file
func writeConfigWithoutDefaults() error {
	defaults := getDefaults()
	allSettings := viper.AllSettings()

	// Filter out default values
	filteredConfig := make(map[string]interface{})
	for key, value := range allSettings {
		defaultValue, hasDefault := defaults[key]
		if !hasDefault || value != defaultValue {
			filteredConfig[key] = value
		}
	}

	// Get config file path from viper
	configFilePath := viper.ConfigFileUsed()
	if configFilePath == "" {
		return fmt.Errorf("no config file path available from viper")
	}

	// Write filtered config
	jsonData, err := json.MarshalIndent(filteredConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling config: %v", err)
	}

	return os.WriteFile(configFilePath, jsonData, 0644)
}

// config set <key> <value>
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]

		// Use viper to set the value
		viper.Set(key, value)

		// Use our custom writer to write only non-default values
		err := writeConfigWithoutDefaults()
		if err != nil {
			fmt.Println("Error writing config:", err)
		} else {
			defaults := getDefaults()
			defaultValue, hasDefault := defaults[key]
			if hasDefault && value == defaultValue {
				fmt.Printf("Set %s = %s (using default, removed from config)\n", key, value)
			} else {
				fmt.Printf("Set %s = %s\n", key, value)
			}
		}
	},
}

// config get <key>
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := viper.Get(key)
		fmt.Printf("%s: %v\n", key, value)
	},
}

// config unset <key>
var configUnsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Remove a configuration key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		// Check if key exists in viper (either from config file or as default)
		if !viper.IsSet(key) {
			fmt.Printf("Key %s not found in config\n", key)
			return
		}

		// Get the default value for this key
		defaults := getDefaults()
		defaultValue, hasDefault := defaults[key]

		// If it has a default, reset to default (which won't be written to file)
		// If no default, we need to remove it completely from viper's overrides
		if hasDefault {
			// Reset to default - viper will use the default value but it won't be in the config file
			viper.Set(key, defaultValue)
		} else {
			// For keys without defaults, we need to remove from viper's settings
			// This is a limitation of viper - we'll use a workaround
			allSettings := viper.AllSettings()
			delete(allSettings, key)

			// Get the current config file path before resetting
			currentConfigFile := viper.ConfigFileUsed()

			// Clear viper and reload
			viper.Reset()
			// Re-initialize defaults
			for dk, dv := range defaults {
				viper.SetDefault(dk, dv)
			}
			// Set the config file path again
			viper.SetConfigFile(currentConfigFile)

			// Reload remaining settings
			for k, v := range allSettings {
				if k != key {
					viper.Set(k, v)
				}
			}
		}

		// Write the updated config
		err := writeConfigWithoutDefaults()
		if err != nil {
			fmt.Println("Error writing config:", err)
		} else {
			fmt.Printf("Unset %s\n", key)
		}
	},
}

// config view
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the config file",
	Run: func(cmd *cobra.Command, args []string) {
		// Get config file path from viper
		configFilePath := viper.ConfigFileUsed()
		if configFilePath == "" {
			fmt.Println("No config file path available from viper")
			return
		}

		// Read config file directly
		fileContent, err := os.ReadFile(configFilePath)
		if err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}

		// Parse and pretty print JSON
		var configData map[string]interface{}
		if err := json.Unmarshal(fileContent, &configData); err != nil {
			fmt.Println("Error parsing config file:", err)
			return
		}

		if len(configData) == 0 {
			fmt.Println("No configuration set")
			return
		}

		for key, value := range configData {
			fmt.Printf("%s: %v\n", key, value)
		}
	},
}
