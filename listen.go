package main

import (
    "io"
    "log"
    "fmt"
    "net/http"
    "encoding/json"
)

//
//    /createAchievement
//    params:     slug (mandatory), title, description, img (optional)
//



type CreateAchievementRequest struct {
    Slug            string  `json:"slug"`
    Title           string  `json:"title"`
    Description     string  `json:"description"`
    Img             string  `json:"img"`
}

type ResponseError struct {
    Result  string  `json:"result"`
    Info    string  `json:"info"`
}

func createAchievementHandler(w http.ResponseWriter, r *http.Request) {
    body := make([]byte, r.ContentLength)
    n, err := r.Body.Read(body)
    if err == io.EOF && n > 0  {
        w.Header().Set("content-type", "application/json")
        var req CreateAchievementRequest
        if err = json.Unmarshal(body, &req); err != nil {
            js, _ := json.Marshal(ResponseError{"error", err.Error()})
            w.Write(js)
            return
        }
        if len(req.Slug) < 1 {
            js, _ := json.Marshal(ResponseError{"error", "slug cannot be empty"})
            w.Write(js)
            return
        }
        fmt.Println(req)
    }
}

func main() {
    http.HandleFunc("/createAchievement", createAchievementHandler)

    log.Fatal(http.ListenAndServe(":4242", nil))
}
