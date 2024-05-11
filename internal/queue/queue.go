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
