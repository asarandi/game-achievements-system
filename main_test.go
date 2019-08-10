package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	initDatabase()
	setRoutes()
}

/* 	ignore these keys when comparing expected response to received response  */
var ignoreKeys = map[string]interface{}{
	"ID"			:nil,	// gorm.Model
	"CreatedAt"		:nil,	// gorm.Model
	"UpdatedAt"		:nil,	// gorm.Model
	"DeletedAt"		:nil,	// gorm.Model
	"game_id"		:nil,	// Stat{}
	"team_id"		:nil,	// Stat{}
	"member_id"		:nil,	// Stat{}
	"is_winner"		:nil,	// Stat{}
}

func isJsonObjectsEqual(expectedBytes []byte,  receivedBytes []byte) bool {
	expectedObject := map[string]interface{}{}
	receivedObject := map[string]interface{}{}
	if json.Unmarshal(expectedBytes, &expectedObject) != nil {	panic("json.Unmarshal() failed")	}
	if json.Unmarshal(receivedBytes, &receivedObject) != nil {	panic("json.Unmarshal() failed")	}
	for key, value := range receivedObject {
		if key == "result" {continue}                                                       /*compare in next for loop*/
		if expectedObject[key] != value {return false}                                      /*objects are different*/
	}
	if receivedObject["result"] == nil && expectedObject["result"] == nil { return true }	/*if result is nil, then objects match*/
	if _, ok := receivedObject["result"].([]interface{}); ok {
			expectedArray := expectedObject["result"].([]interface{})						/*handle json array of objects*/
			receivedArray := receivedObject["result"].([]interface{})
			for i := range receivedArray {
				for key, value := range receivedArray[i].(map[string]interface{}) {
					if _, ok := ignoreKeys[key]; ok {continue}
					if value != expectedArray[i].(map[string]interface{})[key] {
						return false
					}
				}
			}
			return true
	}
	expectedResult := expectedObject["result"].(map[string]interface{})                     /*handle json object*/
	receivedResult := receivedObject["result"].(map[string]interface{})
	for key, value := range receivedResult {
        if _, ok := ignoreKeys[key]; ok {continue}                                          /*ignore CreatedAt, UpdatedAt, etc*/
		if expectedResult[key] != value {return false}                                      /*objects are different*/
	}
	return true																				/*objects are the same*/
}

func requestAndCompareArray(endpoints []string, method string, requestData, expectedResponseData []interface{}, response Response) bool {
	for i := range requestData {
		b, err := json.Marshal(requestData[i])
		if err != nil {
			panic("json.Marshal() failed")
		}
		req, err := http.NewRequest(method, endpoints[i], bytes.NewBuffer(b))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		response.Result = expectedResponseData[i]
		expectedBytes, err := json.Marshal(response)
		if err != nil {	panic("json.Marshal() failed") }
		if verboseTest {
			var buf1, buf2 bytes.Buffer
			json.Indent(&buf1, expectedBytes, "", "\t")
			fmt.Println("EXPECTED\n", buf1.String())
			json.Indent(&buf2, res.Body.Bytes(), "", "\t")
			fmt.Println("RECEIVED\n", buf2.String())
		}
		if !isJsonObjectsEqual(expectedBytes, res.Body.Bytes()) {return false}
	}
	return true
}

func requestAndCompareResponse(endpoint, method string, data []interface{}, response Response) bool {
	for i := range data {
		b, err := json.Marshal(data[i])
		if err != nil { panic("json.Marshal() failed") }
		req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(b))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		response.Result = data[i]
		expectedBytes, err := json.Marshal(response)
		if err != nil {panic("json.Marshal() failed")}
		if !isJsonObjectsEqual(expectedBytes, res.Body.Bytes()) {return false}
	}
	return true
}

