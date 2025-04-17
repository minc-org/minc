package main

import (
	"fmt"
	"github.com/minc-org/minc/pkg/constants"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/minc"
	"github.com/minc-org/minc/pkg/minc/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	provider      string
	logLevel      string
	uShiftVersion string
	uShiftConfig  string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the MicroShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		uShiftConf := viper.GetString("microshift-config")
		if uShiftConf != "" {
			_, err := os.Stat(uShiftConf)
			if os.IsNotExist(err) {
				log.Fatal("config file does not exist", "Config", uShiftConf)
			}
		}
		cType := &types.CreateType{
			Provider:      viper.GetString("provider"),
			UShiftVersion: viper.GetString("microshift-version"),
			UShiftConfig:  uShiftConf,
		}
		err := minc.Create(cType)
		if err != nil {
			log.Fatal("error creating cluster", "err", err)
		}
		log.Info("Cluster created")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the MicroShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		ls, err := minc.List(viper.GetString(provider))
		if err != nil {
			log.Fatal("error listing cluster", "err", err)
		}
		fmt.Printf("%s", ls)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the MicroShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		err := minc.Delete(provider)
		if err != nil {
			log.Fatal("error deleting cluster", "err", err)
		}
		fmt.Println("Item deleted")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of minc",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", constants.Version)
	},
}

func initConfig() {
	appName := "minc"
	configFileName := "config.json"

	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("error getting user config directory", err)
		os.Exit(1)
	}

	appConfigDir := filepath.Join(configDir, appName)
	_ = os.MkdirAll(appConfigDir, 0755)

	configFilePath := filepath.Join(appConfigDir, configFileName)
	// If config file doesn't exist, create an empty one
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		emptyConfig := []byte("{}")
		err := os.WriteFile(configFilePath, emptyConfig, 0644)
		if err != nil {
			fmt.Println("Error creating config file:", err)
			os.Exit(1)
		}
		fmt.Println("Created empty config file: ", configFilePath)
	}

	viper.SetConfigFile(configFilePath)
	// Set defaults
	viper.SetDefault("provider", "podman")
	viper.SetDefault("log-level", "info")
	viper.SetDefault("microshift-version", constants.UShiftVersion)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file", err)
		os.Exit(1)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "minc",
		Short: "MicroShift in Container",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set logger based on user-provided log level
			log.SetLogger(viper.GetString("log-level"))
			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	// create command flags
	createCmd.PersistentFlags().StringVarP(&uShiftVersion, "microshift-version", "m", "",
		fmt.Sprintf("MicroShift version to use, check available tag %s", constants.GetImageRegistry()))
	createCmd.PersistentFlags().StringVarP(&uShiftConfig, "microshift-config", "c", "",
		"MicroShift custom config file")

	rootCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "Specify the provider (e.g., podman, docker)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "", "Log level (e.g., info, debug, warn)")

	// Add config subcommands
	configCmd.AddCommand(configSetCmd, configGetCmd, configUnsetCmd, configViewCmd)

	rootCmd.AddCommand(createCmd, listCmd, deleteCmd, versionCmd, configCmd)

	// Binding with viper
	viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("microshift-version", createCmd.PersistentFlags().Lookup("microshift-version"))
	viper.BindPFlag("microshift-config", createCmd.PersistentFlags().Lookup("microshift-config"))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command: ", err)
		os.Exit(1)
	}
}
