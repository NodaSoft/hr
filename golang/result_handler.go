package main

import "fmt"

// handleProcessingResults print
func handleProcessingResults(
	processed []*Task,
	errs []error,
) {
	// print errors
	fmt.Println()
	println("Errors:")
	for _, err := range errs {
		fmt.Println(err)
	}

	// print ! errors
	fmt.Println()
	println("Done tasks:")
	for _, task := range processed {
		fmt.Printf(
			"task_id: %3d, created_at: %s, result: %s\n",
			task.ID, task.CreatedAt, task.Result,
		)
	}
}
