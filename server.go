package main

import (
    "log"
    "fmt"
    "time"
    "strings"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _"github.com/jinzhu/gorm/dialects/sqlite"
)

type Achievement struct {
    gorm.Model
    Slug            string          `gorm:"unique" json:"slug"`
    Name            string          `gorm:"unique" json:"name"`
    Desc            string          `json:"desc"`
    Img             string          `json:"img"`
    Members         []Member        `gorm:"many2many:member_achievements;" json:"members,omitempty"`
}

type Member struct {
    gorm.Model
    Name            string          `gorm:"unique" json:"name"`
    Img             string          `json:"img"`
    Achievements    []Achievement   `gorm:"many2many:member_achievements;" json:"achievements,omitempty"`
    Teams           []Team          `gorm:"many2many:team_members;" json:"teams,omitempty"`
    Games           []Game          `gorm:"many2many game_members;" json:"games,omitempty"`
    Stats           []Stat          `json:"stats,omitempty"`
}

type Team struct {
    gorm.Model
    Name            string          `gorm:"unique" json:"name"`
    Img             string          `json:"img"`
    Members         []Member        `gorm:"many2many:team_members;" json:"members,omitempty"`
    Games           []Game          `gorm:"many2many game_teams;" json:"games,omitempty"`
    Stats           []Stat          `json:"stats,omitempty"`
}

type GameStatus     int
const (
        NewGame         GameStatus = iota
        StartedGame
        FinishedGame
    )

type Game struct {
    gorm.Model
    Status          GameStatus      `json:"status"`
    Teams           []Team          `gorm:"many2many:game_teams;" json:"teams,omitempty"`
    Members         []Member        `gorm:"many2many:game_members;" json:"members,omitempty"`
    Stats           []Stat          `json:"stats,omitempty"`
    TeamID          int             `json:"team_id,omitempty"`
    Winner          Team            `json:"winner,omitempty"`
}

type Stat struct {
    gorm.Model
    GameID          int             `json:"game_id,omitempty"`
    Game            Game            `json:"game,omitempty"`
    TeamID          int             `json:"team_id,omitempty"`
    Team            Team            `json:"team,omitempty"`
    MemberID        int             `json:"member_id,omitempty"`
    Member          Member          `json:"member,omitempty"`
    Started         time.Time       `json:"started,omitempty"`
    Finished        time.Time       `json:"finished,omitempty"`
    NumAttacks      int             `json:"num_attacks,omitempty"`
    NumHits         int             `json:"num_hits,omitempty"`
    AmountDamage    int             `json:"amount_damage,omitempty"`
    NumKills        int             `json:"num_kills,omitempty"`
    InstantKills    int             `json:"instant_kills,omitempty"`
    NumAssists      int             `json:"num_assists,omitempty"`
    NumSpells       int             `json:"num_spells,omitempty"`
    SpellsDamage    int             `json:"spells_damage,omitempty"`
}


