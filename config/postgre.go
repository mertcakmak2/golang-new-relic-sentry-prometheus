package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var once = sync.Once{}

func ConnectPostgres() *gorm.DB {
	var postgresDb *gorm.DB
	once.Do(func() {
		dsn := getConnectionString()
		var err error
		postgresDb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Creating single postgres db instance now.")
	})
	return postgresDb
}

func ConnectTestPostgres(connStr string) *gorm.DB {
	var postgresDb *gorm.DB
	once.Do(func() {
		var err error
		postgresDb, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Creating single postgres db instance now.")
	})
	return postgresDb
}

func getConnectionString() string {
	host := config().Database.Host
	user := config().Database.Username
	password := config().Database.Password
	dbname := config().Database.DatabaseName
	port := config().Database.Port

	connectionSting := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul", host, user, password, dbname, port)
	return connectionSting
}
