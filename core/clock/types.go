package clock

import "time"

type Clock interface {
	LoadLocation(name string) (*time.Location, error)
	Now() time.Time
	Sleep(d time.Duration)
}