type Response struct {
    Success         bool            `json:"success"`
    Code            int             `json:"code"`
    Message         string          `json:"message"`
    Result          interface{}     `json:"result"`
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

func createEmpty(w http.ResponseWriter, r *http.Request, model interface{}) {
    if err := DB.Create(model).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    jsonResponse(w, Response{true, http.StatusCreated, "ok", model})
}

func createFromData(w http.ResponseWriter, r *http.Request, model interface{}) {
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
    jsonResponse(w, Response{true, http.StatusCreated, "ok", model})
}

func getRecordByID(w http.ResponseWriter, r *http.Request, model interface{}) {
    vars := mux.Vars(r)
    if err := DB.First(model, vars["id"]).Error; err != nil {
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

func getRecordsWhere(w http.ResponseWriter, r *http.Request, model interface{}, where ...interface{}) {
    if err := DB.Find(model, where).Error; err != nil {
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
    if err := DB.First(model, vars["id"]).Error; err != nil {
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
    jsonResponse(w, Response{true, http.StatusAccepted, "ok", model})
}

func deleteRecord(w http.ResponseWriter, r *http.Request, model interface{}) {
    vars := mux.Vars(r)
    if err := DB.First(model, vars["id"]).Error; err != nil {
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
    db.AutoMigrate(&Achievement{}, &Member{}, &Team{}, &Game{}, &Stat{})
    return db
}

func getAssociationRecords(w http.ResponseWriter, r *http.Request, a interface{}, b string, c interface{}) {
    vars := mux.Vars(r)
    if err := DB.First(a, vars["id"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    DB.Model(a).Association(b).Find(c)
    jsonResponse(w, Response{true, http.StatusOK, "ok", c})
}


func addAssociationRecord(w http.ResponseWriter, r *http.Request, a interface{}, b string, c interface{}) {
    vars := mux.Vars(r)
    if err := DB.First(a, vars["id0"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
        return
    }
    if err := DB.First(c, vars["id1"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id1"], errorMessage), nil})
        return
    }
    DB.Model(a).Association(b).Append(c)
    jsonResponse(w, Response{true, http.StatusOK, "ok", nil})
}

func removeAssociationRecord(w http.ResponseWriter, r *http.Request, a interface{}, b string, c interface{}) {
    vars := mux.Vars(r)
    if err := DB.First(a, vars["id0"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
        return
    }
    if err := DB.First(c, vars["id1"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id1"], errorMessage), nil})
        return
    }
    DB.Model(a).Association(b).Delete(c)
    jsonResponse(w, Response{true, http.StatusOK, "ok", nil})
}

func getMemberAchievements(w http.ResponseWriter, r *http.Request) {
    model := Member{}
    getAssociationRecords(w, r, &model, "Achievements", &model.Achievements)
}

func getAchievementMembers(w http.ResponseWriter, r *http.Request) {
    model := Achievement{}
    getAssociationRecords(w, r, &model, "Members", &model.Members)
}

func getMemberTeams(w http.ResponseWriter, r *http.Request) {
    model := Member{}
    getAssociationRecords(w, r, &model, "Teams", &model.Teams)
}

func getTeamMembers(w http.ResponseWriter, r *http.Request) {
    model := Team{}
    getAssociationRecords(w, r, &model, "Members", &model.Members)
}

func addTeamMember(w http.ResponseWriter, r *http.Request) {
    addAssociationRecord(w, r, &Team{}, "Members", &Member{})
}

func removeTeamMember(w http.ResponseWriter, r *http.Request) {
    removeAssociationRecord(w, r, &Team{}, "Members", &Member{})
}

func addMemberTeam(w http.ResponseWriter, r *http.Request) {
    addAssociationRecord(w, r, &Member{}, "Teams", &Team{})
}

func removeMemberTeam(w http.ResponseWriter, r *http.Request) {
    removeAssociationRecord(w, r, &Member{}, "Teams", &Team{})
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

    r.HandleFunc("/memberAchievements/{id:[0-9]+}", getMemberAchievements).Methods("GET")
    r.HandleFunc("/achievementMembers/{id:[0-9]+}", getAchievementMembers).Methods("GET")

    r.HandleFunc("/teamMembers/{id:[0-9]+}", getTeamMembers).Methods("GET")                         // list all members for a given team
    r.HandleFunc("/teamMembers/{id0:[0-9]+}/{id1:[0-9]+}", addTeamMember).Methods("POST")           // a team (id0) adds a member (id1)
    r.HandleFunc("/teamMembers/{id0:[0-9]+}/{id1:[0-9]+}", removeTeamMember).Methods("DELETE")      // a team (id0) removes a member (id1)

    r.HandleFunc("/memberTeams/{id:[0-9]+}", getMemberTeams).Methods("GET")                         // list all teams for a given member
    r.HandleFunc("/memberTeams/{id0:[0-9]+}/{id1:[0-9]+}", addMemberTeam).Methods("POST")           // a member (id0) joins a team (id1)
    r.HandleFunc("/memberTeams/{id0:[0-9]+}/{id1:[0-9]+}", removeMemberTeam).Methods("DELETE")      // a member (id0) leaves a team (id1)


    r.HandleFunc("/games", createGame).Methods("POST")
    r.HandleFunc("/games/all", getAllGames).Methods("GET")
    r.HandleFunc("/games/new", getNewGames).Methods("GET")
    r.HandleFunc("/games/started", getStartedGames).Methods("GET")
    r.HandleFunc("/games/finished", getFinishedGames).Methods("GET")


    // a team (id1) joins a game session (id0)
    r.HandleFunc("/gameTeams/{id0:[0-9]+}/{id1:[0-9]+}", addGameTeam).Methods("POST")

//
//    r.HandleFunc("/gameStats/{id:[0-9]+}", getGameStats).Methods("GET")                             // list all stats for a given game
//    r.HandleFunc("/teamStats/{id:[0-9]+}", getTeamStats).Methods("GET")                             // list all stats for a given team
//    r.HandleFunc("/memberStats/{id:[0-9]+}", getMemberStats).Methods("GET")                         // list all stats for a given member


    r.HandleFunc("/teams", createTeam).Methods("POST")
    r.HandleFunc("/teams/{id:[0-9]+}", getTeam).Methods("GET")
    r.HandleFunc("/teams", getAllTeams).Methods("GET")
    r.HandleFunc("/teams/{id:[0-9]+}", updateTeam).Methods("PUT")
    r.HandleFunc("/teams/{id:[0-9]+}", deleteTeam).Methods("DELETE")

    fmt.Println("listening on :4242")
    log.Fatal(http.ListenAndServe(":4242", r))
}


//
// logic: 
//    - a game can only have two teams
//    - once a new game is created it has 0 teams and its status is `NewGame`
//    - first team to join must have between 3 and 5 members, otherwise error
//    - second team to join must have same number of members as first team, otherwise error
//    - the two teams cannot be the same
//    - the two teams cannot have shared members
//    - once second team joins game status is changed to `StartedGame`
//

const MinNumMembers = 3
const MaxNumMembers = 5

func haveSharedMembers(a []Member, b []Member) bool {
    for _, i := range a {
        for _, j := range b {
            if i.ID == j.ID {
                return true
            }
        }
    }
    return false
}

/// work in progress

func addGameTeam(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    game := Game{}
    if err := DB.First(&game, vars["id0"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
        return
    }
    team := Team{}
    if err := DB.First(&team, vars["id1"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id1"], errorMessage), nil})
        return
    }
    if (game.Status != NewGame) {
        jsonResponse(w, Response{false, http.StatusForbidden, fmt.Sprintf("game status must be %d", NewGame), nil})
        return
    }
    members := []Member{}
    DB.Model(&team).Association("Members").Find(&members)
    if len(members) < MinNumMembers || len(members) > MaxNumMembers {
        s := fmt.Sprintf("team must contain between %d and %d members", MinNumMembers, MaxNumMembers)
        jsonResponse(w, Response{false, http.StatusForbidden, s, nil})
        return
    }

    if  DB.Model(&game).Association("Teams").Count() == 0 {  // this team is first team
        DB.Model(&game).Association("Teams").Append(&team)
        DB.Model(&game).Association("Members").Append(&members)
        jsonResponse(w, Response{true, http.StatusOK, "first team ok", nil})
        return
    }
    prevTeam := Team{}
    prevMembers := []Member{}
    DB.Model(&game).Association("Teams").Find(&prevTeam)
    DB.Model(&game).Association("Members").Find(&prevMembers)
    if prevTeam.ID == team.ID {
        jsonResponse(w, Response{false, http.StatusForbidden, "teams cannot be the same", nil})
        return
    }
    if len(prevMembers) != len(members) {
        jsonResponse(w, Response{false, http.StatusForbidden, "teams must have same number of players", nil})
        return
    }
    if haveSharedMembers(prevMembers, members) {
        jsonResponse(w, Response{false, http.StatusForbidden, "teams cannot have shared members", nil})
        return
    }
    DB.Model(&game).Association("Teams").Append(&team)
    DB.Model(&game).Association("Members").Append(&members)
    game.Status = StartedGame
    DB.Save(&game)
    fmt.Printf("game: +%v\n", game)
    jsonResponse(w, Response{true, http.StatusOK, "second team ok", nil})
}


func createGame(w http.ResponseWriter, r *http.Request) {
    createEmpty(w, r, &Game{})
}

func getAllGames(w http.ResponseWriter, r *http.Request) {
    getAllRecords(w, r, &[]Game{})
}

func getNewGames(w http.ResponseWriter, r *http.Request) {
    getRecordsWhere(w, r, &[]Game{}, "status = ?", NewGame)
}

func getStartedGames(w http.ResponseWriter, r *http.Request) {
    getRecordsWhere(w, r, &[]Game{}, "status = ?", StartedGame)
}

func getFinishedGames(w http.ResponseWriter, r *http.Request) {
    getRecordsWhere(w, r, &[]Game{}, "status = ?", FinishedGame)
}

// create new record
func createAchievement(w http.ResponseWriter, r *http.Request) {
    createFromData(w, r, &Achievement{})
}
func createMember(w http.ResponseWriter, r *http.Request) {
    createFromData(w, r, &Member{})
}
func createTeam(w http.ResponseWriter, r *http.Request) {
    createFromData(w, r, &Team{})
}

// get record by id
func getAchievement(w http.ResponseWriter, r *http.Request) {
    getRecordByID(w, r, &Achievement{})
}
func getMember(w http.ResponseWriter, r *http.Request) {
    getRecordByID(w, r, &Member{})
}
func getTeam(w http.ResponseWriter, r *http.Request) {
    getRecordByID(w, r, &Team{})
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
