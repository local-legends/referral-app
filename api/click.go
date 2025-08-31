package api

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func Click(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/api/click/")
	if slug == "" {
		http.Error(w, "Slug required", 400)
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		http.Error(w, "Database connection failed", 500)
		return
	}
	defer db.Close()

	// Increment click count
	_, err = db.Exec("UPDATE links SET clicks = clicks + 1 WHERE slug = $1", slug)
	if err != nil {
		http.Error(w, "Update failed", 500)
		return
	}

	// Redirect to a destination (you can customize this)
	http.Redirect(w, r, "https://example.com", http.StatusFound)
}
