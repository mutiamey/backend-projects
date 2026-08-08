# Expense Tracker CLI

A simple Command Line Interface (CLI) application built in Go to track and manage personal expenses, generate total summaries, and filter expenses by month.

This project is built to complete the **[Expense Tracker Challenge](https://roadmap.sh/projects/expense-tracker)** from **roadmap.sh**.

## Features

- Add an expense with `--description` and `--amount`.
- Update an existing expense description or amount by `--id`.
- List all recorded expenses.
- Delete an expense by `--id`.
- View total expenses summary.
- View monthly expense summary using `--month <1-12>`.
- Automatically persists data to a local `expenses.json` file.
- Built strictly with the **Go Standard Library** (zero external dependencies).

## Usage

Navigate into this folder and run:

```bash
# Add expenses
go run main.go add --description "Lunch" --amount 20
go run main.go add --description "Dinner" --amount 10

# Update an expense
go run main.go update --id 1 --description "Lunch Box" --amount 25

# View list of expenses
go run main.go list

# View overall summary
go run main.go summary

# View summary for a specific month (e.g. August)
go run main.go summary --month 8

# Delete an expense
go run main.go delete --id 2
```