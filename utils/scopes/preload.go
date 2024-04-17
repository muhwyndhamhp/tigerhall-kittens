package scopes

import "gorm.io/gorm"

type Preload struct {
	Key       string
	Statement string
}

func Preloads(ps ...Preload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, key := range ps {
			db = db.Preload(key.Key, key.Statement)
		}
		return db
	}
}
