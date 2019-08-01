package main

import (
    import "log"
    import "fmt"
    import "net/http"
)

//
//    /createAchievement
//    params:     slug (mandatory), title, description, img (optional)
//




func main() {
    http.Handle("/createAchievement", createAchievementHandler)

    http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    })

    log.Fatal(http.ListenAndServe(":4242", nil))
}
