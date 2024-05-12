package main

import (
	"bufio"
	"fmt"
	"github.com/Redtreatrick/pcclub/event"
	"github.com/Redtreatrick/pcclub/queue"
	"github.com/Redtreatrick/pcclub/table"
	"github.com/Redtreatrick/pcclub/time"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	//file, err := os.Open("test_file.txt")
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
			log.Fatal("tableAmount should be integer:", err)
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

	//fmt.Printf("столов: %d\nоткрытие: %v\nзакрытие: %v\nдоход/час: %d\n",
	//tableAmount, timeOpen, timeClose, hourRate)

	fmt.Println(timeOpen)

	tables := make([]table.Table, tableAmount+1) // first t shall be ignored
	users := make(map[string]int)
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

	for _, t := range tables {
		if t.ClientName != "" {
			fmt.Println(timeClose, event.GotToGoMate, t.ClientName)
		}
	}

	fmt.Println(timeClose)

	for num, t := range tables {
		if num == 0 {
			continue
		}

		if t.ClientName != "" {
			t.Leave(timeClose)
		}
		fmt.Println(num, t.Amount*hourRate, t.TimeTotal)
	}

}
