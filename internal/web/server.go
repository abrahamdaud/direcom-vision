package web

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abrahamdaud/direcom-vision/internal/monitor"
)

func NewHandler(m *monitor.Monitor) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/snapshot", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(m.Snapshot()); err != nil { http.Error(w, "could not encode snapshot", http.StatusInternalServerError) }
	})
	mux.HandleFunc("GET /api/export", func(w http.ResponseWriter, r *http.Request) {
		snapshot := m.Snapshot()
		switch strings.ToLower(r.URL.Query().Get("format")) {
		case "", "json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Content-Disposition", `attachment; filename="direcom-events.json"`)
			if err := json.NewEncoder(w).Encode(snapshot.Events); err != nil { http.Error(w, "could not encode events", http.StatusInternalServerError) }
		case "csv":
			w.Header().Set("Content-Type", "text/csv; charset=utf-8")
			w.Header().Set("Content-Disposition", `attachment; filename="direcom-events.csv"`)
			writer := csv.NewWriter(w)
			_ = writer.Write([]string{"sequence", "timestamp", "type", "path", "is_directory", "size_bytes", "previous_size_bytes", "modified_at"})
			for _, event := range snapshot.Events {
				_ = writer.Write([]string{jsonNumber(event.Sequence), event.At.Format("2006-01-02T15:04:05.000Z07:00"), event.Type, event.Path, jsonBool(event.IsDir), jsonNumber(event.Size), jsonNumber(event.PreviousSize), event.ModifiedAt.Format("2006-01-02T15:04:05.000Z07:00")})
			}
			writer.Flush()
		default:
			http.Error(w, "format must be json or csv", http.StatusBadRequest)
		}
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" { http.NotFound(w, r); return }
		content, err := assets.ReadFile("assets/index.html")
		if err != nil { http.Error(w, "dashboard unavailable", http.StatusInternalServerError); return }
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(content)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) { http.Error(w, "method not allowed", http.StatusMethodNotAllowed) })
	return mux
}

func jsonNumber(value interface{}) string { data, _ := json.Marshal(value); return string(data) }
func jsonBool(value bool) string { if value { return "true" }; return "false" }
