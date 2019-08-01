package main

import (
    "io"
    "log"
    "fmt"
    "net/http"
)

//
//    /createAchievement
//    params:     slug (mandatory), title, description, img (optional)
//



func createAchievementHandler(w http.ResponseWriter, r *http.Request) {
    body := make([]byte, r.ContentLength)
    n, err := r.Body.Read(body)
    if err == io.EOF && n > 0  {
        fmt.Println(string(body))
    }
}

func main() {
    http.HandleFunc("/createAchievement", createAchievementHandler)

    log.Fatal(http.ListenAndServe(":4242", nil))
}
