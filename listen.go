package main

import (
    "log"
    "fmt"
    "strings"
    "net/http"
    "database/sql"
    "encoding/json"
    "github.com/gorilla/mux"
    _"github.com/mattn/go-sqlite3"
)

type Achievement struct {
    Slug            string  `json:"slug"`
    Title           string  `json:"title"`
    Desc            string  `json:"desc"`
    Img             string  `json:"img"`
}

type Response struct {
    Success         string  `json:"success"`
    Message         string  `json:"message"`
}

func dbInit() {
    db, err := sql.Open("sqlite3", "file:./database.sqlite?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
    defer db.Close()

	sqlStmt := `
	create table if not exists achievements        (id integer primary key, slug text, title text, desc text, img text, created datetime);
	create table if not exists members             (id integer primary key, name text, img text, created datetime);
	create table if not exists member_achievements (id integer, achievement_id integer, created datetime);
	create table if not exists teams               (id integer primary key, name text, img text, created datetime);
	create table if not exists team_members        (id integer, member_id integer, created datetime);
	create table if not exists games               (id integer primary key, created datetime);    
	create table if not exists game_stats          (id integer, team_id integer, member_id integer, num_attacks integer, num_hits integer, amount_damage integer, num_kills integer, num_assists integer, num_spells integer, spells_damage integer, finished datetime, is_winner integer);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}


func createAchievement(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var a Achievement
    err := decoder.Decode(&a)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    a.Slug = strings.TrimSpace(a.Slug)
    if a.Slug == "" {
        http.Error(w, "missing \"slug\" field in JSON", http.StatusBadRequest)
        return
    }
    w.Header().Set("content-type", "application/json")
    fmt.Println(a)
}

func main() {
    dbInit()
    r := mux.NewRouter()
    r.HandleFunc("/achievements", createAchievement).Methods("POST")

    log.Fatal(http.ListenAndServe(":4242", r))
}
