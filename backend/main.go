package main

import (
    "log"
    "net/http"
    "os"

    "invoice-backend/internal/server"

    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
    }

    addr := ":8080"
    if v := os.Getenv("PORT"); v != "" {
        addr = ":" + v
    }

    srv := server.New()

    log.Printf("starting server on %s", addr)
    if err := http.ListenAndServe(addr, srv); err != nil {
        log.Fatal(err)
    }
}
