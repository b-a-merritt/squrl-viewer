package main

import (
	"log"
	"net/http"
	"os"

	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/db"
	"github.com/b-a-merritt/squrlviewer/app/handlers"
	"github.com/b-a-merritt/squrlviewer/app/reader"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "github.com/pkg/browser"
)

var port = ":1618"

func main() {
	cfg := config.Config{}
	reader.Init(&cfg, os.Args)

	db, dbErr := db.DbOnce(cfg.ConnStr)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()
	cfg.DB = db

	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/public", "public")
	e.GET("/", handlers.NewListTablesPageHandler(&cfg).ServeHTTP)
	e.GET("/:table", handlers.NewTablePageHandler(&cfg).ServeHTTP)
	e.POST("/:table", handlers.NewAddRowHandler(&cfg).ServeHTTP)
	e.GET("/count/:table", handlers.NewCountTableHandler(&cfg).ServeHTTP)
	e.GET("/table/:table", handlers.NewTableHandler(&cfg).ServeHTTP)
	e.POST("/table/:table", handlers.NewInsertHandler(&cfg).ServeHTTP)
	e.GET("/table/:table/column/:column", handlers.NewColumnSortHandler(&cfg).ServeHTTP)

	// browser.OpenURL("http://localhost" + port)

	if err := e.Start(port); err != http.ErrServerClosed {
		log.Fatal((err))
	}
}
