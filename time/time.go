package time

import (
	"fmt"
	"log"
	"strconv"
)

type Minutes int

func Atoi(timeStr string) Minutes {
	if len(timeStr) != 5 || timeStr[2] != ':' {
		log.Fatalf("time should have format HH:MM\nlen(time)=%d\ntimeStr[2]==':'%v", len(timeStr), timeStr[2] == ':')
	}

	hours, err := strconv.Atoi(timeStr[:2])
	if err != nil || hours < 0 || hours > 23 {
		log.Fatalf("hours should be integer from 0 to 23: %v", err)
	}

	minutes, err := strconv.Atoi(timeStr[3:5])
	if err != nil || minutes < 0 || minutes > 59 {
		log.Fatalf("minutes should be integer from 0 to 59: %v", err)
	}

	return Minutes(hours*60 + minutes)
}

func (t Minutes) String() string {
	return fmt.Sprintf("%02d:%02d", t/60, t%60)
}
