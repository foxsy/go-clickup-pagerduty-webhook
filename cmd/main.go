package main

import (
    "go-clickup-pagerduty-webhook/config"
    "go-clickup-pagerduty-webhook/internal/handler"
    "fmt"
    "log"
    "net/http"
)

func main() {
    // Load YAML configuration for rules
    config.LoadConfig("config/rules.yaml")
    // Load YAML configuration for groups and escalation policy IDs
    config.LoadGroupConfig("config/groups.yaml")

    // Set up the webhook handler
    http.HandleFunc("/webhook", handler.WebhookHandler)

    fmt.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
