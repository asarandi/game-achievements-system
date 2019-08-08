package main

import (
	"errors"
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

func setGameWinners(game *Game) {
	var teams []Team
	var stats []Stat
	var calc = [2]Stat{}
	var idx int
	if db.Model(game).Association("Teams").Find(&teams).Count() != 2 {
		return
	}
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
	if idx == -1 { // game ended in a tie
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

func isMemberAchievement(array []Achievement, ach Achievement) bool {
	for i := range array {
		if array[i].ID == ach.ID {
			return true
		}
	}
	return false
}

func setMemberAchievements(game *Game) {
	var members []Member
	var stats []Stat
	db.Order("id, asc").Model(game).Association("Members").Find(&members)
	/* join statement in case a member was deleted mid-game but the stat still exists */
	db.Order("member_id, asc").
		Joins("JOIN members ON members.id = stats.member_id AND stats.game_id = ?", game.ID).
		Find(&stats)
	if len(members) != len(stats) {		/* should not happen */
		panic("achievements: len(members) != len(stats)")
	}
	for i := range members {			/* preload existing achievements for each member */
		db.Model(&members[i]).Association("Achievements").Find(&members[i].Achievements)
	}
	for i := range asf {
		if db.First(&asf[i].achievement, &Achievement{Slug: asf[i].slug}).Error != nil {
			asf[i].achievement.ID = 0	/* could not find asf slug in database */
		}
	}
	for i := range members {
		for j := range asf {
			if asf[j].achievement.ID == 0 {
				continue
			}
			if isMemberAchievement(members[i].Achievements, asf[j].achievement) {
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
	var game = Game{}
	if err := getRecordByID(&game, mux.Vars(r)["id0"]); err != nil {
		responseJson(w, err, nil, 0)
		return
	}
	if game.Status != startedGame {
		responseJson(w, errors.New("cannot change status of this game"), nil,0)
		return
	}
	game.Status = finishedGame
	setGameWinners(&game)
	setMemberAchievements(&game)
	db.Save(&game)
	responseJson(w, nil, nil,0)
}

/*
	addGameTeam() logic:
     - a game can only have two teams
     - once a new game is created it has 0 teams and its status is `newGame`
     - first team to join must have between 3 and 5 members, otherwise error
     - after first team joins, status is changed to `pendingGame`
     - note: its possible that first team add/removes members while status is `pendingGame`
             if the number of team members becomes < 3 or > 5,
             then team will be removed from game and game status reset to `newGame`
     - second team to join must have same number of members as first team, otherwise error
     - the two teams cannot be the same
     - the two teams cannot have shared members
     - once second team joins game status is changed to `startedGame`
     - empty stats are created for each team member of both teams
*/

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

func createEmptyStats(game *Game) {
	for i := range game.Teams {
		for j := range game.Teams[i].Members {
			db.Create(&Stat{
				GameID: game.ID,
				TeamID: game.Teams[i].ID,
				MemberID: game.Teams[i].Members[j].ID})
		}
	}
}

func addGameTeam(w http.ResponseWriter, r *http.Request) {
	var gameId, teamId =  mux.Vars(r)["id0"], mux.Vars(r)["id1"]
	var game = Game{}
	var team, prevTeam = Team{}, Team{}
	var teamCount = 0

	if err := getRecordByID(&game, gameId); err != nil {
		responseJson(w, err, nil, 0)
		return
	}
	if (game.Status != newGame) && (game.Status != pendingGame) {
		s := fmt.Sprintf("cannot add team to game, status must be %d or %d", newGame, pendingGame)
		responseJson(w, errors.New(s), nil, 0)
		return
	}
	teamCount = db.Model(&game).Association("Teams").Count()
	if teamCount > 0 { // load 1st team
		db.Model(&game).Association("Teams").Find(&prevTeam)
		db.Model(&prevTeam).Association("Members").Find(&prevTeam.Members)
		if (len(prevTeam.Members) < gameMinNumMembers) || (len(prevTeam.Members) > gameMaxNumMembers) {
			db.Model(&game).Association("Teams").Clear()
			game.Status = newGame
			db.Save(&game)
			teamCount = 0
		}
	}
	if err := getRecordByID(&team, teamId); err != nil {
		responseJson(w, err, nil, 0)
		return
	}
	db.Model(&team).Association("Members").Find(&team.Members)
	if len(team.Members) < gameMinNumMembers || len(team.Members) > gameMaxNumMembers {
		s := fmt.Sprintf("team must contain between %d and %d members", gameMinNumMembers, gameMaxNumMembers)
		responseJson(w, errors.New(s), nil, 0)
		return
	}
	if teamCount == 0 {
		db.Model(&game).Association("Teams").Append(&team)			// add first team
		game.Status = pendingGame
		db.Save(&game)
		responseJson(w, nil, nil, 0)
		return
	}
	if prevTeam.ID == team.ID {
		responseJson(w, errors.New("teams cannot be the same"),nil, 0)
		return
	}
	if len(prevTeam.Members) != len(team.Members) {
		responseJson(w, errors.New("teams must have same number of players"),nil, 0)
		return
	}
	if isSharedMembers(prevTeam.Members, team.Members) {
		responseJson(w, errors.New("teams cannot have shared members"),nil, 0)
		return
	}

	db.Model(&game).Association("Teams").Append(&team)				// add second team
	db.Model(&game).Association("Members").Append(&team.Members)
	db.Model(&game).Association("Members").Append(&prevTeam.Members)
	createEmptyStats(&game)
	game.Status = startedGame
	db.Save(&game)
	getGameTeams(w, r)
}
