package main

import (
	_ "bufio"
	"bytes"
	_ "bytes"
	_ "encoding/csv"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "io"
	_ "io/ioutil"
	_ "log"
	"net/http"
	_ "net/http"
	"net/http/httptest"
	_ "net/http/httptest"
	_ "os"
	"testing"
)

func init() {
	initDatabase()
	setRoutes()
}

/* 	ignore these keys when comparing client request to server response  */
var ignoreKeys = map[string]interface{}{
	"ID"			:nil,	// gorm.Model
	"CreatedAt"		:nil,	// gorm.Model
	"UpdatedAt"		:nil,	// gorm.Model
	"DeletedAt"		:nil,	// gorm.Model
	"game_id"		:nil,	// Stat{}
	"team_id"		:nil,	// Stat{}
	"member_id"		:nil,	// Stat{}
	"is_winner"		:nil,	// Stat{}
//	"status"		:nil,	// Game{}
}

func isJsonObjectsEqual(requestBytes []byte,  responseBytes []byte) bool {
	requestObject := map[string]interface{}{}
	responseObject := map[string]interface{}{}
	if json.Unmarshal(requestBytes, &requestObject) != nil {	panic("json.Unmarshal() failed")	}
	if json.Unmarshal(responseBytes, &responseObject) != nil {	panic("json.Unmarshal() failed")	}
	for key, value := range requestObject {
		if key == "result" {continue}						/*compare in next for loop*/
		if responseObject[key] != value {return false}		/*objects are different*/
	}

	if responseObject["result"] == nil && requestObject["result"] == nil { return true }

	if v, ok := responseObject["result"].([]interface{}); ok {
		if len(v) == 0 {	//empty json array
			responseObject["result"] = nil
		} else {
			requestArray := requestObject["result"].([]interface{})
			responseArray := responseObject["result"].([]interface{})
			for i := range responseArray {
				for key, value := range responseArray[i].(map[string]interface{}) {
					if value != requestArray[i].(map[string]interface{})[key] {
						return false
					}
				}
			}
			return true
		}
	}

	if responseObject["result"] == nil && requestObject["result"] == nil { return true }
	
	if _, ok := responseObject["result"].(interface{}); ok {
		if responseObject["result"].(interface{}) == nil && requestObject["result"].(interface{}) == nil {
			return true
		}
	}

	responseResult := responseObject["result"].(map[string]interface{})
	requestResult := requestObject["result"].(map[string]interface{})

	for key, value := range requestResult {
        if _, ok := ignoreKeys[key]; ok {continue}			/*ignore CreatedAt, UpdatedAt, etc*/
		if responseResult[key] != value {return false}		/*objects are different*/
	}
	return true												/*objects are the same*/
}

