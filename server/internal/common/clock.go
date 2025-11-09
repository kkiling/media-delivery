package common

import "time"

// RealClock структура работающая со веременем
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}
