package main

import (
	"encoding/json"
	"fmt"
	"github.com/minc-org/minc/pkg/constants"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/minc"
	"github.com/minc-org/minc/pkg/minc/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
)

var (
	provider      string
	logLevel      string
	uShiftVersion string
	uShiftConfig  string
	httpsPort     string
	httpPort      string
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
		hPort, err := strconv.Atoi(viper.GetString("http-port"))
		if err != nil {
			log.Fatal("http port must be an integer", "port", viper.GetString("http-port"))
		}
		hsPort, err := strconv.Atoi(viper.GetString("https-port"))
		if err != nil {
			log.Fatal("https port must be an integer", "port", viper.GetString("https-port"))
		}

		cType := &types.CreateType{
			Provider:      viper.GetString("provider"),
			UShiftVersion: viper.GetString("microshift-version"),
			UShiftConfig:  uShiftConf,
			HTTPSPort:     hsPort,
			HTTPPort:      hPort,
		}
		err = minc.Create(cType)
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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of MicroShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		status := minc.Status(viper.GetString(provider))
		jsonData, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			log.Fatal("error marshalling status", "err", err)
		}
		fmt.Println(string(jsonData))
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
	createCmd.PersistentFlags().StringVar(&httpsPort, "https-port", "9443",
		"https route port to be exposed by container (default: 9443)")
	createCmd.PersistentFlags().StringVar(&httpPort, "http-port", "9080",
		"http route port to be exposed by container (default: 9080)")

	rootCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "Specify the provider (e.g., podman, docker)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "", "Log level (e.g., info, debug, warn)")

	// Add config subcommands
	configCmd.AddCommand(configSetCmd, configGetCmd, configUnsetCmd, configViewCmd)

	rootCmd.AddCommand(createCmd, listCmd, deleteCmd, versionCmd, statusCmd, configCmd)

	// Binding with viper
	viper.BindPFlag("provider", rootCmd.PersistentFlags().Lookup("provider"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("microshift-version", createCmd.PersistentFlags().Lookup("microshift-version"))
	viper.BindPFlag("microshift-config", createCmd.PersistentFlags().Lookup("microshift-config"))
	viper.BindPFlag("https-port", createCmd.PersistentFlags().Lookup("https-port"))
	viper.BindPFlag("http-port", createCmd.PersistentFlags().Lookup("http-port"))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command: ", err)
		os.Exit(1)
	}
}
