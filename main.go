package main

import (
   "fmt"
   "log"
   "net/http"
   "os"
)

// main sets up the HTTP server, routes, and serves frontend.
func main() {
   http.HandleFunc("/api/generate-qr", qrHandler)
   fs := http.FileServer(http.Dir("./dist"))
   http.Handle("/", fs)
   port := os.Getenv("PORT")
   if port == "" {
       port = "8080"
   }
   fmt.Printf("Server listening on http://localhost:%s\n", port)
   log.Fatal(http.ListenAndServe(":"+port, nil))
}