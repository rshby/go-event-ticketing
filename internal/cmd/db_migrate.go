package cmd

import (
	"github.com/rshby/go-event-ticketing/internal/database"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	DbMigrateCmd = &cobra.Command{
		Use: "db-migrate",
		Run: dbMigrate,
	}
)

func init() {
	RootCmd.AddCommand(DbMigrateCmd)
}

func dbMigrate(cmd *cobra.Command, args []string) {
	dbGorm, err := database.ConnectPostgreSql()
	if err != nil {
		logrus.Fatalf("Gagal terhubung ke database: %v", err)
	}

	dbSql, err := dbGorm.DB()
	if err != nil {
		logrus.Fatalf("Gagal mengekstrak sql.DB: %v", err)
	}

	migrate.SetTable("migrations")

	migrations := &migrate.FileMigrationSource{
		Dir: "./internal/database/migration",
	}

	logrus.Info("Memulai proses migrasi database...")

	n, err := migrate.Exec(dbSql, "postgres", migrations, migrate.Up)
	if err != nil {
		logrus.Fatalf("Migrasi gagal dieksekusi: %v", err)
	}

	if n > 0 {
		logrus.Infof("Berhasil mengeksekusi %d file migrasi ✅", n)
	} else {
		logrus.Info("Database sudah up-to-date. Tidak ada file migrasi baru yang dijalankan ✅")
	}
}
