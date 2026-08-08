package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		Commits []struct {
			Sha string `json:"sha"`
		} `json:"commits"`
		Issue struct {
			Number int `json:"number"`
		} `json:"issue"`
		PullRequest struct {
			Number int `json:"number"`
		} `json:"pull_request"`
	} `json:"payload"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		os.Exit(1)
	}

	username := os.Args[1]
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Set User-Agent as required by GitHub API guidelines
	req.Header.Set("User-Agent", "Go-GitHub-Activity-CLI")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching data from GitHub: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Error: User '%s' not found.\n", username)
		os.Exit(1)
	} else if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: GitHub API returned status code %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Printf("No recent activity found for user '%s'.\n", username)
		return
	}

	displayActivity(events)
}

func displayActivity(events []Event) {
	fmt.Println("Output:")
	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			commitCount := len(event.Payload.Commits)
			commitStr := "commit"
			if commitCount > 1 {
				commitStr = "commits"
			}
			fmt.Printf("- Pushed %d %s to %s\n", commitCount, commitStr, event.Repo.Name)

		case "IssuesEvent":
			action := event.Payload.Action
			if action == "" {
				action = "interacted with"
			}
			fmt.Printf("- %s an issue in %s\n", capitalize(action), event.Repo.Name)

		case "WatchEvent":
			fmt.Printf("- Starred %s\n", event.Repo.Name)

		case "CreateEvent":
			fmt.Printf("- Created a new resource in %s\n", event.Repo.Name)

		case "PullRequestEvent":
			action := event.Payload.Action
			if action == "" {
				action = "updated"
			}
			fmt.Printf("- %s a pull request in %s\n", capitalize(action), event.Repo.Name)

		case "IssueCommentEvent":
			fmt.Printf("- Commented on an issue in %s\n", event.Repo.Name)

		case "ForkEvent":
			fmt.Printf("- Forked %s\n", event.Repo.Name)

		default:
			fmt.Printf("- %s in %s\n", event.Type, event.Repo.Name)
		}
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := s[0]
	if firstChar >= 'a' && firstChar <= 'z' {
		firstChar -= 32
	}
	return string(firstChar) + s[1:]
}
