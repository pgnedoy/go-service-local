package clock

import "time"

type RealClock struct{}

func New() *RealClock {
	return &RealClock{}
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

func (c *RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (c *RealClock) LoadLocation(name string) (*time.Location, error) {
	return time.LoadLocation(name)
}
