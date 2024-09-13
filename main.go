package main

import (
	"Dandelion/handler"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := handler.NewServer()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Router,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdown(ctx, srv)

}

func shutdown(ctx context.Context, srv *http.Server) {
	// 启动HTTP服务器
	go func() {
		fmt.Println("Starting server...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("Error starting server: ", err)
		}
	}()

	// 用于接收退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	/// 等待退出信号
	<-quit
	fmt.Println("Received exit signal, shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown: ", err)
	}

	fmt.Println("Server closed...")
}
