package queue

import (
	"fmt"
	"testing"
)

func TestPush(t *testing.T) {
	cb := NewCircularBuffer(3)
	cb.Push("1")
	cb.Push("2")
	cb.Push("3")
	if cb.count != 3 || cb.buffer[0] != "1" || cb.buffer[1] != "2" || cb.buffer[2] != "3" {
		t.Errorf("Push failed")
	}
	if ok := cb.Push("4"); ok {
		t.Errorf("Push succeeded on full buffer")
	}
}

func TestPop(t *testing.T) {
	cb := NewCircularBuffer(3)
	_, ok := cb.Pop()
	if ok {
		t.Errorf("Expected Pop to fail on empty buffer")
	}

	cb.Push("1")
	cb.Push("2")
	cb.Push("3")
	popped, ok := cb.Pop()
	if !ok || popped != "1" {
		t.Errorf("Pop failed")
	}
	popped, ok = cb.Pop()
	if !ok || popped != "2" {
		t.Errorf("Pop failed")
	}
	popped, ok = cb.Pop()
	if !ok || popped != "3" {
		t.Errorf("Pop failed")
	}
	if cb.count != 0 {
		t.Errorf("Pop did not decrement count")
	}
	ok = cb.Push("4")
	if !ok {
		t.Errorf("Push failed")
	}
	ok = cb.Push("5")
	if !ok {
		t.Errorf("Push failed")
	}
	popped, ok = cb.Pop()
	if !ok || popped != "4" {
		t.Errorf("Pop failed")
	}
	ok = cb.Push("6")
	if !ok {
		t.Errorf("Push failed")
	}
	ok = cb.Push("7")
	if !ok {
		t.Errorf("Push failed")
	}
	ok = cb.Push("8")
	if ok {
		t.Errorf("Push succeeded on full buffer")
	}

}

func TestContains(t *testing.T) {
	cb := NewCircularBuffer(3)
	cb.Push("1")
	cb.Push("2")
	cb.Push("3")

	if !cb.Contains("1") || !cb.Contains("2") || !cb.Contains("3") {
		t.Errorf("Contains failed")
	}

	if cb.Contains("") {
		t.Errorf("Contains succeeded on wrong case")
	}
	if cb.Contains("4") {
		t.Errorf("Contains succeeded on wrong case")
	}

	cb.Pop()
	if cb.Contains("1") {
		t.Errorf("Contains succeeded on wrong case")
	}
	if !cb.Contains("2") {
		t.Errorf("Contains failed on 2")
	}
	if !cb.Contains("3") {
		t.Errorf("Contains failed on 3")
	}

	cb2 := NewCircularBuffer(3)
	if cb.Contains("1") || cb.Contains("") {
		t.Errorf("Contains succeeded on wrong case")
	}

	cb2.Push("1")
	cb2.Push("2")
	if !cb2.Contains("1") || !cb2.Contains("2") || cb2.Contains("3") {
		t.Errorf("Contains failed")
	}
	cb2.Push("3")
	if !cb2.Contains("3") || !cb2.Contains("1") || !cb2.Contains("2") {
		t.Errorf("Contains failed")
	}
}

func TestFind(t *testing.T) {
	cb := NewCircularBuffer(3)
	if cb.Find("") != -1 {
		t.Errorf("Find succeeded on empty buffer")
	}

	cb.Push("1")
	cb.Push("2")
	if cb.Find("2") != 1 {
		t.Errorf("Find failed on tail > head non full buffer")
	}
	cb.Push("3")
	if cb.Find("2") != 1 {
		t.Errorf("Find failed on full buffer")
	}

	cb.Pop()
	if cb.Find("2") != 1 {
		t.Errorf("Find failed on tail < head non full buffer, got %d", cb.Find("2"))
	}

	cb2 := NewCircularBuffer(100)
	for i := 0; i < 99; i++ {
		cb2.Push(fmt.Sprintf("%d", i))
	}

	for i := 0; i < 50; i++ {
		cb2.Pop()
	}

	for i := 0; i < 99; i++ {
		if cb2.Find(fmt.Sprintf("%d", i)) != i {
			t.Errorf("Find failed, want %d got %d", i, cb2.Find(fmt.Sprintf("%d", i)))
		}
	}

}

func TestLength(t *testing.T) {
	cb := NewCircularBuffer(5)
	cb.Push("1")
	cb.Push("2")
	cb.Push("3")
	cb.Push("4")
	if cb.Length() != 4 {
		t.Errorf("Length failed, want %d got %d", 4, cb.Length())
	}
	cb.Push("5")
	if cb.Length() != 5 {
		t.Errorf("Length failed, want %d got %d", 5, cb.Length())
	}

	for i := 4; i >= 0; i-- {
		cb.Pop()
		if cb.Length() != i {
			t.Errorf("Length failed, want %d got %d", 4, cb.Length())
		}
	}

}

