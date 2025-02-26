package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	logFile, err := os.OpenFile("minc.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	rootCmd := &cobra.Command{
		Use:   "minc",
		Short: "MicroShift in Container",
	}

	rootCmd.AddCommand(createCmd, listCmd, deleteCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the MicroShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Create command executed")
		fmt.Println("Item created")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the MicroShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("List command executed")
		fmt.Println("Listing all items")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the MicroShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Delete command executed")
		fmt.Println("Item deleted")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of minc",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Version command executed")
		fmt.Println("version: 0.0.1")
	},
}
