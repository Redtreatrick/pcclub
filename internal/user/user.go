package user

import (
	"regexp"
)

type Name string

type Data struct {
	Table int
	//TimeMinutes time.Minutes
}

func (n Name) Ok() bool {
	pattern := `^[a-z0-9_-]+$`
	// Compile the regular expression.
	re := regexp.MustCompile(pattern)
	// Test the input string against the compiled regular expression.
	return re.MatchString(string(n))
}
