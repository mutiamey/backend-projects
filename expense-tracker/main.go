package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Expense struct {
	ID          int     `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

const fileName = "expenses.json"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd()
	case "list":
		handleList()
	case "delete":
		handleDelete()
	case "summary":
		handleSummary()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  expense-tracker add --description \"Lunch\" --amount 20")
	fmt.Println("  expense-tracker list")
	fmt.Println("  expense-tracker delete --id 1")
	fmt.Println("  expense-tracker summary [--month 8]")
}

func handleAdd() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	desc := addCmd.String("description", "", "Expense description")
	amount := addCmd.Float64("amount", 0, "Expense amount")
	addCmd.Parse(os.Args[2:])

	if *desc == "" || *amount <= 0 {
		fmt.Println("Error: --description and a positive --amount are required.")
		return
	}

	expenses := loadExpenses()
	newID := 1
	if len(expenses) > 0 {
		newID = expenses[len(expenses)-1].ID + 1
	}

	expense := Expense{
		ID:          newID,
		Date:        time.Now().Format("2006-01-02"),
		Description: *desc,
		Amount:      *amount,
	}

	expenses = append(expenses, expense)
	saveExpenses(expenses)

	fmt.Printf("# Expense added successfully (ID: %d)\n", newID)
}

func handleList() {
	expenses := loadExpenses()
	if len(expenses) == 0 {
		fmt.Println("No expenses recorded yet.")
		return
	}

	fmt.Printf("%-4s %-12s %-20s %-10s\n", "ID", "Date", "Description", "Amount")
	for _, e := range expenses {
		fmt.Printf("%-4d %-12s %-20s $%.2f\n", e.ID, e.Date, e.Description, e.Amount)
	}
}

func handleDelete() {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	id := deleteCmd.Int("id", 0, "Expense ID to delete")
	deleteCmd.Parse(os.Args[2:])

	if *id <= 0 {
		fmt.Println("Error: Please provide a valid --id.")
		return
	}

	expenses := loadExpenses()
	updatedExpenses := []Expense{}
	found := false

	for _, e := range expenses {
		if e.ID == *id {
			found = true
			continue
		}
		updatedExpenses = append(updatedExpenses, e)
	}

	if !found {
		fmt.Printf("Error: Expense with ID %d not found.\n", *id)
		return
	}

	saveExpenses(updatedExpenses)
	fmt.Println("# Expense deleted successfully")
}

func handleSummary() {
	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	month := summaryCmd.Int("month", 0, "Filter summary by month (1-12)")
	summaryCmd.Parse(os.Args[2:])

	expenses := loadExpenses()
	var total float64

	if *month > 0 {
		if *month < 1 || *month > 12 {
			fmt.Println("Error: Month must be between 1 and 12.")
			return
		}

		for _, e := range expenses {
			t, err := time.Parse("2006-01-02", e.Date)
			if err == nil && int(t.Month()) == *month && t.Year() == time.Now().Year() {
				total += e.Amount
			}
		}

		monthName := time.Month(*month).String()
		fmt.Printf("# Total expenses for %s: $%.2f\n", monthName, total)
	} else {
		for _, e := range expenses {
			total += e.Amount
		}
		fmt.Printf("# Total expenses: $%.2f\n", total)
	}
}

func loadExpenses() []Expense {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return []Expense{}
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return []Expense{}
	}

	var expenses []Expense
	json.Unmarshal(data, &expenses)
	return expenses
}

func saveExpenses(expenses []Expense) {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		fmt.Printf("Error saving expenses: %v\n", err)
		return
	}
	os.WriteFile(fileName, data, 0644)
}
