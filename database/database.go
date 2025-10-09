package database

import (
	"database/sql"
	"errors"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func GetDatabasePath() (string, error) {
	_, fullPath, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("error getting database path")
	}

	normalizedPath := filepath.ToSlash(fullPath)
	rootFolderName := "bot-spot-v3"
	index := strings.LastIndex(normalizedPath, rootFolderName)

	if index == -1 {
		log.Printf("Folder %s not found in %s", rootFolderName, normalizedPath)
		return "", errors.New("error getting database path")
	}

	endIndex := index + len(rootFolderName)
	projectRootPath := normalizedPath[:endIndex]

	dbDir := filepath.Join(projectRootPath, "db")
	if _, err := os.Stat(dbDir); errors.Is(err, os.ErrNotExist) {
		if mkErr := os.MkdirAll(dbDir, os.ModePerm); mkErr != nil {
			return "", mkErr
		}
	}

	dbFile := filepath.Join(dbDir, "bot.db")

	return dbFile, nil

}

func InitDatabase() (err error) {
	dbPath, err := GetDatabasePath()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Ping or create
	if err := db.Ping(); err != nil {
		return err
	}

	// Utility function
	execAndIgnoreDuplicateColumn := func(stmt string) error {
		_, err := db.Exec(stmt)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate column name") {
				return nil // Ignore the error, the column is already there.
			}
		}
		return err
	}

	// Create table cycles
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cycles (id INTEGER PRIMARY KEY)")
	if err != nil {
		return err
	}

	if err = execAndIgnoreDuplicateColumn("ALTER TABLE cycles ADD COLUMN exchange TEXT"); err != nil {
		return err
	}
	if err = execAndIgnoreDuplicateColumn("ALTER TABLE cycles ADD COLUMN status TEXT"); err != nil {
		return err
	}
	if err = execAndIgnoreDuplicateColumn("ALTER TABLE cycles ADD COLUMN quantity REAL"); err != nil {
		return err
	}
	if err = execAndIgnoreDuplicateColumn("ALTER TABLE cycles ADD COLUMN buyPrice REAL"); err != nil {
		return err
	}
	if err = execAndIgnoreDuplicateColumn("ALTER TABLE cycles ADD COLUMN buyId TEXT"); err != nil {
		return err
	}
	if err = execAndIgnoreDuplicateColumn("ALTER TABLE cycles ADD COLUMN sellPrice REAL"); err != nil {
		return err
	}
	if err = execAndIgnoreDuplicateColumn("ALTER TABLE cycles ADD COLUMN sellId TEXT"); err != nil {
		return err
	}

	// Create table cfg_items
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cfg_items (key TEXT PRIMARY KEY, value TEXT)")
	if err != nil {
		return err
	}

	return nil
}

func GetDB() (sqlDB *sql.DB, err error) {
	path, err := GetDatabasePath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, err
}
