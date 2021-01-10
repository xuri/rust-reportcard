package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"flag"
	"html/template"

	"github.com/dgraph-io/badger/v2"
)

var domain = flag.String("domain", "rust-reportcard.xuri.me", "Domain used for your rust-reportcard installation")
var googleAnalyticsKey = flag.String("google_analytics_key", "", "Google Analytics Account Id")

// ReportHandler handles the report page
func ReportHandler(w http.ResponseWriter, r *http.Request, db *badger.DB, repo string) {
	log.Printf("Displaying report: %q", repo)
	t := template.Must(template.New("report.html").Delims("[[", "]]").ParseFiles("templates/report.html", "templates/footer.html"))
	resp, err := getFromCache(db, repo)
	needToLoad := false
	if err != nil {
		switch err.(type) {
		case notFoundError:
			// don't bother logging - we already log in getFromCache. continue
		default:
			log.Println("ERROR ReportHandler:", err) // log error, but continue
		}
		needToLoad = true
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("ERROR ReportHandler: could not marshal JSON: ", err)
		http.Error(w, "Failed to load cache object", 500)
		return
	}

	t.Execute(w, map[string]interface{}{
		"repo":                 repo,
		"response":             string(respBytes),
		"loading":              needToLoad,
		"domain":               domain,
		"google_analytics_key": googleAnalyticsKey,
	})
}
