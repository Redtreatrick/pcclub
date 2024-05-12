package event

import (
	"github.com/Redtreatrick/pcclub/queue"
	"github.com/Redtreatrick/pcclub/table"
	"testing"
)

import (
	"fmt"
)

func TestClientEntered(t *testing.T) {
	users := make(map[string]int)

	event1 := &Event{
		TimeMinutes: 1,
		ID:          ClientEntered,
		ClientName:  "1",
	}

	event2 := &Event{
		TimeMinutes: 2,
		ID:          ClientEntered,
		ClientName:  "1",
	}

	HandleClientEntered(users, event1, 0, 100)
	if len(users) != 1 {
		t.Errorf("len(users) = %d, want 1", len(users))
	}
	HandleClientEntered(users, event2, 0, 100)
	if len(users) != 1 {
		t.Errorf("len(users) = %d, want 1", len(users))
	}

	event3 := &Event{
		TimeMinutes: 3,
		ID:          ClientEntered,
		ClientName:  "2",
	}
	HandleClientEntered(users, event3, 0, 100)
	if len(users) != 2 {
		t.Errorf("len(users) = %d, want 1", len(users))
	}

}

func TestClientSat(t *testing.T) {
	users := make(map[string]int)
	users["1"] = 0
	users["2"] = 0
	tables := make([]table.Table, 4)

	event1 := &Event{
		TimeMinutes: 1,
		ID:          ClientSat,
		ClientName:  "1",
		Table:       1,
	}

	HandleClientSat(users, tables, event1)
	if users["1"] != 1 || tables[1].ClientName != "1" {
		t.Errorf("users 1 table = %d, want %d\ntable 1 user = %s, want %s\n",
			users["1"], 1, tables[1].ClientName, "1")
	}

	event2 := &Event{
		TimeMinutes: 2,
		ID:          ClientSat,
		ClientName:  "2",
		Table:       1,
	}
	HandleClientSat(users, tables, event2)
	if users["2"] != 0 || tables[1].ClientName != "1" {
		t.Errorf("users 2 table = %d, want %d\ntable 1 user = %s, want %s\n",
			users["2"], 0, tables[1].ClientName, "1")
	}

	event3 := &Event{
		TimeMinutes: 3,
		ID:          ClientSat,
		ClientName:  "1",
		Table:       2,
	}
	HandleClientSat(users, tables, event3)
	if tables[1].ClientName != "" || tables[1].TimeTotal != 2 {
		t.Errorf("table 1 should have no user, got %s\ntable1 timeTotal should be 2, got %d",
			tables[1].ClientName, tables[1].TimeTotal)
	}
	if users["1"] != 2 {
		t.Errorf("users 1 table = %d, want %d", users["1"], 2)
	}

	event4 := &Event{
		TimeMinutes: 4,
		ID:          ClientSat,
		ClientName:  "1",
		Table:       0,
	}
	HandleClientSat(users, tables, event4)
	if tables[0].ClientName != "" {
		t.Errorf("table 0 should have no user, got %s", tables[0].ClientName)
	}
	if users["1"] != 2 || tables[2].ClientName != "1" {
		t.Errorf("users 1 table = %d, want %d\ntable 1 user = %s, want %s\n",
			users["1"], 2, tables[2].ClientName, "1")
	}
}

func TestClientLeft(t *testing.T) {
	users := make(map[string]int)
	users["1"] = 0
	users["2"] = 0
	users["3"] = 0
	users["4"] = 0
	users["5"] = 0
	users["6"] = 0
	users["7"] = 0
	users["8"] = 0
	users["9"] = 0
	users["10"] = 0

	tables := make([]table.Table, 4)
	q := queue.NewCircularBuffer(3)

	eventsSit := []*Event{{
		TimeMinutes: 1,
		ID:          ClientSat,
		ClientName:  "1",
		Table:       1,
	}, {
		TimeMinutes: 2,
		ID:          ClientSat,
		ClientName:  "2",
		Table:       2,
	}, {
		TimeMinutes: 3,
		ID:          ClientSat,
		ClientName:  "3",
		Table:       3,
	},
	}
	for _, e := range eventsSit {
		HandleClientSat(users, tables, e)
	}

	eventLeftFromIdle := &Event{
		TimeMinutes: 4,
		ID:          ClientLeft,
		ClientName:  "4",
	}
	HandleClientLeft(users, tables, q, eventLeftFromIdle)
	if len(users) != 9 {
		t.Errorf("len(users) = %d, want %d", len(users), 9)
	}

	eventLeftFromPCWithEmptyQ := &Event{
		TimeMinutes: 10,
		ID:          ClientLeft,
		ClientName:  "1",
	}
	HandleClientLeft(users, tables, q, eventLeftFromPCWithEmptyQ)
	if tables[1].Amount != 1 || tables[1].TimeTotal != 9 || tables[1].ClientName != "" {
		t.Errorf("error on client left from pc with empty q")
	}
	if len(users) != 8 {
		t.Errorf("len(users) = %d, want %d", len(users), 8)
	}

	eventEnterPC := &Event{
		TimeMinutes: 7,
		ID:          ClientSat,
		ClientName:  "6",
		Table:       1,
	}
	HandleClientSat(users, tables, eventEnterPC)

	eventsWaiting := []*Event{{
		TimeMinutes: 5,
		ID:          ClientWaiting,
		ClientName:  "4",
	}, {
		TimeMinutes: 8,
		ID:          ClientWaiting,
		ClientName:  "7",
	}, {
		TimeMinutes: 9,
		ID:          ClientWaiting,
		ClientName:  "8",
	}, {
		TimeMinutes: 10,
		ID:          ClientWaiting,
		ClientName:  "9",
	}, {
		TimeMinutes: 11,
		ID:          ClientWaiting,
		ClientName:  "10",
	}}
	for _, e := range eventsWaiting {
		HandleClientWaiting(users, tables, q, e)
	}
	if len(users) != 7 {
		t.Errorf("len(users) = %d, want 7", len(users))
	}
	fmt.Println(tables)
	fmt.Println(q)
	fmt.Println(users)

	eventsLeft := []*Event{{
		TimeMinutes: 20,
		ID:          ClientLeft,
		ClientName:  "2",
	}, {
		TimeMinutes: 30,
		ID:          ClientLeft,
		ClientName:  "3",
	}, {
		TimeMinutes: 40,
		ID:          ClientLeft,
		ClientName:  "6",
	}, {
		TimeMinutes: 50,
		ID:          ClientLeft,
		ClientName:  "5",
	}, {
		TimeMinutes: 60,
		ID:          ClientLeft,
		ClientName:  "4",
	}}
	for _, e := range eventsLeft {
		HandleClientLeft(users, tables, q, e)
		if tables[0].ClientName != "" {
			t.Errorf("table 0 should have no user, got %s", tables[0].ClientName)
		}
		if tables[0].TimeTotal != 0 || tables[0].Timer != 0 || tables[0].Amount != 0 {
			t.Errorf("table 0 should have no timeTotal, got %s\nno timer, got %s\nno amount, got %d", tables[0].TimeTotal, tables[0].Timer, tables[0].Amount)
		}
	}

	fmt.Println(tables)
	fmt.Println(q)
	fmt.Println(users)

}
