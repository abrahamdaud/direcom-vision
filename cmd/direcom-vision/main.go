package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abrahamdaud/direcom-vision/internal/monitor"
	webui "github.com/abrahamdaud/direcom-vision/internal/web"
)

func main() {
	root := flag.String("root", ".", "directory tree to monitor")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	interval := flag.Duration("interval", 250*time.Millisecond, "filesystem scan interval (minimum 100ms)")
	flag.Parse()

	m, err := monitor.New(*root, *interval)
	if err != nil {
		log.Fatalf("initialize monitor: %v", err)
	}
	go m.Run()

	server := &http.Server{
		Addr:              *addr,
		Handler:           webui.NewHandler(m),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown: %v", err)
		}
	}()

	log.Printf("Direcom Vision monitoring %s at http://localhost%s (scan interval %s)", m.Root(), *addr, m.Interval())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server: %v", err)
	}
}