func TestInsertMembers(t *testing.T) {
	endpoints := []string{}
	for i := 0; i < 10; i++ {endpoints = append(endpoints, "/members") }
	method := http.MethodPost
	response := Response{Success: true, Code: 201, Message: "ok", Result: nil}
	requestData := []interface{}{
		Member{Name: "Mario", Img: "https://upload.wikimedia.org/wikipedia/en/a/a9/MarioNSMBUDeluxe.png"},
		Member{Name: "Luigi", Img: "https://upload.wikimedia.org/wikipedia/en/f/f1/LuigiNSMBW.png"},
		Member{Name: "Princess Peach", Img: "https://upload.wikimedia.org/wikipedia/en/d/d5/Peach_%28Super_Mario_3D_World%29.png"},
		Member{Name: "Toad", Img: "https://upload.wikimedia.org/wikipedia/en/d/d1/Toad_3D_Land.png"},
		Member{Name: "Yoshi", Img: "https://upload.wikimedia.org/wikipedia/en/d/d9/YoshiMarioParty10.png"},
		Member{Name: "Bowser", Img: "https://upload.wikimedia.org/wikipedia/en/1/11/BowserNSMBUD.png"},
		Member{Name: "Bowser Jr.", Img: "https://upload.wikimedia.org/wikipedia/en/d/d2/Bowser_Jr.png"},
		Member{Name: "Princess Daisy", Img: "https://upload.wikimedia.org/wikipedia/en/b/bd/Daisy_%28Super_Mario_Party%29.png"},
		Member{Name: "Wario", Img: "https://upload.wikimedia.org/wikipedia/en/8/81/Wario.png"},
		Member{Name: "Waluigi", Img: "https://upload.wikimedia.org/wikipedia/en/4/46/Waluigi.png"},
	}
	resposeData := requestData
	if !requestAndCompareArray(endpoints, method, requestData, resposeData, response) {
		t.Fatal("TestInsertMembers() failed")
	}
}

func TestInsertTeams(t *testing.T) {
	endpoints := []string{}
	for i := 0; i < 4; i++ {endpoints = append(endpoints, "/teams") }
	method := http.MethodPost
	response := Response{Success: true, Code: 201, Message: "ok", Result: nil}
	requestData := []interface{}{
		Team{Name: "The Sluggers Team", Img: "https://upload.wikimedia.org/wikipedia/en/0/0c/MarioSuperSluggers.png",},
		Team{Name: "Strikers Group Inc", Img: "https://upload.wikimedia.org/wikipedia/en/5/50/Mario_Strikers_Charged.jpg",},
		Team{Name: "Superstar Alliance Corp", Img: "https://upload.wikimedia.org/wikipedia/en/2/2f/Mario_Superstar_Baseball.jpg",},
		Team{Name: "Advanced Tourists Ltd", Img: "https://upload.wikimedia.org/wikipedia/en/0/0d/Advance_Tour_Cover.jpg",},
	}
	responseData := requestData
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestInsertTeams() failed")
	}
}

func TestInsertAchievements(t *testing.T) {
	endpoints := []string{}
	for i := 0; i < 4; i++ {endpoints = append(endpoints, "/achievements") }
	method := http.MethodPost
	response := Response{Success: true, Code: 201, Message: "ok", Result: nil}
	requestData := []interface{}{
        Achievement{Slug: "sharpshooter", Name: "“Sharpshooter” Award", Desc: "Land at least 75% of all your attacks, assuming you attacked at least once.", Img: "http://images.clipartpanda.com/sharpshooter-clipart-gg58876635.jpg",},
        Achievement{Slug: "bruiser", Name: "“Bruiser” Award", Desc: "Do more than 500 points of damage during one game.", Img: "https://i1.wp.com/azwildlife.org/wp-content/uploads/2018/07/trophy.png", },
        Achievement{Slug: "veteran", Name: "“Veteran” Award", Desc: "Play more than 1000 games.", Img: "http://icons.iconarchive.com/icons/google/noto-emoji-activities/1024/52725-trophy-icon.png",},
        Achievement{Slug: "bigwinner", Name: "“Big Winner” Award", Desc: "Have over 200 wins.", Img: "https://png.pngtree.com/element_origin_min_pic/17/09/21/91633641bc0263293bce7cab1593e41c.jpg",},
	}
	responseData := requestData
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestInsertAchievements() failed")
	}
}

// create two teams with 4 members each
// create two more teams with 3 members each
func TestMembersJoinTeams(t *testing.T) {
	endpoints := []string{
		"/members/1/teams/1",	//mario joins sluggers
		"/members/2/teams/1",	//luigi joins sluggers
		"/members/3/teams/1",	//princess peach joins sluggers
		"/members/4/teams/1",	//toad joins sluggers
		"/members/5/teams/2",	//yoshi joins strikers
		"/members/6/teams/2",	//bowser joins strikers
		"/members/7/teams/2",	//bowser jr joins strikers
		"/members/8/teams/2",	//princess daisy joins strikers
		"/members/1/teams/3",	//mario joins superstar
		"/members/2/teams/3",	//luigi joins superstar
		"/members/3/teams/3",	//princess peach joins superstar
		"/members/5/teams/4",	//yoshi joins tourists
		"/members/6/teams/4",	//bowser joins tourists
		"/members/7/teams/4",	//bowser jr joins tourists
	}
	method := http.MethodPost
	requestData := make([]interface{}, 14)
	responseData := make([]interface{}, 14)
	response := Response{Success: true, Code: 200, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestMembersJoinTeams() failed")
	}
}

