package main

import (
	"dsacli/cmd/complete"
	"dsacli/cmd/list"
	"dsacli/cmd/seed"
	"dsacli/cmd/status"
	"dsacli/cmd/today"
	"dsacli/config"
	"dsacli/db"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

const (
	Version = "1.0.0"
)

func versionCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("dsacli version %s\n", Version)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "dsacli",
		Short: "A CLI tool to practice DSA questions using spaced repetition",
		Long:  `A CLI tool to practice DSA questions using a spaced repetition algorithm with difficulty progression.`,
	}

	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Show the current version of the DSA CLI tool.`,
		Run:   versionCmd,
	}

	db, err := db.NewSQLDatabase(config.NewDefaultConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(today.GetCommand(db))
	rootCmd.AddCommand(complete.GetCommand(db))
	rootCmd.AddCommand(complete.GetProgressCommand(db))
	rootCmd.AddCommand(list.GetCommand(db))
	rootCmd.AddCommand(seed.GetCommand(db))
	rootCmd.AddCommand(status.GetCommand(db))
	rootCmd.AddCommand(versionCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
