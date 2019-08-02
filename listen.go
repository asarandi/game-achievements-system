package main

import (
    "log"
    "fmt"
    _"reflect"
    _"strings"
    _"unsafe"
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

type Member struct {
    gorm.Model
    Name            string  `gorm:"unique" json:"name"`
    Img             string  `json:"img"`
}

type Team struct {
    gorm.Model
    Name            string  `gorm:"unique" json:"name"`
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

func createRecord(w http.ResponseWriter, r *http.Request, model interface{}) {
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(model); err != nil {
        jsonResponse(w, Response{false, http.StatusBadRequest, err.Error(), nil})
        return
    }
    if err := DB.Create(model).Error; err != nil {
        jsonResponse(w, Response{false, http.StatusInternalServerError, err.Error(), nil})
        return
    }
    jsonResponse(w, Response{true, http.StatusOK, "ok", nil})
}

func createAchievement(w http.ResponseWriter, r *http.Request) {
    createRecord(w, r, &Achievement{})
}

func createMember(w http.ResponseWriter, r *http.Request) {
    createRecord(w, r, &Member{})
}

func createTeam(w http.ResponseWriter, r *http.Request) {
    createRecord(w, r, &Team{})
}

func dbInit() *gorm.DB {
    db, err := gorm.Open("sqlite3", "database.sqlite")
    if err != nil {
        panic("failed to connect to database")
    }
    db.AutoMigrate(&Achievement{}, &Member{}, &Team{})
    return db
}

func main() {
    DB = dbInit()
    r := mux.NewRouter()
    r.HandleFunc("/achievements", createAchievement).Methods("POST")
    r.HandleFunc("/members", createMember).Methods("POST")
    r.HandleFunc("/teams", createTeam).Methods("POST")
    fmt.Println("listening on :4242")
    log.Fatal(http.ListenAndServe(":4242", r))
}
