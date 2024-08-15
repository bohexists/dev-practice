package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}
	return output
}

// Task represents a unit of work for a worker.
type Task struct {
	user User
}

// worker is a function that will be run as a goroutine.
// It processes tasks from the jobs channel and sends the results to the results channel.
func worker(id int, jobs <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range jobs {
		saveUserInfo(task.user)
		fmt.Printf("Worker %d processed user %d\n", id, task.user.id)
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	users := generateUsers(100)
	jobs := make(chan Task, len(users))
	var wg sync.WaitGroup

	// Determine the number of workers
	numWorkers := 100
	fmt.Printf("Using %d workers\n", numWorkers)

	// Start the workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Distribute the work to the workers
	for _, user := range users {
		jobs <- Task{user: user}
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(user User) {

	if err := os.MkdirAll("refactoring/users", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("refactoring/users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
}

func generateUsers(count int) []User {
	users := make([]User, count)

	for i := 0; i < count; i++ {
		users[i] = User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", i+1)
		time.Sleep(time.Millisecond * 100)
	}

	return users
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
