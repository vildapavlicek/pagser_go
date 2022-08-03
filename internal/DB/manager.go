package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	ds "github.com/vildapavlicek/go-pagser/internal/data_structures"
	"github.com/vildapavlicek/go-pagser/internal/tracer"
	"github.com/vildapavlicek/go-pagser/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

/// https://bun.uptrace.dev/postgres/#pgdriver
func Connect(dsn string) DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true), bundebug.FromEnv("BUNDEBUG")))

	return DB{db}

}

type DB struct {
	inner *bun.DB
}

func (db *DB) SelectCustomer(id int32, ctx *gin.Context, tracer_ctx context.Context) (*ds.CustomerDBO, error) {
	_, span := otel.Tracer(tracer.RequestTracerContext).Start(tracer_ctx, "SelectCustomer")
	span.SetAttributes(attribute.Int("id", int(id)))
	defer span.End()
	fmt.Printf("[SelectCustomer] span ID %s\n", span.SpanContext().SpanID())

	logger := utils.LoggerWithContext(ctx)
	customer := new(ds.CustomerDBO)

	if err := db.inner.NewSelect().Model(customer).Where("customer_id = ?", id).Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("requested customer not found", zap.Int32("customer_id", id))
		} else {
			logger.Error("failed to query customer", zap.String("error", err.Error()))
		}
		return nil, err
	}
	return customer, nil
}

/// Inserts new customer data into DB
func (db *DB) InsertCustomer(customer ds.CustomerDBO, ctx *gin.Context) error {
	logger := utils.LoggerWithContext(ctx)
	result, err := db.inner.NewInsert().Model(&customer).Exec(ctx)
	if err != nil {
		logger.Error("failed to insert new customer", zap.String("error", err.Error()))
		return err
	}

	lastInsterId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	logger.Info("inserted new customer to db", zap.String("customer", fmt.Sprintf("%#v", customer)), zap.Int64("lastInsertId", lastInsterId), zap.Int64("rowsAffected", rowsAffected))

	return nil

}