func requestAndCompareArray(endpoints []string, method string, requestData, expectedResponseData []interface{}, response Response) bool {
	fmt.Println("response", response)
	for i := range requestData {
		b, err := json.Marshal(requestData[i])
		if err != nil { panic("json.Marshal() failed") }
		req, err := http.NewRequest(method, endpoints[i], bytes.NewBuffer(b))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		response.Result = expectedResponseData[i]
		expectedBytes, err := json.Marshal(response)
		if err != nil {panic("json.Marshal() failed")}
		fmt.Println("expected", string(expectedBytes))
		fmt.Println("  actual", res.Body.String())
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
	endpoint := "/members"
	method := http.MethodPost
	response := Response{Success: true, Code: 201, Message: "ok", Result: nil}
	data := []interface{}{
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
	if !requestAndCompareResponse(endpoint, method, data, response) {
		t.Fatal("TestInsertMembers() failed")
	}
}

func TestInsertTeams(t *testing.T) {
	endpoint := "/teams"
	method := http.MethodPost
	response := Response{Success: true, Code: 201, Message: "ok", Result: nil}
	data := []interface{}{
		Team{Name: "The Sluggers Team", Img: "https://upload.wikimedia.org/wikipedia/en/0/0c/MarioSuperSluggers.png",},
		Team{Name: "Strikers Group Inc", Img: "https://upload.wikimedia.org/wikipedia/en/5/50/Mario_Strikers_Charged.jpg",},
		Team{Name: "Superstar Alliance Corp", Img: "https://upload.wikimedia.org/wikipedia/en/2/2f/Mario_Superstar_Baseball.jpg",},
		Team{Name: "Advanced Tourists Ltd", Img: "https://upload.wikimedia.org/wikipedia/en/0/0d/Advance_Tour_Cover.jpg",},
	}
	if !requestAndCompareResponse(endpoint, method, data, response) {
		t.Fatal("TestInsertTeams() failed")
	}
}

func TestInsertAchievements(t *testing.T) {
	endpoint := "/achievements"
	method := http.MethodPost
	response := Response{Success: true, Code: 201, Message: "ok", Result: nil}
	data := []interface{}{
        Achievement{Slug: "sharpshooter", Name: "“Sharpshooter” Award", Desc: "Land at least 75% of all your attacks, assuming you attacked at least once.", Img: "http://images.clipartpanda.com/sharpshooter-clipart-gg58876635.jpg",},
        Achievement{Slug: "bruiser", Name: "“Bruiser” Award", Desc: "Do more than 500 points of damage during one game.", Img: "https://i1.wp.com/azwildlife.org/wp-content/uploads/2018/07/trophy.png", },
        Achievement{Slug: "veteran", Name: "“Veteran” Award", Desc: "Play more than 1000 games.", Img: "http://icons.iconarchive.com/icons/google/noto-emoji-activities/1024/52725-trophy-icon.png",},
        Achievement{Slug: "bigwinner", Name: "“Big Winner” Award", Desc: "Have over 200 wins.", Img: "https://png.pngtree.com/element_origin_min_pic/17/09/21/91633641bc0263293bce7cab1593e41c.jpg",},
	}
	if !requestAndCompareResponse(endpoint, method, data, response) {
		t.Fatal("TestInsertAchievements() failed")
	}
}
//
//func TestErrorInsertDuplicate(t *testing.T) {
//	endpoint := "/members"
//	method := http.MethodPost
//	response := Response{Success: false, Code: 406, Message: "record already exists", Result: nil}
//	data := []interface{}{
//		Member{Name: "Mario", Img: "https://upload.wikimedia.org/wikipedia/en/a/a9/MarioNSMBUDeluxe.png"},
//	}
//	if !requestAndCompareResponse(endpoint, method, data, response) {
//		t.Fatal("TestErrorInsertDuplicate failed")
//	}
//}
//


// create two teams with 4 members each
// create another two teams with 3 members each
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

/*
	at this point we have 2 games with status `gameStarted`

	game_id 1:
		team_id 1:
			members id: 1,2,3,4
		team_id 2:
			members id 5,6,7,8

	game_id 2:
		team_id 3:
			members id: 1,2,3
		team_id 4:
			members id: 5,6,7
*/


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

/*
		now we have "rigged" the game,
		before we close the game,
		lets get a list of all achievements holders,
		the list should be empty
		after we close the game, luigi should have the "sharpshooter" achievement
		and princess peach should have the "bruiser" achievement
 */
func TestGetAchievementMembersBefore(t *testing.T) {
	endpoints := []string{
		"/achievements/1/members",
		"/achievements/2/members",
		"/achievements/3/members",
		"/achievements/4/members",
	}
	method := http.MethodGet		/* GET */
	requestData := make([]interface{}, 4)
	responseData := requestData
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
	responseData := requestData
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
		[]interface{}{
			Member{},
		},
		nil,
		nil,
		nil,
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
	responseData := requestData
	response := Response{Success: true, Code: 200, Message: "ok", Result: nil}
	if !requestAndCompareArray(endpoints, method, requestData, responseData, response) {
		t.Fatal("TestGetGameWinnerAfter() failed")
	}
}







//
//func TestInsertMembers(t *testing.T){
//	endpoint := "/members"
//	method := http.MethodPost
//	data := []Member{
//		{Name: "Mario",          Img: "https://upload.wikimedia.org/wikipedia/en/a/a9/MarioNSMBUDeluxe.png"},
//        {Name: "Luigi",          Img: "https://upload.wikimedia.org/wikipedia/en/f/f1/LuigiNSMBW.png"},
//        {Name: "Princess Peach", Img: "https://upload.wikimedia.org/wikipedia/en/d/d5/Peach_%28Super_Mario_3D_World%29.png"},
//        {Name: "Toad",           Img: "https://upload.wikimedia.org/wikipedia/en/d/d1/Toad_3D_Land.png"},
//        {Name: "Yoshi",          Img: "https://upload.wikimedia.org/wikipedia/en/d/d9/YoshiMarioParty10.png"},
//        {Name: "Bowser",         Img: "https://upload.wikimedia.org/wikipedia/en/1/11/BowserNSMBUD.png"},
//        {Name: "Bowser Jr.",     Img: "https://upload.wikimedia.org/wikipedia/en/d/d2/Bowser_Jr.png"},
//        {Name: "Princess Daisy", Img: "https://upload.wikimedia.org/wikipedia/en/b/bd/Daisy_%28Super_Mario_Party%29.png"},
//        {Name: "Wario",          Img: "https://upload.wikimedia.org/wikipedia/en/8/81/Wario.png"},
//        {Name: "Waluigi",        Img: "https://upload.wikimedia.org/wikipedia/en/4/46/Waluigi.png"},
//	}
//	for i := range data {
//		b, err := json.Marshal(data[i])
//		if err != nil {
//			t.Fatal("json.Marshal() failed")
//		}
//		fmt.Println(b)
//		req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(b))
//		res := httptest.NewRecorder()
//		r.ServeHTTP(res, req)
//		expectedBytes, err := json.Marshal(Response{Success: true, Code: 201, Message: "ok", Result: data[i]})
//		if err != nil {t.Fatal("json.Marshal() failed")}
//		if !isJsonObjectsEqual(expectedBytes, res.Body.Bytes()) {t.Fatal("isJsonObjectsEqual() returns false")}
//	}
//}



//
//func TestInsertSampleData(t *testing.T) {
//	expected := []byte(`{"success":true,"code":201,"message":"ok","result":{"ID":`)
//	csvFiles := []struct{filename string; endpoint string; f func(s []string) interface{}}{
//		{"data/achievements.csv", "/achievements", func (s []string)interface{}{return Achievement{Slug:s[0],Name:s[1],Desc:s[2],Img:s[3]}}},
//		{"data/members.csv", "/members", func (s []string)interface{}{return Member{Name:s[0],Img:s[1]}}},
//		{"data/teams.csv", "/teams", func (s []string)interface{}{return Team{Name:s[0],Img:s[1]}}},
//	}
//	for _, data := range csvFiles {
//		fd, err := os.Open(data.filename);
//		if err != nil {
//			t.Error(err)
//		}
//		rows, _ := csv.NewReader(bufio.NewReader(fd)).ReadAll()
//		_ = fd.Close()
//		for _, s := range rows {
//			j, _ := json.Marshal(data.f(s))
//			req, err := http.NewRequest("POST", data.endpoint, bytes.NewBuffer(j))
//			if err != nil { t.Error(err) }
//			res := httptest.NewRecorder()
//			r.ServeHTTP(res, req)
//			if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected){
//				t.Fatal("TestInsertSampleData() failed") }
//		}
//	}
//}
//
///*
//	set up 10 teams with five members each
//		pattern:
//			team 1: members: 11,12,13,14,15
//			team 2: members: 21,21,23,24,25
//			...
// */
//func TestTeamsAddMembers(t *testing.T) {
//	expected := []byte(`{"success":true,"code":200,"message":"ok","result":null}`)
//	for i := 1; i <= 10; i++ {
//		for j := i*10+1; j < i*10+6; j++ {
//			endpoint := fmt.Sprintf("/teams/%d/members/%d",i,j)
//			req, err := http.NewRequest("POST", endpoint, nil)
//			if err != nil { t.Error(err) }
//			res := httptest.NewRecorder()
//			r.ServeHTTP(res, req)
//			if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected){
//				t.Fatal("TestTeamsAddMembers() failed") }
//		}
//	}
//}
//
///*
//	create games # 1..81
// */
//func TestCreateGames(t *testing.T) {
//	expected := []byte(`{"success":true,"code":201,"message":"ok","result":{"ID":`)
//	for i := 1; i <= 81; i++ {
//		req, err := http.NewRequest("POST", "/games", nil)
//		if err != nil { t.Error(err) }
//		res := httptest.NewRecorder()
//		r.ServeHTTP(res, req)
//		if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected) { t.Fatal("TestCreateGames() failed") }
//	}
//}
//
//func TestTeamsJoinGames(t *testing.T) {
//	expected := []byte(`{"success":true,"code":202,"message":"ok","result":null}`)
//	var numGame = 1
//	for i := 1; i <= 9; i++ {
//		for j := 1; j <= 9; j++ {
//			if i == j { continue }
//			endpoint := fmt.Sprintf("/games/%d/teams/%d", numGame, i)	/* team i */
//			req, err := http.NewRequest("POST", endpoint, nil)
//			if err != nil { t.Error(err) }
//			res := httptest.NewRecorder()
//			r.ServeHTTP(res, req)
//			if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected) { t.Fatal("TestTeamsJoinGames() failed") }
//			endpoint = fmt.Sprintf("/games/%d/teams/%d", numGame, j)		/* team j */
//			req, err = http.NewRequest("POST", endpoint, nil)
//			if err != nil { t.Error(err) }
//			res = httptest.NewRecorder()
//			r.ServeHTTP(res, req)
//			if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected) { t.Fatal("TestTeamsJoinGames() failed") }
//			numGame += 1														/* next game */
//		}
//	}
//}
//
////fixme
///*
//func TestUpdateMemberStats(t *testing.T) {
//	var numGame = 1
//	for i := 1; i <= 10; i++ {
//		for j := i*10+1; j < i*10+6; j++ {
//			body := `{"num_attacks":12,"num_hits":23,"amount_damage":500,"num_kills":34,"instant_kills":45,"num_assists":56,"num_spells":67,"spells_damage":78}`
//			endpoint := fmt.Sprintf("/games/%d/members/%d/stats",numGame,j)
//			req, err := http.NewRequest("PUT", endpoint, bytes.NewReader([]byte(body)))
//			if err != nil { t.Error(err) }
//			res := httptest.NewRecorder()
//			r.ServeHTTP(res, req)
//			fmt.Println(res.Code, res.Body.String())
//			if res.Code != http.StatusOK {
//				t.Error("expecting server to return 200")
//			}
//			numGame += 1
//		}
//	}
//}*/
