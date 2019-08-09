package main

import (
	"bufio"
	_"bufio"
	"bytes"
	_"bytes"
	"encoding/csv"
	_"encoding/csv"
	"encoding/json"
	_"encoding/json"
	"fmt"
	_"io"
	_ "io"
	_"io/ioutil"
	_ "io/ioutil"
	_"log"
	"net/http"
	_"net/http"
	"net/http/httptest"
	_"net/http/httptest"
	"os"
	_"os"
	"testing"
)

func init() {
	initDatabase()
	setRoutes()
}


func TestInsertSampleData(t *testing.T) {
	expected := []byte(`{"success":true,"code":201,"message":"ok","result":{"ID":`)
	csvFiles := []struct{filename string; endpoint string; f func(s []string) interface{}}{
		{"data/achievements.csv", "/achievements", func (s []string)interface{}{return Achievement{Slug:s[0],Name:s[1],Desc:s[2],Img:s[3]}}},
		{"data/members.csv", "/members", func (s []string)interface{}{return Member{Name:s[0],Img:s[1]}}},
		{"data/teams.csv", "/teams", func (s []string)interface{}{return Team{Name:s[0],Img:s[1]}}},
	}
	for _, data := range csvFiles {
		fd, err := os.Open(data.filename);
		if err != nil {
			t.Error(err)
		}
		rows, _ := csv.NewReader(bufio.NewReader(fd)).ReadAll()
		_ = fd.Close()
		for _, s := range rows {
			j, _ := json.Marshal(data.f(s))
			req, err := http.NewRequest("POST", data.endpoint, bytes.NewBuffer(j))
			if err != nil { t.Error(err) }
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)
			if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected){
				t.Fatal("TestInsertSampleData() failed") }
		}
	}
}

/*
	set up 10 teams with five members each
		pattern:
			team 1: members: 11,12,13,14,15
			team 2: members: 21,21,23,24,25
			...
 */
func TestTeamsAddMembers(t *testing.T) {
	expected := []byte(`{"success":true,"code":200,"message":"ok","result":null}`)
	for i := 1; i <= 10; i++ {
		for j := i*10+1; j < i*10+6; j++ {
			endpoint := fmt.Sprintf("/teams/%d/members/%d",i,j)
			req, err := http.NewRequest("POST", endpoint, nil)
			if err != nil { t.Error(err) }
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)
			if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected){
				t.Fatal("TestTeamsAddMembers() failed") }
		}
	}
}

/*
	create games # 1..81
 */
func TestCreateGames(t *testing.T) {
	expected := []byte(`{"success":true,"code":201,"message":"ok","result":{"ID":`)
	for i := 1; i <= 81; i++ {
		req, err := http.NewRequest("POST", "/games", nil)
		if err != nil { t.Error(err) }
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected) { t.Fatal("TestCreateGames() failed") }
	}
}

func TestTeamsJoinGames(t *testing.T) {
	expected := []byte(`{"success":true,"code":202,"message":"ok","result":null}`)
	var numGame = 1
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			if i == j { continue }
			endpoint := fmt.Sprintf("/games/%d/teams/%d", numGame, i)	/* team i */
			req, err := http.NewRequest("POST", endpoint, nil)
			if err != nil { t.Error(err) }
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)
			if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected) { t.Fatal("TestTeamsJoinGames() failed") }
			endpoint = fmt.Sprintf("/games/%d/teams/%d", numGame, j)		/* team j */
			req, err = http.NewRequest("POST", endpoint, nil)
			if err != nil { t.Error(err) }
			res = httptest.NewRecorder()
			r.ServeHTTP(res, req)
			if 0 != bytes.Compare(res.Body.Bytes()[:len(expected)], expected) { t.Fatal("TestTeamsJoinGames() failed") }
			numGame += 1														/* next game */
		}
	}
}

//fixme
/*
func TestUpdateMemberStats(t *testing.T) {
	var numGame = 1
	for i := 1; i <= 10; i++ {
		for j := i*10+1; j < i*10+6; j++ {
			body := `{"num_attacks":12,"num_hits":23,"amount_damage":500,"num_kills":34,"instant_kills":45,"num_assists":56,"num_spells":67,"spells_damage":78}`
			endpoint := fmt.Sprintf("/games/%d/members/%d/stats",numGame,j)
			req, err := http.NewRequest("PUT", endpoint, bytes.NewReader([]byte(body)))
			if err != nil { t.Error(err) }
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)
			fmt.Println(res.Code, res.Body.String())
			if res.Code != http.StatusOK {
				t.Error("expecting server to return 200")
			}
			numGame += 1
		}
	}
}*/