# Expense Tracker

A simple Command Line Interface (CLI) application built in Go to manage your expenses and display summary reports.

This project is built to complete the **[Expense Tracker Challenge](https://roadmap.sh/projects/expense-tracker)** from **roadmap.sh**.

## Features

- Add an expense with `--description` and `--amount`.
- List all recorded expenses.
- Delete an expense by `--id`.
- View total expenses summary.
- View monthly summary using `--month <1-12>`.
- Automatically saves data to `expenses.json`.

## Usage

```bash
# Add expenses
go run main.go add --description "Lunch" --amount 20
go run main.go add --description "Dinner" --amount 15

# View list
go run main.go list

# Summary
go run main.go summary
go run main.go summary --month 8

# Delete
go run main.go delete --id 1
```