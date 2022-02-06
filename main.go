package main

import (
	"encoding/json"
	"log"
	"net/http"
    "fmt"
	"github.com/gorilla/mux"
    "database/sql"
    _ "github.com/lib/pq"
    "context"


	"github.com/georgysavva/scany/sqlscan"
)

// "Item type" (tipo um objeto)
type Item struct {
    ID        string   `json:"id"`
    Date string   `json:"date"`
    Category  string   `json:"category"`
    Title   string `json:"title"`
    Value float32 `json:"value"`
}

var people []Item

// Get mostra todos os contatos da variável people
func Get(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()

    db := OpenConnection()

	var items []*Item
	sqlscan.Select(ctx, db, &items, `SELECT * FROM item`)

	peopleBytes, _ := json.MarshalIndent(items, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer db.Close()
}

// GetItem mostra apenas um contato
func GetItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    ctx := context.Background()

    db := OpenConnection()

	var items []*Item
	sqlscan.Select(ctx, db, &items, `SELECT * FROM item where id = ` + params["id"])

	peopleBytes, _ := json.MarshalIndent(items, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer db.Close()
}

// CreateItem cria um novo contato
func CreateItem(w http.ResponseWriter, r *http.Request) {
    db := OpenConnection()

	var p Item
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO item (date, category, title, value) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, p.Date, p.Category, p.Title, p.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

// DeleteItem deleta um contato
func DeleteItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    db := OpenConnection()

	var p Item
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `DELETE FROM item where id = ` + params["id"]
	_, err = db.Exec(sqlStatement)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
}

const (
	host     = "backend-postgres-12.6"
	port     = 5432
	user     = "admin"
	password = "admin"
	dbname   = "financas_db"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
    user, password, host, port, dbname)


	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

// função principal para executar a api
func main() {
    router := mux.NewRouter()

    router.HandleFunc("/", Get).Methods("GET")
    router.HandleFunc("/{id}", GetItem).Methods("GET")
    router.HandleFunc("/", CreateItem).Methods("POST")
    router.HandleFunc("/{id}", DeleteItem).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8000", router))
}