func TestEvacuate(t *testing.T) {
	cb := NewCircularBuffer(5)
	cb.Push("1")
	cb.Push("2")
	cb.Push("3")
	cb.Push("4")
	if cb.Length() != 4 {
		t.Errorf("Length failed, want %d got %d", 4, cb.Length())
	}

	// test evacuate on itemPos == c.head
	cb.Evacuate("1")
	cb.Evacuate("2")
	if cb.Length() != 2 {
		t.Errorf("Length failed, want %d got %d", 2, cb.Length())
	}
	if !cb.Contains("3") && !cb.Contains("4") {
		t.Errorf("Evacuate failed")
	}

	// test evacuate on itemPos outside [head:tail]
	cb.Push("5")
	cb.Pop()         // len = 2
	cb.Evacuate("1") // len = 2 still
	if cb.Length() != 2 {
		t.Errorf("Length failed, want %d got %d", 2, cb.Length())
	}

	// test for head = 0, tail = 2, x = 3
	cb.Push("1")
	cb.Push("2")
	cb.Pop()
	cb.Pop()
	if cb.Contains("3") || cb.Contains("4") || cb.Contains("5") ||
		!cb.Contains("1") || !cb.Contains("2") || cb.Length() != 2 {
		t.Errorf("something went wrong!")
	}
	cb.Evacuate("4")
	if cb.Length() != 2 {
		t.Errorf("Length failed, want %d got %d", 2, cb.Length())
	}
	if cb.buffer[3] != "F" {
		t.Errorf("Evacuate failed, want %s got %s", "F", cb.buffer[4])
	}

	// test for head = 1, tail = 4, x = 0
	cb2 := NewCircularBuffer(5)
	for i := 0; i < 4; i++ {
		cb2.Push(fmt.Sprintf("%d", i))
	}
	cb2.Pop()
	if cb2.Length() != 3 {
		t.Errorf("Evacuate failed, want %d got %d", 3, cb.Length())
	}
	cb2.Evacuate("0")
	if cb2.Length() != 3 {
		t.Errorf("Evacuate failed, want %d got %d", 3, cb.Length())
	}
	if cb2.buffer[0] != "F" {
		t.Errorf("Evacuate failed, want %s got %s", "F", cb2.buffer[0])
	}

	// test for head = 4, tail = 0, x = 2
	cb3 := NewCircularBuffer(5)
	for i := 0; i < 5; i++ {
		cb3.Push(fmt.Sprintf("%d", i))
	}
	for i := 0; i < 4; i++ {
		cb3.Pop()
	}
	if cb3.Length() != 1 {
		t.Errorf("Evacuate failed, want %d got %d", 1, cb3.Length())
	}
	cb3.Evacuate("2")
	if cb3.Length() != 1 {
		t.Errorf("Evacuate failed, want %d got %d", 1, cb3.Length())
	}
	if cb3.buffer[2] != "F" {
		t.Errorf("Evacuate failed, want %s got %s", "F", cb3.buffer[2])
	}

	// now tests for itemPos contains in [head:tail], these require actual evacuation
	cb4 := NewCircularBuffer(5)
	for i := 0; i < 4; i++ {
		cb4.Push(fmt.Sprintf("%d", i))
	}
	cb4.Pop()
	cb4.Evacuate("3")
	if cb4.buffer[1] != "F" {
		t.Errorf("Evacuate failed, want %s, %s, %s got %s, %s, %s",
			"F", "1", "2",
			cb4.buffer[1], cb4.buffer[2], cb4.buffer[3])
	}
	if cb4.Length() != 2 {
		t.Errorf("Evacuate failed, want %d got %d", 2, cb4.Length())
	}

	cb5 := NewCircularBuffer(5)
	for i := 0; i < 5; i++ {
		cb5.Push(fmt.Sprintf("%d", i))

	}
	for i := 0; i < 4; i++ {
		cb5.Pop()
	}
	for i := 0; i < 3; i++ {
		cb5.Push(fmt.Sprintf("%d", i))
	}
	cb5.Evacuate("0")
	if cb5.Length() != 3 {
		t.Errorf("Evacuate failed, want %d got %d", 3, cb5.Length())
	}
	if cb5.buffer[0] != "1" || cb5.buffer[1] != "2" || cb5.buffer[2] != "3" || cb5.buffer[3] != "F" || cb5.buffer[4] != "4" {
		t.Errorf("Evacuate failed, want %s, %s,%s,%s, %s, got %s, %s, %s, %s, %s",
			"1", "2", "3", "F", "4",
			cb5.buffer[0], cb5.buffer[1], cb5.buffer[2], cb5.buffer[3], cb5.buffer[4])
	}
}
