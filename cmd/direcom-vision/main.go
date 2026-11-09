package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/abrahamdaud/direcom-vision/internal/monitor"
	webui "github.com/abrahamdaud/direcom-vision/internal/web"
)

func main() {
	root := flag.String("root", ".", "directory tree to monitor")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	interval := flag.Duration("interval", 250*time.Millisecond, "filesystem scan interval (minimum 100ms)")
	mode := flag.String("mode", "web", "run mode: web or cli")
	flag.Parse()

	m, err := monitor.New(*root, *interval)
	if err != nil {
		log.Fatalf("initialize monitor: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	stopCh := make(chan struct{})
	go func() {
		defer close(stopCh)
		m.Run(ctx.Done())
	}()

	if strings.EqualFold(*mode, "cli") {
		runCLI(m)
		<-stopCh
		return
	}
	if !strings.EqualFold(*mode, "web") {
		log.Fatalf("unknown mode %q; choose web or cli", *mode)
	}

	server := &http.Server{
		Addr:              *addr,
		Handler:           webui.NewHandler(m),
		ReadHeaderTimeout: 5 * time.Second,
	}
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
	<-stopCh
}

func runCLI(m *monitor.Monitor) {
	go func() {
		<-time.After(100 * time.Millisecond)
		fmt.Println("DIRECOM VISION · CLI MODE")
		fmt.Printf("Root: %s\nScan interval: %s\n\n", m.Root(), m.Interval())
		fmt.Printf("%-10s %-8s %-12s %s\n", "TYPE", "KIND", "SIZE", "PATH")
		for _, entry := range m.Snapshot().Entries {
			kind := "file"
			if entry.IsDir {
				kind = "dir"
			}
			fmt.Printf("%-10s %-8s %-12s %s\n", "current", kind, formatBytes(entry.Size), filepath.ToSlash(entry.Path))
		}
	}()

	lastSequence := uint64(0)
	ticker := time.NewTicker(350 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		snapshot := m.Snapshot()
		for _, event := range snapshot.Events {
			if event.Sequence <= lastSequence {
				continue
			}
			lastSequence = event.Sequence
			kind := "file"
			if event.IsDir {
				kind = "dir"
			}
			fmt.Printf("%-10s %-8s %-12s %s [%s]\n", strings.ToUpper(event.Type), kind, formatBytes(event.Size), filepath.ToSlash(event.Path), event.At.Local().Format("15:04:05"))
		}
	}
}

func formatBytes(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	return fmt.Sprintf("%.1f KB", float64(n)/1024)
}
