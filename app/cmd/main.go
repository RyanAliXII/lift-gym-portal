package main

import (
	"context"
	"crypto/tls"
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
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
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

	e.AutoTLSManager.Cache = autocert.DirCache("var/www/.cache") 
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

func customHTTPServer() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<h1>Welcome to Echo!</h1>
			<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
		`)
	})

	autoTLSManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
		Cache: autocert.DirCache("/var/www/.cache"),
		//HostPolicy: autocert.HostWhitelist("<DOMAIN>"),
	}
	s := http.Server{
		Addr:    ":443",
		Handler: e, // set Echo as handler
		TLSConfig: &tls.Config{
			//Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
		//ReadTimeout: 30 * time.Second, // use custom timeouts
	}
	if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}