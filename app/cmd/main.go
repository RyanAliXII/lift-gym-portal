package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/http/middlewares"
	"lift-fitness-gym/app/pkg/applog"
	"lift-fitness-gym/app/pkg/browser"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/handlers"
	"net/http"
	"path/filepath"
	"slices"
	"time"

	_ "time/tzdata"

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
	if data == nil {
		data = handlers.Data{}
	}
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	mapData, ok := data.(handlers.Data)
	if ok {      
		mapData["currentUser"] = sessionData.User	
	}
	_, hasCSRFToken := mapData["csrf"]
	
	if !hasCSRFToken {
		mapData["csrf"] = c.Get("csrf")
	}
	err := t.templates.ExecuteTemplate(w, name, mapData)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
var logger * zap.Logger = applog.Get()
func main() {
	godotenv.Load(".env")
	store := mysqlsession.GetMySQLStore()
	db.GetConnection()
	db.CreateRootAccount()
	browser, err  := browser.NewBrowser()
	if err != nil{
		logger.Error(err.Error())
	}
	defer browser.GetBrowser().Close()
	defer browser.GetLauncher().Close()
	e := echo.New()
	defer e.Shutdown(context.Background())
	logger := applog.Get()
	defer logger.Sync()
	e.Use(middlewares.LoggerMiddleware)
	e.Use(session.Middleware(store))
	e.Use(middleware.CSRF())
	e.Static("/", "./assets")
	e.Renderer = &TemplateRegistry{
		templates: loadTemplates("./views"),
	}
	handlers.RegisterHandlers(e)
	registerNotFoundHandler(e)
	e.Logger.Fatal(e.Start(":80"))	
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
	
	tmpls, err := template.New("views").Funcs(template.FuncMap{
		"hasPermission": func(requiredPermission string, permissions []string )bool {
			return slices.Contains(permissions, requiredPermission)
		},
		"toReadableDate": func(dateStr string) string {
			t, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				return dateStr
			}
			TextDate := "January 2, 2006"
			return t.Format(TextDate)
		},
	}).ParseFiles(templateList...)
	
	if err != nil {
        panic(err)
    }
	return tmpls
}

func registerNotFoundHandler(router *  echo.Echo) {
	router.RouteNotFound("/*", noRouteFunc)
	router.RouteNotFound("/app/*", noRouteFunc)
	router.RouteNotFound("/clients/*", noRouteFunc)
	router.RouteNotFound("/coaches/*", noRouteFunc)
	
}

func noRouteFunc (c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		return c.JSON(http.StatusNotFound, handlers.JSONResponse{
			Status: http.StatusNotFound,
			Message: "Not found",
		})
	}
	
	return c.Render(http.StatusNotFound,"partials/error/404-page", nil)
}