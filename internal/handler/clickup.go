package handler

import (
    "go-clickup-pagerduty-webhook/config"
    "go-clickup-pagerduty-webhook/internal/pagerduty"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "bytes"
    "os"
)

// Define structs based on the payload structure for task updates and creation
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"username"`
    Email string `json:"email"`
}

type HistoryItem struct {
    ID     string `json:"id"`
    Field  string `json:"field"`
    Before struct {
        Priority string `json:"priority,omitempty"`
    } `json:"before,omitempty"`
    After struct {
        Priority string `json:"priority,omitempty"`
        Status   string `json:"status,omitempty"`
    } `json:"after,omitempty"`
}

type ClickUpEvent struct {
    Event        string        `json:"event"`
    HistoryItems []HistoryItem `json:"history_items"`
    TaskID       string        `json:"task_id"`
    WebhookID    string        `json:"webhook_id"`
    Task         map[string]interface{} `json:"task"`  // Handle dynamic task data
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	authToken := os.Getenv("CLICKUP_WEBHOOK_TOKEN")
	token := r.URL.Query().Get("token")
    if token != authToken {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        fmt.Println("token error")
        return
    }

    // Read and print the raw request body
    body, ioErr := io.ReadAll(r.Body)
    if ioErr != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }
    fmt.Println(string(body))  // Debugging: print the raw JSON payload

    // Reset the request body so it can be read again
    r.Body = io.NopCloser(bytes.NewBuffer(body))

    // Parse the body into the ClickUpEvent struct
    var event ClickUpEvent
    err := json.NewDecoder(r.Body).Decode(&event)
    if err != nil {
        http.Error(w, "Error decoding JSON", http.StatusBadRequest)
        return
    }

    // Construct the link to the ClickUp task
    clickUpTaskLink := fmt.Sprintf("https://app.clickup.com/t/%s", event.TaskID)

    // Handle task updates and task creations
    if event.Event == "taskUpdated" {
        for _, item := range event.HistoryItems {
            if item.Field == "priority" {
                fmt.Printf("Task ID: %s - Priority changed from %s to %s\n", event.TaskID, item.Before.Priority, item.After.Priority)
                // If priority is set to urgent, take action
                if item.After.Priority == "urgent" {
                    fmt.Println("Urgent priority detected, sending alert!")
                    // Call the PagerDuty alert function
                    pagerduty.SendPagerDutyAlert("Priority changed to urgent", "devops", clickUpTaskLink)
                }
            }
        }
    } else if event.Event == "taskCreated" {
        // For task creation, check the history items for priority
        for _, item := range event.HistoryItems {
            if item.Field == "priority" && item.After.Priority == "urgent" {
                fmt.Printf("Task ID: %s created with urgent priority\n", event.TaskID)
                // Send an alert for the new urgent task
                pagerduty.SendPagerDutyAlert("New task created with urgent priority", "devops", clickUpTaskLink)
            }
        }
    }

    // Iterate over rules and check for matches (based on your existing logic)
    for _, rule := range config.AppConfig.Rules {
        if rule.Event == event.Event {
            // Look up the condition dynamically using the key from YAML
            if value, exists := event.Task[rule.Condition.Key]; exists && fmt.Sprintf("%v", value) == rule.Condition.Value {
                // Check if the list or space matches the rule
                listMatches := rule.List == "" || rule.List == fmt.Sprintf("%v", event.Task["list_id"])
                spaceMatches := rule.Space == "" || rule.Space == fmt.Sprintf("%v", event.Task["project"].(map[string]interface{})["id"])

                if listMatches || spaceMatches {
                    fmt.Println("Matching rule found! Sending alert to group:", rule.Group)
                    pagerduty.SendPagerDutyAlert(
                        fmt.Sprintf("ClickUp task matched condition: %s = %s", rule.Condition.Key, rule.Condition.Value),
                        rule.Group,
                        clickUpTaskLink,
                    )
                }
            }
        }
    }

    // Respond with a success message
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Webhook received"))
}
