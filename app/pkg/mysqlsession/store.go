package mysqlsession

import (
	"lift-fitness-gym/app/db"
	"os"
	"sync"

	"github.com/srinathgs/mysqlstore"
)

var store * mysqlstore.MySQLStore

var once sync.Once

func GetMySQLStore() * mysqlstore.MySQLStore{
	once.Do(func() {
		sessionSecret := os.Getenv("SESSION_SECRET")
		s, storeErr := mysqlstore.NewMySQLStore(db.GetConnectionEnvs().DSN, "session", "/", 3600 * 24, []byte(sessionSecret))
		if storeErr != nil {
			panic(storeErr.Error())
		}
		store = s
	})
	
	return store
}