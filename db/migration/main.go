package main

import (
	"os"

	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"gorm.io/gorm"
)

func main() {
	d := db.GetDB()
	isDryRun := len(os.Args) > 1 && os.Args[1] == "--dry-run"

	if isDryRun {
		d = d.Session(&gorm.Session{DryRun: true})
	} else {
		d = d.Debug()
	}

	runAutoMigrate(d)
}

func runAutoMigrate(d *gorm.DB) {
	err := d.AutoMigrate(&entities.User{})
	if err != nil {
		panic(err)
	}

	err = d.AutoMigrate(&entities.Tiger{})
	if err != nil {
		panic(err)
	}

	err = d.AutoMigrate(&entities.Sighting{})
	if err != nil {
		panic(err)
	}
}
