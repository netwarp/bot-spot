package database

import (
	"log"
	"testing"
)

func TestRootDir(t *testing.T) {
	dir := RootDir()
	log.Println(dir)
}

func TestGetDatabasePath(t *testing.T) {
	path, err := GetDatabasePath()
	if err != nil {
	}
	log.Println(path)
}

func TestGetDB(t *testing.T) {
	db, err := GetDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)
}

func TestDeleteRowsByRange(t *testing.T) {
	start := 1
	end := 3

	for i := start; i <= end; i++ {
		//DeleteByIdInt(int32(i))
	}
}
