package table

import "pcclub/time"

type Table struct {
	Amount     int
	TimeTotal  time.Minutes
	Timer      time.Minutes
	ClientName string
}

func (t *Table) Enter(clientName string, start time.Minutes) {
	t.ClientName = clientName
	t.Timer = start
}

func (t *Table) Leave(stop time.Minutes) {
	t.ClientName = ""
	timeUsed := stop - t.Timer
	t.Amount += int(timeUsed+60) / 60
	t.TimeTotal += timeUsed
	t.Timer = stop
}
