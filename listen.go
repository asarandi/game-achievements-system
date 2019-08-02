package main

import (
    "log"
    "fmt"
    "strings"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _"github.com/jinzhu/gorm/dialects/sqlite"
)

type Achievement struct {
    gorm.Model
    Slug            string  `gorm:"unique" json:"slug"`
    Title           string  `json:"title"`
    Desc            string  `json:"desc"`
    Img             string  `json:"img"`
}

type Response struct {
    Success         bool        `json:"success"`
    Code            int         `json:"code"`
    Message         string      `json:"message"`
    Result          interface{} `json:"result"`
}

var DB *gorm.DB

func jsonResponse(w http.ResponseWriter, r Response) {
    w.Header().Set("content-type", "application/json")
    w.WriteHeader(r.Code)
    json.NewEncoder(w).Encode(r)
}

func createAchievement(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    a := Achievement{}
    err := decoder.Decode(&a)
    if err != nil {
        jsonResponse(w, Response{false, http.StatusBadRequest, err.Error(), nil})
        return
    }
    a.Slug = strings.TrimSpace(a.Slug)
    if a.Slug == "" {
        jsonResponse(w, Response{false, http.StatusBadRequest, "missing slug field in json", nil})
        return
    }
    if err = DB.Save(&a).Error; err != nil {
        jsonResponse(w, Response{false, http.StatusBadRequest, err.Error(), nil})
        return
    }
    jsonResponse(w, Response{true, http.StatusOK, "ok", nil})
}

func dbInit() *gorm.DB {
    db, err := gorm.Open("sqlite3", "database.sqlite")
    if err != nil {
        panic("failed to connect to database")
    }
    db.AutoMigrate(&Achievement{})
    return db
}

func main() {
    DB = dbInit()
    r := mux.NewRouter()
    r.HandleFunc("/achievements", createAchievement).Methods("POST")
    fmt.Println("listening on :4242")
    log.Fatal(http.ListenAndServe(":4242", r))
}