func TestCreateGames(t *testing.T) {
	endpoints := []string{"/games", "/games",}
	method := http.MethodPost
	requestData := make([]interface{}, 2)
	responseData := []interface{}{
		Game{Status:newGame},
		Game{Status:newGame},
	}
	response := Response{Success: true, Code: 201, Message: "ok", Result:nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestCreateGames() failed")
	}
}

func TestTeamsJoinGames(t *testing.T) {
	endpoints := []string{
		"/games/1/teams/1",
		"/games/1/teams/2",
		"/games/2/teams/3",
		"/games/2/teams/4",
	}
	method := http.MethodPost
	requestData := make([]interface{}, 4)
	responseData := requestData
	response := Response{Success: true, Code: 202, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestTeamsJoinGames() failed")
	}
}

/*	at this point we have 2 games with status `gameStarted`

	game_id 1:
		team_id 1:
			members id: 1,2,3,4
		team_id 2:
			members id 5,6,7,8

	game_id 2:
		team_id 3:
			members id: 1,2,3
		team_id 4:
			members id: 5,6,7				*/

func TestUpdateGameMemberStats(t *testing.T) {
	endpoints := []string{
		"/games/1/members/1/stats",	//mario
		"/games/1/members/2/stats",	//luigi
		"/games/1/members/3/stats",	//princess peach
		"/games/1/members/4/stats",	//toad
		"/games/1/members/5/stats",	//yoshi
		"/games/1/members/6/stats",	//bowser
		"/games/1/members/7/stats",	//bowser jr
		"/games/1/members/8/stats",	//princess daisy
		"/games/2/members/1/stats",	//mario
		"/games/2/members/2/stats",	//luigi
		"/games/2/members/3/stats",	//princess peach
		"/games/2/members/5/stats",	//yoshi
		"/games/2/members/6/stats",	//bowser
		"/games/2/members/7/stats",	//bowser jr
	}
	method := http.MethodPut		/* PUT */
	requestData := []interface{}{
        Stat{NumHits:1,     NumAttacks:12,  AmountDamage:13,    NumKills:14,  InstantKills:15,  NumAssists:16,  NumSpells:17,  SpellsDamage:18},   //game 1: team 1: member 1: mario
		Stat{NumHits:99999, NumAttacks:22,  AmountDamage:23,    NumKills:24,  InstantKills:25,  NumAssists:26,  NumSpells:27,  SpellsDamage:28},   //game 1: team 1: member 2: luigi				<<-- sharpshooter
		Stat{NumHits:3,     NumAttacks:32,  AmountDamage:33,    NumKills:34,  InstantKills:35,  NumAssists:36,  NumSpells:37,  SpellsDamage:38},   //game 1: team 1: member 3: princess peach
		Stat{NumHits:4,     NumAttacks:42,  AmountDamage:43,    NumKills:44,  InstantKills:45,  NumAssists:46,  NumSpells:47,  SpellsDamage:48},   //game 1: team 1: member 4: toad
		Stat{NumHits:5,     NumAttacks:52,  AmountDamage:53,    NumKills:54,  InstantKills:55,  NumAssists:56,  NumSpells:57,  SpellsDamage:58},   //game 1: team 2: member 5: yoshi
		Stat{NumHits:6,     NumAttacks:62,  AmountDamage:63,    NumKills:64,  InstantKills:65,  NumAssists:66,  NumSpells:67,  SpellsDamage:68},   //game 1: team 2: member 6: bowser
		Stat{NumHits:7,     NumAttacks:72,  AmountDamage:73,    NumKills:74,  InstantKills:75,  NumAssists:76,  NumSpells:77,  SpellsDamage:78},   //game 1: team 2: member 7: bowser jr
		Stat{NumHits:8,     NumAttacks:82,  AmountDamage:83,    NumKills:84,  InstantKills:85,  NumAssists:86,  NumSpells:87,  SpellsDamage:88},   //game 1: team 2: member 8: princess daisy
		Stat{NumHits:9,     NumAttacks:92,  AmountDamage:93,    NumKills:94,  InstantKills:95,  NumAssists:96,  NumSpells:97,  SpellsDamage:98},   //game 2: team 3: member 1: mario
		Stat{NumHits:10,    NumAttacks:102, AmountDamage:103,   NumKills:104, InstantKills:105, NumAssists:106, NumSpells:107, SpellsDamage:108},  //game 2: team 3: member 2: luigi
		Stat{NumHits:11,    NumAttacks:112, AmountDamage:99999, NumKills:114, InstantKills:115, NumAssists:116, NumSpells:117, SpellsDamage:118},  //game 2: team 3: member 3: princess peach	    <<-- bruiser
		Stat{NumHits:12,    NumAttacks:122, AmountDamage:123,   NumKills:124, InstantKills:125, NumAssists:126, NumSpells:127, SpellsDamage:128},  //game 2: team 4: member 5: yoshi
		Stat{NumHits:13,    NumAttacks:132, AmountDamage:133,   NumKills:134, InstantKills:135, NumAssists:136, NumSpells:137, SpellsDamage:138},  //game 2: team 4: member 6: bowser
		Stat{NumHits:14,    NumAttacks:142, AmountDamage:143,   NumKills:144, InstantKills:145, NumAssists:146, NumSpells:147, SpellsDamage:148},  //game 2: team 1: member 7: bowser jr
	}
	responseData := requestData
	response := Response{Success: true, Code: 202, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestUpdateGameMemberStats() failed")
	}
}

