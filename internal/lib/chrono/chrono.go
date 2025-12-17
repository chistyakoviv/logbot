package chrono

import "time"

type Chrono interface {
	Now() time.Time
}

type chrono struct{}

func New() Chrono {
	return &chrono{}
}

func (c *chrono) Now() time.Time {
	return time.Now().UTC()
}
