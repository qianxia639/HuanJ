package main

import (
	"HuanJ/config"
	db "HuanJ/db/sqlc"
	"HuanJ/handler"
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"HuanJ/logs"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	conf, err := config.LoadConfig("config/.")
	if err != nil {
		logs.Fatalf("Can't load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connPool, err := pgxpool.New(ctx, conf.Postgres.DatabaseSource())
	if err != nil {
		logs.Fatalf("Cannot connect to database: %v\n", err)
	}

	rdb := initRedisClient(conf.Redis.Address())
	if err := rdb.Ping(ctx).Err(); err != nil {
		logs.Fatalf("Connect Redis error: %v", err)
	}

	runDBMigration(conf.Postgres.MigrationUrl, conf.Postgres.DatabaseUrl())

	store := db.NewStore(connPool)

	router := handler.NewHandler(conf, store, rdb)

	srv := &http.Server{
		Addr:    conf.Http.Address(),
		Handler: router.Router,
	}

	shutdown(ctx, srv)

}

func shutdown(ctx context.Context, srv *http.Server) {
	// 启动HTTP服务器
	go func() {
		logs.Infof("Listening and serving HTTP on %s", strings.Split(srv.Addr, ":")[1])
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
