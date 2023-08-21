package db

import (
	"fmt"
	"lift-fitness-gym/app/pkg/applog"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db * sqlx.DB

var once sync.Once


var logger = applog.Get()

type ConnectionEnvs struct {
	Name string
	User string
	Port string
	Password string
	Host string
	DSN string
}

func GetConnection () * sqlx.DB{
	envs := GetConnectionEnvs()
	var err error
	once.Do(func() {
		db, err = sqlx.Open("mysql", envs.DSN)
		if err != nil {
			log.Fatal(err)
		}
		// Check if the database connection is successful
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
	})

	return db
}

func GetConnectionEnvs()ConnectionEnvs{
	var name = os.Getenv("DB_NAME")
	var user = os.Getenv("DB_USER")
	var port = os.Getenv("DB_PORT")
	var password = os.Getenv("DB_PASSWORD")	
	var	host = os.Getenv("DB_HOST")
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",user, password, host, port, name)
	return  ConnectionEnvs{
		Name: name,
		User: user,
		Port: port,
		Password: password,
		Host: host,
		DSN: dsn,
	}
}