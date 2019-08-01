package main

import (
    "log"
    _"fmt"
    "strings"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _"github.com/jinzhu/gorm/dialects/sqlite"
)

type Achievement struct {
    gorm.Model
    Slug            string  `gorm:"unique" gorm:"default:'galeone'" json:"slug"`
    Title           string  `gorm:"default:'galeone'"  json:"title"`
    Desc            string  `gorm:"default:'galeone'" json:"desc"`
    Img             string  `gorm:"default:'galeone'" json:"img"`
}

type Response struct {
    Success         string  `json:"success"`
    Message         string  `json:"message"`
}

var DB *gorm.DB

func createAchievement(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    a := Achievement{}
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
    if err = DB.Save(&a).Error; err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

}

func main() {
    db, err := gorm.Open("sqlite3", "database.sqlite")
    if err != nil {
        panic("failed to connect database")
    }
    DB = db
    defer db.Close()
    db.AutoMigrate(&Achievement{})

    r := mux.NewRouter()
    r.HandleFunc("/achievements", createAchievement).Methods("POST")

    log.Fatal(http.ListenAndServe(":4242", r))
}
