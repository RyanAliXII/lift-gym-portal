package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db * sqlx.DB

var once sync.Once

func GetConnection () * sqlx.DB{
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")	
	host := os.Getenv("DB_HOST")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",user, password, host, port, name)
	fmt.Println(dsn)
	var err error
	once.Do(func() {
		db, err = sqlx.Open("mysql", dsn)
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