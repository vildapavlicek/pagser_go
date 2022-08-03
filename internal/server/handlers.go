package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/vildapavlicek/go-pagser/internal/DB"
	ds "github.com/vildapavlicek/go-pagser/internal/data_structures"
	"github.com/vildapavlicek/go-pagser/internal/tracer"
	"github.com/vildapavlicek/go-pagser/internal/utils"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func getUser(db *db.DB, c *gin.Context) {
	logger := utils.LoggerWithContext(c)
	user := new(ds.UserId)
	tracer_ctx, _ := c.MustGet(tracer.RequestTracerContext).(context.Context)
	newCtx, span := otel.Tracer(tracer.AppName).Start(tracer_ctx, "getUser")
	defer span.End()

	fmt.Printf("[getUser] span ID %s\n", span.SpanContext().SpanID())

	if err := c.ShouldBindUri(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	logger.Info("querying for user", zap.Int32("userId", user.Id))
	customer, err := db.SelectCustomer(user.Id, c, newCtx)
	if err != nil {

		if err == sql.ErrNoRows {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, customer.ToResponse())
}

func postUser(db *db.DB, c *gin.Context) {
	customer := new(ds.CustomerInsertRequest)

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := db.InsertCustomer(customer.ToCustomerDBO(), c); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)

}
