package adapters

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationManager struct {
	databaseURL   string
	migrationPath string
}

func NewMigrationManager(databaseURL, migrationPath string) *MigrationManager {
	return &MigrationManager{databaseURL: databaseURL, migrationPath: migrationPath}
}

func (m *MigrationManager) LoadMigrations() error {
	return m.load(m.migrationPath)
}

func (m *MigrationManager) RollbackMigration() error {
	return m.rollback(m.migrationPath)
}

func (m *MigrationManager) GetStatus() (version uint, dirty bool, err error) {
	return m.status(m.migrationPath)
}

func (m *MigrationManager) load(dir string) error {
	migrator, err := m.newMigrator(dir)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database schema is already up to date")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Successfully applied all migrations")
	return nil
}

func (m *MigrationManager) rollback(dir string) error {
	migrator, err := m.newMigrator(dir)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer migrator.Close()

	version, _, err := migrator.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Println("No migrations to roll back - database is at base version")
			return nil
		}
		return fmt.Errorf("failed to check migration version: %w", err)
	}

	if err := migrator.Steps(-1); err != nil {
		return fmt.Errorf("failed to roll back migration: %w", err)
	}

	log.Printf("Successfully rolled back from version %d", version)
	return nil
}

func (m *MigrationManager) status(dir string) (version uint, dirty bool, err error) {
	migrator, err := m.newMigrator(dir)
	if err != nil {
		return 0, false, fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer migrator.Close()

	ver, dirty, err := migrator.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("failed to check version: %w", err)
	}

	return ver, dirty, nil
}

func (m *MigrationManager) newMigrator(dir string) (*migrate.Migrate, error) {
	sourceURL := fmt.Sprintf("file://%s", dir)
	databaseURL := m.databaseURL

	migrator, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return migrator, nil
}
