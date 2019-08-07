package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

// get records where ... condition a b c
func getGameStats(w http.ResponseWriter, r *http.Request) {
    v := mux.Vars(r)
    getRecordsWhereABC(w, r, &[]Stat{}, "game_id = ?", v["id0"], nil)
}

func getGameTeamStats(w http.ResponseWriter, r *http.Request) {
    v := mux.Vars(r)
    getRecordsWhereABC(w, r, &[]Stat{}, "game_id = ? AND team_id = ?", v["id0"], v["id1"])
}

func getGameMemberStats(w http.ResponseWriter, r *http.Request) {
    v := mux.Vars(r)
    getRecordsWhereABC(w, r, &[]Stat{}, "game_id = ? AND member_id = ?", v["id0"], v["id1"])
}

// restriction: member stats can only be updated if game status == StartedGame
func updateGameMemberStats(w http.ResponseWriter, r *http.Request) {
    v := mux.Vars(r)
    game := Game{}
    if err := DB.First(&game, v["id0"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, errorMessage, nil})
        return
    }
    if game.Status != StartedGame {
        jsonResponse(w, Response{false, http.StatusForbidden, "cannot update stats for this game", nil})
        return
    }
    updateRecordWhereABC(w, r, &Stat{}, "game_id = ? AND member_id = ?", v["id0"], v["id1"])
}

// create records by model type
func createAchievement(w http.ResponseWriter, r *http.Request) {
    createFromData(w, r, &Achievement{})
}

func createMember(w http.ResponseWriter, r *http.Request) {
    createFromData(w, r, &Member{})
}

func createTeam(w http.ResponseWriter, r *http.Request) {
    createFromData(w, r, &Team{})
}

func createGame(w http.ResponseWriter, r *http.Request) {
    createEmpty(w, r, &Game{Status: NewGame})
}

// get all records by model type
func getAllAchievements(w http.ResponseWriter, r *http.Request) {
    getAllRecords(w, r, &[]Achievement{})
}

func getAllMembers(w http.ResponseWriter, r *http.Request) {
    getAllRecords(w, r, &[]Member{})
}

func getAllTeams(w http.ResponseWriter, r *http.Request) {
    getAllRecords(w, r, &[]Team{})
}

func getAllGames(w http.ResponseWriter, r *http.Request) {
    getAllRecords(w, r, &[]Game{})
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

// update record by id
func updateAchievement(w http.ResponseWriter, r *http.Request) {
    updateRecordByID(w, r, &Achievement{})
}
func updateMember(w http.ResponseWriter, r *http.Request) {
    updateRecordByID(w, r, &Member{})
}
func updateTeam(w http.ResponseWriter, r *http.Request) {
    updateRecordByID(w, r, &Team{})
}

// delete record by id
func deleteAchievement(w http.ResponseWriter, r *http.Request) {
    deleteRecordByID(w, r, &Achievement{})
}
func deleteMember(w http.ResponseWriter, r *http.Request) {
    deleteRecordByID(w, r, &Member{})
}
func deleteTeam(w http.ResponseWriter, r *http.Request) {
    deleteRecordByID(w, r, &Team{})
}

// get association records
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

func getGameMembers(w http.ResponseWriter, r *http.Request) {
    model := Game{}
    getAssociationRecords(w, r, &model, "Members", &model.Members)
}

func getWinnersByGameID(gameID uint) []Member {
    var ids []uint
    stats := []Stat{}
    members := []Member{}
    DB.Where("game_id = ? AND is_winner = ?", gameID, true).Find(&stats)
    for _, stat := range stats {
        ids = append(ids, stat.MemberID)
    }
    DB.Where("ID IN (?)", ids).Find(&members)
    return members
}

func getGameWinners(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    game := Game{}
    if err := DB.First(&game, vars["id0"]).Error; err != nil {
        errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
        jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
        return
    }
    winners := getWinnersByGameID(game.ID)
    jsonResponse(w, Response{true, http.StatusOK, "ok", &winners})
}

func getGameTeams(w http.ResponseWriter, r *http.Request) {
    model := Game{}
    getAssociationRecords(w, r, &model, "Teams", &model.Teams)
}

// add association records
func addTeamMember(w http.ResponseWriter, r *http.Request) {
    addAssociationRecord(w, r, &Team{}, "Members", &Member{})
}

func addMemberTeam(w http.ResponseWriter, r *http.Request) {
    addAssociationRecord(w, r, &Member{}, "Teams", &Team{})
}

// remove association records
func removeTeamMember(w http.ResponseWriter, r *http.Request) {
    removeAssociationRecord(w, r, &Team{}, "Members", &Member{})
}

func removeMemberTeam(w http.ResponseWriter, r *http.Request) {
    removeAssociationRecord(w, r, &Member{}, "Teams", &Team{})
}
