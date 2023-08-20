package main

import (
	"encoding/gob"
	"html/template"
	"io"
	"io/fs"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/handlers"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)



type TemplateRegistry struct {
	templates *template.Template
}
  
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	gob.Register(mysqlsession.SessionData{})
	godotenv.Load(".env")
	store := mysqlsession.GetMySQLStore()
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