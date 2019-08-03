package main

import (
    "bufio"
    "encoding/csv"
    "encoding/json"
    "net/http"
    "bytes"
    "fmt"
    _"io"
    "io/ioutil"
    "log"
    "os"
)

type Achievement struct {
//    gorm.Model
    Slug            string      `gorm:"unique" json:"slug"`
    Name            string      `gorm:"unique" json:"name"`
    Desc            string      `json:"desc"`
    Img             string      `json:"img"`
}

type Member struct {
//    gorm.Model
    Name            string      `gorm:"unique" json:"name"`
    Img             string      `json:"img"`
}

type Team struct {
//    gorm.Model
    Name            string      `gorm:"unique" json:"name"`
    Img             string      `json:"img"`
}


type sampleData struct {
    filename        string
    endpoint        string
    f               func (s []string) interface{}
}

func insertSampleData(server string, array []sampleData) {
    for _, data := range array {
        fd, err := os.Open(data.filename);
        if err != nil {
            log.Fatal(err)
        }
        strs, _ := csv.NewReader(bufio.NewReader(fd)).ReadAll()
        defer fd.Close()
        for _, s := range strs {
            a := data.f(s)
            j, _ := json.Marshal(a)
            resp, err := http.Post(server + data.endpoint, "application/json", bytes.NewBuffer(j))
            if (err != nil) {
                log.Fatal(err)
            }
            body, _ := ioutil.ReadAll(resp.Body)
            fmt.Println(resp.Status, string(body))
            if resp.StatusCode != 201 {
                fmt.Println(a)
            }
        }
    }
}

func main() {
    data := []sampleData{
        {"data/achievements.csv", "/achievements", func (s []string)interface{}{return Achievement{s[0],s[1],s[2],s[3]}}},
        {"data/members.csv", "/members", func (s []string)interface{}{return Member{s[0],s[1]}}},
        {"data/teams.csv", "/teams", func (s []string)interface{}{return Team{s[0],s[1]}}},
    }
    insertSampleData("http://0.0.0.0:4242", data)
}
