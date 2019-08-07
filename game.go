package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const gameMinNumMembers = 3
const gameMaxNumMembers = 5

type gameStatus int

const (
	newGame gameStatus = iota + 1
	pendingGame
	startedGame
	finishedGame
)

func setGameWinner(game *Game) {
	teams := []Team{}
	stats := []Stat{}
	calc := [2]Stat{}
	var idx int
	db.Model(game).Association("Teams").Find(&teams)
	db.Model(game).Association("Stats").Find(&stats)
	for _, stat := range stats {
		idx = 0
		if stat.TeamID == teams[1].ID {
			idx = 1
		}
		calc[idx].NumAttacks += stat.NumAttacks
		calc[idx].NumHits += stat.NumHits
		calc[idx].AmountDamage += stat.AmountDamage
		calc[idx].NumKills += stat.NumKills
		calc[idx].InstantKills += stat.InstantKills
		calc[idx].NumAssists += stat.NumAssists
		calc[idx].NumSpells += stat.NumSpells
		calc[idx].SpellsDamage += stat.SpellsDamage
	}
	idx = -1
	switch {
	case calc[0].NumKills > calc[1].NumKills:
		idx = 0
	case calc[0].NumKills < calc[1].NumKills:
		idx = 1
	case calc[0].AmountDamage > calc[1].AmountDamage:
		idx = 0
	case calc[0].AmountDamage < calc[1].AmountDamage:
		idx = 1
	case calc[0].NumHits > calc[1].NumHits:
		idx = 0
	case calc[0].NumHits < calc[1].NumHits:
		idx = 1
	case calc[0].NumAttacks > calc[1].NumAttacks:
		idx = 0
	case calc[0].NumAttacks < calc[1].NumAttacks:
		idx = 1
	case calc[0].InstantKills > calc[1].InstantKills:
		idx = 0
	case calc[0].InstantKills < calc[1].InstantKills:
		idx = 1
	case calc[0].NumAssists > calc[1].NumAssists:
		idx = 0
	case calc[0].NumAssists < calc[1].NumAssists:
		idx = 1
	case calc[0].SpellsDamage > calc[1].SpellsDamage:
		idx = 0
	case calc[0].SpellsDamage < calc[1].SpellsDamage:
		idx = 1
	case calc[0].NumSpells > calc[1].NumSpells:
		idx = 0
	case calc[0].NumSpells < calc[1].NumSpells:
		idx = 1
	}
	if idx == -1 { // tie
		return
	}
	for _, stat := range stats {
		if stat.TeamID == teams[idx].ID {
			stat.IsWinner = true
		} else {
			stat.IsWinner = false
		}
		db.Save(&stat)
	}
}

func setMemberAchievements(game *Game) {
	members := getWinnersByGameID(game.ID)
	stats := []Stat{}
	for i := range members {
		db.Model(&members[i]).Association("Achievements").Find(&members[i].Achievements)
		stat := Stat{GameID: game.ID, MemberID: members[i].ID}
		db.First(&stat, &stat)
		stats = append(stats, stat)
	}
	for i := range asf {
		if err := db.First(&asf[i].achievement, &Achievement{Slug: asf[i].slug}).Error; err != nil {
			asf[i].achievement.ID = 0
		}
	}
	for i := range members {
		for j := range asf {
			if asf[j].achievement.ID == 0 {
				continue
			} // failed to preload
			if isAwardedAlready(members[i].Achievements, asf[j].achievement) {
				continue
			}
			if !asf[j].function(stats[i]) {
				continue
			}
			db.Model(&members[i]).Association("Achievements").Append(asf[j].achievement)
		}
	}
}

func endGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	game := Game{}
	if err := db.First(&game, vars["id0"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
		return
	}
	if game.Status != startedGame {
		jsonResponse(w, Response{false, http.StatusForbidden, "cannot change status of this game", nil})
		return
	}
	game.Status = finishedGame
	setGameWinner(&game)
	setMemberAchievements(&game)
	jsonResponse(w, Response{true, http.StatusOK, "ok", game})
}

