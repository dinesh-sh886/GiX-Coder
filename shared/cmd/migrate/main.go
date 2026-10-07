package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var (
	dsn            string
	migrationsPath string
)

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration tool",
	Long:  `Database migration tool for GiX-Coder`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.MarkFlagRequired("dsn"); err != nil {
			return fmt.Errorf("failed to mark dsn flag as required: %w", err)
		}
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := migrate.New(migrationsPath, dsn)
		if err != nil {
			return fmt.Errorf("create migrate instance: %w", err)
		}

		version, dirty, err := m.Version()
		if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
			return fmt.Errorf("get version: %w", err)
		}

		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "PostgreSQL DSN (required)")
	rootCmd.PersistentFlags().StringVar(&migrationsPath, "path", "./migrations", "Migrations directory path")
}

func runMigrations(direction string) error {
	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("run up migrations: %w", err)
		}
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("run down migrations: %w", err)
		}
	}

	fmt.Printf("Migrations %s completed successfully\n", direction)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
