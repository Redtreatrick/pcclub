package table

import (
	"pc_club/internal/time"
	"pc_club/internal/user"
)

type Table struct {
	Amount     int
	TimeTotal  time.Minutes
	Timer      time.Minutes
	ClientName user.Name
}

func (t *Table) Enter(clientName user.Name, start time.Minutes) {
	t.ClientName = clientName
	t.Timer = start
}

func (t *Table) Leave(stop time.Minutes) {
	//fmt.Printf("from %v to %v client %v, in total: ", t.Timer, stop, t.ClientName)
	t.ClientName = ""
	timeUsed := stop - t.Timer
	t.Amount += int(timeUsed+60) / 60
	t.TimeTotal += timeUsed
	t.Timer = stop
}
