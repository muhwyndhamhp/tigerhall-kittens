package db

import (
	"fmt"

	"github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
	libsql "github.com/renxzen/gorm-libsql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	url := config.Get(config.LIBSQL_URL)
	auth := config.Get(config.LIBSQL_TOKEN)

	d, err := gorm.Open(libsql.Open(fmt.Sprintf("%s?authToken=%s", url, auth)), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to database: %s\n", d.Name())

	db = d
}

func GetDB() *gorm.DB {
	return db
}
