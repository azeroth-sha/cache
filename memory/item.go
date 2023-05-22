package memory

import "time"

type item struct {
	v interface{}
	e *time.Time
}

func (i *item) Expired() bool {
	if i.e == nil {
		return false
	}
	return time.Now().After(*i.e)
}
