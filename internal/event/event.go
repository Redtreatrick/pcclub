package event

import (
	"fmt"
	"log"
	"pc_club/internal/queue"
	"pc_club/internal/table"
	"pc_club/internal/time"
	"pc_club/internal/user"
	"strconv"
)

type Event struct {
	TimeMinutes time.Minutes
	ID          ID
	ClientName  user.Name
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
		if !user.Name(data[2]).Ok() {
			panic(fmt.Sprintf("Error in client name in: %s %s %s %s", data[0], data[1], data[2], data[3]))
		}

		return &Event{
			TimeMinutes: time.Atoi(data[0]),
			ID:          ID(id),
			ClientName:  user.Name(data[2]),
			Table:       tableNumber,
		}
	}

	if !user.Name(data[2]).Ok() {
		panic(fmt.Sprintf("Error in client name in: %s %s %s", data[0], data[1], data[2]))
	}

	return &Event{
		TimeMinutes: time.Atoi(data[0]),
		ID:          ID(id),
		ClientName:  user.Name(data[2]),
	}
}

func HandleClientEntered(users map[user.Name]user.Data, event *Event, timeOpen, timeClose time.Minutes) {
	if _, ok := users[event.ClientName]; ok {
		fmt.Println(event.TimeMinutes, Error, "YouShallNotPass")
		return
	}

	if event.TimeMinutes < timeOpen || event.TimeMinutes > timeClose {
		fmt.Println(event.TimeMinutes, Error, "NotOpenYet")
		return
	}

	users[event.ClientName] = user.Data{
		Table: 0,
		//TimeMinutes: 0,
	}
	fmt.Printf("%s %d %s\n", event.TimeMinutes, event.ID, event.ClientName)
}

func HandleClientSat(users map[user.Name]user.Data, tables []table.Table, event *Event) {
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

	if users[event.ClientName].Table != 0 {
		tables[users[event.ClientName].Table].Leave(event.TimeMinutes)
	}

	tables[event.Table].Enter(event.ClientName, event.TimeMinutes)
	users[event.ClientName] = user.Data{Table: event.Table /*TimeMinutes: event.TimeMinutes*/}
	fmt.Println(event.TimeMinutes, event.ID, event.ClientName, event.Table)
}

func HandleClientWaiting(users map[user.Name]user.Data, tables []table.Table, q *queue.CircularBuffer, event *Event) {
	if _, ok := users[event.ClientName]; !ok {
		fmt.Println(event.TimeMinutes, Error, "ClientUnknown")
		return
	}

	for i, t := range tables {
		if t.ClientName == "" && i != 0 {
			fmt.Println(event.TimeMinutes, Error, "ICanWaitNoLonger")
			return
		}
	}

	if q.Full() {
		fmt.Println(event.TimeMinutes, GotToGoMate, event.ClientName)
		delete(users, event.ClientName)
		return
	}

	q.Push(string(event.ClientName))
	fmt.Println(event.TimeMinutes, event.ID, event.ClientName)
}

func HandleClientLeft(users map[user.Name]user.Data, tables []table.Table, q *queue.CircularBuffer, event *Event) {
	// если клиент не в клубе - ошибка
	uName := event.ClientName
	uData, ok := users[uName]
	if !ok {
		fmt.Println(event.TimeMinutes, Error, "ClientUnknown")
		return
	}
	//fmt.Println("client", uName, "with table", uData.Table, "left")
	// если клиент не находится ни за компом ни в очереди, его уход ни на что не влияет
	if !q.Contains(string(uName)) && uData.Table == 0 {
		delete(users, uName)
		return
	}

	// если клиент уходит из очереди, то она подвигается
	if q.Contains(string(uName)) && uData.Table == 0 {
		//panic("non specified case")
		return
	}

	// если клиент уходит из-за стола при пустой очереди - он свободен идти
	if q.Empty() {
		//fmt.Println("client", uName, "with table ", uData.Table, "is free to go")
		tables[uData.Table].Leave(event.TimeMinutes)
		delete(users, uName)
		return
	}

	// если клиент уходит из-за стола при непустой очереди - его стол занимает первый в очереди
	//fmt.Println("client", uName, "with table ", uData.Table, "has to give his place to someone else")
	tables[uData.Table].Leave(event.TimeMinutes)
	delete(users, uName)
	uNameFromQ, ok := q.Pop()
	tables[uData.Table].Enter(user.Name(uNameFromQ), event.TimeMinutes)
	users[user.Name(uNameFromQ)] = user.Data{Table: uData.Table}

}
