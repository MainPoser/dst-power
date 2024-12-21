package options

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
}

var database Database

type myDatabase struct {
	db *gorm.DB
}

func (m *myDatabase) GetDB() *gorm.DB {
	return m.db
}

func GetDB() *gorm.DB {
	return database.GetDB()
}

type DbOptions struct {
	Url string
}

func (o *DbOptions) InitDatabase() error {
	dsn := o.Url
	conf := &gorm.Config{
		PrepareStmt: true,
	}
	db, err := gorm.Open(mysql.Open(dsn), conf)
	if err != nil {
		return err
	}
	database = &myDatabase{db: db}
	return nil
}
