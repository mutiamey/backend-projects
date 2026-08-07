package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// GitHubEvent represents the simplified structure of a GitHub Event API response.
type GitHubEvent struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		RefType string `json:"ref_type"`
		Size    int    `json:"size"` // Number of commits in a PushEvent
	} `json:"payload"`
}

func main() {
	// Ensure a username argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		os.Exit(1)
	}

	username := os.Args[1]
	events, err := fetchUserEvents(username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if the user has any recent public activity
	if len(events) == 0 {
		fmt.Printf("No recent activity found for user '%s'.\n", username)
		return
	}

	// Display formatted events
	fmt.Printf("Recent activity for %s:\n", username)
	for _, event := range events {
		formattedMessage := formatEvent(event)
		if formattedMessage != "" {
			fmt.Printf("- %s\n", formattedMessage)
		}
	}
}

// fetchUserEvents calls the GitHub REST API to retrieve public events for a given user.
func fetchUserEvents(username string) ([]GitHubEvent, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	// Create an HTTP client with a 10-second timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Prepare the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header as required by GitHub API guidelines
	req.Header.Set("User-Agent", "GitHub-User-Activity-CLI")

	// Execute the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GitHub API: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP status codes
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user '%s' not found", username)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON into slice of GitHubEvent
	var events []GitHubEvent
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return events, nil
}

// formatEvent formats a GitHub event into a human-readable string based on its type.
func formatEvent(event GitHubEvent) string {
	switch event.Type {
	case "PushEvent":
		commits := event.Payload.Size
		if commits == 1 {
			return fmt.Sprintf("Pushed 1 commit to %s", event.Repo.Name)
		}
		return fmt.Sprintf("Pushed %d commits to %s", commits, event.Repo.Name)

	case "IssuesEvent":
		action := event.Payload.Action
		if action == "" {
			action = "interacted with"
		}
		return fmt.Sprintf("%s an issue in %s", capitalize(action), event.Repo.Name)

	case "WatchEvent":
		return fmt.Sprintf("Starred %s", event.Repo.Name)

	case "CreateEvent":
		refType := event.Payload.RefType
		if refType == "" {
			refType = "repository"
		}
		return fmt.Sprintf("Created a new %s in %s", refType, event.Repo.Name)

	case "ForkEvent":
		return fmt.Sprintf("Forked %s", event.Repo.Name)

	case "IssueCommentEvent":
		return fmt.Sprintf("Commented on an issue in %s", event.Repo.Name)

	case "PullRequestEvent":
		action := event.Payload.Action
		if action == "" {
			action = "updated"
		}
		return fmt.Sprintf("%s a pull request in %s", capitalize(action), event.Repo.Name)

	default:
		// Return a general fallback message for other unhandled event types
		return fmt.Sprintf("%s in %s", event.Type, event.Repo.Name)
	}
}

// capitalize turns the first character of a string to uppercase.
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return fmt.Sprintf("%c%s", s[0]-32, s[1:])
}
