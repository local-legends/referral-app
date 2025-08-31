package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type Link struct {
	ID         int    `json:"id"`
	Slug       string `json:"slug"`
	ShortURL   string `json:"shortUrl"`
	ReferrerID string `json:"referrerId"`
	Clicks     int    `json:"clicks"`
}

func Links(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		http.Error(w, "Database connection failed", 500)
		return
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id SERIAL PRIMARY KEY,
			slug VARCHAR(255) UNIQUE NOT NULL,
			referrer_id VARCHAR(255) NOT NULL,
			clicks INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		http.Error(w, "Table creation failed", 500)
		return
	}

	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT id, slug, referrer_id, clicks FROM links")
		if err != nil {
			http.Error(w, "Query failed", 500)
			return
		}
		defer rows.Close()

		var links []Link
		for rows.Next() {
			var link Link
			err := rows.Scan(&link.ID, &link.Slug, &link.ReferrerID, &link.Clicks)
			if err != nil {
				continue
			}
			link.ShortURL = "https://" + r.Host + "/" + link.Slug
			links = append(links, link)
		}

		json.NewEncoder(w).Encode(links)

	case "POST":
		var link Link
		err := json.NewDecoder(r.Body).Decode(&link)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		err = db.QueryRow(
			"INSERT INTO links (slug, referrer_id) VALUES ($1, $2) RETURNING id",
			link.Slug, link.ReferrerID,
		).Scan(&link.ID)
		if err != nil {
			http.Error(w, "Insert failed", 500)
			return
		}

		link.ShortURL = "https://" + r.Host + "/" + link.Slug
		link.Clicks = 0
		json.NewEncoder(w).Encode(link)
	}
}
