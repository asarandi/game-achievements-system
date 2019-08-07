package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// get records where ... condition a b c
func getGameStats(w http.ResponseWriter, r *http.Request) {
	var model []Stat
	var id string = mux.Vars(r)["id0"]
	responseDb(w, getAllRecordsWhereABC(&model, "game_id = ?", id, nil), &model, http.StatusOK)
}

func getGameTeamStats(w http.ResponseWriter, r *http.Request) {
	var model []Stat
	var i, j string = mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseDb(w, getAllRecordsWhereABC(&model, "game_id = ? AND team_id = ?", i, j), &model, http.StatusOK)
}

func getGameMemberStats(w http.ResponseWriter, r *http.Request) {
	var model Stat
	var i, j string = mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseDb(w, getRecordWhereABC(&model, "game_id = ? AND member_id = ?", i, j), &model, http.StatusOK)
}

// restriction: member stats can only be updated if game status == startedGame
func updateGameMemberStats(w http.ResponseWriter, r *http.Request) {
	var game Game
	var oldRecord, newRecord Stat
	var i, j = mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	if err := db.First(&game, i).Error; err != nil {
		responseDb(w, err, nil, 0)
		return
	}
	if game.Status != startedGame {
		responseDb(w, errors.New("cannot update stats for this game"), nil, 0)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&newRecord); err != nil {
		responseDb(w, err, nil, 0);
		return
	}
	responseDb(w, updateRecordWhereABC(&oldRecord, &newRecord, "game_id = ? AND member_id = ?", i, j), &newRecord, http.StatusAccepted)
}

func createFromRequest(w http.ResponseWriter, r *http.Request, model interface{}) {
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		responseDb(w, err, nil, 0);
		return
	}
	responseDb(w, createFromModel(model), model, http.StatusCreated)
}

func updateFromRequest(w http.ResponseWriter, r *http.Request, model, newRecord, id interface{}) {
	if err := json.NewDecoder(r.Body).Decode(&newRecord); err != nil {
		responseDb(w, err, nil, 0);
		return
	}
	responseDb(w, updateRecordByID(model, newRecord, id), &newRecord, http.StatusAccepted)
}

// create records by model type
func createAchievement(w http.ResponseWriter, r *http.Request) {
	createFromRequest(w, r, &Achievement{})
}

func createMember(w http.ResponseWriter, r *http.Request) {
	createFromRequest(w, r, &Member{})
}

func createTeam(w http.ResponseWriter, r *http.Request) {
	createFromRequest(w, r, &Team{})
}

func createGame(w http.ResponseWriter, r *http.Request) {
	var model Game = Game{Status: newGame}
	responseDb(w, createFromModel(&model), &model, http.StatusCreated)
}

// get all records by model type
func getAllAchievements(w http.ResponseWriter, r *http.Request) {
	var model []Achievement
	responseDb(w, getAllRecords(&model), &model, 0)
}

func getAllMembers(w http.ResponseWriter, r *http.Request) {
	var model []Member
	responseDb(w, getAllRecords(&model), &model, 0)
}

func getAllTeams(w http.ResponseWriter, r *http.Request) {
	var model []Team
	responseDb(w, getAllRecords(&model), &model, 0)
}

func getAllGames(w http.ResponseWriter, r *http.Request) {
	var model []Game
	responseDb(w, getAllRecords(&model), &model, 0)
}

// get record by id
func getAchievement(w http.ResponseWriter, r *http.Request) {
	var model Achievement
	responseDb(w, getRecordByID(&model, mux.Vars(r)["id0"]), &model, 0);
}
func getMember(w http.ResponseWriter, r *http.Request) {
	var model Member
	responseDb(w, getRecordByID(&model, mux.Vars(r)["id0"]), &model, 0);
}
func getTeam(w http.ResponseWriter, r *http.Request) {
	var model Team
	responseDb(w, getRecordByID(&model, mux.Vars(r)["id0"]), &model, 0);
}

// update record by id
func updateAchievement(w http.ResponseWriter, r *http.Request) {
	updateFromRequest(w, r, &Achievement{}, &Achievement{}, mux.Vars(r)["id0"])
}

