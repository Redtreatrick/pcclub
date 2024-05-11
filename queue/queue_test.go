package queue

import (
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
