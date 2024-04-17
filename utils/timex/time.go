package timex

import "time"

// Mock time.Now() function
// This function is used to mock time.Now() function in the test.

var tn *time.Time

func init() {
	tn = nil
}

func SetTestTime(t time.Time) {
	tn = &t
}

func RemoveTestTime() {
	tn = nil
}

func Now() time.Time {
	if tn != nil {
		return *tn
	}

	return time.Now()
}
