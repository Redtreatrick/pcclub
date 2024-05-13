package queue

type CircularBuffer struct {
	buffer []string
	size   int
	head   int
	tail   int
	count  int
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buffer: make([]string, size),
		size:   size,
		head:   0,
		tail:   0,
		count:  0, // counter to check if its full or empty
	}
}
func (c *CircularBuffer) Push(data string) bool {
	if c.count == c.size {
		// required by task
		return false
	} else {
		c.count++
	}
	c.buffer[c.tail] = data
	c.tail = (c.tail + 1) % c.size
	return true
}

func (c *CircularBuffer) Pop() (string, bool) {
	if c.count == 0 {
		return "", false
	}

	oldHead := c.head

	c.head = (c.head + 1) % c.size

	c.count--

	return c.buffer[oldHead], true
}

func (c *CircularBuffer) Contains(username string) bool {
	if c.count == 0 {
		return false
	}

	if c.count == c.size {
		for _, v := range c.buffer {
			if v == username {
				return true
			}
		}
	}
	if c.tail > c.head {
		for _, v := range c.buffer[c.head:c.tail] {
			if v == username {
				return true
			}
		}
	}

	if c.tail < c.head {
		for _, v := range c.buffer[c.head : len(c.buffer)-c.tail] {
			if v == username {
				return true
			}
		}
	}

	return false
}

// Find returns position of an element in buffer. If not found returns -1
func (c *CircularBuffer) Find(username string) int {
	if c.count == 0 {
		return -1
	}

	for i, v := range c.buffer {
		if v == username {
			return i
		}
	}

	return -1
}

// Evacuate solely exists to solve the problem of client leaving from waiting queue
func (c *CircularBuffer) Evacuate(username string) {
	//if c.tail > c.head
	itemPos := c.Find(username)
	if itemPos == -1 {
		return
	}

	if c.head == itemPos {
		c.buffer[itemPos] = "F"
		c.Pop()
		return
	}

	if c.head > c.tail && itemPos < c.head && itemPos >= c.tail {
		c.buffer[itemPos] = "F"
		return
	}
	if c.tail > c.head && itemPos > c.tail || c.tail > c.head && itemPos < c.head {
		c.buffer[itemPos] = "F"
		return
	}

	if c.tail >= c.head && itemPos > c.head && itemPos <= c.tail {
		for i := itemPos; i > c.head; i-- {
			c.buffer[i] = c.buffer[(i - 1)]
		}
		c.buffer[c.head] = "F"
		c.head++
		c.count--
		return
	}

	if c.head > c.tail && itemPos < c.tail {
		for i := itemPos; i < c.tail; i++ {
			c.buffer[i] = c.buffer[(i + 1)]
		}
		c.buffer[c.tail] = "F"
		c.tail--
		c.count--
		return
	}

	panic("case not covered")
}

func (c *CircularBuffer) Full() bool {
	return c.count == c.size
}

func (c *CircularBuffer) Empty() bool {
	return c.count == 0
}

func (c *CircularBuffer) Count() int {
	return c.count
}

func (c *CircularBuffer) Peek() string {
	return c.buffer[c.head]
}

func (c *CircularBuffer) Length() int {
	if c.count == 0 {
		return 0
	}

	if c.count == c.size {
		return c.size
	}

	if c.tail > c.head {
		return c.tail - c.head
	}

	return c.tail + c.size - c.head
}
