package main

import (
	db "Dandelion/db/service"
	"Dandelion/handler"
	"Dandelion/internal/config"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Dandelion/internal/logs"

	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conf := config.LoadConfig("internal/config/.")

	dbConnect, err := sqlx.Connect(conf.Postgres.Driver, conf.Postgres.DatabaseSource())
	if err != nil {
		logs.Fatalf("Connect database err: %v\n", err)
	}
	defer dbConnect.Close()

	rdb := initRedisClient(conf.RedisClient.Addr)
	if err := rdb.Ping(ctx).Err(); err != nil {
		logs.Fatalf("Connect Redis error: %v", err)
	}

	queries := db.NewQueries(dbConnect)

	router := handler.NewHandler(*conf, queries, rdb)

	srv := &http.Server{
		Addr:    conf.Http.Address(),
		Handler: router.Router,
	}

	runDBMigration(conf.Postgres.MigrationUrl, conf.Postgres.DatabaseUrl())

	shutdown(ctx, srv)

}

func shutdown(ctx context.Context, srv *http.Server) {
	// 启动HTTP服务器
	go func() {
		logs.Info("Starting server...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logs.Fatal("Error starting server: ", err)
		}
	}()

	// 用于接收退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	/// 等待退出信号
	<-quit
	logs.Info("Received exit signal, shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		logs.Fatal("Server shutdown: ", err)
	}

	logs.Info("Server closed...")
}

// sql file migrate
func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logs.Fatalf("cannot create new migrate instance, Err: %v", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		logs.Fatalf("failed to run migrate up, Err: %v", err)
	}

	logs.Info("db migrated successfully...")
}

func initRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
