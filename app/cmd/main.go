package main

import (
	"encoding/gob"
	"html/template"
	"io"
	"io/fs"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/handlers"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/srinathgs/mysqlstore"
)



type TemplateRegistry struct {
	templates *template.Template
}
  
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
var store *mysqlstore.MySQLStore

func main() {
	gob.Register(model.User{})
	godotenv.Load(".env")
	dbENVS := db.GetConnectionEnvs()
	sessionSecret := os.Getenv("SESSION_SECRET")
	store, storeErr := mysqlstore.NewMySQLStore(dbENVS.DSN, "session", "/", 3600 * 24, []byte(sessionSecret))
	if storeErr != nil {
		panic(storeErr.Error())
	}
	db.GetConnection()
	db.CreateRootAccount()
	e := echo.New()
	e.Use(session.Middleware(store))
	e.Static("/", "/assets")
	e.Renderer = &TemplateRegistry{
		templates: loadTemplates("./views"),
	}
	handlers.RegisterHandlers(e)
	e.Logger.Fatal(	e.Start(":5000"))

}
func loadTemplates(path string) * template.Template{
	templateList := []string{}
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			fileExtension := filepath.Ext(path)
			if fileExtension == ".html" {
				templateList = append(templateList, path)
				
			}
		}
		return nil
	})
	tmpls, err := template.ParseFiles(templateList...)
	if err != nil {
        panic(err)
    }
	return tmpls
}