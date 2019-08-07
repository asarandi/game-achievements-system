package main

import (
	"github.com/gorilla/mux"
)

var router *mux.Router

func setRoutes() {
	router = mux.NewRouter()

	router.HandleFunc("/achievements", getAllAchievements).Methods("GET")                // get all records
	router.HandleFunc("/achievements", createAchievement).Methods("POST")                // create new record
	router.HandleFunc("/achievements/{id0:[0-9]+}", getAchievement).Methods("GET")       // get record by id
	router.HandleFunc("/achievements/{id0:[0-9]+}", updateAchievement).Methods("PUT")    // update record
	router.HandleFunc("/achievements/{id0:[0-9]+}", deleteAchievement).Methods("DELETE") // delete record
	router.HandleFunc("/achievements/{id0:[0-9]+}/members", getAchievementMembers).Methods("GET")

	router.HandleFunc("/members", getAllMembers).Methods("GET")
	router.HandleFunc("/members", createMember).Methods("POST")
	router.HandleFunc("/members/{id0:[0-9]+}", getMember).Methods("GET")
	router.HandleFunc("/members/{id0:[0-9]+}", updateMember).Methods("PUT")
	router.HandleFunc("/members/{id0:[0-9]+}", deleteMember).Methods("DELETE")
	router.HandleFunc("/members/{id0:[0-9]+}/achievements", getMemberAchievements).Methods("GET")
	router.HandleFunc("/members/{id0:[0-9]+}/teams", getMemberTeams).Methods("GET")                   // list all teams for a given member
	router.HandleFunc("/members/{id0:[0-9]+}/teams/{id1:[0-9]+}", addMemberTeam).Methods("POST")      // a member (id0) joins a team (id1)
	router.HandleFunc("/members/{id0:[0-9]+}/teams/{id1:[0-9]+}", removeMemberTeam).Methods("DELETE") // a member (id0) leaves a team (id1)

	router.HandleFunc("/teams", getAllTeams).Methods("GET")
	router.HandleFunc("/teams", createTeam).Methods("POST")
	router.HandleFunc("/teams/{id0:[0-9]+}", getTeam).Methods("GET")
	router.HandleFunc("/teams/{id0:[0-9]+}", updateTeam).Methods("PUT")
	router.HandleFunc("/teams/{id0:[0-9]+}", deleteTeam).Methods("DELETE")
	router.HandleFunc("/teams/{id0:[0-9]+}/members", getTeamMembers).Methods("GET")                   // list all members for a given team
	router.HandleFunc("/teams/{id0:[0-9]+}/members/{id1:[0-9]+}", addTeamMember).Methods("POST")      // a team (id0) adds a member (id1)
	router.HandleFunc("/teams/{id0:[0-9]+}/members/{id1:[0-9]+}", removeTeamMember).Methods("DELETE") // a team (id0) removes a member (id1)

	router.HandleFunc("/games", getAllGames).Methods("GET")
	router.HandleFunc("/games", createGame).Methods("POST")
	router.HandleFunc("/games/{id0:[0-9]+}", endGame).Methods("DELETE") // set status as `finishedGame`
	router.HandleFunc("/games/{id0:[0-9]+}/stats", getGameStats).Methods("GET")
	router.HandleFunc("/games/{id0:[0-9]+}/teams", getGameTeams).Methods("GET")
	router.HandleFunc("/games/{id0:[0-9]+}/teams/{id1:[0-9]+}", addGameTeam).Methods("POST") // a team (id1) joins a game (id0)
	router.HandleFunc("/games/{id0:[0-9]+}/members", getGameMembers).Methods("GET")
	router.HandleFunc("/games/{id0:[0-9]+}/winners", getGameWinners).Methods("GET")
	router.HandleFunc("/games/{id0:[0-9]+}/members/{id1:[0-9]+}/stats", getGameMemberStats).Methods("GET")
	router.HandleFunc("/games/{id0:[0-9]+}/members/{id1:[0-9]+}/stats", updateGameMemberStats).Methods("PUT") // update game member stats
}
