# GitHub User Activity

A Command Line Interface (CLI) application built in Go that fetches and displays recent public activities of a specified GitHub user using the GitHub REST API.

This project is built to complete the **[GitHub User Activity Challenge](https://roadmap.sh/projects/github-user-activity)** from **roadmap.sh**.

## Features

- Fetches recent public activities using the official GitHub REST API (`https://api.github.com/users/<username>/events`).
- Formats event outputs clearly (e.g., `- Pushed 3 commits to ...`, `- Starred ...`, `- Opened an issue in ...`).
- Graceful error handling for non-existent users, empty activity feeds, and network/API errors.
- Built strictly with the **Go Standard Library** (zero external dependencies).

## Usage

Navigate into this folder and run:

```bash
go run main.go <username>

# Example
go run main.go mutiamey
```