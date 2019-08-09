package main

import (
	_ "bufio"
	"bytes"
	_ "bytes"
	_ "encoding/csv"
	"encoding/json"
	_ "encoding/json"
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
	"status"		:nil,	// Game{}
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
	if requestObject["result"] == nil || responseObject["result"] == nil {return true} /*both have empty results*/
	requestResult := requestObject["result"].(map[string]interface{})
	responseResult := responseObject["result"].(map[string]interface{})
	for key, value := range requestResult {
        if _, ok := ignoreKeys[key]; ok {continue}			/*ignore CreatedAt, UpdatedAt, etc*/
		if responseResult[key] != value {return false}		/*objects are different*/
	}
	return true												/*objects are the same*/
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
		t.Fatal("TestInsertMembers failed")
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
	}
	if !requestAndCompareResponse(endpoint, method, data, response) {
		t.Fatal("TestInsertTeams failed")
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
		t.Fatal("TestInsertAchievements failed")
	}
}

func TestErrorInsertDuplicate(t *testing.T) {
	endpoint := "/members"
	method := http.MethodPost
	response := Response{Success: false, Code: 406, Message: "record already exists", Result: nil}
	data := []interface{}{
		Member{Name: "Mario", Img: "https://upload.wikimedia.org/wikipedia/en/a/a9/MarioNSMBUDeluxe.png"},
	}
	if !requestAndCompareResponse(endpoint, method, data, response) {
		t.Fatal("TestErrorInsertDuplicate failed")
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
