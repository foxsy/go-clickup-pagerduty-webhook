package pagerduty

import (
	"go-clickup-pagerduty-webhook/config"
    "github.com/go-resty/resty/v2"
    "fmt"
    "os"
)

func getEscalationPolicyIDForGroup(groupName string) (string, error) {
    for _, group := range config.GroupAppConfig.Groups {
        if group.Name == groupName {
            return group.EscalationPolicyID, nil
        }
    }
    return "", fmt.Errorf("Group %s not found", groupName)
}

func SendPagerDutyAlert(summary string, group string, clickUpTaskLink string) {
    client := resty.New()

	// API key should be stored securely in an environment variable
    apiKey := os.Getenv("PAGERDUTY_API_KEY")

    // You will need a From address when using this API call, it has to be a valid address of someone in PagerDuty
    userEmail := os.Getenv("PAGERDUTY_USER_EMAIL")

    serviceId := os.Getenv("PAGERDUTY_SERVICE_ID")

	if apiKey == "" || userEmail == "" {
        fmt.Println("PAGERDUTY_API_KEY or PAGERDUTY_USER_EMAIL environment variables are not set")
        return
    }

    // Find the escalation policy ID for the group
    escalationPolicyID, err := getEscalationPolicyIDForGroup(group)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Incident payload, using the escalation policy ID based on the group
    payload := map[string]interface{}{
        "incident": map[string]interface{}{
            "type":  "incident",
            "title": summary,
            "service": map[string]interface{}{
                "id":   serviceId,  // Replace with your service ID
                "type": "service_reference",
            },
            "body": map[string]interface{}{
                "type":    "incident_body",
                "details": fmt.Sprintf("ClickUp task link: %s\nGroup: %s", clickUpTaskLink, group),
            },
            "escalation_policy": map[string]interface{}{
                "id":   escalationPolicyID,  // Escalation policy ID from the group
                "type": "escalation_policy_reference",
            },
        },
    }

    // Send the API request
    resp, err := client.R().
        SetBody(payload).
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", "Token token=" + apiKey).
        SetHeader("From", userEmail).
        Post("https://api.pagerduty.com/incidents")

    if err != nil {
        fmt.Printf("Error sending alert to PagerDuty: %v\n", err)
        return
    }

    fmt.Printf("PagerDuty response: %v\n", resp.String())
}
