package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Constants for task status and default storage filename
const (
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
	FileName         = "tasks.json"
)

// Task represents the data structure for a single task item
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func main() {
	// Parse CLI arguments (os.Args[0] is the program name)
	args := os.Args[1:]

	if len(args) < 1 {
		printUsage()
		return
	}

	command := args[0]

	switch command {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: Missing task description. Usage: task-cli add \"Buy groceries\"")
			return
		}
		addTask(args[1])

	case "update":
		if len(args) < 3 {
			fmt.Println("Error: Missing arguments. Usage: task-cli update <id> \"New description\"")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Task ID must be a valid integer")
			return
		}
		updateTask(id, args[2])

	case "delete":
		if len(args) < 2 {
			fmt.Println("Error: Missing task ID. Usage: task-cli delete <id>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Task ID must be a valid integer")
			return
		}
		deleteTask(id)

	case "mark-in-progress":
		if len(args) < 2 {
			fmt.Println("Error: Missing task ID. Usage: task-cli mark-in-progress <id>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Task ID must be a valid integer")
			return
		}
		updateTaskStatus(id, StatusInProgress)

	case "mark-done":
		if len(args) < 2 {
			fmt.Println("Error: Missing task ID. Usage: task-cli mark-done <id>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Task ID must be a valid integer")
			return
		}
		updateTaskStatus(id, StatusDone)

	case "list":
		filterStatus := ""
		if len(args) >= 2 {
			filterStatus = args[1]
		}
		listTasks(filterStatus)

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
	}
}

// FILE HELPER FUNCTIONS

// loadTasks reads the tasks.json file. Returns an empty slice if the file does not exist.
func loadTasks() ([]Task, error) {
	if _, err := os.Stat(FileName); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(FileName)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// saveTasks writes the current slice of tasks to tasks.json formatted with indentation
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(FileName, data, 0644)
}

// TASK OPERATION FUNCTIONS

func addTask(description string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error reading tasks file: %v\n", err)
		return
	}

	// Auto-increment ID implementation
	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1
	}

	now := time.Now()
	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)
	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Error saving task: %v\n", err)
		return
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newID)
}

func updateTask(id int, newDescription string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error reading tasks file: %v\n", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: Task with ID %d not found.\n", id)
		return
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Printf("Error updating task: %v\n", err)
		return
	}

	fmt.Printf("Task %d updated successfully\n", id)
}

func deleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error reading tasks file: %v\n", err)
		return
	}

	indexToDelete := -1
	for i, t := range tasks {
		if t.ID == id {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		fmt.Printf("Error: Task with ID %d not found.\n", id)
		return
	}

	// Remove element from slice
	tasks = append(tasks[:indexToDelete], tasks[indexToDelete+1:]...)

	if err := saveTasks(tasks); err != nil {
		fmt.Printf("Error deleting task: %v\n", err)
		return
	}

	fmt.Printf("Task %d deleted successfully\n", id)
}

func updateTaskStatus(id int, newStatus string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error reading tasks file: %v\n", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = newStatus
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: Task with ID %d not found.\n", id)
		return
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Printf("Error updating status: %v\n", err)
		return
	}

	fmt.Printf("Task %d marked as %s\n", id, newStatus)
}

func listTasks(filterStatus string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error reading tasks file: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	count := 0
	fmt.Printf("%-4s | %-12s | %s\n", "ID", "Status", "Description")
	fmt.Println("--------------------------------------------------")

	for _, t := range tasks {
		if filterStatus != "" && t.Status != filterStatus {
			continue
		}
		fmt.Printf("%-4d | %-12s | %s\n", t.ID, t.Status, t.Description)
		count++
	}

	if count == 0 {
		fmt.Printf("No tasks found with status '%s'.\n", filterStatus)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  task-cli add \"<description>\"")
	fmt.Println("  task-cli update <id> \"<description>\"")
	fmt.Println("  task-cli delete <id>")
	fmt.Println("  task-cli mark-in-progress <id>")
	fmt.Println("  task-cli mark-done <id>")
	fmt.Println("  task-cli list [done|todo|in-progress]")
}
