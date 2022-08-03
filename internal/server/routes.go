package server

import (
	"github.com/gin-gonic/gin"
	db "github.com/vildapavlicek/go-pagser/internal/DB"
)

func registerGetUser(engine *gin.Engine, db *db.DB) {
	engine.GET("/user/:id", func(c *gin.Context) { getUser(db, c) })
}

func registerPostUser(engine *gin.Engine, db *db.DB) {
	engine.POST("/user", func(c *gin.Context) { postUser(db, c) })
}