func updateMember(w http.ResponseWriter, r *http.Request) {
	updateFromRequest(w, r, &Member{}, &Member{}, mux.Vars(r)["id0"])
}
func updateTeam(w http.ResponseWriter, r *http.Request) {
	updateFromRequest(w, r, &Team{}, &Team{}, mux.Vars(r)["id0"])
}

// delete record by id
func deleteAchievement(w http.ResponseWriter, r *http.Request) {
	responseDb(w, deleteRecordByID(&Achievement{}, mux.Vars(r)["id0"]), nil, http.StatusNoContent)
}
func deleteMember(w http.ResponseWriter, r *http.Request) {
	responseDb(w, deleteRecordByID(&Member{}, mux.Vars(r)["id0"]), nil, http.StatusNoContent)
}
func deleteTeam(w http.ResponseWriter, r *http.Request) {
	responseDb(w, deleteRecordByID(&Team{}, mux.Vars(r)["id0"]), nil, http.StatusNoContent)
}

///// ?? XXX wtf

func getWinnersByGameID(gameID uint) []Member {
	var ids []uint
	stats := []Stat{}
	members := []Member{}
	db.Where("game_id = ? AND is_winner = ?", gameID, true).Find(&stats)
	for _, stat := range stats {
		ids = append(ids, stat.MemberID)
	}
	db.Where("ID IN (?)", ids).Find(&members)
	return members
}

func getGameWinners(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	game := Game{}
	if err := db.First(&game, vars["id0"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, &Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
		return
	}
	winners := getWinnersByGameID(game.ID)
	jsonResponse(w, &Response{true, http.StatusOK, "ok", &winners})
}
/////// fixme XXX

// get association records
func getMemberAchievements(w http.ResponseWriter, r *http.Request) {
	m := Member{}
	id := mux.Vars(r)["id0"]
	responseDb(w, findAssociationRecords(&m, id, "Achievements", &m.Achievements), &m.Achievements, 0)
}

func getAchievementMembers(w http.ResponseWriter, r *http.Request) {
	m := Achievement{}
	id := mux.Vars(r)["id0"]
	responseDb(w, findAssociationRecords(&m, id, "Members", &m.Members), &m.Members, 0)
}

func getMemberTeams(w http.ResponseWriter, r *http.Request) {
	m := Member{}
	id := mux.Vars(r)["id0"]
	responseDb(w, findAssociationRecords(&m, id, "Teams", &m.Teams), &m.Teams, 0)
}

func getTeamMembers(w http.ResponseWriter, r *http.Request) {
	m := Team{}
	id := mux.Vars(r)["id0"]
	responseDb(w, findAssociationRecords(&m, id, "Members", &m.Members), &m.Members, 0)
}

func getGameMembers(w http.ResponseWriter, r *http.Request) {
	m := Game{}
	id := mux.Vars(r)["id0"]
	responseDb(w, findAssociationRecords(&m, id, "Members", &m.Members), &m.Members, 0)
}

func getGameTeams(w http.ResponseWriter, r *http.Request) {
	m := Game{}
	id  := mux.Vars(r)["id0"]
	responseDb(w, findAssociationRecords(&m, id, "Teams", &m.Teams), &m.Teams, 0)
}

// add association records
func addTeamMember(w http.ResponseWriter, r *http.Request) {
	i, j := mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseDb(w, appendAssociationRecord(&Team{}, i, "Members", &Member{}, j), nil, 0)
}

func addMemberTeam(w http.ResponseWriter, r *http.Request) {
	i, j := mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseDb(w, appendAssociationRecord(&Member{}, i, "Teams", &Team{}, j), nil, 0)
}

// remove association records
func removeTeamMember(w http.ResponseWriter, r *http.Request) {
	i, j := mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseDb(w, appendAssociationRecord(&Team{}, i, "Members", &Member{}, j), nil, 0)
}

func removeMemberTeam(w http.ResponseWriter, r *http.Request) {
	i, j := mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseDb(w, deleteAssociationRecord(&Member{}, i, "Teams", &Team{}, j), nil, 0)
}
