# Task Tracker

A lightweight Command Line Interface (CLI) application written in Go to track and manage your daily tasks.

This project is built to complete the **https://roadmap.sh/projects/task-tracker** from **roadmap.sh**.

## Features

- Add, update, and delete tasks.
- Change task status to `in-progress` or `done`.
- List all tasks or filter by status (`todo`, `in-progress`, `done`).
- Automatic local JSON file storage (`tasks.json`).
- Built purely using the **Go standard library** (zero external dependencies).

## Usage

Navigate into this folder and run:

```bash
# Add a new task
go run main.go add "Buy groceries"

# Update a task description
go run main.go update 1 "Buy groceries and cook dinner"

# Mark status
go run main.go mark-in-progress 1
go run main.go mark-done 1

# List tasks
go run main.go list
go run main.go list todo
go run main.go list in-progress
go run main.go list done

# Delete a task
go run main.go delete 1
```
