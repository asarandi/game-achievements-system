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
    Slug            string      `gorm:"unique" json:"slug"`
    Name            string      `gorm:"unique" json:"name"`
    Desc            string      `json:"desc"`
    Img             string      `json:"img"`
}

type Member struct {
    gorm.Model
    Name            string      `gorm:"unique" json:"name"`
    Img             string      `json:"img"`
}

type Team struct {
    gorm.Model
    Name            string      `gorm:"unique" json:"name"`
    Img             string      `json:"img"`
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

func translateError(code int, msg string) (int, string) {
    switch {
    case strings.Contains(msg, "UNIQUE constraint failed"): msg = "record already exists"; code = http.StatusNotAcceptable
    case msg == "record not found": code = http.StatusNotFound
    }
    return code, msg
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

func getAllRecords(w http.ResponseWriter, r *http.Request, model interface{}) {
    if err := DB.Find(model).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    jsonResponse(w, Response{true, http.StatusOK, "ok", model})
}

// first find existing record to populate struct with all extra data: ID, CreatedAt, UpdatedAt, etc
// load updated data from request body into same struct and save 
func updateRecord(w http.ResponseWriter, r *http.Request, model interface{}) {
    vars := mux.Vars(r)
    id := vars["id"]
    if err := DB.First(model, id).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(model); err != nil {
        jsonResponse(w, Response{false, http.StatusBadRequest, err.Error(), nil})
        return
    }
    if err := DB.Save(model).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    jsonResponse(w, Response{true, http.StatusAccepted, "ok", nil})
}

func deleteRecord(w http.ResponseWriter, r *http.Request, model interface{}) {
    vars := mux.Vars(r)
    id := vars["id"]
    if err := DB.First(model, id).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    if err := DB.Delete(model).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    jsonResponse(w, Response{true, http.StatusNoContent, "ok", nil})
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

    r.HandleFunc("/achievements", createAchievement).Methods("POST")                // create new record
    r.HandleFunc("/achievements/{id:[0-9]+}", getAchievement).Methods("GET")        // get record by id
    r.HandleFunc("/achievements", getAllAchievements).Methods("GET")                // get all records
    r.HandleFunc("/achievements/{id:[0-9]+}", updateAchievement).Methods("PUT")     // update record
    r.HandleFunc("/achievements/{id:[0-9]+}", deleteAchievement).Methods("DELETE")  // delete record

    r.HandleFunc("/members", createMember).Methods("POST")
    r.HandleFunc("/members/{id:[0-9]+}", getMember).Methods("GET")
    r.HandleFunc("/members", getAllMembers).Methods("GET")
    r.HandleFunc("/members/{id:[0-9]+}", updateMember).Methods("PUT")
    r.HandleFunc("/members/{id:[0-9]+}", deleteMember).Methods("DELETE")

    r.HandleFunc("/teams", createTeam).Methods("POST")
    r.HandleFunc("/teams/{id:[0-9]+}", getTeam).Methods("GET")
    r.HandleFunc("/teams", getAllTeams).Methods("GET")
    r.HandleFunc("/teams/{id:[0-9]+}", updateTeam).Methods("PUT")
    r.HandleFunc("/teams/{id:[0-9]+}", deleteTeam).Methods("DELETE")

    fmt.Println("listening on :4242")
    log.Fatal(http.ListenAndServe(":4242", r))
}



// create new record
func createAchievement(w http.ResponseWriter, r *http.Request) {
    createRecord(w, r, &Achievement{})
}
func createMember(w http.ResponseWriter, r *http.Request) {
    createRecord(w, r, &Member{})
}
func createTeam(w http.ResponseWriter, r *http.Request) {
    createRecord(w, r, &Team{})
}

// get record by id
func getAchievement(w http.ResponseWriter, r *http.Request) {
    getRecord(w, r, &Achievement{})
}
func getMember(w http.ResponseWriter, r *http.Request) {
    getRecord(w, r, &Member{})
}
func getTeam(w http.ResponseWriter, r *http.Request) {
    getRecord(w, r, &Team{})
}

// get all records
func getAllAchievements(w http.ResponseWriter, r *http.Request) {
    getAllRecords(w, r, &[]Achievement{})
}
func getAllMembers(w http.ResponseWriter, r *http.Request) {
    getAllRecords(w, r, &[]Member{})
}
func getAllTeams(w http.ResponseWriter, r *http.Request) {
    getAllRecords(w, r, &[]Team{})
}

// update record by id
func updateAchievement(w http.ResponseWriter, r *http.Request) {
    updateRecord(w, r, &Achievement{})
}
func updateMember(w http.ResponseWriter, r *http.Request) {
    updateRecord(w, r, &Member{})
}
func updateTeam(w http.ResponseWriter, r *http.Request) {
    updateRecord(w, r, &Team{})
}

// delete record by id
func deleteAchievement(w http.ResponseWriter, r *http.Request) {
    deleteRecord(w, r, &Achievement{})
}
func deleteMember(w http.ResponseWriter, r *http.Request) {
    deleteRecord(w, r, &Member{})
}
func deleteTeam(w http.ResponseWriter, r *http.Request) {
    deleteRecord(w, r, &Team{})
}
