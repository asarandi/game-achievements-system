package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

/*
	updateGameMemberStats():
		restriction: member stats can only be updated if game status == startedGame
*/
func updateGameMemberStats(w http.ResponseWriter, r *http.Request) {
	var game Game
	var oldRecord, newRecord Stat
	var i, j = mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	if err := db.First(&game, i).Error; err != nil {
		responseJson(w, err, nil, 0)
		return
	}
	if game.Status != startedGame {
		responseJson(w, errors.New("cannot update stats for this game"), nil, 0)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&newRecord); err != nil {
		responseJson(w, err, nil, 0)
		return
	}
	if err := db.Where("game_id = ? AND member_id = ?", i, j).First(&oldRecord).Error; err != nil {
		responseJson(w, err, nil, 0)
		return
	}
	responseJson(w, db.Model(&oldRecord). /* restrict access to some fields via Omit() */
		Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt", "GameID", "TeamID", "MemberID", "IsWinner").
		Updates(&newRecord).Error,&newRecord, http.StatusAccepted)
}

func getGameWinners(w http.ResponseWriter, r *http.Request) {
	var gameID = mux.Vars(r)["id0"]
	var members []Member
	responseJson(w,
		db.Joins("JOIN stats ON stats.member_id = members.id AND stats.game_id = ? AND stats.is_winner = ?",
			gameID, true).Find(&members).Error,
		&members, 0)
}

/*
	get records where ... condition a b c
*/
func getGameStats(w http.ResponseWriter, r *http.Request) {
	var model []Stat
	responseJson(w, getAllRecordsWhereABC(&model, "game_id = ?", mux.Vars(r)["id0"], nil), &model, http.StatusOK)
}

func getGameTeamStats(w http.ResponseWriter, r *http.Request) {
	var model []Stat
	var i, j = mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseJson(w, getAllRecordsWhereABC(&model, "game_id = ? AND team_id = ?", i, j), &model, http.StatusOK)
}

func getGameMemberStats(w http.ResponseWriter, r *http.Request) {
	var model Stat
	var i, j = mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseJson(w, getRecordWhereABC(&model, "game_id = ? AND member_id = ?", i, j), &model, http.StatusOK)
}

/*
	create records by model type
*/
func createFromRequest(w http.ResponseWriter, r *http.Request, model interface{}) {
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		responseJson(w, err, nil, 0)
		return
	}
	responseJson(w, createFromModel(model), model, http.StatusCreated)
}

func createAchievement(w http.ResponseWriter, r *http.Request) {
	createFromRequest(w, r, &Achievement{})
}

func createMember(w http.ResponseWriter, r *http.Request) {
	createFromRequest(w, r, &Member{})
}

func createTeam(w http.ResponseWriter, r *http.Request) {
	createFromRequest(w, r, &Team{})
}

func createGame(w http.ResponseWriter, _ *http.Request) {
	var model = Game{Status: newGame}
	responseJson(w, createFromModel(&model), &model, http.StatusCreated)
}

/*
	get all records by model type
*/
func getAllAchievements(w http.ResponseWriter, _ *http.Request) {
	var model []Achievement
	responseJson(w, getAllRecords(&model), &model, 0)
}

func getAllMembers(w http.ResponseWriter, _ *http.Request) {
	var model []Member
	responseJson(w, getAllRecords(&model), &model, 0)
}

func getAllTeams(w http.ResponseWriter, _ *http.Request) {
	var model []Team
	responseJson(w, getAllRecords(&model), &model, 0)
}

func getAllGames(w http.ResponseWriter, _ *http.Request) {
	var model []Game
	responseJson(w, getAllRecords(&model), &model, 0)
}

/*
	get record by id
*/
func getAchievement(w http.ResponseWriter, r *http.Request) {
	var model Achievement
	responseJson(w, getRecordByID(&model, mux.Vars(r)["id0"]), &model, 0)
}
func getMember(w http.ResponseWriter, r *http.Request) {
	var model Member
	responseJson(w, getRecordByID(&model, mux.Vars(r)["id0"]), &model, 0)
}
func getTeam(w http.ResponseWriter, r *http.Request) {
	var model Team
	responseJson(w, getRecordByID(&model, mux.Vars(r)["id0"]), &model, 0)
}