/*	lets get a list of all achievements holders,
	the list should be empty,
	after we close the game, luigi should have the "sharpshooter" achievement
	and princess peach should have the "bruiser" achievement	*/

func TestGetAchievementMembersBefore(t *testing.T) {
	endpoints := []string{
		"/achievements/1/members",
		"/achievements/2/members",
		"/achievements/3/members",
		"/achievements/4/members",
	}
	method := http.MethodGet		/* GET */
	requestData := make([]interface{}, 4)
	responseData := []interface{}{
		[]interface{}{},
		[]interface{}{},
		[]interface{}{},
		[]interface{}{},
	}
	response := Response{Success: true, Code: 200, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestGetAchievementMembersBefore() failed")
	}
}

func TestGetGameWinnersBefore(t *testing.T) {
	endpoints := []string{
		"/games/1/winners",
		"/games/2/winners",
	}
	method := http.MethodGet		/* GET */
	requestData := make([]interface{}, 2)
	responseData := []interface{}{
		[]interface{}{},
		[]interface{}{},
	}
	response := Response{Success: true, Code: 200, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestGetGameWinnersBefore() failed")
	}
}

func TestCloseGame(t *testing.T) {
	endpoints := []string{
		"/games/1",
		"/games/2",
	}
	method := http.MethodDelete		/* DELETE */
	requestData := make([]interface{}, 2)
	responseData := requestData
	response := Response{Success: true, Code: 200, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestCloseGame() failed")
	}
}

func TestGetAchievementMembersAfter(t *testing.T) {
	endpoints := []string{
		"/achievements/1/members",
		"/achievements/2/members",
		"/achievements/3/members",
		"/achievements/4/members",
	}
	method := http.MethodGet		/* GET */
	requestData := make([]interface{}, 4)
	responseData := []interface{}{
		[]interface{}{Member{Name: "Luigi", Img: "https://upload.wikimedia.org/wikipedia/en/f/f1/LuigiNSMBW.png"},},
		[]interface{}{Member{Name: "Princess Peach", Img: "https://upload.wikimedia.org/wikipedia/en/d/d5/Peach_%28Super_Mario_3D_World%29.png"},},
		[]interface{}{},
		[]interface{}{},
	}
	response := Response{Success: true, Code: 200, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestGetAchievementMembersAfter() failed")
	}
}

func TestGetGameWinnersAfter(t *testing.T) {
	endpoints := []string{
		"/games/1/winners",
		"/games/2/winners",
	}
	method := http.MethodGet		/* GET */
	requestData := make([]interface{}, 2)
	responseData := []interface{}{
		/*for game #1, team #2 had better combined stats than team #1; team #2 wins*/
		[]interface{}{
			Member{Name: "Yoshi", Img: "https://upload.wikimedia.org/wikipedia/en/d/d9/YoshiMarioParty10.png"},
			Member{Name: "Bowser", Img: "https://upload.wikimedia.org/wikipedia/en/1/11/BowserNSMBUD.png"},
			Member{Name: "Bowser Jr.", Img: "https://upload.wikimedia.org/wikipedia/en/d/d2/Bowser_Jr.png"},
			Member{Name: "Princess Daisy", Img: "https://upload.wikimedia.org/wikipedia/en/b/bd/Daisy_%28Super_Mario_Party%29.png"},
		},
		/*for game #2, team #4 had better combined stats than team #3; team #4 wins*/
		[]interface{}{
			Member{Name: "Yoshi", Img: "https://upload.wikimedia.org/wikipedia/en/d/d9/YoshiMarioParty10.png"},
			Member{Name: "Bowser", Img: "https://upload.wikimedia.org/wikipedia/en/1/11/BowserNSMBUD.png"},
			Member{Name: "Bowser Jr.", Img: "https://upload.wikimedia.org/wikipedia/en/d/d2/Bowser_Jr.png"},
		},
	}
	response := Response{Success: true, Code: 200, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestGetGameWinnersAfter() failed")
	}
}
