package table

import (
	"pc_club/internal/user"
	"testing"
)

func TestEnter(t *testing.T) {
	tables := make([]Table, 3)
	clients := []string{"1", "2", "3", "4"}

	tables[1].Enter(user.Name(clients[0]), 0)
	tables[1].Leave(100)
	if tables[1].Amount != 2 {
		t.Errorf("wrong amount on table: want 2 got %d", tables[1].Amount)
	}
	if tables[1].TimeTotal != 100 {
		t.Errorf("wrong time total: want 100 got %d", tables[1].TimeTotal)
	}

	tables[2].Enter(user.Name(clients[0]), 100)
	tables[1].Enter(user.Name(clients[1]), 200)

	tables[1].Leave(201)
	if tables[1].Amount != 3 {
		t.Errorf("wrong amount on table: want 3 got %d", tables[1].Amount)
	}
	if tables[1].TimeTotal != 101 {
		t.Errorf("wrong time total: want 101 got %d", tables[1].TimeTotal)
	}

	tables[2].Leave(201)
	if tables[2].Amount != 2 {
		t.Errorf("wrong amount on table: want 2 got %d", tables[2].Amount)
	}
	if tables[2].TimeTotal != 101 {
		t.Errorf("wrong time total: want 101 got %d", tables[2].TimeTotal)
	}
}
