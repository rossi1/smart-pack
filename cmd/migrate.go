package cmd

import (
	"errors"
	"strings"

	"github.com/rossi1/smart-pack/adapters"
	"github.com/spf13/cobra"
)

var (
	direction string
)

const (
	up   = "up"
	down = "down"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "postgres migration command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please specify direction: usage migrate up or migrate down")
		}
		direction = strings.TrimSpace(args[0])
		if direction != up && direction != down {
			return errors.New("not a valid direction: usage migrate up or migrate down")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		migrationRepo := adapters.NewMigrationManager(cfg.DatabaseURL, cfg.DatabaseMigrationPath)

		switch direction {
		case up:
			migrateUp(migrationRepo)
		case down:
			migrateDown(migrationRepo)
		}
	},
}

func migrateUp(m *adapters.MigrationManager) {
	err := m.LoadMigrations()
	if err != nil {
		panic(err)
	}
	println("Migration up completed")
}

func migrateDown(m *adapters.MigrationManager) {
	err := m.RollbackMigration()
	if err != nil {
		panic(err)
	}
	println("Migration down completed")
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
