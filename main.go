package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	phone    string
}

// use this kind of struct to collect all the errors
type FullOutput struct {
	User
	err error
}

func addPhone(usr User, phones map[string]string) User {
	// emulates some response/processing time
	time.Sleep(1 * time.Second)
	usr.phone = phones[usr.Username]

	return usr
}

func processUsers(wg *sync.WaitGroup, inputCh <-chan User, outputCh chan<- FullOutput, phones map[string]string) {
	defer wg.Done()

	for usr := range inputCh {
		usr = addPhone(usr, phones)

		outputCh <- FullOutput{
			User: usr,
			err:  nil, // in case if you have an error, you can collect this way and handle it later, otherwise you can just log it here
		}
	}
}

func main() {
	users, phones, err := openDataset()
	if err != nil {
		log.Fatal("cannot open testing dataset: %w", err)
	}

	const workerNum = 3

	inputCh := make(chan User)
	outputCh := make(chan FullOutput)
	wg := &sync.WaitGroup{}

	// here we "produce" data
	go func() {
		defer close(inputCh)

		for i := range users {
			inputCh <- users[i]
		}
	}()

	// worker pool itself
	go func() {
		for i := 0; i < workerNum; i++ {
			wg.Add(1)
			go processUsers(wg, inputCh, outputCh, phones)
		}
		wg.Wait()
		close(outputCh)
	}()

	outputUsers := make([]User, 0)

	// here we "collect" data
	for res := range outputCh {
		outputUsers = append(outputUsers, res.User)
	}

	for i := range outputUsers {
		fmt.Println(outputUsers[i])
	}
}
