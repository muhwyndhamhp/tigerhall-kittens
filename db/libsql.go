package db

import (
	"fmt"

	"github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
	libsql "github.com/renxzen/gorm-libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	if db == nil {

		url := config.Get(config.LIBSQL_URL)
		auth := config.Get(config.LIBSQL_TOKEN)

		str := ""
		if url == "" {
			str = "file:db/kittens.db"
		} else {
			str = fmt.Sprintf("%s?authToken=%s", url, auth)
		}

		d, err := gorm.Open(libsql.Open(str), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Connected to database: %s\n", d.Name())

		db = d
	}
	return db
}

func GetTestDB() *gorm.DB {
	d, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to test database: %s\n", d.Name())

	return d
}
