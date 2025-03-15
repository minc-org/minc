package main

import (
	"fmt"
	"github.com/minc-org/minc/pkg/constants"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/minc"
	"github.com/spf13/cobra"
	"os"
)

var (
	provider string
	logLevel string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the MicroShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		err := minc.Create(provider)
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
		err := minc.List(provider)
		if err != nil {
			log.Fatal("error listing cluster", "err", err)
		}
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

func main() {
	rootCmd := &cobra.Command{
		Use:   "minc",
		Short: "MicroShift in Container",
	}

	rootCmd.AddCommand(createCmd, listCmd, deleteCmd, versionCmd)
	rootCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "podman", "Specify the provider (e.g., podman, docker)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (e.g., info, debug, warn)")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		log.SetLogger(logLevel) // Apply log level after parsing flags
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command: ", err)
		os.Exit(1)
	}
}
