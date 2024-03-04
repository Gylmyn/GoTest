package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Item struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price int       `json:"price"`
}

func main() {
	connStr := "postgresql://Gylmyn:kLUA2E6lYWam@ep-delicate-pond-24139553.ap-southeast-1.aws.neon.tech/mobile_db?sslmode=require"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Setup routes
	http.HandleFunc("/", getAllItems)
	http.HandleFunc("/mytabel/add", addItem)
	http.HandleFunc("/mytabel/delete", deleteItem)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM myitems")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	json.NewEncoder(w).Encode(items)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM myitems WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Item with ID %s deleted successfully\n", id)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New()

	_, err := db.Exec("INSERT INTO myitems(id, name, price) VALUES($1, $2, $3)", id, item.Name, item.Price)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Item added: %+v\n", item)
}