// addGameTeam() logic:
//      - a game can only have two teams
//      - once a new game is created it has 0 teams and its status is `newGame`
//      - first team to join must have between 3 and 5 members, otherwise error
//      - after first team joins, status is changed to `pendingGame`
//      - note: its possible that first team add/removes/updates members while status is `pendingGame`
//              if the number of team members becomes < 3 or > 5,
//              then team will be removed from game and game status reset to `newGame`
//      - second team to join must have same number of members as first team, otherwise error
//      - the two teams cannot be the same
//      - the two teams cannot have shared members
//      - once second team joins game status is changed to `startedGame`
//      - empty stats are created for each team member of both teams

func addGameTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	game := Game{}
	team, prevTeam := Team{}, Team{}

	if err := db.First(&game, vars["id0"]).Error; err != nil { // game not found
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id0"], errorMessage), nil})
		return
	}
	if (game.Status != newGame) && (game.Status != pendingGame) {
		s := fmt.Sprintf("cannot join game, status must be %d or %d", newGame, pendingGame)
		jsonResponse(w, Response{false, http.StatusForbidden, s, nil})
		return
	}

	if db.Model(&game).Association("Teams").Count() != 0 { // load 1st team
		db.Model(&game).Association("Teams").Find(&prevTeam)
		db.Model(&game).Association("Members").Clear()
		db.Model(&prevTeam).Association("Members").Find(&prevTeam.Members)
		db.Model(&game).Association("Members").Append(&prevTeam.Members) // reload game members
		db.Save(&game)
	}

	if len(prevTeam.Members) < gameMinNumMembers || len(prevTeam.Members) > gameMaxNumMembers {
		db.Model(&game).Association("Members").Clear()
		db.Model(&game).Association("Teams").Clear()
		game.Status = newGame
		db.Save(&game)
	}

	if err := db.First(&team, vars["id1"]).Error; err != nil {
		errorCode, errorMessage := translateError(http.StatusInternalServerError, err.Error())
		jsonResponse(w, Response{false, errorCode, fmt.Sprintf("%s: %s", vars["id1"], errorMessage), nil})
		return
	}

	db.Model(&team).Association("Members").Find(&team.Members)
	if len(team.Members) < gameMinNumMembers || len(team.Members) > gameMaxNumMembers {
		s := fmt.Sprintf("team must contain between %d and %d members", gameMinNumMembers, gameMaxNumMembers)
		jsonResponse(w, Response{false, http.StatusForbidden, s, nil})
		return
	}

	if db.Model(&game).Association("Teams").Count() == 0 { // this team is first team
		db.Model(&game).Association("Teams").Append(&team)
		db.Model(&game).Association("Members").Append(&team.Members)
		game.Status = pendingGame
		db.Save(&game)
		getGameTeams(w, r)
		return
	}

	if prevTeam.ID == team.ID {
		jsonResponse(w, Response{false, http.StatusForbidden, "teams cannot be the same", nil})
		return
	}
	if len(prevTeam.Members) != len(team.Members) {
		jsonResponse(w, Response{false, http.StatusForbidden, "teams must have same number of players", nil})
		return
	}
	if isSharedMembers(prevTeam.Members, team.Members) {
		jsonResponse(w, Response{false, http.StatusForbidden, "teams cannot have shared members", nil})
		return
	}
	db.Model(&game).Association("Teams").Append(&team) // add 2nd team to game
	db.Model(&game).Association("Members").Append(&team.Members)
	game.Status = startedGame
	createEmptyStats(&game, &team, &team.Members)
	createEmptyStats(&game, &prevTeam, &prevTeam.Members)
	db.Save(&game)
	getGameTeams(w, r)
}

func isSharedMembers(teamA, teamB []Member) bool {
	for i := range teamA {
		for j := range teamB {
			if teamA[i].ID == teamB[j].ID {
				return true
			}
		}
	}
	return false
}

func createEmptyStats(game *Game, team *Team, members *[]Member) {
	for _, member := range *members {
		db.Create(&Stat{GameID: game.ID, TeamID: team.ID, MemberID: member.ID})
	}
}

func isAwardedAlready(array []Achievement, ach Achievement) bool {
	for i := range array {
		if array[i].ID == ach.ID {
			return true
		}
	}
	return false
}
