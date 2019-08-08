package main

import (
	"github.com/gorilla/mux"
)

var r *mux.Router

func setRoutes() {
	r = mux.NewRouter()

	r.HandleFunc("/achievements", getAllAchievements).Methods("GET")
	r.HandleFunc("/achievements", createAchievement).Methods("POST")
	r.HandleFunc("/achievements/{id0:[0-9]+}", getAchievement).Methods("GET")
	r.HandleFunc("/achievements/{id0:[0-9]+}", updateAchievement).Methods("PUT")
	r.HandleFunc("/achievements/{id0:[0-9]+}", deleteAchievement).Methods("DELETE")
	r.HandleFunc("/achievements/{id0:[0-9]+}/members", getAchievementMembers).Methods("GET")

	r.HandleFunc("/members", getAllMembers).Methods("GET")
	r.HandleFunc("/members", createMember).Methods("POST")
	r.HandleFunc("/members/{id0:[0-9]+}", getMember).Methods("GET")
	r.HandleFunc("/members/{id0:[0-9]+}", updateMember).Methods("PUT")
	r.HandleFunc("/members/{id0:[0-9]+}", deleteMember).Methods("DELETE")
	r.HandleFunc("/members/{id0:[0-9]+}/achievements", getMemberAchievements).Methods("GET")
	r.HandleFunc("/members/{id0:[0-9]+}/teams", getMemberTeams).Methods("GET")
	r.HandleFunc("/members/{id0:[0-9]+}/teams/{id1:[0-9]+}", addMemberTeam).Methods("POST")
	r.HandleFunc("/members/{id0:[0-9]+}/teams/{id1:[0-9]+}", removeMemberTeam).Methods("DELETE")

	r.HandleFunc("/teams", getAllTeams).Methods("GET")
	r.HandleFunc("/teams", createTeam).Methods("POST")
	r.HandleFunc("/teams/{id0:[0-9]+}", getTeam).Methods("GET")
	r.HandleFunc("/teams/{id0:[0-9]+}", updateTeam).Methods("PUT")
	r.HandleFunc("/teams/{id0:[0-9]+}", deleteTeam).Methods("DELETE")
	r.HandleFunc("/teams/{id0:[0-9]+}/members", getTeamMembers).Methods("GET")
	r.HandleFunc("/teams/{id0:[0-9]+}/members/{id1:[0-9]+}", addTeamMember).Methods("POST")
	r.HandleFunc("/teams/{id0:[0-9]+}/members/{id1:[0-9]+}", removeTeamMember).Methods("DELETE")

	r.HandleFunc("/games", getAllGames).Methods("GET")
	r.HandleFunc("/games", createGame).Methods("POST")
	r.HandleFunc("/games/{id0:[0-9]+}", endGame).Methods("DELETE") 			/* close game */
	r.HandleFunc("/games/{id0:[0-9]+}/stats", getGameStats).Methods("GET")
	r.HandleFunc("/games/{id0:[0-9]+}/teams", getGameTeams).Methods("GET")
	r.HandleFunc("/games/{id0:[0-9]+}/teams/{id1:[0-9]+}", addGameTeam).Methods("POST")
	r.HandleFunc("/games/{id0:[0-9]+}/teams/{id1:[0-9]+}/stats", getGameTeamStats).Methods("GET")
	r.HandleFunc("/games/{id0:[0-9]+}/members", getGameMembers).Methods("GET")
	r.HandleFunc("/games/{id0:[0-9]+}/winners", getGameWinners).Methods("GET")
	r.HandleFunc("/games/{id0:[0-9]+}/members/{id1:[0-9]+}/stats",getGameMemberStats).Methods("GET")
	r.HandleFunc("/games/{id0:[0-9]+}/members/{id1:[0-9]+}/stats", updateGameMemberStats).Methods("PUT")
}