/*
	update record by id
*/
func updateFromRequest(w http.ResponseWriter, r *http.Request, model, newRecord, id interface{}) {
	if err := json.NewDecoder(r.Body).Decode(&newRecord); err != nil {
		responseJson(w, err, nil, 0)
		return
	}
	responseJson(w, updateRecordByID(model, newRecord, id), &newRecord, http.StatusAccepted)
}

func updateAchievement(w http.ResponseWriter, r *http.Request) {
	updateFromRequest(w, r, &Achievement{}, &Achievement{}, mux.Vars(r)["id0"])
}

func updateMember(w http.ResponseWriter, r *http.Request) {
	updateFromRequest(w, r, &Member{}, &Member{}, mux.Vars(r)["id0"])
}
func updateTeam(w http.ResponseWriter, r *http.Request) {
	updateFromRequest(w, r, &Team{}, &Team{}, mux.Vars(r)["id0"])
}

/*
	delete record by id
*/
func deleteAchievement(w http.ResponseWriter, r *http.Request) {
	responseJson(w, deleteRecordByID(&Achievement{}, mux.Vars(r)["id0"]), nil, http.StatusNoContent)
}
func deleteMember(w http.ResponseWriter, r *http.Request) {
	responseJson(w, deleteRecordByID(&Member{}, mux.Vars(r)["id0"]), nil, http.StatusNoContent)
}
func deleteTeam(w http.ResponseWriter, r *http.Request) {
	responseJson(w, deleteRecordByID(&Team{}, mux.Vars(r)["id0"]), nil, http.StatusNoContent)
}

/*
	get association records
*/
func getMemberAchievements(w http.ResponseWriter, r *http.Request) {
	m := Member{}
	id := mux.Vars(r)["id0"]
	responseJson(w,
		findAssociationRecords(&m, id, "Achievements", &m.Achievements), &m.Achievements, 0)
}

func getAchievementMembers(w http.ResponseWriter, r *http.Request) {
	m := Achievement{}
	id := mux.Vars(r)["id0"]
	responseJson(w, findAssociationRecords(&m, id, "Members", &m.Members), &m.Members, 0)
}

func getMemberTeams(w http.ResponseWriter, r *http.Request) {
	m := Member{}
	id := mux.Vars(r)["id0"]
	responseJson(w, findAssociationRecords(&m, id, "Teams", &m.Teams), &m.Teams, 0)
}

func getTeamMembers(w http.ResponseWriter, r *http.Request) {
	m := Team{}
	id := mux.Vars(r)["id0"]
	responseJson(w, findAssociationRecords(&m, id, "Members", &m.Members), &m.Members, 0)
}

func getGameMembers(w http.ResponseWriter, r *http.Request) {
	m := Game{}
	id := mux.Vars(r)["id0"]
	responseJson(w, findAssociationRecords(&m, id, "Members", &m.Members), &m.Members, 0)
}

func getGameTeams(w http.ResponseWriter, r *http.Request) {
	m := Game{}
	id  := mux.Vars(r)["id0"]
	responseJson(w, findAssociationRecords(&m, id, "Teams", &m.Teams), &m.Teams, 0)
}

/*
	add association records
*/
func addTeamMember(w http.ResponseWriter, r *http.Request) {
	i, j := mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseJson(w, appendAssociationRecord(&Team{}, i, "Members", &Member{}, j), nil, 0)
}

func addMemberTeam(w http.ResponseWriter, r *http.Request) {
	i, j := mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseJson(w, appendAssociationRecord(&Member{}, i, "Teams", &Team{}, j), nil, 0)
}

/*
	remove association records
*/
func removeTeamMember(w http.ResponseWriter, r *http.Request) {
	i, j := mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseJson(w, appendAssociationRecord(&Team{}, i, "Members", &Member{}, j), nil, 0)
}

func removeMemberTeam(w http.ResponseWriter, r *http.Request) {
	i, j := mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	responseJson(w, deleteAssociationRecord(&Member{}, i, "Teams", &Team{}, j), nil, 0)
}
