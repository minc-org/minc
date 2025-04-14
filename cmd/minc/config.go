package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration for minc",
}

// config set <key> <value>
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]
		viper.Set(key, value)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error writing config:", err)
		} else {
			fmt.Printf("Set %s = %s\n", key, value)
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
		configMap := viper.AllSettings()
		delete(configMap, key)
		encodedConfig, _ := json.MarshalIndent(configMap, "", " ")
		if err := viper.ReadConfig(bytes.NewReader(encodedConfig)); err != nil {
			fmt.Println("Error reading config:", err)
		}
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error updating config:", err)
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
		for _, key := range viper.AllKeys() {
			fmt.Printf("%s: %v\n", key, viper.Get(key))
		}
	},
}
