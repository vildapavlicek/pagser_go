package main

import (
	"context"
	"os"

	db "github.com/vildapavlicek/go-pagser/internal/DB"
	"github.com/vildapavlicek/go-pagser/internal/server"
	"github.com/vildapavlicek/go-pagser/internal/tracer"
	"github.com/vildapavlicek/go-pagser/internal/utils"
	"go.uber.org/zap"
)

func main() {
	utils.Initialize()
	defer utils.Logger.Sync()

	trace_file, _ := os.OpenFile("traces.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer trace_file.Close()

	shutdown, _ := tracer.InitTracer("http://localhost:9411/api/v2/spans")

	defer func() {
		if err := shutdown(context.Background()); err != nil {
			utils.Logger.Error("failed to close trace exporter", zap.String("error", err.Error()))
		}
	}()

	utils.Logger.Info("Hello World")
	db := db.Connect("postgres://pagser:pagser1234@127.0.0.1:5432/pagila?sslmode=disable")

	server.StartServer(&db, context.Background())
}
