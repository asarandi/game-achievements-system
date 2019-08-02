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
    Name            string  `gorm:"unique" json:"name"`
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
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    jsonResponse(w, Response{true, http.StatusCreated, "ok", nil})
}

func translateError(code int, msg string) (int, string) {
    switch {
    case strings.Contains(msg, "UNIQUE constraint failed"): msg = "record already exists"; code = http.StatusNotAcceptable
    case msg == "record not found": code = http.StatusNotFound
    }
    return code, msg
}

func getRecord(w http.ResponseWriter, r *http.Request, model interface{}) {
    vars := mux.Vars(r)
    id := vars["id"]
    if err := DB.First(model, id).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    jsonResponse(w, Response{true, http.StatusOK, "ok", model})
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

func getAchievement(w http.ResponseWriter, r *http.Request) {
    getRecord(w, r, &Achievement{})
}

func getMember(w http.ResponseWriter, r *http.Request) {
    getRecord(w, r, &Member{})
}

func getTeam(w http.ResponseWriter, r *http.Request) {
    getRecord(w, r, &Team{})
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
    r.HandleFunc("/achievements/{id:[0-9]+}", getAchievement).Methods("GET")
    r.HandleFunc("/members", createMember).Methods("POST")
    r.HandleFunc("/members/{id:[0-9]+}", getMember).Methods("GET")
    r.HandleFunc("/teams", createTeam).Methods("POST")
    r.HandleFunc("/teams/{id:[0-9]+}", getTeam).Methods("GET")

    fmt.Println("listening on :4242")
    log.Fatal(http.ListenAndServe(":4242", r))
}
