package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	CreateMigrationFileCmd = &cobra.Command{
		Use:  "create-migration-file [migration_name]",
		Args: cobra.ExactArgs(1),
		Run:  createMigrationFile,
	}
)

func init() {
	RootCmd.AddCommand(CreateMigrationFileCmd)
}

func createMigrationFile(cmd *cobra.Command, args []string) {
	migrationName := args[0]
	migrationDir := "./internal/database/migration"

	if err := os.MkdirAll(migrationDir, os.ModePerm); err != nil {
		logrus.Fatalf("gagal membuat direktori migrasi: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, migrationName)
	filepath := filepath.Join(migrationDir, filename)

	// 3. Template wajib dari library sql-migrate
	template := `-- +migrate Up
-- SQL untuk UP (Create/Alter table) ditulis di bawah baris ini


-- +migrate Down
-- SQL untuk DOWN (Drop table/Rollback) ditulis di bawah baris ini

`

	err := os.WriteFile(filepath, []byte(template), 0644)
	if err != nil {
		logrus.Fatalf("gagal menulis file migrasi: %v", err)
	}

	logrus.Infof("File migrasi berhasil dibuat: %s ✅", filepath)
}
