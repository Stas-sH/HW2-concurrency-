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

func generateUsers(count int) []User {
	users := make([]User, count)
	wg3 := new(sync.WaitGroup)
	for i := 0; i < count; i++ {
		wg3.Add(1)
		go func(i int, wg3 *sync.WaitGroup) {
			users[i] = User{
				id:    i + 1,
				email: fmt.Sprintf("user%d@company.com", i+1),
				logs:  generateLogs(rand.Intn(1000)),
			}
			fmt.Printf("generated user %d\n", i+1)
			time.Sleep(time.Millisecond * 100)
			wg3.Done()
		}(i, wg3)
	}
	wg3.Wait()
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

const countUsers = 100

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	users := generateUsers(countUsers)

	tasks := make(chan User, countUsers)
	wg := new(sync.WaitGroup)

	for i := 0; i < countUsers; i++ {
		wg.Add(1)
		go WorkerForSaveUserInfo(tasks, wg)
	}
	for i := 0; i < countUsers; i++ {
		tasks <- users[i]
	}
	wg.Wait()
	close(tasks)
	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
	fmt.Println("Len(users) - ", len(users))

}

func WorkerForSaveUserInfo(tasks <-chan User, wg *sync.WaitGroup) {
	for t := range tasks {
		fmt.Printf("WRITING FILE FOR UID %d\n", t.id)

		filename := fmt.Sprintf("users/uid%d.txt", t.id)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		wg1 := new(sync.WaitGroup)
		resultFromGetActivityInfo := make(chan string, 1)
		wg1.Add(1)
		go WorkerByGetActivityInfo(t, resultFromGetActivityInfo, wg1)

		information := <-resultFromGetActivityInfo
		wg1.Wait()
		file.WriteString(information)

		time.Sleep(time.Second)
		wg.Done()
	}
}

func WorkerByGetActivityInfo(t User, resultFromGetActivityInfo chan<- string, wg1 *sync.WaitGroup) {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", t.id, t.email)
	for index, item := range t.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	resultFromGetActivityInfo <- output
	wg1.Done()
}
