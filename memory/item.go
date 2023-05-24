package memory

import "time"

type item struct {
	v interface{}
	e *time.Time
}

func (i *item) Expired() bool {
	return i.e != nil && time.Now().After(*i.e)
}
