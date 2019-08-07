package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    initDatabase()
    setRoutes()
    fmt.Printf("server ready at: %s\n", serverAddress)
    log.Fatal(http.ListenAndServe(serverAddress, router))
}
