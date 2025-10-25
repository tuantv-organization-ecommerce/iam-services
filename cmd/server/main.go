package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/tvttt/iam-services/internal/app"
)

func main() {
	// Global panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("FATAL: Service crashed with panic: %v\n", r)
			fmt.Printf("Stack trace:\n%s\n", debug.Stack())
			os.Exit(1)
		}
	}()

	// Create application
	application, err := app.New()
	if err != nil {
		fmt.Printf("Failed to create application: %v\n", err)
		os.Exit(1)
	}

	// Initialize dependencies
	if err := application.Initialize(); err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Run application
	if err := application.Run(); err != nil {
		fmt.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
}
