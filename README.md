# Backend Projects

A collection of backend practice projects built with **Go**, based on challenges from [roadmap.sh/projects](https://roadmap.sh/projects).

## Project List

| No | Project Name | Description | Folder Link |
|---|---|---|---|
| 1 | **Task Tracker CLI** | A CLI application to track and manage tasks (CRUD) using local JSON storage. | [`/task-tracker`](./task-tracker) |
| 2 | **GitHub User Activity CLI** | A CLI application to fetch recent GitHub user activity via the GitHub REST API. | [`/github-user-activity`](./github-user-activity) |

## Getting Started

To run any project, navigate into its corresponding directory:

```bash
# Task Tracker
cd task-tracker
go run main.go list

# GitHub User Activity
cd github-user-activity
go run main.go mutiamey  # ./githuh-activity username
```