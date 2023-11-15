package mysqlsession

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/srinathgs/mysqlstore"
)

var store * mysqlstore.MySQLStore

var once sync.Once

func GetMySQLStore() *mysqlstore.MySQLStore {
	once.Do(func() {
		var (
			name     = os.Getenv("DB_NAME")
			user     = os.Getenv("DB_USER")
			port     = os.Getenv("DB_PORT")
			password = os.Getenv("DB_PASSWORD")
			host     = os.Getenv("DB_HOST")
			sessionSecret = os.Getenv("SESSION_SECRET")
			retries       = 10
			storeErr      error
		)

		var dsn string
		for i := 1; i <= retries; i++ {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, password, host, port, name)
			store, storeErr = mysqlstore.NewMySQLStore(dsn, "session", "/", 3600*24, []byte(sessionSecret))
			if storeErr == nil {
				break
			}

			fmt.Printf("Error connecting to MySQL (attempt %d): %v\n", i, storeErr)
			time.Sleep(2 * time.Second) // Add a delay before the next retry
		}

		if storeErr != nil {
			panic(fmt.Sprintf("Failed to connect to MySQL after %d attempts: %v", retries, storeErr))
		}
	})

	return store
}