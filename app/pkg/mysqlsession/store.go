package mysqlsession

import (
	"fmt"
	"os"
	"sync"

	"github.com/srinathgs/mysqlstore"
)

var store * mysqlstore.MySQLStore

var once sync.Once

func GetMySQLStore() * mysqlstore.MySQLStore{
	once.Do(func() {
		var name = os.Getenv("DB_NAME")
		var user = os.Getenv("DB_USER")
		var port = os.Getenv("DB_PORT")
		var password = os.Getenv("DB_PASSWORD")	
		var	host = os.Getenv("DB_HOST")
		var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",user, password, host, port, name)
		sessionSecret := os.Getenv("SESSION_SECRET")
		s, storeErr := mysqlstore.NewMySQLStore(dsn, "session", "/", 3600 * 24, []byte(sessionSecret))
		if storeErr != nil {
			panic(storeErr.Error())
		}
		store = s
	})
	
	return store
}