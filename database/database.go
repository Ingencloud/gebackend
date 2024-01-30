package database

import (
	"ge/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("ge.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}

	log.Println("Connected to the database successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")
	//TODO: Add migrations
	db.AutoMigrate(
		&models.User{},
		&models.Playlist{},
		&models.InviteCode{},
		&models.PlaylistDetails{},
	)

	Database = DbInstance{Db: db}
}
