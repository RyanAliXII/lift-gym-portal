package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/pkg/applog"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/handlers"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
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
	var logger = applog.Get()
	defer logger.Sync()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI: true,
		LogStatus: true,
		LogMethod: true,
		LogLatency: true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			contentType := c.Request().Header.Get("content-type")
			if v.Status >= 400{
				logger.Error("Error", zap.String("URI", v.URI), zap.Int("status", v.Status), zap.String("duration",fmt.Sprint(v.Latency.Milliseconds()," ms")), zap.String("method", v.Method), zap.String("contentType", contentType), zap.String("IP", v.RemoteIP))
			} 
			if v.Status >= 300 && v.Status < 400 {
				logger.Info("Redirect", zap.String("URI", v.URI), zap.Int("status", v.Status), zap.String("duration",fmt.Sprint(v.Latency.Milliseconds()," ms")), zap.String("method", v.Method), zap.String("contentType", contentType), zap.String("IP", v.RemoteIP))
			}
			if v.Status  >= 200 && v.Status < 300 {
				logger.Info("Request", zap.String("URI", v.URI), zap.Int("status", v.Status), zap.String("duration",fmt.Sprint(v.Latency.Milliseconds()," ms")), zap.String("method", v.Method), zap.String("contentType", contentType), zap.String("IP", v.RemoteIP))
			}
			
			return nil
		},

	}))
	e.Use(session.Middleware(store))
	e.Use(middleware.CSRF())
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