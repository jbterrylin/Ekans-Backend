package app

import (
	"ekans/app/model"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// RegisterModels :
func RegisterModels(db *pg.DB) error {

	models := []interface{}{
		&model.Group{},
		&model.Leave{},
		&model.Record{},
		&model.Roster{},
		&model.ShiftRequest{},
		&model.User{},
		&model.Area{},
	}

	options := &orm.CreateTableOptions{IfNotExists: true}
	for _, model := range models {
		if err := db.Model(model).CreateTable(options); err != nil {
			return err
		}
	}
	return nil
}
