package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pc_club/internal/event"
	"pc_club/internal/queue"
	"pc_club/internal/table"
	"pc_club/internal/time"
	"pc_club/internal/user"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("test_file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var tableAmount int
	var timeOpen, timeClose time.Minutes
	var hourRate int

	if scanner.Scan() {
		tableAmount, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("tableAmount should be integer:", err)
			return
		}
	}

	if scanner.Scan() {
		timeOpenStr, timeCloseStr, ok := strings.Cut(scanner.Text(), " ")
		if !ok {
			log.Fatal("Error parsing time\n", timeOpenStr)
		}

		timeOpen, timeClose = time.Atoi(timeOpenStr), time.Atoi(timeCloseStr)
		if timeOpen > timeClose {
			log.Fatal("open time should be less than close time\n", timeCloseStr)
		}
	}

	if scanner.Scan() {
		hourRate, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Error converting hourRate:", err)
			return
		}
	}

	fmt.Printf("столов: %d\nоткрытие: %v\nзакрытие: %v\nдоход/час: %d\n",
		tableAmount, timeOpen, timeClose, hourRate)

	tables := make([]table.Table, tableAmount+1) // first t shall be ignored
	users := make(map[user.Name]user.Data)
	q := queue.NewCircularBuffer(tableAmount)

	for scanner.Scan() {
		currEvent := event.ReadEvent(strings.Fields(scanner.Text()))
		switch currEvent.ID {
		case event.ClientEntered:
			event.HandleClientEntered(users, currEvent, timeOpen, timeClose)
		case event.ClientSat:
			event.HandleClientSat(users, tables, currEvent)
		case event.ClientWaiting:
			event.HandleClientWaiting(users, tables, q, currEvent)
		case event.ClientLeft:
			event.HandleClientLeft(users, tables, q, currEvent)

		}
	}

	for num, t := range tables {
		if t.ClientName != "" {
			t.Leave(timeClose)
		}
		fmt.Println(num, t.Amount*hourRate, t.TimeTotal)
	}

}
