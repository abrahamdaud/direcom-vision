package web

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/abrahamdaud/direcom-vision/internal/monitor"
)

func TestExportCSVAndJSON(t *testing.T) {
	m, err := monitor.New(t.TempDir(), time.Second)
	if err != nil { t.Fatal(err) }
	handler := NewHandler(m)
	for _, format := range []string{"csv", "json"} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/export?format="+format, nil)
		handler.ServeHTTP(recorder, request)
		if recorder.Code != 200 { t.Fatalf("format %s: got status %d", format, recorder.Code) }
		if format == "csv" && !strings.Contains(recorder.Body.String(), "sequence,timestamp,type,path") { t.Fatalf("CSV header missing: %s", recorder.Body.String()) }
		if format == "json" { var events []json.RawMessage; if err := json.Unmarshal(recorder.Body.Bytes(), &events); err != nil { t.Fatalf("invalid JSON export: %v", err) } }
	}
}
