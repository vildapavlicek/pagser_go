package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	db "github.com/vildapavlicek/go-pagser/internal/DB"
	"golang.org/x/net/context"
)

func StartServer(db_manager *db.DB, appCtx context.Context) {
	gin.ForceConsoleColor()
	gin_router := gin.New()
	gin_router.Use(logger_mw(&appCtx))
	gin_router.Use(gin.Recovery())

	ginprometheus.NewPrometheus("gin").Use(gin_router)

	registerGetUser(gin_router, db_manager)
	registerPostUser(gin_router, db_manager)

	server := &http.Server{
		Handler:           gin_router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		Addr:              "127.0.0.1:8081",
	}

	server.ListenAndServe()

}
