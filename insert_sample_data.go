package main

import (
    "bufio"
    "encoding/csv"
    _"encoding/json"
    "fmt"
    "io"
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

func main() {
    csvFile, _ := os.Open("data/achievements.csv")
    reader := csv.NewReader(bufio.NewReader(csvFile))
    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }
        a := Achievement{
            Slug: line[0],
            Name: line[1],
            Desc: line[2],
            Img:  line[3],
        }
        fmt.Println(a)
    }
}
