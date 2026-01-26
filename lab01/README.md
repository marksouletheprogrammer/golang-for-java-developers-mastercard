# Lab 1: Environment Setup & Go Basics

Welcome to Go! In this lab, you'll create your first Go module and build a working program. By the end, you'll have a project structure that you'll build upon throughout this course.

### Part 1: Create Your First Module
You'll create a project called "hello-go" - a simple program to get familiar with Go's tooling.

1. Use this lab directory as the root of your project: `cd lab01`
2. Initialize a new Go module using the command `go mod init` with an appropriate module path (e.g., `go mod init github.com/yourusername/hello-go`)
3. Examine the `go.mod` file that was created - it should contain the module path and the Go version. This file tracks your module's dependencies (though you have none yet)

### Part 2: Write Your First Go Program
1. Create a file named `main.go` in your project root (i.e. lab01/main.go).
2. Write a program that prints a welcome message and the current timestamp to the console.

### Part 3: Build and Run
1. Run your program directly using `go run main.go`
2. Build a binary executable using `go build`
3. Run the binary you just created.
4. Note the difference between the two approaches.

### Part 4: Adding External Dependencies
Now you'll learn how to add and use external packages in Go. You'll add a popular logging library and observe how Go manages dependencies.

1. Add the logrus logging library to your project:
   ```bash
   go get github.com/sirupsen/logrus
   ```

2. After running `go get`, examine the changes:
   - Open `go.mod` - you'll see the logrus dependency has been added
   - Notice `go.sum` was created - this file contains checksums for your dependencies and their transitive dependencies (dependencies of dependencies)
   - Look for both **direct** dependencies (logrus) and **indirect** dependencies (packages that logrus needs)

3. Update your `main.go` to use logrus. Replace your existing imports and print statements with:
   
   **Import section** (at the top of the file):
   ```go
   import (
       "time"
       "github.com/sirupsen/logrus"
   )
   ```
   
   **In your main function**, replace your `fmt.Println` calls with:
   ```go
   logrus.Info("Welcome to Go!")
   
   currentTime := time.Now()
   logrus.Infof("Current timestamp: %s", currentTime.Format("2006-01-02 15:04:05 MST"))
   logrus.Infof("ISO 8601 format: %s", currentTime.Format(time.RFC3339))
   ```

4. Run your program again with `go run main.go` and observe the different output format from the logging library.

### Part 5: Explore Go Commands
Experiment with these commands to understand Go's tooling:
1. Run `go mod tidy` - removes unused dependencies and adds missing ones. This ensures your `go.mod` and `go.sum` files accurately reflect the dependencies your code actually uses.
2. Run `go fmt` or `go fmt .` - automatically formats your code to Go standards.
3. Run `go vet` - checks for common mistakes and suspicious code.
4. Run `go help` to see all available commands, and try `go help build` or `go help run` for detailed help on specific commands.
5. Observe what each command does and when you might use them in development.