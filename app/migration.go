package app

import (
	"ekans/app/model"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"golang.org/x/crypto/bcrypt"
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

	count, _ := db.Model((*model.User)(nil)).Count()
	if count > 0 {
		user := new(model.User)
		user.Name = "Wang Lin Lee"
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte("1"), bcrypt.MinCost)
		encodePW := string(hashPassword)
		user.Password = encodePW
		user.PhoneNumber = "60167560969"
		user.Job = "Admin"
		user.RegisterDateTime = time.Now()
		db.Model(user).Insert()
	}

	count1, _ := db.Model((*model.Area)(nil)).Count()
	if count1 == 0 {
		area := new(model.Area)
		area.AreaName = "Heart"
		db.Model(area).Insert()
		area1 := new(model.Area)
		area1.AreaName = "Bone"
		db.Model(area1).Insert()
	}

	return nil
}
