package main

import (
	"Dandelion/handler"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	router := handler.NewServer()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Router,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbConn, err := pgx.Connect(ctx, `postgres://postgres:postgres@localhost:5432/dandelion?sslmode=disable`)
	if err != nil {
		log.Fatalln("Can't connect: ", err)
	}
	defer dbConn.Close(ctx)

	runDBMigration("file://db/migration", `postgres://postgres:postgres@localhost:5432/dandelion?sslmode=disable`)

	shutdown(ctx, srv)

}

func shutdown(ctx context.Context, srv *http.Server) {
	// 启动HTTP服务器
	go func() {
		log.Println("Starting server...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("Error starting server: ", err)
		}
	}()

	// 用于接收退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	/// 等待退出信号
	<-quit
	log.Println("Received exit signal, shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("Server shutdown: ", err)
	}

	log.Println("Server closed...")
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalf("cannot create new migrate instance, Err: %v", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrate up, Err: %v", err)
	}

	log.Print("db migrated successfully...")
}
