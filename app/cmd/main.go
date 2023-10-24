package main

import (
	"context"
	"html/template"
	"io"
	"io/fs"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/http/middlewares"
	"lift-fitness-gym/app/pkg/applog"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/handlers"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)



type TemplateRegistry struct {
	templates *template.Template
}
  
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	mapData, ok := data.(handlers.Data)
	if ok {      
		mapData["currentUser"] = sessionData.User
	}
	return t.templates.ExecuteTemplate(w, name, mapData)
}

func main() {
	godotenv.Load(".env")
	store := mysqlsession.GetMySQLStore()
	db.GetConnection()
	db.CreateRootAccount()
	e := echo.New()
	defer e.Shutdown(context.Background())
	logger := applog.Get()
	defer logger.Sync()
	e.Use(middlewares.LoggerMiddleware)
	e.Use(session.Middleware(store))
	e.Use(middleware.CSRF())
	e.Static("/", "/assets")
	e.Renderer = &TemplateRegistry{
		templates: loadTemplates("./views"),
	}
	handlers.RegisterHandlers(e)
	e.Logger.Fatal(e.Start(":5000"))	
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