package main

import (
    "fmt"
    "os"
    "log"
    "net/http"
)

func main() {

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s", os.Getenv("USER_NAME"))
    })

    log.Fatal(http.ListenAndServe(":8081", nil))

}
