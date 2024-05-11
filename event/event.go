package event

import (
	"fmt"
	"log"
	"pc_club/queue"
	"pc_club/table"
	"pc_club/time"
	"regexp"
	"strconv"
	"strings"
)

type Event struct {
	TimeMinutes time.Minutes
	ID          ID
	ClientName  string
	Table       int
}

type ID int

const (
	ClientEntered ID = 1
	ClientSat     ID = 2
	ClientWaiting ID = 3
	ClientLeft    ID = 4

	GotToGoMate ID = 11
	TableAwaits ID = 12
	Error       ID = 13
)

func ReadEvent(data []string) *Event {
	id, err := strconv.Atoi(data[1])
	if err != nil {
		log.Fatalf("Error converting id in: %s", data)
	}

	if ID(id) == ClientSat {
		tableNumber, err := strconv.Atoi(data[3])
		if err != nil {
			log.Fatal("Error converting table number:", err)
		}

		if !valid(data[2]) {
			panic(fmt.Sprintf("Error in client name in: %s %s %s %s", data[0], data[1], data[2], data[3]))
		}

		return &Event{
			TimeMinutes: time.Atoi(data[0]),
			ID:          ID(id),
			ClientName:  data[2],
			Table:       tableNumber,
		}
	}

	if !valid(data[2]) {
		panic(fmt.Sprintf("Error in client name in: %s %s %s", data[0], data[1], data[2]))
	}
	fmt.Println(strings.Join(data, " "))
	return &Event{
		TimeMinutes: time.Atoi(data[0]),
		ID:          ID(id),
		ClientName:  data[2],
	}
}

func valid(str string) bool {
	pattern := `^[a-z0-9_-]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

func HandleClientEntered(users map[string]int, event *Event, timeOpen, timeClose time.Minutes) {
	if _, ok := users[event.ClientName]; ok {
		fmt.Println(event.TimeMinutes, Error, "YouShallNotPass")
		return
	}

	if event.TimeMinutes < timeOpen || event.TimeMinutes > timeClose {
		fmt.Println(event.TimeMinutes, Error, "NotOpenYet")
		return
	}

	users[event.ClientName] = 0

}

func HandleClientSat(users map[string]int, tables []table.Table, event *Event) {
	fmt.Println(event.TimeMinutes, event.ID, event.ClientName, event.Table)

	if _, ok := users[event.ClientName]; !ok {
		fmt.Println(event.TimeMinutes, Error, "ClientUnknown")
		return
	}

	if event.Table == 0 {
		fmt.Println(event.TimeMinutes, Error, "TableUnknown")
		return
	}

	if tables[event.Table].ClientName != "" {
		fmt.Println(event.TimeMinutes, Error, "PlaceIsBusy")
		return
	}

	if users[event.ClientName] != 0 {
		tables[users[event.ClientName]].Leave(event.TimeMinutes)
	}

	tables[event.Table].Enter(event.ClientName, event.TimeMinutes)
	users[event.ClientName] = event.Table
}

func HandleClientWaiting(users map[string]int, tables []table.Table, q *queue.CircularBuffer, event *Event) {
	if _, ok := users[event.ClientName]; !ok {
		fmt.Println(event.TimeMinutes, Error, "ClientUnknown")
		return
	}

	for i, t := range tables {
		if t.ClientName == "" && i != 0 {
			fmt.Println(event.TimeMinutes, Error, "ICanWaitNoLonger!")
			return
		}
	}

	if q.Full() {
		fmt.Println(event.TimeMinutes, GotToGoMate, event.ClientName)
		delete(users, event.ClientName)
		return
	}

	q.Push(event.ClientName)
}

func HandleClientLeft(users map[string]int, tables []table.Table, q *queue.CircularBuffer, event *Event) {
	uName := event.ClientName
	uData, ok := users[uName]
	if !ok {
		fmt.Println(event.TimeMinutes, Error, "ClientUnknown")
		return
	}

	if !q.Contains(uName) && uData == 0 {
		delete(users, uName)
		return
	}

	if q.Contains(uName) && uData == 0 {
		return
	}

	if q.Empty() {
		tables[uData].Leave(event.TimeMinutes)
		delete(users, uName)
		return
	}

	tables[uData].Leave(event.TimeMinutes)
	delete(users, uName)
	uNameFromQ, ok := q.Pop()
	tables[uData].Enter(uNameFromQ, event.TimeMinutes)
	users[uNameFromQ] = uData
	fmt.Println(event.TimeMinutes, TableAwaits, uNameFromQ, uData)
}
