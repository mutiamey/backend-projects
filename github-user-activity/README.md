# GitHub User Activity

A CLI application built in Go that fetches and displays recent public activities of a specified GitHub user using the GitHub REST API.

This project is built to complete the **[GitHub User Activity Challenge](https://roadmap.sh/projects/github-user-activity)** from **roadmap.sh**.

## Features

- Fetches recent events directly from the official GitHub REST API.
- Formats activity types (`PushEvent`, `WatchEvent`, `IssuesEvent`, `CreateEvent`, etc.).
- Robust error handling for non-existent users, empty activity feeds, and network issues.
- Built strictly with the **Go Standard Library** (zero external dependencies).

## Usage

Navigate into this folder and run:

```bash
go run main.go <username>

# Example
go run main.go mutiamey
```