package web

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/abrahamdaud/direcom-vision/internal/monitor"
)

//go:embed assets/*
var assets embed.FS

import "embed"

func NewHandler(m *monitor.Monitor) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/snapshot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(m.Snapshot()); err != nil {
			http.Error(w, "could not encode snapshot", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("GET /api/export", func(w http.ResponseWriter, r *http.Request) {
		snapshot := m.Snapshot()
		format := strings.ToLower(r.URL.Query().Get("format"))
		switch format {
		case "", "json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Content-Disposition", `attachment; filename="direcom-vision-events.json"`)
			if err := json.NewEncoder(w).Encode(snapshot.Events); err != nil {
				http.Error(w, "could not encode events", http.StatusInternalServerError)
			}
		case "csv":
			w.Header().Set("Content-Type", "text/csv; charset=utf-8")
			w.Header().Set("Content-Disposition", `attachment; filename="direcom-vision-events.csv"`)
			writer := csv.NewWriter(w)
			if err := writer.Write([]string{"sequence", "timestamp", "type", "path", "is_dir", "size_bytes", "previous_size_bytes", "modified_at"}); err != nil {
				http.Error(w, "failed to write CSV header", http.StatusInternalServerError)
				return
			}
			for _, event := range snapshot.Events {
				record := []string{
					strconv.FormatUint(event.Sequence, 10),
					event.At.Format(time.RFC3339Nano),
					event.Type,
					event.Path,
					strconv.FormatBool(event.IsDir),
					strconv.FormatInt(event.Size, 10),
					strconv.FormatInt(event.PreviousSize, 10),
					event.ModifiedAt.Format(time.RFC3339Nano),
				}
				if err := writer.Write(record); err != nil {
					http.Error(w, "failed to write CSV row", http.StatusInternalServerError)
					return
				}
			}
			writer.Flush()
		default:
			http.Error(w, "format must be json or csv", http.StatusBadRequest)
		}
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		content, err := assets.ReadFile("assets/index.html")
		if err != nil {
			http.Error(w, "dashboard unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(content)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
	return mux
}
